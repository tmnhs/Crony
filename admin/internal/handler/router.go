package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouters(c *gin.Engine) {
	hello := c.Group("/ping")
	{
		hello.GET("", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
	}
	job := c.Group("/job")
	{
		job.POST("add", defaultJobRouter.CreateOrUpdate)
		job.POST("del", defaultJobRouter.Delete)
		job.POST("find", defaultJobRouter.FindById)
		job.POST("search", defaultJobRouter.Search)
	}
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS = 200
	ERROR   = 1000

	ErrorRequestParameter = 1001
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "operation success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "operation success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func FailWithMessage(code int, message string, c *gin.Context) {
	Result(code, map[string]interface{}{}, message, c)
}
func FailWithCode(code int, c *gin.Context) {
	Result(code, map[string]interface{}{}, "operation failed", c)
}
func FailWithDetailed(code int, data interface{}, message string, c *gin.Context) {
	Result(code, data, message, c)
}
