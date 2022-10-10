package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/middlerware"
)

func RegisterRouters(c *gin.Engine) {
	//跨域
	c.Use(middlerware.Cors())
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
	job.Use(middlerware.JWTAuth())
	{
		job.POST("add", defaultJobRouter.CreateOrUpdate)
		job.POST("del", defaultJobRouter.Delete)
		job.GET("find", defaultJobRouter.FindById)
		job.POST("search", defaultJobRouter.Search)
		job.POST("log", defaultJobRouter.SearchLog)
		job.POST("once", defaultJobRouter.Once)
	}

	user := c.Group("/user")
	user.Use(middlerware.JWTAuth())
	{
		user.POST("del", defaultUserRouter.Delete)
		user.POST("update", defaultUserRouter.Update)
		user.POST("change_pw", defaultUserRouter.ChangePassword)
		user.GET("find", defaultUserRouter.FindById)
		user.POST("search", defaultUserRouter.Search)
		user.GET("get_by_group", defaultUserRouter.GetByGroupId)
		user.POST("join", defaultUserRouter.JoinGroup)
		user.POST("kick", defaultUserRouter.KickGroup)
	}
	node := c.Group("/node")
	node.Use(middlerware.JWTAuth())
	{
		node.POST("search", defaultNodeRouter.Search)
		node.POST("join", defaultNodeRouter.JoinGroup)
		node.POST("kick", defaultNodeRouter.KickGroup)
		node.GET("get_by_group", defaultNodeRouter.GetByGroupId)
	}

	group := c.Group("/group")
	//group.Use(middlerware.JWTAuth())
	{
		group.POST("update", defaultGroupRouter.CreateOrUpdate)
		group.POST("del", defaultGroupRouter.Delete)
		group.GET("find", defaultGroupRouter.FindById)
		group.GET("search", defaultGroupRouter.Search)
	}
}
