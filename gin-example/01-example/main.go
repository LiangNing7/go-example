package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的路由引擎.
	r := gin.Default()
	// GET: 请求方式；/hello: 请求的路径.
	r.GET("/hello", func(c *gin.Context) {
		// c.JSON: 返回 JSON 格式的数据.
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})
	// 启动 HTTP 服务，默认在 0.0.0.0:8080 启动服务.
	r.Run()
}
