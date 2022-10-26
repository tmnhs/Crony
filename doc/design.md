## 1. Etcd中各key值的作用

|                  key                   |   value   |          备注           |
| :------------------------------------: | :-------: | :-------------------: |
|        /crony/node/<node_uuid>         |   node    | 监听节点node的状态,用于服务注册与发现 |
|    /crony/job/<node_uuid>/<job_id>     |    job    |   监听node节点需要执行的job    |
|          /crony/once/<job_id>          | node_uuid |      监听job是否立即执行      |
| /crony/proc/<node_uuid>/<job_id>/<pid> |  procVal  |     监听正在执行的任务job      |
|    /crony/system/<node_uuid>/switch    |  "alive"  |    监听是否需要获取节点服务器状态    |
|     /crony/system/<node_uuid>/get      | nodeInfo  |    获取节点node的服务器状态     |

## 2. 服务注册与发现

> 使用etcd实现服务注册与发现，key值为/crony/node/<node_uuid>

#### 2.1 node实现服务注册

```go
type ServerReg struct {
	Client        *Client
	stop          chan error
	leaseId       clientv3.LeaseID
	cancelFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	//time-to-live
	Ttl int64
}

func NewServerReg(ttl int64) *ServerReg {
	return &ServerReg{
		Client: _defalutEtcd,
		Ttl:    ttl,
		stop:   make(chan error),
	}
}

func (s *ServerReg) Register(key string, value string) error {
	if err := s.setLease(s.Ttl); err != nil {
		return err
	}
	go s.keepAlive()
	if err := s.putService(key, value); err != nil {
		return err
	}
	return nil
}

func (s *ServerReg) setLease(ttl int64) error {
	leaseResp, err := Grant(ttl)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := s.Client.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}
	s.leaseId = leaseResp.ID
	s.cancelFunc = cancelFunc
	s.keepAliveChan = leaseRespChan
	return nil
}
func (s *ServerReg) Stop() {
	s.stop <- nil
}

//Monitor the lease renewal
func (s *ServerReg) keepAlive() {
	for {
		select {
		case <-s.stop:
			return
		case leaseKeepResp := <-s.keepAliveChan:
			if leaseKeepResp == nil {
				logger.GetLogger().Info("the lease renewal function has been turned off\n")
				return
			}
		}
	}
}

func (s *ServerReg) putService(key, val string) error {
	kv := clientv3.NewKV(s.Client.Client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(s.leaseId))
	return err
}

func (s *ServerReg) RevokeLease() error {
	s.cancelFunc()
	time.Sleep(2 * time.Second)
	_, err := Revoke(s.leaseId)
	return err
}
```

#### 2.2 admin实现服务发现

```go
func (n *NodeWatcherService) Watch() error {
	resp, err := n.client.Get(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	_ = n.extractNodes(resp)

	go n.watcher()
	return nil
}

func (n *NodeWatcherService) watcher() {
	rch := n.client.Watch(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				n.setNodeList(n.GetUUID(string(ev.Kv.Key)), string(ev.Kv.Value))
			case mvccpb.DELETE:
				uuid := n.GetUUID(string(ev.Kv.Key))
				n.delNodeList(uuid)
				node := &models.Node{UUID: uuid}
				err := node.FindByUUID()
				if err != nil {
					return
				}
				//FailOver 故障转移
				success, fail, err := n.FailOver(uuid)
				if err != nil {
					return
				}
				// if the failover is all successful, delete the node in the database
				if fail.Count() == 0 {
					err = node.Delete()
					if err != nil {
					}
				}
				//Node inactivation information defaults to email.
				msg := &notify.Message{
					Type:      notify.NotifyTypeMail,
					IP:        fmt.Sprintf("%s:%s", node.IP, node.PID),
					Subject:   "节点失活报警",
					Body:      fmt.Sprintf("[Crony Warning]crony node[%s] in the cluster has failed,，fail over success count:%d jobID are :%s ,fail count:%d jobID are :%s ", uuid, success.Count(), success.String(), fail.Count(), fail.String()),
					To:        config.GetConfigModels().Email.To,
					OccurTime: time.Now().Format(utils.TimeFormatSecond),
				}
				//send 发送告警通知
				go notify.Send(msg)
			}
		}
	}
}
```

## 3. 任务自动分配

> 支持http回调方法自动分配,shell任务可通过相应的预设脚本来配置服务器的环境，从而实现自动分配
>
> 默认分配给任务数量最少的节点

```go
//Give priority to the node with the least number of tasks
func (j *JobService) AutoAllocateNode() string {
	//Get all the living nodes
	nodeList := DefaultNodeWatcher.List2Array()
	resultCount, resultNodeUUID := MaxJobCount, ""
	for _, nodeUUID := range nodeList {
		//Check the database to see if it is alive
		node := &models.Node{UUID: nodeUUID}
		err := node.FindByUUID()
		if err != nil {
			continue
		}
		if node.Status == models.NodeConnFail {
			//The node has failed
			delete(DefaultNodeWatcher.nodeList, nodeUUID)
			continue
		}
		count, err := DefaultNodeWatcher.GetJobCount(nodeUUID)
		if err != nil {
			continue
		}
		if resultCount > count {
			resultCount, resultNodeUUID = count, nodeUUID
		}
	}
	return resultNodeUUID
}
```



## 4. 故障转移

> 支持http回调方法故障转移,shell任务可通过相应的预设脚本来配置服务器的环境，从而实现故障转移

