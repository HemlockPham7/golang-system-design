package main

import "github.com/HemlockPham7/golang-system-design/internal/api"

func main() {
	app := api.NewEngine()
	err := app.Start()
	if err != nil {
		panic(err)
	}
}
