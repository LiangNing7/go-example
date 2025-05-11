package main

import (
	"net/http"

	proto "github.com/LiangNing7/go-example/gin-example/05-render/proto/protoexample"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 引擎
	r := gin.Default()

	// 注册路由，返回 protobuf 内容
	r.GET("/someProtoBuf", func(c *gin.Context) {
		// 构造示例数据
		reps := []int64{1, 2}
		label := "test"
		data := &proto.Test{
			Label: label,
			Reps:  reps,
		}

		// 使用 c.ProtoBuf 序列化并返回二进制 protobuf
		c.ProtoBuf(http.StatusOK, data)
	})

	// 启动服务
	r.Run(":8080") // 默认在 0.0.0.0:8080
}
