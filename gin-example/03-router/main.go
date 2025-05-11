package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// 定义普通路由组.
	router.GET("/index", func(c *gin.Context) {
		c.String(http.StatusOK, "%s - %s", c.Request.Method, c.Request.RequestURI)
	})

	router.POST("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "%s - %s", c.Request.Method, c.Request.RequestURI)
	})

	router.Any("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "%s - %s", c.Request.Method, c.Request.RequestURI)
	})

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "%s - %s NotFound\n", c.Request.Method, c.Request.RequestURI)
	})

	// 路由参数.
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		c.String(http.StatusOK, "%s is %s", name, action)
	})

	// 路由组.
	shopGroup := router.Group("/shop")
	{
		shopGroup.GET("/index", func(c *gin.Context) {
			c.String(http.StatusOK, "%s - %s", c.Request.Method, c.Request.RequestURI)
		})
		shopGroup.POST("/checkout", func(c *gin.Context) {
			c.String(http.StatusOK, "%s - %s", c.Request.Method, c.Request.RequestURI)
		})
		xx := shopGroup.Group("/xx")
		{
			xx.GET("/oo", func(c *gin.Context) {
				c.String(http.StatusOK, "%s - %s", c.Request.Method, c.Request.RequestURI)
			})
		}
	}

	router.Run(":8080")
}
