package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MyHeader struct {
	Rate   int    `header:"Rate" binding:"required"`
	Domain string `header:"Domain" binding:"required"`
}

func main() {
	r := gin.New()

	r.GET("/UAgent", func(c *gin.Context) {
		userAgent := c.Request.Header.Get("User-Agent")
		c.JSON(http.StatusOK, gin.H{
			"User-Agent": userAgent,
		})
	})

	r.GET("/token", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		c.String(http.StatusOK, "Token: %s", token)
	})

	r.GET("/up", func(c *gin.Context) {
		var h MyHeader
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Rate":   h.Rate,
			"Domain": h.Domain,
		})
	})

	r.Run()
}
