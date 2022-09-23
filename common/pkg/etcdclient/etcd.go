package etcdclient

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"time"
)

var _defalutEtcd *clientv3.Client

func Init(e models.Etcd) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: time.Duration(e.DialTimeout)*time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	_defalutEtcd=cli
}

func GetEtcdClient() *clientv3.Client {
	if _defalutEtcd==nil{
		logger.Errorf("mysql database is not initialized")
		return nil
	}
	return _defalutEtcd
}