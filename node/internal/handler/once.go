package handler

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
)

// Execute the job immediately
// value
// If a single node is executed, the value is nodeUUID
// If the node where the job is located needs to be executed, the value is null ""
func WatchOnce() clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdOnceProfile, clientv3.WithPrefix())
}
