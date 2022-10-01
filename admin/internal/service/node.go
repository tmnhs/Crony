package service

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"strings"
	"sync"
)

type NodeWatcher struct {
	client     *etcdclient.Client
	serverList map[string]string
	lock       sync.Mutex
}

var DefaultNodeWatcher = NewNodeWatcher()

func NewNodeWatcher() *NodeWatcher {
	return &NodeWatcher{
		client:     etcdclient.GetEtcdClient(),
		serverList: make(map[string]string),
	}
}

func (s *NodeWatcher) Watch() error {
	resp, err := s.client.Get(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	_ = s.extractAddrs(resp)

	go s.watcher()
	return nil
}

func (s *NodeWatcher) watcher() {
	rch := s.client.Watch(context.Background(), etcdclient.KeyEtcdNode, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				//todo insert or update
				/*node:=&models.Node{
					UUID:s.GetUUID(string(ev.Kv.Key)),
				}*/
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				fmt.Println("server delete")
				//todo
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (s *NodeWatcher) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			s.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (s *NodeWatcher) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	logger.GetLogger().Debug(fmt.Sprintf("set data key : %s val:%s", key, val))
}

func (s *NodeWatcher) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	logger.GetLogger().Debug(fmt.Sprintf("del data key: %s", key))
}

func (s *NodeWatcher) SerList2Array() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *NodeWatcher) Close() error {
	return nil
}

func (s *NodeWatcher) GetUUID(key string) string {
	// /crony/node/<node_uuid>
	index := strings.LastIndex(key, "/")
	if index == -1 {
		return ""
	}
	logger.GetLogger().Debug(fmt.Sprintf("key_index:%s key_index+1%s", key[index:], key[index+1:]))
	return key[index+1:]
}
