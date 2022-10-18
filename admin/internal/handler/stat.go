package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"time"
)

type StatRouter struct {
}

var defaultStatRouter = new(StatRouter)

func (s *StatRouter) GetTodayStatistics(c *gin.Context) {
	jobExcSuccess, err := service.DefaultJobService.GetTodayJobExcCount(models.JobExcSuccess)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_statisitcs] GetTodayJobExcCount(successs)  error:%s", err.Error()))
	}
	jobExcFail, err := service.DefaultJobService.GetTodayJobExcCount(models.JobExcFail)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_statisitcs] GetTodayJobExcCount(fail) error:%s", err.Error()))
	}
	jobRunningCount, err := service.DefaultJobService.GetRunningJobCount()
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_statisitcs] GetRunningJobCount error:%s", err.Error()))
	}
	normalNodeCount, err := service.DefaultNodeWatcher.GetNodeCount(models.NodeConnSuccess)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_statisitcs] GetNodeCount(success) error:%s", err.Error()))
	}
	failNodeCount, err := service.DefaultNodeWatcher.GetNodeCount(models.NodeConnFail)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_statisitcs] GetNodeCount(fail) error:%s", err.Error()))
	}
	resp.OkWithDetailed(resp.RspSystemStatistics{
		NormalNodeCount:    normalNodeCount,
		FailNodeCount:      failNodeCount,
		JobExcSuccessCount: jobExcSuccess,
		JobRunningCount:    jobRunningCount,
		JobExcFailCount:    jobExcFail,
	}, "ok", c)
}

func (s *StatRouter) GetWeekStatistics(c *gin.Context) {
	t := time.Now()
	jobExcSuccess, err := service.DefaultJobService.GetJobExcCount(t.Unix()-60*60*24*7, t.Unix(), models.JobExcSuccess)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_week_statisitcs] GetTodayJobExcCount(successs)  error:%s", err.Error()))
	}
	jobExcFail, err := service.DefaultJobService.GetJobExcCount(t.Unix()-60*60*24*7, t.Unix(), models.JobExcFail)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("[get_week_statisitcs] GetTodayJobExcCount(fail) error:%s", err.Error()))
	}
	resp.OkWithDetailed(resp.RspDateCount{
		SuccessDateCount: jobExcSuccess,
		FailDateCount:    jobExcFail,
	}, "ok", c)
}

func (s *StatRouter) GetSystemInfo(c *gin.Context) {
	var req request.ByUUID
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[get_system_info] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[get_system_info] request parameter error", c)
		return
	}
	var server *utils.Server
	var err error
	if req.UUID == "" {
		//Get native information of admin
		server, err = utils.GetServerInfo()
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("[get_system_info]  error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[get_system_info]  error", c)
			return
		}
	} else {
		//Set the survival time to 30 seconds
		_, err := etcdclient.PutWithTtl(fmt.Sprintf(etcdclient.KeyEtcdSystemSwitch, req.UUID), models.NodeSystemInfoSwitch, 30)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] etcd put error: %s", req.UUID, err.Error()))
			resp.FailWithMessage(resp.ERROR, "[get_system_info]  error", c)
			return
		}
		//There will be a delay. By default, wait 2s.
		time.Sleep(2 * time.Second)
		server, err = service.GetNodeSystemInfo(req.UUID)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("get system info from node[%s] watch key error: %s", req.UUID, err.Error()))
			resp.FailWithMessage(resp.ERROR, "[get_system_info]  error", c)
			return
		}
	}
	resp.OkWithDetailed(server, "ok", c)

}
