package main

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/notify"
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
	logger.GetLogger().Debug(fmt.Sprintf("nodeServer:%#v", *nodeServer))
	logger.GetLogger().Debug(fmt.Sprintf("node:%#v", *nodeServer.Node))
	if err = nodeServer.Register(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("register node into etcd error:%s", err.Error()))
		os.Exit(1)
	}
	if err = nodeServer.Run(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node run error: %s", err.Error()))
		os.Exit(1)
	}
	//邮件相关操作
	go notify.Serve()
	logger.GetLogger().Info(fmt.Sprintf("crony node %s service started, Ctrl+C or send kill sign to exit", nodeServer.String()))
	// 注册退出事件
	event.OnEvent(event.EXIT, nodeServer.Stop /*,stopwatcher()*/)
	// 监听退出信号
	event.WaitEvent()
	// 处理退出事件
	event.EmitEvent(event.EXIT, nil)
	logger.GetLogger().Info("exit success")
}
