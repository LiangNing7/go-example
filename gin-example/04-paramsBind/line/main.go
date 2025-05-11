package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// URI 绑定结构体
type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

// Query 绑定结构体
type Search struct {
	Query  string `form:"q" binding:"required"`
	Page   int    `form:"page,default=1"`
	Limits []int  `form:"limit"`
}

func main() {
	router := gin.Default()

	// 1. 演示请求方法、原始 URI、解析后 URI、路径、协议版本
	router.GET("/info/*any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method":             c.Request.Method,
			"raw_request_uri":    c.Request.RequestURI,
			"parsed_request_uri": c.Request.URL.RequestURI(),
			"path":               c.Request.URL.Path,
			"proto":              c.Request.Proto,
			"proto_major":        c.Request.ProtoMajor,
			"proto_minor":        c.Request.ProtoMinor,
		})
	})

	// 2. URI 参数绑定
	router.GET("/user/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"name": person.Name,
			"uuid": person.ID,
		})
	})

	// 3. 获取所有 Query 参数，及单个参数和默认值
	router.GET("/params", func(c *gin.Context) {
		allQuery := c.Request.URL.Query()
		lastname := c.Query("lastname")
		firstname := c.DefaultQuery("firstname", "Guest")
		c.JSON(http.StatusOK, gin.H{
			"all_query": allQuery,
			"lastname":  lastname,
			"firstname": firstname,
		})
	})

	// 4. Query 中的数组与 Map
	// 示例： /tags?tags=go&tags=gin
	//        /opts?opts[a]=1&opts[b]=2
	router.GET("/extras", func(c *gin.Context) {
		tags := c.QueryArray("tags")
		opts := c.QueryMap("opts")
		c.JSON(http.StatusOK, gin.H{
			"tags": tags,
			"opts": opts,
		})
	})

	// 5. 将 Query 绑定到结构体
	// 示例： /search?q=gin&page=2&limit=10&limit=20
	router.GET("/search", func(c *gin.Context) {
		var s Search
		if err := c.ShouldBindQuery(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	router.Run(":8080") // 启动服务，监听 8080 端口
}
