package handler

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
)

// 马上执行 job 任务
// 注册到 /cronsun/once/group/<jobID>
// value
// 若执行单个结点，则值为 NodeID
// 若 job 所在的结点都需执行，则值为空 ""
func PutOnce(groupId, jobID int, nodeUUID string) error {
	_, err := etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdOnce, groupId, jobID), nodeUUID)
	return err
}

func WatchOnce() clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdOnceProfile, clientv3.WithPrefix())
}
