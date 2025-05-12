package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 为 multipart forms 设置较低的内存限制(默认 32 MiB.)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/public", "./public")
	router.POST("/upload", func(c *gin.Context) {
		// Source.
		file, err := c.FormFile("f1")
		if err != nil {
			c.String(http.StatusOK, "get form err: %s", err.Error())
			return
		}

		filename := filepath.Base(file.Filename)
		dst := "./public/" + filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		c.String(http.StatusOK, "File %s uploaded successfully!", file.Filename)
	})

	router.POST("/uploads", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["files[]"]
		log.Println(files)

		for index, file := range files {
			log.Println(file.Filename)
			dst := fmt.Sprintf("./public/%d_%s", index, file.Filename)

			// 上传文件至指定的目录.
			c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})
	router.Run(":8080")
}
