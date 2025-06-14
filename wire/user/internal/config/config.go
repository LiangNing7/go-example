package config

import "github.com/LiangNing7/go-example/wire/user/pkg/db"

type Config struct {
	MySQL db.MySQLOptions `json:"mysql" yaml:"mysql"`
}
