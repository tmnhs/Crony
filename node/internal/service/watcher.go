package service

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"github.com/tmnhs/crony/node/internal/handler"
	"strings"
)

func (srv *NodeServer) watchJobs() {
	rch := handler.WatchJobs(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate():
				var job handler.Job
				if err := json.Unmarshal(ev.Kv.Value, &job); err != nil {
					err = fmt.Errorf("watch job[%s] create json umarshal err: %s", string(ev.Kv.Key), err.Error())
					continue
				}
				srv.jobs[job.ID] = &job
				job.InitNodeInfo(models.JobStatusAssigned, srv.UUID, srv.Hostname, srv.IP)
				srv.addJob(&job)
			case ev.IsModify():
				var job handler.Job
				if err := json.Unmarshal(ev.Kv.Value, &job); err != nil {
					err = fmt.Errorf("watch job[%s] modify json umarshal err: %s", string(ev.Kv.Key), err.Error())
					continue
				}
				job.InitNodeInfo(models.JobStatusAssigned, srv.UUID, srv.Hostname, srv.IP)
				srv.modifyJob(&job)
			case ev.Type == mvccpb.DELETE:
				srv.deleteJob(handler.GetJobIDFromKey(string(ev.Kv.Key)))
			default:
				logger.GetLogger().Warn(fmt.Sprintf("watch job unknown event type[%v] from job[%s]", ev.Type, string(ev.Kv.Key)))
			}
		}
	}
}

//fixme kill executing job
/*
func (srv *NodeServer) watchKilledProc() {
	rch := handler.WatchProc(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
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
					if err := syscall.Kill(proc.ID, syscall.SIGKILL); err != nil {
						logger.GetLogger().Error(fmt.Sprintf("process:[%d] force kill failed, error:[%s]", proc.ID, err))
						return
					}
				}

			}
		}
	}
}*/

func (srv *NodeServer) watchSystemInfo() {
	rch := handler.WatchSystem(srv.UUID)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsCreate() || ev.IsModify():
				key := string(ev.Kv.Key)
				if string(ev.Kv.Value) != models.NodeSystemInfoSwitch || srv.Node.UUID != getUUID(key) {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] ,switch is not alive ", srv.UUID))
					continue
				}
				s, err := utils.GetServerInfo()
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] error: %s", srv.UUID, err.Error()))
					continue
				}
				b, err := json.Marshal(s)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] json marshal error: %s", srv.UUID, err.Error()))
					continue
				}
				_, err = etcdclient.PutWithTtl(fmt.Sprintf(etcdclient.KeyEtcdSystemGet, getUUID(key)), string(b), 5*60)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] etcd put error: %s", srv.UUID, err.Error()))
					continue
				}

			}
		}
	}
}

func getUUID(key string) string {
	// /crony/node/<node_uuid>
	index := strings.LastIndex(key, "/")
	if index == -1 {
		return ""
	}
	return key[index+1:]
}

func (srv *NodeServer) watchOnce() {
	rch := handler.WatchOnce()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch {
			case ev.IsModify(), ev.IsCreate():
				// is not executed on this node
				if len(ev.Kv.Value) != 0 && string(ev.Kv.Value) != srv.UUID {
					continue
				}
				j, ok := srv.jobs[handler.GetJobIDFromKey(string(ev.Kv.Key))]
				if !ok {
					continue
				}
				go j.RunWithRecovery()
			}
		}
	}
}
