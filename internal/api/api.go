package api

import (
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/handler"
	"github.com/HemlockPham7/golang-system-design/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
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

// Start starts the application
func (e *engine) Start() error {
	return e.app.Run(":8080")
}

// ServeHTTP to test the API endpoint
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	genPassService := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	e.app.GET("/genpass", genPassHandler.GeneratePassword)
}
