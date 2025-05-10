package main

import (
	"net/http"

	"github.com/LiangNing7/go-example/gee/04-group/gee"
)

/*
(1) index
curl -i http://localhost:8080/index
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sat, 10 May 2025 07:16:45 GMT
Content-Length: 19

<h1>Index Page</h1>

(2) v1
curl -i http://localhost:8080/v1/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sat, 10 May 2025 07:17:19 GMT
Content-Length: 18

<h1>Hello Gee</h1>

(3)
curl "http://localhost:8080/v1/hello?name=liangning"
hello liangning, you're at /v1/hello

(4)
curl "http://localhost:8080/v2/hello/liangning"
hello liangning, you're at /v2/hello/liangning

(5)
curl "http://localhost:8080/v2/login" -X POST -d 'username=liangning&password=1234'
{"password":"1234","username":"liangning"}

(6)
curl "http://localhost:8080/hello"
404 NOT FOUND: /hello
*/

func main() {
	r := gee.New()
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=liangning
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/liangning
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":8080")
}
