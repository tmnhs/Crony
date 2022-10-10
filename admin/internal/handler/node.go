package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
)

type NodeRouter struct{}

var defaultNodeRouter = new(NodeRouter)

func (n *NodeRouter) Search(c *gin.Context) {
	var req request.ReqNodeSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_node] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[search_node] request parameter error", c)
		return
	}
	req.Check()
	nodes, total, err := service.DefaultNodeWatcher.Search(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_node] search node error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[search_node] search node  error", c)
		return
	}
	var resultNodes []resp.RspNodeSearch
	for _, node := range nodes {
		resultNode := resp.RspNodeSearch{
			Node: node,
		}
		resultNode.JobCount, _ = service.DefaultNodeWatcher.GetJobCount(node.UUID)
		resultNodes = append(resultNodes, resultNode)
	}

	resp.OkWithDetailed(resp.PageResult{
		List:     resultNodes,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "search success", c)
}

func (n *NodeRouter) GetStatistics(c *gin.Context) {
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
		NormalNodeCount: normalNodeCount,
		FailNodeCount:   failNodeCount,
		JobExcCount:     jobExcFail + jobExcSuccess,
		JobRunningCount: jobRunningCount,
		JobExcFailCount: jobExcFail,
	}, "ok", c)
}
