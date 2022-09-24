package server

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron/v3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"net/http"
	"os"
	"strconv"
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

// Node 执行 cron 命令服务的结构体
type NodeServer struct {
	*etcdclient.Client
	*models.Node
	*cron.Cron

	jobs   models.Jobs // 和结点相关的任务
	groups models.Groups
	cmds   map[string]*models.Cmd

	models.Link
	// 删除的 job id，用于 group 更新
	delIDs map[string]bool

	timeout int64
	lID     clientv3.LeaseID // lease id
	done    chan struct{}
}

func NewNodeServer(serverName string, inits ...func()) (*NodeServer, error) {
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
		Client: etcdclient.GetEtcdClient(),
		Node: &models.Node{
			ID:       uuid,
			PID:      strconv.Itoa(os.Getpid()),
			IP:       ip.String(),
			Hostname: hostname,
		},
		Cron: cron.New(),

		jobs: make(models.Jobs, 8),
		cmds: make(map[string]*models.Cmd),

		Link:   make(models.Link, 8),
		delIDs: make(map[string]bool, 8),

		timeout: defaultConfig.System.NodeTimeout,
		done:    make(chan struct{}),
	}, nil
}
