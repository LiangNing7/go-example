package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed file/*
var embeddedAssets embed.FS

func main() {
	fileContent, err := fs.Sub(embeddedAssets, "file")
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/file", http.FS(fileContent))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
