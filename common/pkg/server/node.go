package server

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jessevdk/go-flags"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"net/http"
	"os"
	"time"
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

// 执行 cron cmd 的进程
// 注册到 /cronsun/node/<id>
type Node struct {
	ID       string `json:"id"`  // machine id
	PID      string `json:"pid"` // 进程 pid
	PIDFile  string `json:"-"`
	IP       string `json:"ip"` // node ip
	Hostname string `json:"hostname"`

	Version  string    `json:"version"`
	UpTime   time.Time `json:"up"`   // 启动时间
	DownTime time.Time `json:"down"` // 上次关闭时间

	Alived    bool `json:"alived"`   // 是否可用
	Connected bool `son:"connected"` // 当 Alived 为 true 时有效，表示心跳是否正常
}

type NodeServer struct {
	*etcdclient.Client
	*Node
	//*cron.Cron

	//jobs   Jobs // 和结点相关的任务
	//groups Groups
	//cmds   map[string]*cronsun.Cmd

	//link
	// 删除的 job id，用于 group 更新
	delIDs map[string]bool

	ttl  int64
	lID  clientv3.LeaseID // lease id
	done chan struct{}
}

func InitNodeServer(serverName string, inits ...func()) error {
	var parser = flags.NewParser(&NodeOptions, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		return err
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
			return err
		}
	}

	var configFile = NodeOptions.ConfigFileName
	if configFile == "" {
		configFile = "main"
	}
	defaultConfig, err := config.LoadConfig(env.String(), serverName, configFile)
	if err != nil {
		fmt.Printf("node-server:init config error:%s", err.Error())
		return err
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

	return nil
}
