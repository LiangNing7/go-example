package main

import (
	"net/http"
	"time"

	"github.com/LiangNing7/go-example/gin-example/10-self/selfmiddleware/cost"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(cost.StatCost())

	router.GET("/", func(c *gin.Context) {
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "Hello World")
	})

	router.Run(":8080")
}
