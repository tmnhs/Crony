package service

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/node/internal/handler"
)

func (srv *NodeServer) watchJobs() {
	rch := handler.WatchJobs(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate():
				job, err := handler.GetJobFromKv(ev.Kv.Key, ev.Kv.Value)
				if err != nil {
					logger.GetLogger().Warn(fmt.Sprintf("watch job err: %s, kv: %s", err.Error(), ev.Kv.String()))
					continue
				}
				srv.jobs[job.ID] = job
				job.InitNodeInfo(srv.UUID, srv.Hostname, srv.IP)
				srv.addJob(job)
			case ev.IsModify():
				job, err := handler.GetJobFromKv(ev.Kv.Key, ev.Kv.Value)
				if err != nil {
					logger.GetLogger().Warn(fmt.Sprintf("watch job err: %s, kv: %s", err.Error(), ev.Kv.String()))
					continue
				}
				job.InitNodeInfo(srv.UUID, srv.Hostname, srv.IP)
				srv.modifyJob(job)
			case ev.Type == mvccpb.DELETE:
				srv.deleteJob(handler.GetJobIDFromKey(string(ev.Kv.Key)))
			default:
				logger.GetLogger().Warn(fmt.Sprintf("watch job unknown event type[%v] from job[%s]", ev.Type, string(ev.Kv.Key)))
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
					logger.GetLogger().Error(fmt.Sprintf("watch killed proc error:%s kv:%s", err.Error(), ev.Kv.String()))
					continue
				}
				procVal := &models.JobProcVal{}
				err = json.Unmarshal(ev.Kv.Value, procVal)
				if err != nil {
					logger.GetLogger().Warn(fmt.Sprintf("watch killed proc json warn:%s kv:%s", err.Error(), ev.Kv.String()))
					continue
				}
				proc.JobProcVal = *procVal
				if proc.Killed {
					//fixme this is only used for linux/mac/unix system
					/*if err := syscall.Kill(-proc.ID, syscall.SIGKILL); err != nil {
						logger.Errorf("process:[%d] force kill failed, error:[%s]\n", proc.ID, err)
						return
					}*/

				}
			}
		}
	}
}

func (srv *NodeServer) watchOnce() {
	rch := handler.WatchOnce()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsModify(), ev.IsCreate():
				//不是在该node节点执行
				if len(ev.Kv.Value) != 0 && string(ev.Kv.Value) != srv.UUID {
					continue
				}
				//
				j, ok := srv.jobs[handler.GetJobIDFromKey(string(ev.Kv.Key))]
				if !ok {
					continue
				}
				//todo 判断是否正在运行，如果正在运行，则continue

				//立即执行
				go j.RunWithRecovery()
			}
		}
	}
}
