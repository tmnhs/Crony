package service

import (
	"github.com/ouqiang/goutil"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/node/internal/group"
	"github.com/tmnhs/crony/node/internal/job"
)

func (srv *NodeServer) Run() (err error) {
	//todo defer

	if err = srv.loadJobs(); err != nil {
		return
	}
	//start cron
	srv.Cron.Start()
	//todo watchJobs
	//go n.watchJobs()
	//todo watchExcutingProc
	//go n.watchExcutingProc()
	//todo watchGroups
	//go n.watchGroups()
	//todo watchOnce
	//go n.watchOnce()
	// todo node into mysql
	//n.Node.On()
	return
}

func (srv *NodeServer) loadJobs() (err error) {
	//先获取所有的分组
	if srv.Groups, err = group.GetGroups(""); err != nil {
		return
	}
	//再获取本机分配的定时任务
	jobs, err := job.GetJobs(srv.ID)
	if err != nil {
		return
	}

	if len(jobs) == 0 {
		return
	}
	srv.jobs = jobs

	for _, j := range jobs {
		j.InitNodeInfo(srv.ID, srv.Hostname, srv.IP)
		srv.addJob(j, false)
	}

	return
}

//todo notice
func (srv *NodeServer) addJob(j *job.Job, notice bool) {
	taskFunc := job.CreateJob(j)
	if taskFunc == nil {
		logger.Errorf("创建任务处理Job失败,不支持的任务协议#%s", j.JobType)
		return
	}
	err := goutil.PanicToError(func() {
		srv.Cron.AddFunc(j.Spec, taskFunc)
	})
	if err != nil {
		logger.Errorf("添加任务到调度器失败#%v", err.Error())
	}

	return
}
