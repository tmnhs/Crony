package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouters(c *gin.Engine) {
	hello := c.Group("/ping")
	{
		hello.GET("", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
	}
	base := c.Group("")
	{
		base.POST("register", defaultUserRouter.Register)
		base.POST("login", defaultUserRouter.Login)
	}
	job := c.Group("/job")
	{
		job.POST("add", defaultJobRouter.CreateOrUpdate)
		job.POST("del", defaultJobRouter.Delete)
		job.GET("find", defaultJobRouter.FindById)
		job.GET("search", defaultJobRouter.Search)
		job.GET("log", defaultJobRouter.SearchLog)
	}
	user := c.Group("/user")
	{
		user.POST("del", defaultUserRouter.Delete)
		user.POST("find", defaultUserRouter.FindById)
		user.POST("search", defaultUserRouter.Search)
	}
}
