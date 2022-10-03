package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"time"
)

type GroupRouter struct {
}

var defaultGroupRouter = new(GroupRouter)

func (g *GroupRouter) CreateOrUpdate(c *gin.Context) {
	var req models.Group
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[create_group] request parameter error", c)
		return
	}
	logger.GetLogger().Debug(fmt.Sprintf("create group req:%#v", req))
	var err error
	t := time.Now()
	if req.ID > 0 {
		//update
		req.Updated = t.Unix()
		err = req.Update()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[update_group] into db  error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[update_group] into db id error", c)
			return
		}
		resp.OkWithMessage("update success", c)
	} else {
		//create
		req.Created = t.Unix()
		_, err = req.Insert()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[create_group] insert group into db error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[create_group] insert group into db error", c)
			return
		}
		resp.OkWithMessage("create success", c)

	}
}

func (g *GroupRouter) Delete(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_group] request parameter error", c)
		return
	}
	group := models.Group{ID: req.ID}
	err := group.Delete()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_group]  db error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_group]  db error", c)
		return
	}
	resp.OkWithMessage("delete success", c)
}

func (g *GroupRouter) FindById(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_group] request parameter error", c)
		return
	}
	group := models.Group{ID: req.ID}
	err := group.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_group] find group by id :%d error:%s", req.ID, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_group] find group by id error", c)
		return
	}
	resp.OkWithDetailed(group, "find success", c)
}

func (g *GroupRouter) Search(c *gin.Context) {
	var req request.ReqGroupSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[search_group] request parameter error", c)
		return
	}
	req.Check()
	groups, total, err := service.DefaultGroupService.Search(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_group] search group error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[search_group] search group error", c)
		return
	}
	resp.OkWithDetailed(resp.PageResult{
		List:     groups,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "search success", c)
}
