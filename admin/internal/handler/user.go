package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/pkg/logger"
)

type UserRouter struct {
}

func (u *UserRouter) Login(c *gin.Context) {
	var req model.ReqUserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_login] request parameter error:%s", err.Error()))
		FailWithMessage(ErrorRequestParameter, "[user_login] request parameter error", c)
		return
	}
	user, err := service.DefaultUserService.Login(req.UserName, req.Password)
	if err != nil || user == nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_login] db error:%v", err))
		FailWithMessage(ERROR, "[user_login] username or password is incorrect", c)
		return
	}
	OkWithDetailed(user, "login success", c)
}

func (u *UserRouter) Register(c *gin.Context) {
	var req model.ReqUserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_login] request parameter error:%s", err.Error()))
		FailWithMessage(ErrorRequestParameter, "[user_login] request parameter error", c)
		return
	}
	user, err := service.DefaultUserService.Login(req.UserName, req.Password)
	if err != nil || user == nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_login] db error:%v", err))
		FailWithMessage(ERROR, "[user_login] username or password is incorrect", c)
		return
	}
	OkWithDetailed(user, "login success", c)
}
