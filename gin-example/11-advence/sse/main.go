package main

import (
	"io"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 简单实现
	r.GET("/stream1", func(c *gin.Context) {
		// 设置 SSE 必要头
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")

		// 每秒推送一次数据
		c.Stream(func(w io.Writer) bool {
			c.SSEvent("time", gin.H{"now": time.Now().Format(time.RFC3339)})
			time.Sleep(1 * time.Second)
			return true // 返回 false 时结束该流
		})
	})

	// 使用第三方库.
	r.GET("/stream2", func(c *gin.Context) {
		// 标准 SSE 头
		c.Writer.Header().Set("Content-Type", sse.ContentType)
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")

		// 使用 sse.Encode 输出多种类型数据
		c.Stream(func(w io.Writer) bool {
			sse.Encode(w, sse.Event{
				Event: "msg",
				Data:  "hello world",
			})
			sse.Encode(w, sse.Event{
				Id:    "123",
				Event: "json",
				Data:  gin.H{"user": "Alice", "time": time.Now().Unix()},
			})
			time.Sleep(2 * time.Second)
			return true
		})
	})
	r.Run(":8080")
}
