package main

import (
	"fmt"
	"github.com/tmnhs/crony/admin/internal/handler"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/notify"
	"github.com/tmnhs/crony/common/pkg/server"
	"os"
	"time"
)

const (
	ServerName = "admin"
)

func main() {
	srv, err := server.NewApiServer(ServerName)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("new api server error:%s", err.Error()))
		os.Exit(1)
	}
	// Register the API routing service
	srv.RegisterRouters(handler.RegisterRouters)
	service.DefaultNodeWatcher = service.NewNodeWatcherService()
	err = service.DefaultNodeWatcher.Watch()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("resolver  error:%#v", err))
	}
	// Notify operation
	go notify.Serve()
	// log cleaner
	var closeChan chan struct{}
	period := config.GetConfigModels().System.LogCleanPeriod
	if period > 0 {
		closeChan = service.RunLogCleaner(time.Duration(period)*time.Minute, config.GetConfigModels().System.LogCleanExpiration)
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("startup api server error:%v", err.Error()))
		close(closeChan)
		os.Exit(1)
	}
	os.Exit(0)
}
