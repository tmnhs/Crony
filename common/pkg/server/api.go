package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/notify"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	shutdownMaxAge = 15 * time.Second
	shutdownWait   = 1000 * time.Millisecond
)
const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var (
	ApiOptions struct {
		flags.Options
		Environment       string `short:"e" long:"env" description:"Use ApiServer environment" default:"testing"`
		Version           bool   `short:"v" long:"verbose"  description:"Show ApiServer version"`
		EnablePProfile    bool   `short:"p" long:"enable-pprof"  description:"enable pprof"`
		PProfilePort      int    `short:"d" long:"pprof-port"  description:"pprof port" default:"8188"`
		EnableHealthCheck bool   `short:"a" long:"enable-health-check"  description:"enable health check"`
		HealthCheckURI    string `short:"i" long:"health-check-uri"  description:"health check uri" default:"/health" `
		HealthCheckPort   int    `short:"f" long:"health-check-port"  description:"health check port" default:"8186"`
		ConfigFileName    string `short:"c" long:"config" description:"Use ApiServer config file" default:"main"`
		EnableDevMode     bool   `short:"m" long:"enable-dev-mode"  description:"enable dev mode"`
	}
)

type healthCheckHttpServer struct {
}

func (server *healthCheckHttpServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	io.WriteString(response, "ok\n")
}

var healthCheckServer = &healthCheckHttpServer{}

type ApiServer struct {
	Engine      *gin.Engine
	HttpServer  *http.Server
	Addr        string
	mu          sync.Mutex
	doneChan    chan struct{}
	Routers     []func(*gin.Engine)
	Middlewares []func(*gin.Engine)
	Shutdowns   []func(*ApiServer)
	Services    []func(*ApiServer)
}

//get close Chan
func (srv *ApiServer) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.getDoneChanLocked()
}

func (srv *ApiServer) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *ApiServer) Shutdown(ctx context.Context) {
	//Give priority to business shutdown Hook
	if len(srv.Shutdowns) > 0 {
		for _, shutdown := range srv.Shutdowns {
			shutdown(srv)
		}
	}
	//wait for registry shutdown
	select {
	case <-time.After(shutdownWait):
	}
	// close the HttpServer
	srv.HttpServer.Shutdown(ctx)
}

// ApiRecovery recovery any panics and writes a 500 if there was one.
func (srv *ApiServer) apiRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				stack := stack(3)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}

				if brokenPipe {
					logger.GetLogger().Error(fmt.Sprintf("%s\n%s%s", err, string(httpRequest), reset))
				} else {
					logger.GetLogger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s%s",
						formatTime(time.Now()), err, stack, reset))
				}

				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}

func (srv *ApiServer) setupSignal() {
	go func() {
		var sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan /*syscall.SIGUSR1,*/, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownMaxAge)
		defer shutdownCancel()

		for sig := range sigChan {
			if sig == syscall.SIGINT || sig == syscall.SIGHUP || sig == syscall.SIGTERM {
				logger.GetLogger().Error(fmt.Sprintf("Graceful shutdown:signal %v to stop api-server ", sig))
				srv.Shutdown(shutdownCtx)
			} else {
				logger.GetLogger().Info(fmt.Sprintf("Caught signal %v", sig))
			}
		}
		logger.Shutdown()
	}()
}

func NewApiServer(serverName string, inits ...func()) (*ApiServer, error) {
	var parser = flags.NewParser(&ApiOptions, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		return nil, err
	}

	if ApiOptions.Version {
		fmt.Printf("%s Version:%s\n", ApiModule, Version)
		os.Exit(0)
	}

	if ApiOptions.EnablePProfile {
		go func() {
			fmt.Printf("enable pprof http server at:%d\n", ApiOptions.PProfilePort)
			fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", ApiOptions.PProfilePort), nil))
		}()
	}

	if ApiOptions.EnableHealthCheck {
		go func() {
			fmt.Printf("enable healthcheck http server at:%d\n", ApiOptions.HealthCheckPort)
			fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", ApiOptions.HealthCheckPort), healthCheckServer))
		}()
	}
	var env = config.Environment(ApiOptions.Environment)
	if env.Invalid() {
		var err error
		env, err = config.NewGlobalEnvironment()
		if err != nil {
			return nil, err
		}
	}

	var configFile = ApiOptions.ConfigFileName
	if configFile == "" {
		configFile = "main"
	}
	defaultConfig, err := config.LoadConfig(env.String(), serverName, configFile)
	if err != nil {
		fmt.Printf("api-server:init config error:%s", err.Error())
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
		logger.GetLogger().Error(fmt.Sprintf("api-server:init mysql failed , error:%s", err.Error()))
	} else {
		logger.GetLogger().Info("api-server:init mysql success")
	}
	//etcd
	_, err = etcdclient.Init(etcdConfig.Endpoints, etcdConfig.DialTimeout, etcdConfig.ReqTimeout)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("api-server:init etcd failed , error:%s", err.Error()))
	} else {
		logger.GetLogger().Info("api-server:init etcd success")
	}

	if len(inits) > 0 {
		for _, init := range inits {
			init()
		}
	}

	apiServer := &ApiServer{
		Addr: fmt.Sprintf(":%d", defaultConfig.System.Addr),
	}

	apiServer.setupSignal()
	//set gin mode
	switch env {
	case config.EnvProduction:
		gin.SetMode(gin.ReleaseMode)
	case config.EnvTesting:
		gin.SetMode(gin.DebugMode)
	}
	return apiServer, nil
}

// ListenAndServe Listen And Serve()
func (srv *ApiServer) ListenAndServe() error {
	srv.Engine = gin.New()
	srv.Engine.Use(srv.apiRecoveryMiddleware())

	for _, service := range srv.Services {
		service(srv)
	}

	for _, middleware := range srv.Middlewares {
		middleware(srv.Engine)
	}

	for _, c := range srv.Routers {
		c(srv.Engine)
	}

	srv.HttpServer = &http.Server{
		Handler:        srv.Engine,
		Addr:           srv.Addr,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := srv.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Register Shutdown Handler
func (srv *ApiServer) RegisterShutdown(handlers ...func(*ApiServer)) {
	srv.Shutdowns = append(srv.Shutdowns, handlers...)
}

// Register Service Handler
func (srv *ApiServer) RegisterService(handlers ...func(*ApiServer)) {
	srv.Services = append(srv.Services, handlers...)
}

// Register Middleware Middleware
func (srv *ApiServer) RegisterMiddleware(middlewares ...func(engine *gin.Engine)) {
	srv.Middlewares = append(srv.Middlewares, middlewares...)
}

// RegisterRouters
func (srv *ApiServer) RegisterRouters(routers ...func(engine *gin.Engine)) *ApiServer {
	srv.Routers = append(srv.Routers, routers...)
	return srv
}
