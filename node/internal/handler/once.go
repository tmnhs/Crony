package handler

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
)

// 马上执行 job 任务
// value
// 若执行单个结点，则值为 nodeUUID
// 若 job 所在的结点都需执行，则值为空 ""

func WatchOnce() clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdOnceProfile, clientv3.WithPrefix())
}
