package service

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"github.com/ouqiang/goutil"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"github.com/tmnhs/crony/node/internal/handler"
	"os"
	"strconv"
	"syscall"
)

// Node 执行 cron 命令服务的结构体
type NodeServer struct {
	*etcdclient.ServerReg
	*models.Node
	*cron.Cron

	jobs handler.Jobs // 和结点相关的任务
	//fixme
	//Groups handler.Groups

	// 删除的 job id，用于 group 更新
	delIDs map[string]bool
}

func NewNodeServer() (*NodeServer, error) {
	uuid, err := utils.UUID()
	if err != nil {
		return nil, err
	}
	ip, err := utils.LocalIP()
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = uuid
		err = nil
	}
	return &NodeServer{
		Node: &models.Node{
			UUID:     uuid,
			PID:      strconv.Itoa(os.Getpid()),
			IP:       ip.String(),
			Hostname: hostname,
		},
		Cron: cron.New(),

		jobs: make(handler.Jobs, 8),

		delIDs: make(map[string]bool, 8),

		ServerReg: etcdclient.NewServerReg(config.GetConfigModels().System.NodeTtl),
	}, nil

}

// Check whether the node is registered with ETCD
// If yes, PID is returned. If no, -1 is returned
func (srv *NodeServer) exist(nodeUUID string) (pid int, err error) {
	resp, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdNode, nodeUUID))
	if err != nil {
		return
	}

	if len(resp.Kvs) == 0 {
		return -1, nil
	}

	if pid, err = strconv.Atoi(string(resp.Kvs[0].Value)); err != nil {
		if _, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdNode, nodeUUID)); err != nil {
			return
		}
		return -1, nil
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return -1, nil
	}

	// TODO: 暂时不考虑 linux/unix 以外的系统
	if p != nil && p.Signal(syscall.Signal(0)) == nil {
		return
	}
	return -1, nil
}

// Register into ETCD with /crony/node/<node_id>
func (srv *NodeServer) Register() error {
	pid, err := srv.exist(srv.UUID)
	if err != nil {
		return err
	}
	if pid != -1 {
		return fmt.Errorf("node[%s] with pid[%d] exist", srv.UUID, pid)
	}
	//creates a new lease
	if err := srv.ServerReg.Register(fmt.Sprintf(etcdclient.KeyEtcdNode, srv.UUID), srv.PID); err != nil {
		return err
	}
	return nil
}

// 停止服务
func (srv *NodeServer) Stop(i interface{}) {
	etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdNode, srv.UUID))
	srv.Down()

	srv.Client.Close()
	srv.Cron.Stop()
}

//todo On 结点实例停用后，在 mongoDB 中去掉存活信息
func (srv *NodeServer) Down() {
	/*	n.Alived, n.DownTime = false, time.Now()
		n.SyncToMgo()*/
}

func (srv *NodeServer) Run() (err error) {
	defer func() {
		if err != nil {
			srv.Stop(err)
		}
	}()

	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		if r := recover(); r != nil {

		}
	}()
	if err = srv.loadJobs(); err != nil {
		return
	}
	//start cron
	srv.Cron.Start()
	go srv.watchJobs()
	go srv.watchKilledProc()
	go srv.watchOnce()
	//n.Node.On()
	return
}

func (srv *NodeServer) loadJobs() (err error) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		if r := recover(); r != nil {
			logger.GetLogger().Warn(fmt.Sprintf("load jobs panic:%v", r))
		}
	}()
	//再获取本机分配的定时任务
	jobs, err := handler.GetJobs(srv.UUID)
	if err != nil {
		return
	}

	if len(jobs) == 0 {
		return
	}
	srv.jobs = jobs

	for _, j := range jobs {
		j.InitNodeInfo(srv.UUID, srv.Hostname, srv.IP)
		srv.addJob(j, false)
	}

	return
}

//todo notice
func (srv *NodeServer) addJob(j *handler.Job, notice bool) {
	if err := j.Check(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("job check error :%s", err.Error()))
		return
	}
	taskFunc := handler.CreateJob(j)
	if taskFunc == nil {
		logger.GetLogger().Error(fmt.Sprintf("创建任务处理Job失败,不支持的任务协议#%s", j.Type))
		return
	}
	err := goutil.PanicToError(func() {
		srv.Cron.AddFunc(j.Spec, taskFunc, srv.jobCronName(j.ID))
	})
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("添加任务到调度器失败#%v", err.Error()))
	}

	return
}
func (srv *NodeServer) jobCronName(jobId int) string {
	return fmt.Sprintf(srv.UUID+"/%d", jobId)
}

func (srv *NodeServer) modifyJob(j *handler.Job) {
	oldJob, ok := srv.jobs[j.ID]
	// 之前此任务没有在当前结点执行，直接增加任务
	if !ok {
		srv.addJob(j, true)
		return
	}
	//先删除
	srv.deleteJob(oldJob.ID)
	//再
	srv.addJob(j, true)
	return
}

func (srv *NodeServer) deleteJob(jobId int) {
	if _, ok := srv.jobs[jobId]; ok {
		//存在则删除并且移除任务
		srv.Cron.RemoveJob(srv.jobCronName(jobId))
		delete(srv.jobs, jobId)
		return
	}
	return
}
