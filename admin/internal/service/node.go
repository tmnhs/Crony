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
	client   *etcdclient.Client
	nodeList map[string]models.Node
	lock     sync.Mutex
}

var DefaultNodeWatcher *NodeWatcherService

func NewNodeWatcherService() *NodeWatcherService {
	return &NodeWatcherService{
		client:   etcdclient.GetEtcdClient(),
		nodeList: make(map[string]models.Node),
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
				n.setNodeList(n.GetUUID(string(ev.Kv.Key)), string(ev.Kv.Value))
			case mvccpb.DELETE:
				uuid := n.GetUUID(string(ev.Kv.Key))
				n.delNodeList(uuid)
				logger.GetLogger().Warn(fmt.Sprintf("crony node[%s] DELETE event detected", uuid))
				//todo 可能node的状态还是alive 需要更新
				node := &models.Node{UUID: uuid}
				err := node.FindByUUID()
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("crony node[%s] find by uuid  error:%s", uuid, err.Error()))
					return
				}

				success, fail, err := n.FailOver(uuid)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("crony node[%s] fail over error:%s", uuid, err.Error()))
					return
				}
				// if the failover is all successful, delete the node in the database
				if fail.Count() == 0 {
					err = node.Delete()
					if err != nil {
						logger.GetLogger().Error(fmt.Sprintf("crony node[%s] delete by uuid  error:%s", uuid, err.Error()))
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
			n.setNodeList(n.GetUUID(string(resp.Kvs[i].Key)), string(resp.Kvs[i].Value))
			nodes = append(nodes, string(v))
		}
	}
	return nodes
}

func (n *NodeWatcherService) setNodeList(key, val string) {
	var node models.Node
	err := json.Unmarshal([]byte(val), &node)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("discover node[%s] json error:%s", key, err.Error()))
		return
	}
	n.lock.Lock()
	n.nodeList[key] = node
	n.lock.Unlock()
	logger.GetLogger().Debug(fmt.Sprintf("discover node node[%s] with pid[%s]", key, val))
	//Wait for the node to be fully started and assign the node
	time.Sleep(5 * time.Second)
	//find unassigned job
	jobs, err := DefaultJobService.GetNotAssignedJob()
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("discover node[%s],pid[%s] and get not assigned job err:%s", key, val, err.Error()))
		return
	}
	for _, job := range jobs {
		if job.Type == models.JobTypeCmd && !config.GetConfigModels().System.CmdAutoAllocation {
			logger.GetLogger().Warn(fmt.Sprintf("assign unassigned job[%d]  don't support cmd type", job.ID))
			continue
		}
		oldUUID := job.RunOn
		nodeUUID := DefaultJobService.AutoAllocateNode()
		if nodeUUID == "" {
			//If automatic allocation fails, it will be directly assigned to the new node.
			nodeUUID = key
		}
		err = n.assignJob(nodeUUID, &job)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("assign unassigned job[%d]  error:%s", job.ID, err.Error()))
			continue
		}
		//Delete the key value if the transfer is successful
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, oldUUID, job.ID))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over etcd delete job error:%s", nodeUUID, job.ID, err.Error()))
			continue
		}
	}
}

func (n *NodeWatcherService) delNodeList(key string) {
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

func (n *NodeWatcherService) assignJob(nodeUUID string, job *models.Job) (err error) {
	if nodeUUID == "" {
		return fmt.Errorf("node uuid can't be null")
	}
	node, ok := n.nodeList[nodeUUID]
	if !ok {
		return fmt.Errorf("assign unassigned job[%d] but  node[%s] not exist ", job.ID, nodeUUID)
	}
	job.InitNodeInfo(models.JobStatusAssigned, node.UUID, node.Hostname, node.IP)

	b, err := json.Marshal(job)
	if err != nil {
		return
	}
	_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdJob, nodeUUID, job.ID), string(b))
	if err != nil {
		return
	}
	err = job.Update()
	if err != nil {
		return
	}
	return
}

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
		//Determine whether shell command failover is supported
		if job.Type == models.JobTypeCmd && !config.GetConfigModels().System.CmdAutoAllocation {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over don't support cmd type", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}
		oldUUID := job.RunOn
		autoUUID := DefaultJobService.AutoAllocateNode()
		if autoUUID == "" {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over auto allocate node error", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}
		err = n.assignJob(autoUUID, &job)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over assign job error", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}
		//Delete the key value if the transfer is successful
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

//get all the job under a node
func (n *NodeWatcherService) GetJobs(nodeUUID string) (jobs []models.Job, err error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
	if err != nil {
		return
	}
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
