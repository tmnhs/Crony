package service

import (
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/node/internal/handler"
)

//fixme 不进行数据库操作
func (srv *NodeServer) watchJobs() {
	rch := handler.WatchJobs(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate():
				job, err := handler.GetJobFromKv(ev.Kv.Key, ev.Kv.Value)
				if err != nil {
					logger.Warnf("err: %s, kv: %s", err.Error(), ev.Kv.String())
					continue
				}
				//todo insert into srv.jobs
				srv.jobs[job.ID] = job
				job.InitNodeInfo(srv.UUID, srv.Hostname, srv.IP)
				srv.addJob(job, true)
			case ev.IsModify():
				job, err := handler.GetJobFromKv(ev.Kv.Key, ev.Kv.Value)
				if err != nil {
					logger.Warnf("err: %s, kv: %s", err.Error(), ev.Kv.String())
					continue
				}
				job.InitNodeInfo(srv.UUID, srv.Hostname, srv.IP)
				//todo modify job
				srv.modifyJob(job)
			case ev.Type == mvccpb.DELETE:
				//todo delete job
				srv.deleteJob(handler.GetJobIDFromKey(string(ev.Kv.Key)))
			default:
				logger.Warnf("unknown event type[%v] from job[%s]", ev.Type, string(ev.Kv.Key))
			}
		}
	}
}

func (srv *NodeServer) watchKilledProc() {
	rch := handler.WatchProc(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			//监控是否被修改
			case ev.IsModify():
				proc, err := handler.GetProcFromKey(string(ev.Kv.Key))
				if err != nil {
					logger.Errorf("killed proc error:%s kv:%s", err.Error(), ev.Kv.String())
					continue
				}
				procVal := &models.JobProcVal{}
				err = json.Unmarshal(ev.Kv.Value, procVal)
				if err != nil {
					logger.Warnf("killed proc json warn:%s kv:%s", err.Error(), ev.Kv.String())
					continue
				}
				proc.JobProcVal = *procVal
				if proc.Killed {
					//todo killed
				}
			}
		}
	}
}

func (srv *NodeServer) watchOnce() {

}
