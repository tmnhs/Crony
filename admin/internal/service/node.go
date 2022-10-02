package service

import (
	"context"
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
	client     *etcdclient.Client
	serverList map[string]string
	lock       sync.Mutex
}

var DefaultNodeWatcher = NewNodeWatcherService()

func NewNodeWatcherService() *NodeWatcherService {
	return &NodeWatcherService{
		client:     etcdclient.GetEtcdClient(),
		serverList: make(map[string]string),
	}
}

func (n *NodeWatcherService) Watch() error {
	resp, err := n.client.Get(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	_ = n.extractAddrs(resp)

	go n.watcher()
	return nil
}

func (n *NodeWatcherService) watcher() {
	rch := n.client.Watch(context.Background(), etcdclient.KeyEtcdNode, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				n.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				uuid := n.GetUUID(string(ev.Kv.Key))
				logger.GetLogger().Warn(fmt.Sprintf("crony node[%s] DELETE event detected", uuid))
				nodeModel := models.Node{UUID: uuid}
				err := nodeModel.FindById()
				if err != nil {
					logger.GetLogger().Warn(fmt.Sprintf("failed to fetch node[%s] from db: %s", uuid, err.Error()))
					continue
				}
				//todo notice
				/*if node.Alived {
					n.Send(&Message{
						Subject: fmt.Sprintf("[Cronsun Warning] Node[%s] break away cluster at %s",
							node.Hostname, time.Now().Format(time.RFC3339)),
						Body: fmt.Sprintf("Cronsun Node breaked away cluster, this might happened when node crash or network problems.\nUUID: %s\nHostname: %s\nIP: %s\n", id, node.Hostname, node.IP),
						To:   conf.Config.Mail.To,
					})
				}*/
				n.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//todo 是否需要
func (n *NodeWatcherService) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			n.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (n *NodeWatcherService) SetServiceList(key, val string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.serverList[key] = val
	logger.GetLogger().Debug(fmt.Sprintf("set data key : %s val:%s", key, val))
}

func (n *NodeWatcherService) DelServiceList(key string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.serverList, key)
	logger.GetLogger().Debug(fmt.Sprintf("del data key: %s", key))
}

func (n *NodeWatcherService) SerList2Array() []string {
	n.lock.Lock()
	defer n.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range n.serverList {
		addrs = append(addrs, v)
	}
	return addrs
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
