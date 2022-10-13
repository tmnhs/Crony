package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/notify"
	"github.com/tmnhs/crony/common/pkg/utils"
	"strings"
	"sync"
	"time"
)

type NodeWatcherService struct {
	client *etcdclient.Client
	//<uuid> <pid>
	nodeList map[string]string
	lock     sync.Mutex
}

var DefaultNodeWatcher *NodeWatcherService

func NewNodeWatcherService() *NodeWatcherService {
	return &NodeWatcherService{
		client:   etcdclient.GetEtcdClient(),
		nodeList: make(map[string]string),
	}
}

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
				n.SetNodeList(n.GetUUID(string(ev.Kv.Key)), string(ev.Kv.Value))
			case mvccpb.DELETE:
				uuid := n.GetUUID(string(ev.Kv.Key))
				logger.GetLogger().Warn(fmt.Sprintf("crony node[%s] DELETE event detected", uuid))
				node := &models.Node{UUID: uuid}
				err := node.FindByUUID()
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("crony node[%s] find by uuid  error:%s", uuid, err.Error()))
					return
				}
				// 先删除再故障转移
				n.DelNodeList(n.GetUUID(string(ev.Kv.Key)))
				success, fail, err := n.FailOver(uuid)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("crony node[%s] fail over error:%s", uuid, err.Error()))
					return
				}
				//如果故障转移全部成功则在数据库里删除节点
				if fail.Count() == 0 {
					err = node.Delete()
					if err != nil {
						logger.GetLogger().Error(fmt.Sprintf("crony node[%s] delete by uuid  error:%s", uuid, err.Error()))
					}
				}
				//节点失活信息默认使用邮件
				msg := &notify.Message{
					Type:      notify.NotifyTypeMail,
					IP:        fmt.Sprintf("%s:%s", node.IP, node.PID),
					Subject:   "节点失活报警",
					Body:      fmt.Sprintf("[Crony Warning]crony node[%s] in the cluster has failed,，fail over success count:%d jobID are :%s ,fail count:%d jobID are :%s ", uuid, success.Count(), success.String(), fail.Count(), fail.String()),
					To:        config.GetConfigModels().Email.To,
					OccurTime: time.Now().Format(utils.TimeFormatSecond),
				}

				go notify.Send(msg)
			}
		}
	}
}

func (n *NodeWatcherService) extractNodes(resp *clientv3.GetResponse) []string {
	nodes := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return nodes
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			n.SetNodeList(n.GetUUID(string(resp.Kvs[i].Key)), string(resp.Kvs[i].Value))
			nodes = append(nodes, string(v))
		}
	}
	return nodes
}

func (n *NodeWatcherService) SetNodeList(key, val string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nodeList[key] = val
	logger.GetLogger().Debug(fmt.Sprintf("discover node[%s],pid[%s]", key, val))
}

func (n *NodeWatcherService) DelNodeList(key string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.nodeList, key)
	logger.GetLogger().Debug(fmt.Sprintf("delelte node[%s]", key))
}

func (n *NodeWatcherService) List2Array() []string {
	n.lock.Lock()
	defer n.lock.Unlock()
	nodes := make([]string, 0)

	for k, _ := range n.nodeList {
		nodes = append(nodes, k)
	}
	return nodes
}

func (n *NodeWatcherService) Close() error {
	return nil
}

func (n *NodeWatcherService) GetUUID(key string) string {
	// /crony/node/<node_uuid>
	index := strings.LastIndex(key, "/")
	if index == -1 {
		return ""
	}
	return key[index+1:]
}

func (n *NodeWatcherService) Search(s *request.ReqNodeSearch) ([]models.Node, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyNodeTableName)
	if len(s.UUID) > 0 {
		db = db.Where("uuid = ?", s.UUID)
	}
	if len(s.IP) > 0 {
		db.Where("ip = ?", s.IP)
	}
	if s.Status > 0 {
		db.Where("status = ?", s.Status)
	}
	if s.UpTime > 0 {
		db.Where("up > ?", s.UpTime)
	}
	nodes := make([]models.Node, 2)
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Order("up desc").Find(&nodes).Error
	if err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}

func (n *NodeWatcherService) GetJobCount(nodeUUID string) (int, error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}
	return int(resps.Count), nil
}

type Result []int

func (r Result) Count() (count int) {
	for _, v := range r {
		if v != 0 {
			count++
		}
	}
	return
}
func (r Result) String() (str string) {
	str = "["
	for _, v := range r {
		if v != 0 {
			str += fmt.Sprintf("%d,", v)
		}
	}
	str += "]"
	return
}

//故障转移
func (n *NodeWatcherService) FailOver(nodeUUID string) (success Result, fail Result, err error) {
	jobs, err := n.GetJobs(nodeUUID)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node[%s] fail over get jobs error:%s", nodeUUID, err.Error()))
		return
	}
	if len(jobs) == 0 {
		return
	}
	for _, job := range jobs {
		//fixme
		/*if job.Type == models.JobTypeCmd {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over don't support cmd type", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}*/
		oldUUID := job.RunOn
		autoUUID := DefaultJobService.AutoAllocateNode()
		if autoUUID == "" {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over auto allocate node error", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}
		node := &models.Node{UUID: autoUUID}
		err = node.FindByUUID()
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over auto allocate node db find error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		job.InitNodeInfo(node.UUID, node.Hostname, node.IP)
		err = job.Update()
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over auto allocate node db update error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		b, err := json.Marshal(job)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over json marshal job error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}

		_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdJob, job.RunOn, job.ID), string(b))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over etcd put job error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		//如果转移成功则删除key值
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, oldUUID, job.ID))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over etcd delete job error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		success = append(success, job.ID)
	}
	return
}

//获取某个node下的所有job
func (n *NodeWatcherService) GetJobs(nodeUUID string) (jobs []models.Job, err error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
	if err != nil {
		return
	}
	jobs = make([]models.Job, 2)
	count := len(resps.Kvs)
	if count == 0 {
		return
	}
	for _, j := range resps.Kvs {
		var job models.Job
		if err := json.Unmarshal(j.Value, &job); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("job[%s] umarshal err: %s", string(j.Key), err.Error()))
			continue
		}
		jobs = append(jobs, job)
	}
	return
}

//获取某个node的数量
func (n *NodeWatcherService) GetNodeCount(status int) (int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyNodeTableName)
	if status > 0 {
		db = db.Where("status = ?", status)
	}
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetNodeSystemInfo(uuid string) (s *utils.Server, err error) {
	defer func() {
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdSystemSwitch, uuid))
	}()
	s = new(utils.Server)
	res, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdSystemGet, uuid), clientv3.WithPrefix())
	if err != nil || len(res.Kvs) == 0 {
		return
	}
	err = json.Unmarshal(res.Kvs[0].Value, s)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("json error:%v", err))
	}
	return
}
