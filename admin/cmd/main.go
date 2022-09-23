package main

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/server"
	"os"
)

const  (
	ServerName="admin"
)
func main() {
	srv,err:=server.NewApiServer(ServerName)
	if err != nil {
		fmt.Printf("new api server error:%s",err.Error())
		os.Exit(1)
	}
	logger.Infof("hello logger")
	//todo 邮件相关操作
	//todo 定时清理日志
	err = srv.ListenAndServe()
	if err != nil {
		logger.Errorf("startup api server error:%v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
