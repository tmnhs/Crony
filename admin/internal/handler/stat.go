package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
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
