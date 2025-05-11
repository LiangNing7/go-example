package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Uri struct {
	ID   string `uri:"id"`
	Name string `uri:"name"`
}

func main() {
	router := gin.New()

	router.GET("/user/:id/:name", func(c *gin.Context) {
		var u Uri
		if err := c.ShouldBindUri(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, u)
	})
	router.Run(":8080")
}
