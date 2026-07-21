package api

import (
	"fmt"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/docs"
	_ "github.com/HemlockPham7/golang-system-design/docs"
	"github.com/HemlockPham7/golang-system-design/internal/handler"
	"github.com/HemlockPham7/golang-system-design/internal/repository"
	"github.com/HemlockPham7/golang-system-design/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine interface for starting the application
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct for starting the application
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

// NewEngine creates a new engine
func NewEngine(cfg *Config, redisClient *redis.Client) Engine {
	app := &engine{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	app.initRoutes()
	return app
}

// Start starts the application
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP to test the API endpoint
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	genPassService := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	urlStorage := repository.NewUrlStorage(e.redisClient)
	urlService := service.NewShortenUrl(urlStorage, genPassService)
	urlHandler := handler.NewShortenUrl(urlService)

	// genpass
	e.app.GET("/genpass", genPassHandler.GeneratePassword)

	// Init swagger routes
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init v1 routes
	v1Routes := e.app.Group("/v1")
	{
		v1Routes.POST("/links/shorten", urlHandler.ShortenLink)
		v1Routes.GET("/links/redirect/:code", urlHandler.Redirect)
	}
}
