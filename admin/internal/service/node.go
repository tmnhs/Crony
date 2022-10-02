package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"strings"
	"sync"
)

type NodeWatcherService struct {
	client *etcdclient.Client
	//<uuid> <pid>
	nodeList map[string]string
	lock     sync.Mutex
}

var DefaultNodeWatcher = NewNodeWatcherService()

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
	rch := n.client.Watch(context.Background(), etcdclient.KeyEtcdNode, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				n.SetNodeList(n.GetUUID(string(ev.Kv.Key)), string(ev.Kv.Value))
			case mvccpb.DELETE:
				uuid := n.GetUUID(string(ev.Kv.Key))
				logger.GetLogger().Warn(fmt.Sprintf("crony node[%s] DELETE event detected", uuid))
				//先删除在故障转移
				n.DelNodeList(n.GetUUID(string(ev.Kv.Key)))
				success, fail, err := n.FailOver(uuid)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("crony node[%s] fail over error:%s", uuid, err.Error()))
				}
				logger.GetLogger().Info(fmt.Sprintf("crony node[%s] fail over success count:%d they are :%#v  ,fail count:%d they are :%#v ", uuid, len(success), success, len(fail), fail))
				//todo notice
				// 故障转移
				/*if node.Alived {
					n.Send(&Message{
						Subject: fmt.Sprintf("[Cronsun Warning] Node[%s] break away cluster at %s",
							node.Hostname, time.Now().Format(time.RFC3339)),
						Body: fmt.Sprintf("Cronsun Node breaked away cluster, this might happened when node crash or network problems.\nUUID: %s\nHostname: %s\nIP: %s\n", id, node.Hostname, node.IP),
						To:   conf.Config.Mail.To,
					})
				}*/
			}
		}
	}
}

//todo 是否需要
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
	logger.GetLogger().Debug(fmt.Sprintf("set data key : %s val:%s", key, val))
}

func (n *NodeWatcherService) DelNodeList(key string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.nodeList, key)
	logger.GetLogger().Debug(fmt.Sprintf("del data key: %s", key))
}

func (n *NodeWatcherService) List2Array() []string {
	n.lock.Lock()
	defer n.lock.Unlock()
	nodes := make([]string, 0)

	for _, v := range n.nodeList {
		nodes = append(nodes, v)
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
	logger.GetLogger().Debug(fmt.Sprintf("key_index:%s key_index+1%s", key[index:], key[index+1:]))
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
	err := db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&nodes).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return nodes, total, nil
}

//todo 获取某节点的job数量
func (n *NodeWatcherService) GetJobCount(nodeUUID string) (int, error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return len(resps.Kvs), nil
}

func (n *NodeWatcherService) FindByGroupId(groupId int) ([]models.Node, error) {
	var nodes []models.Node
	sql := fmt.Sprintf("select n.* from %s ng join %s n on ng.group_id = ? and ng.node_uuid = n.uuid", models.CronyNodeGroupTableName, models.CronyNodeTableName)
	err := dbclient.GetMysqlDB().Raw(sql, groupId).Scan(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

//故障转移
func (n *NodeWatcherService) FailOver(nodeUUID string) (success []int, fail []int, err error) {
	jobs, err := n.GetJobs(nodeUUID)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node[%s] fail over get jobs error:%s", nodeUUID, err.Error()))
		return
	}
	if len(jobs) == 0 {
		return
	}
	success = make([]int, 2)
	fail = make([]int, 2)
	for _, job := range jobs {
		if job.Type == models.JobTypeCmd {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over don't support cmd type", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}
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

		_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdJob, job.RunOn, job.GroupId, job.ID), string(b))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over etcd put job error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		success = append(success, job.ID)
	}
	//todo 自动分配成功的有几个，失败的有几个
	return
}

//获取某个node下的所有job
func (n *NodeWatcherService) GetJobs(nodeUUID string) (jobs []*models.Job, err error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
	if err != nil {
		return
	}
	jobs = make([]*models.Job, 2)
	count := len(resps.Kvs)
	if count == 0 {
		return
	}
	for _, j := range resps.Kvs {
		job := new(models.Job)
		if e := json.Unmarshal(j.Value, job); e != nil {
			logger.GetLogger().Warn(fmt.Sprintf("job[%s] umarshal err: %s", string(j.Key), e.Error()))
			continue
		}
		jobs = append(jobs, job)
	}
	return
}
