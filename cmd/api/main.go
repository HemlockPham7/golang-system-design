package main

import "github.com/HemlockPham7/golang-system-design/internal/api"

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

	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
