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
	if err = nodeServer.Register(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("register node into etcd error:%s", err.Error()))
		os.Exit(1)
	}
	if err = nodeServer.Run(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node run error: %s", err.Error()))
		os.Exit(1)
	}
	//notification operation
	go notify.Serve()
	logger.GetLogger().Info(fmt.Sprintf("crony node %s service started, Ctrl+C or send kill sign to exit", nodeServer.String()))
	// Register the logout event
	event.OnEvent(event.EXIT, nodeServer.Stop)
	// Listen for exit signals
	event.WaitEvent()
	event.EmitEvent(event.EXIT, nil)
	logger.GetLogger().Info("exit success")
}
