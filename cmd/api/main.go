package main

import (
	"github.com/HemlockPham7/golang-system-design/internal/api"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0.0
// @description API for managing bookmarks
// @BasePath /
func main() {
	// create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	// create redis client
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
