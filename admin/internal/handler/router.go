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
}
