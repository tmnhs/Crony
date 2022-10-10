package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"strconv"
	"strings"
	"sync/atomic"
)

//继承
// 当前执行中的任务信息
// key: /cronsun/proc/node/group/jobId/pid
// value: 开始执行时间
//key 会自动过期，防止进程意外退出后没有清除相关 key，过期时间可配置
type JobProc struct {
	*models.JobProc
}

func GetProcFromKey(key string) (proc *JobProc, err error) {
	ss := strings.Split(key, "/")
	var sslen = len(ss)
	if sslen < 5 {
		err = fmt.Errorf("invalid proc key [%s]", key)
		return
	}
	id, err := strconv.Atoi(ss[sslen-1])
	if err != nil {
		return
	}
	jobId, err := strconv.Atoi(ss[sslen-2])
	if err != nil {
		return
	}
	proc = &JobProc{
		JobProc: &models.JobProc{
			ID:       id,
			JobID:    jobId,
			NodeUUID: ss[sslen-3],
		},
	}
	return
}

func (p *JobProc) Key() string {
	return fmt.Sprintf(etcdclient.KeyEtcdProc, p.NodeUUID, p.JobID, p.ID)
}

func (p *JobProc) Val() (string, error) {
	b, err := json.Marshal(&p.JobProcVal)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
		logger.GetLogger().Warn(fmt.Sprintf("proc del[%s] err: %s", p.Key(), err.Error()))
	}
}

func WatchProc(nodeUUID string) clientv3.WatchChan {
	return etcdclient.Watch(fmt.Sprintf(etcdclient.KeyEtcdProcProfile, nodeUUID), clientv3.WithPrefix())
}

func (p *JobProc) Start() error {
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
	_, err = etcdclient.PutWithTtl(p.Key(), string(val), config.GetConfigModels().System.JobProcTtl)
	if err != nil {
		return err
	}
	p.Wg.Done()
	return nil
}
