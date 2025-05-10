package main

import (
	"net/http"

	"github.com/LiangNing7/go-example/gee/07-recover/gee"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello LiangNing\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"liangning"}
		c.String(http.StatusOK, "%s", names[100])
	})

	r.Run(":8080")
}
