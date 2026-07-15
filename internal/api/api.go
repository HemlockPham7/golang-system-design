package api

import (
	"github.com/HemlockPham7/golang-system-design/internal/handler"
	"github.com/HemlockPham7/golang-system-design/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
}

type engine struct {
	app *gin.Engine
}

func NewEngine() Engine {
	app := &engine{
		app: gin.Default(),
	}
	app.initRoutes()
	return app
}

func (e *engine) Start() error {
	return e.app.Run(":8080")
}

func (e *engine) initRoutes() {
	genPassService := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	e.app.GET("/genpass", genPassHandler.GeneratePassword)
}
