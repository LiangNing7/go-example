package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/JSONP", func(c *gin.Context) {
		data := map[string]any{
			"foo": "bar",
		}

		// /JSONP?callback=x
		// 将输出：x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
