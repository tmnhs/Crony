package service

import (
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/node/internal/handler"
)

func (srv *NodeServer) watchJobs() {
	rch := handler.WatchJobs(srv.ID)
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
				job.InitNodeInfo(srv.ID, srv.Hostname, srv.IP)
				srv.addJob(job, true)
			case ev.IsModify():
				job, err := handler.GetJobFromKv(ev.Kv.Key, ev.Kv.Value)
				if err != nil {
					logger.Warnf("err: %s, kv: %s", err.Error(), ev.Kv.String())
					continue
				}

				job.InitNodeInfo(srv.ID, srv.Hostname, srv.IP)
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
