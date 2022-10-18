package handler

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
)

type NodeRouter struct{}

var defaultNodeRouter = new(NodeRouter)

func (n *NodeRouter) Delete(c *gin.Context) {
	var req request.ByUUID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_node] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_node] request parameter error", c)
		return
	}
	node := &models.Node{UUID: req.UUID}
	err := node.FindByUUID()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_node] find node by uuid :%s error:%s", req.UUID, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_node] db find error", c)
		return
	}
	if node.Status == models.NodeConnSuccess {
		resp.FailWithMessage(resp.ERROR, "[delete_node] can't delete a node that is already alive ", c)
		return
	}
	_, _ = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, req.UUID), clientv3.WithPrefix())
	err = node.Delete()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_node] into db error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_node] db delete error", c)
		return
	}
	resp.OkWithMessage("delete success", c)
}

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
