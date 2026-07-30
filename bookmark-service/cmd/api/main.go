package main

import (
	"github.com/HemlockPham7/golang-system-design/internal/api"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/logger"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/HemlockPham7/golang-system-design/pkg/sqldb"
)

// @title Bookmark Management API
// @version 1.0.7
// @description API for managing bookmarks
// @BasePath /
func main() {
	// create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	// set log level
	logger.SetLogLevel(cfg.LogLevel)

	// create redis client
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	// Init db
	db, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg, redisClient, db)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