```go

func (n *NodeWatcherService) FailOver(nodeUUID string) (success Result, fail Result, err error) {
	jobs, err := n.GetJobs(nodeUUID)
	if err != nil {
		return
	}
	if len(jobs) == 0 {
		return
	}
	for _, job := range jobs {
		//Determine whether shell command failover is supported
		if job.Type == models.JobTypeCmd && !config.GetConfigModels().System.CmdAutoAllocation {
			fail = append(fail, job.ID)
			continue
		}
		oldUUID := job.RunOn
		autoUUID := DefaultJobService.AutoAllocateNode()
		if autoUUID == "" {
			fail = append(fail, job.ID)
			continue
		}
		err = n.assignJob(autoUUID, &job)
		if err != nil {
			fail = append(fail, job.ID)
			continue
		}
		//Delete the key value if the transfer is successful
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, oldUUID, job.ID))
		if err != nil {
			fail = append(fail, job.ID)
			continue
		}
		success = append(success, job.ID)
	}
	return
}
```

## 5. 告警通知

> 支持邮件和webhook告警，提供邮件模板和飞书告警模板

#### 5.1 Email

```go
var mailTemplate = `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <title></title>
    <meta charset="utf-8"/>

</head>
<body>
<div class="cap" style="
            border: 2px solid black;
            background-color: whitesmoke;
            height: 500px"
>
    <div class="content" style="
            background-color: white;
            background-clip: padding-box;
            color:black;
            font-size: 13px;
            padding: 25px 25px;
            margin: 25px 25px;
        ">
        <div class="hello" style="text-align: center; color: #FF3333;font-size: 18px;font-weight: bolder">
            {{.Subject}}
        </div>
        <br>
        <div>
            <table border="1"  bordercolor="black" cellspacing="0px" cellpadding="4px" style="margin: 0 auto;">
                <tr >
                    <td>告警主机</td>
                    <td >{{.IP}}</td>
                </tr>

                <tr>
                    <td>告警时间</td>
                    <td>{{.OccurTime}}</td>
                </tr>

                <tr>
                    <td>告警信息</td>
                    <td>{{.Body}}</td>
                </tr>

            </table>
        </div>
        <br><br>
    </div>
</div>
<br>

</body>
</html>
`
func (mail *Mail) SendMsg(msg *Message) {
	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(_defaultMail.From, _defaultMail.Nickname)) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	msgData := parseMailTemplate(msg)
	m.SetBody("text/html", msgData)

	d := gomail.NewDialer(_defaultMail.Host, _defaultMail.Port, _defaultMail.From, _defaultMail.Secret)
	if err := d.DialAndSend(m); err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("smtp send msg[%+v] err: %s", msg, err.Error()))
	}
}
```

#### 5.2 WebHook

```go

var feiShuTemplateCard = `{
  "msg_type": "interactive",
  "card": {
    "config": {
      "wide_screen_mode": true
    },
    "header": {
      "title": {
        "tag": "plain_text",
        "content": "subjectSlot - Crony定时任务平台报警" 
      },
      "template": "red" 
    },
    "elements": [
				{
			  "fields": [
				{
				  "is_short": true,
				  "text": {
					"content": "**🕐 时间：**\ntimeSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": true,
				  "text": {
					"content": "**📋报警主机：**\nipSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": true,
				  "text": {
					"content": "**👤 值班：**\nuserSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": false,
				  "text": {
					"content": "**报警信息:**\nmsgSlot",
					"tag": "lark_md"
				  }
				}
			  ],
			  "tag": "div"
			},
			{
			  "actions": [
				{
				  "tag": "button",
				  "text": {
					"content": "跟进处理",
					"tag": "plain_text"
				  },
				  "type": "primary",
				  "value": {
					"key1": "https://cloud.tencent.com/developer/article/1467743"
				  },
					"url":"https://cloud.tencent.com/developer/article/1467743"
				}
			  ],
			  "tag": "action"
			},
			{
			  "tag": "hr"
			}
    ]
  }
}`

func (w *WebHook) SendMsg(msg *Message) {
	switch _defaultWebHook.Kind {
	case "feishu":
		var sendData = feiShuTemplateCard
		sendData = strings.Replace(sendData, "timeSlot", msg.OccurTime, 1)
		sendData = strings.Replace(sendData, "ipSlot", msg.IP, 1)

		userSlot := ""
		for _, to := range msg.To {
			userSlot += fmt.Sprintf("<at email='' >%s</at>", to)
		}
		sendData = strings.Replace(sendData, "userSlot", userSlot, 1)
		sendData = strings.Replace(sendData, "msgSlot", msg.Body, 1)
		sendData = strings.Replace(sendData, "subjectSlot", msg.Subject, 1)
		_, err := httpclient.PostJson(_defaultWebHook.Url, sendData, 0)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("feishu  send msg[%+v] err: %s", msg, err.Error()))
		}
	default:
		b, err := json.Marshal(msg)
		if err != nil {
			return
		}
		_, err = httpclient.PostJson(_defaultWebHook.Url, string(b), 0)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("web hook api send msg[%+v] err: %s", msg, err.Error()))
		}
	}
	return
}

```

#### 5.3 运行

```go
func Serve() {
   for {
      select {
      case msg := <-msgQueue:
         if msg == nil {

         }
         switch msg.Type {
         case 1:
            //Mail
            msg.Check()
            _defaultMail.SendMsg(msg)
         case 2:
            //webhook
            msg.Check()
            go _defaultWebHook.SendMsg(msg)
         }
      }
   }
}


```