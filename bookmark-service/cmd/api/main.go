package main

import "github.com/HemlockPham7/golang-system-design/internal/infrastructure"

// @title Bookmark Management API
// @version 1.0.7
// @description API for managing bookmarks
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Init api
	app := infrastructure.CreateAPI()

	// Run api
	err := app.Start()
	if err != nil {
		panic(err)
	}
}
