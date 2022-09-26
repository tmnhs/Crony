package etcdclient

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

//创建租约注册服务
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

//设置租约
func (s *ServerReg) setLease(ttl int64) error {
	//设置租约时间
	leaseResp, err := Grant(ttl)
	if err != nil {
		return err
	}

	//设置续租
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

//监听 续租情况
func (s *ServerReg) keepAlive() {
	for {
		select {
		case <-s.stop:
			return
		case leaseKeepResp := <-s.keepAliveChan:
			if leaseKeepResp == nil {
				fmt.Printf("the lease renewal function has been turned off\n")
				return
			} else {
				fmt.Printf("renew the lease successfully.\n")
			}
		}
	}
}

//通过租约 注册服务
func (s *ServerReg) putService(key, val string) error {
	kv := clientv3.NewKV(s.Client.Client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(s.leaseId))
	return err
}

//撤销租约
func (s *ServerReg) RevokeLease() error {
	s.cancelFunc()
	time.Sleep(2 * time.Second)
	_, err := Revoke(s.leaseId)
	return err
}
