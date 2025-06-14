package main

import (
	user "github.com/LiangNing7/go-example/wire/user/internal"
	"github.com/LiangNing7/go-example/wire/user/internal/config"
	"github.com/LiangNing7/go-example/wire/user/pkg/db"
)

func main() {
	cfg := &config.Config{
		MySQL: db.MySQLOptions{
			Address:  "127.0.0.1:3306",
			Database: "user",
			Username: "root",
			Password: "root",
		},
	}

	app, cleanup, err := user.NewApp(cfg)
	if err != nil {
		panic(err)
	}

	defer cleanup()
	app.Run()
}
