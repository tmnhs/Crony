package main

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/server"
	"github.com/tmnhs/crony/common/pkg/utils/event"
	"github.com/tmnhs/crony/node/internal/service"
	"os"
)

const ServerName = "node"

func main() {
	if _, err := server.InitNodeServer(ServerName); err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	nodeServer, err := service.NewNodeServer()
	if err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	logger.Debugf("nodeServer:%#v", *nodeServer)
	logger.Debugf("node:%#v", *nodeServer.Node)
	//todo register to etcd
	if err = nodeServer.Register(); err != nil {
		fmt.Println("register node into etcd error:", err.Error())
		os.Exit(1)
	}
	//todo run

	logger.Infof("crony node %s service started, Ctrl+C or send kill sign to exit", nodeServer.String())
	// 注册退出事件
	event.OnEvent(event.EXIT, nodeServer.Stop /*,stopwatcher()*/)
	// 监听退出信号
	event.WaitEvent()
	// 处理退出事件
	event.EmitEvent(event.EXIT, nil)
	logger.Infof("exit success")
}
