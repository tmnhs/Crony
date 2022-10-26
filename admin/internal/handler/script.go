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

type ScriptRouter struct {
}

var defaultScriptRouter = new(ScriptRouter)

func (s *ScriptRouter) CreateOrUpdate(c *gin.Context) {
	var req models.Script
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_script] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[create_script] request parameter error", c)
		return
	}
	var err error
	t := time.Now()
	if req.ID > 0 {
		//update
		req.Updated = t.Unix()
		err = req.Update()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[update_script] into db  error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[update_script] into db id error", c)
			return
		}
	} else {
		//create
		req.Created = t.Unix()
		_, err = req.Insert()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[create_script] insert script into db error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[create_script] insert script into db error", c)
			return
		}
	}
	resp.OkWithDetailed(req, "operate success", c)
}

func (s *ScriptRouter) Delete(c *gin.Context) {
	var req request.ByIDS
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_script] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_script] request parameter error", c)
		return
	}
	for _, id := range req.IDS {
		script := models.Script{ID: id}
		err := script.FindById()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[delete_script] find script by id :%d error:%s", id, err.Error()))
			continue
		}
		err = script.Delete()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[delete_script] into db error:%s", err.Error()))
			continue
		}
	}
	resp.OkWithMessage("delete success", c)
}

func (s *ScriptRouter) FindById(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[find_script] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[find_script] request parameter error", c)
		return
	}
	script := models.Script{ID: req.ID}
	err := script.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[find_script] find script by id :%d error:%s", req.ID, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[find_script] find script by id error", c)
		return
	}
	resp.OkWithDetailed(script, "find success", c)
}

func (s *ScriptRouter) Search(c *gin.Context) {
	var req request.ReqScriptSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_script] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[search_script] request parameter error", c)
		return
	}
	req.Check()
	scripts, total, err := service.DefaultScriptService.Search(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_script] search script error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[search_script] search script error", c)
		return
	}
	resp.OkWithDetailed(resp.PageResult{
		List:     scripts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "search success", c)
}
