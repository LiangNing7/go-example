package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]any{
			"lang": "Go 语言",
			"tag":  "<br>",
		}
		// 输出: {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}

		c.AsciiJSON(http.StatusOK, data)
	})
	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
