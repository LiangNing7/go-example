package main

import (
	"log"
	"net/http"
	"time"

	"github.com/LiangNing7/go-example/gee/05-middleware/gee"
)

/*
(1) global middleware Logger
curl http://localhost:8080/
<h1>Hello Gee</h1>
>>> log
2025/05/10 17:31:42 [200] / in 6.955µs

(2) global + group middleware
$ curl http://localhost:8080/v2/hello/liangning
{"message":"Internal Server Error"}
>>> log
2025/05/10 17:32:43 [500] /v2/hello/liangning in 59.364µs for group v2
2025/05/10 17:32:43 [500] /v2/hello/liangning in 81.815µs
*/

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer.
		t := time.Now()
		// if a server err occurred.
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time.
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // global midlleware
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/liangning
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":8080")
}
