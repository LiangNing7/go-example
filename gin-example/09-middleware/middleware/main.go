package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 全局注册路由
	r.Use(gin.Logger())

	// 为单个路由定义中间件
	r.GET("/index", gin.ErrorLogger(), func(c *gin.Context) {
		c.String(http.StatusOK, "helloworld")
	})

	// 路由组注册路由
	rv := r.Group("/recover", gin.Recovery())
	{
		rv.GET("/index", func(c *gin.Context) {
			c.String(http.StatusOK, "hello world")
		})
	}
}
