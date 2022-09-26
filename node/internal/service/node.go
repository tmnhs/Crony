package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/utils"
	"github.com/tmnhs/crony/node/internal/group"
	"github.com/tmnhs/crony/node/internal/job"
	"os"
	"strconv"
	"syscall"
)

// Node 执行 cron 命令服务的结构体
type NodeServer struct {
	*etcdclient.ServerReg
	*models.Node
	*cron.Cron

	jobs   job.Jobs // 和结点相关的任务
	Groups group.Groups

	models.Link
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
			ID:       uuid,
			PID:      strconv.Itoa(os.Getpid()),
			IP:       ip.String(),
			Hostname: hostname,
		},
		Cron: cron.New(),

		jobs: make(job.Jobs, 8),

		Link:   make(models.Link, 8),
		delIDs: make(map[string]bool, 8),

		ServerReg: etcdclient.NewServerReg(config.GetConfigModels().System.NodeTtl),
	}, nil

}

// Check whether the node is registered with ETCD
// If yes, PID is returned. If no, -1 is returned
func (srv *NodeServer) exist(nodeId string) (pid int, err error) {
	resp, err := etcdclient.Get(etcdclient.KeyEtcdNode + nodeId)
	if err != nil {
		return
	}

	if len(resp.Kvs) == 0 {
		return -1, nil
	}

	if pid, err = strconv.Atoi(string(resp.Kvs[0].Value)); err != nil {
		if _, err = etcdclient.Delete(etcdclient.KeyEtcdNode + nodeId); err != nil {
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
	pid, err := srv.exist(srv.ID)
	if err != nil {
		return err
	}
	if pid != -1 {
		return fmt.Errorf("node[%s] with pid[%d] exist", srv.ID, pid)
	}
	//creates a new lease
	if err := srv.ServerReg.Register(etcdclient.KeyEtcdNode+srv.ID, srv.PID); err != nil {
		return err
	}
	return nil
}

// 停止服务
func (srv *NodeServer) Stop(i interface{}) {
	//n.Node.Down()
	//todo 删除key值
	srv.Client.Close()
	srv.Cron.Stop()
}
