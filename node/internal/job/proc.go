package job

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"strings"
	"sync/atomic"
)

//继承
// 当前执行中的任务信息
// key: /cronsun/proc/node/group/jobId/pid
// value: 开始执行时间
//todo key 会自动过期，防止进程意外退出后没有清除相关 key，过期时间可配置
type JobProc struct {
	*models.JobProc
	*etcdclient.ServerReg
}

func GetProcFromKey(key string) (proc *JobProc, err error) {
	ss := strings.Split(key, "/")
	var sslen = len(ss)
	if sslen < 5 {
		err = fmt.Errorf("invalid proc key [%s]", key)
		return
	}
	proc = &JobProc{
		JobProc: &models.JobProc{
			ID:     ss[sslen-1],
			JobID:  ss[sslen-2],
			Group:  ss[sslen-3],
			NodeID: ss[sslen-4],
		},
	}
	return
}

func (p *JobProc) Key() string {
	return etcdclient.KeyEtcdProc + p.NodeID + "/" + p.Group + "/" + p.JobID + "/" + p.ID
}

func (p *JobProc) Val() (string, error) {
	b, err := json.Marshal(&p.JobProcVal)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 获取节点正在执行任务的数量
func (j *Job) CountRunning() (int64, error) {
	resp, err := etcdclient.Get(etcdclient.KeyEtcdProc+j.RunOn+"/"+j.Group+"/"+j.ID, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}
func (p *JobProc) del() error {
	if atomic.LoadInt32(&p.HasPut) != 1 {
		return nil
	}

	_, err := etcdclient.Delete(p.Key())
	return err
}

func (p *JobProc) Stop() {
	if p == nil {
		return
	}
	if !atomic.CompareAndSwapInt32(&p.Runnig, 1, 0) {
		return
	}
	p.Wg.Wait()

	if err := p.del(); err != nil {
		logger.Warnf("proc del[%s] err: %s", p.Key(), err.Error())
	}
}

func WatchProc(nid string) clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdProc+nid, clientv3.WithPrefix())
}

//todo 注册
func (p *JobProc) Start() error {
	//todo new regserver in node server
	if p == nil {
		return nil
	}

	if !atomic.CompareAndSwapInt32(&p.Runnig, 0, 1) {
		return nil
	}

	p.Wg.Add(1)
	//creates a new lease
	val, err := json.Marshal(p.JobProcVal)
	if err != nil {
		return err
	}
	if err := p.ServerReg.Register(p.Key(), string(val)); err != nil {
		return err
	}
	p.Wg.Done()
	return nil
}
