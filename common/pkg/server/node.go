package server

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"net/http"
	"os"
)

var (
	NodeOptions struct {
		flags.Options
		Environment    string `short:"e" long:"env" description:"Use nodeServer environment" default:"testing"`
		Version        bool   `short:"v" long:"verbose"  description:"Show nodeServer version"`
		EnablePProfile bool   `short:"p" long:"enable-pprof"  description:"enable pprof"`
		PProfilePort   int    `short:"d" long:"pprof-port"  description:"pprof port" default:"8188"`
		ConfigFileName string `short:"c" long:"config" description:"Use nodeServer config file" default:"main"`
		EnableDevMode  bool   `short:"m" long:"enable-dev-mode"  description:"enable dev mode"`
	}
)

func InitNodeServer(serverName string, inits ...func()) (*models.Config, error) {
	var parser = flags.NewParser(&NodeOptions, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		return nil, err
	}

	if NodeOptions.Version {
		fmt.Printf("%s Version:%s\n", NodeModule, Version)
		os.Exit(0)
	}

	if NodeOptions.EnablePProfile {
		go func() {
			fmt.Printf("enable pprof http server at:%d\n", NodeOptions.PProfilePort)
			fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", NodeOptions.PProfilePort), nil))
		}()
	}
	var env = config.Environment(NodeOptions.Environment)
	if env.Invalid() {
		var err error
		env, err = config.NewGlobalEnvironment()
		if err != nil {
			return nil, err
		}
	}

	var configFile = NodeOptions.ConfigFileName
	if configFile == "" {
		configFile = "main"
	}
	defaultConfig, err := config.LoadConfig(env.String(), serverName, configFile)
	if err != nil {
		fmt.Printf("node-server:init config error:%s", err.Error())
		return nil, err
	}
	//todo
	logger.Init(&defaultConfig.Log, serverName)

	//初始化数据层服务
	_, err = dbclient.Init(defaultConfig.Mysql)
	if err != nil {
		logger.Errorf("node-server:init mysql failed , error:%s", err.Error())
	} else {
		logger.Info("node-server:init mysql success")
	}
	//初始化etcd
	_, err = etcdclient.Init(defaultConfig.Etcd)
	if err != nil {
		logger.Errorf("node-server:init etcd failed , error:%s", err.Error())
	} else {
		logger.Info("node-server:init etcd success")
	}
	if len(inits) > 0 {
		for _, init := range inits {
			init()
		}
	}
	return defaultConfig, nil
}
