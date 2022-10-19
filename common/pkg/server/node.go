package server

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/notify"
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
	logConfig := defaultConfig.Log
	mysqlConfig := defaultConfig.Mysql
	etcdConfig := defaultConfig.Etcd
	//log
	logger.Init(serverName, logConfig.Level, logConfig.Format, logConfig.Prefix, logConfig.Director, logConfig.ShowLine, logConfig.EncodeLevel, logConfig.StacktraceKey, logConfig.LogInConsole)
	//notify
	notify.Init(&notify.Mail{
		Port:     defaultConfig.Email.Port,
		From:     defaultConfig.Email.From,
		Host:     defaultConfig.Email.Host,
		Secret:   defaultConfig.Email.Secret,
		Nickname: defaultConfig.Email.Nickname,
	}, &notify.WebHook{
		Url:  defaultConfig.WebHook.Url,
		Kind: defaultConfig.WebHook.Kind,
	})
	//db
	dsn := mysqlConfig.EmptyDsn()
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 ;", mysqlConfig.Dbname)
	if err := dbclient.CreateDatabase(dsn, "mysql", createSql); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("create mysql database failed , error:%s", err.Error()))
	}
	_, err = dbclient.Init(mysqlConfig.Dsn(), mysqlConfig.LogMode, mysqlConfig.MaxIdleConns, mysqlConfig.MaxOpenConns)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node-server:init mysql failed , error:%s", err.Error()))
	} else {
		logger.GetLogger().Info("node-server:init mysql success")
	}
	//etcd
	_, err = etcdclient.Init(etcdConfig.Endpoints, etcdConfig.DialTimeout, etcdConfig.ReqTimeout)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node-server:init etcd failed , error:%s", err.Error()))
	} else {
		logger.GetLogger().Info("node-server:init etcd success")
	}
	if len(inits) > 0 {
		for _, init := range inits {
			init()
		}
	}
	return defaultConfig, nil
}
