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

func (n *NodeRouter) JoinGroup(c *gin.Context) {
	var req models.NodeGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[node_join_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[node_join_group] request parameter error", c)
		return
	}
	_, err := req.Insert()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[node_join_group] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[node_join_group] db error", c)
		return
	}
	resp.OkWithMessage("join success", c)
}

func (n *NodeRouter) KickGroup(c *gin.Context) {
	var req models.NodeGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[node_kick_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[node_kick_group] request parameter error", c)
		return
	}
	err := req.Delete()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[node_kick_group] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[node_kick_group] db error", c)
		return
	}
	resp.OkWithMessage("kick success", c)
}

func (n *NodeRouter) GetByGroupId(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[node_get_by_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[node_get_by_group] request parameter error", c)
		return
	}
	nodes, err := service.DefaultNodeWatcher.FindByGroupId(req.ID)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[node_get_by_group] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[node_get_by_group] db error", c)
		return
	}
	resp.OkWithDetailed(nodes, "get success", c)
}
