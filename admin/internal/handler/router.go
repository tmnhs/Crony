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
		user.GET("find", defaultUserRouter.FindById)
		user.GET("search", defaultUserRouter.Search)
		user.GET("get_by_group", defaultUserRouter.GetByGroupId)
		user.POST("join", defaultUserRouter.JoinGroup)
		user.POST("kick", defaultUserRouter.KickGroup)
	}
	node := c.Group("/node")
	{
		node.GET("search", defaultNodeRouter.Search)
		node.POST("join", defaultNodeRouter.JoinGroup)
		node.POST("kick", defaultNodeRouter.KickGroup)
		node.GET("get_by_group", defaultNodeRouter.GetByGroupId)
	}

	group := c.Group("/group")
	{
		group.POST("update", defaultGroupRouter.CreateOrUpdate)
		group.POST("del", defaultGroupRouter.Delete)
		group.GET("find", defaultGroupRouter.FindById)
		group.GET("search", defaultGroupRouter.Search)
	}
}
