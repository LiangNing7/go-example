package cost

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// StatCost 是一个统计请求耗时的中间件.
func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer.
		start := time.Now()

		// 创建一个用于标记请求是否已正常完成的 channel.
		finished := make(chan struct{}, 1)

		// 启动一个协程，监听上下文的 Done 信号.
		go func() {
			select {
			case <-c.Request.Context().Done():
				// 客户端断开连接或超时导致 Context 取消.
				elapsed := time.Since(start)
				log.Printf("[GIN] %s %s TERMINATED after %v\n",
					c.Request.Method,
					c.Request.URL.RequestURI(),
					elapsed,
				)
			case <-finished:
				// 请求已正常退出.
			}
		}()

		// 继续处理后续请求.
		c.Next()

		// 请求正常结束后记录耗时.
		elapsed := time.Since(start)
		log.Printf("[GIN] %s %s completed in %v\n",
			c.Request.Method,
			c.Request.URL.RequestURI(),
			elapsed,
		)

		// 通知监听协程不再输出 TERMINATED 日志.
		close(finished)
	}
}
