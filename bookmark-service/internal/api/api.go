package api

import (
	"fmt"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/docs"
	_ "github.com/HemlockPham7/golang-system-design/docs"
	"github.com/HemlockPham7/golang-system-design/internal/api/middleware"
	"github.com/HemlockPham7/golang-system-design/internal/handler"
	bookmarkHdl "github.com/HemlockPham7/golang-system-design/internal/handler/bookmark"
	userHdl "github.com/HemlockPham7/golang-system-design/internal/handler/user"
	"github.com/HemlockPham7/golang-system-design/internal/repository"
	bookmarkRepo "github.com/HemlockPham7/golang-system-design/internal/repository/bookmark"
	userRepo "github.com/HemlockPham7/golang-system-design/internal/repository/user"
	"github.com/HemlockPham7/golang-system-design/internal/service"
	bookmarkSvc "github.com/HemlockPham7/golang-system-design/internal/service/bookmark"
	userSvc "github.com/HemlockPham7/golang-system-design/internal/service/user"
	"github.com/HemlockPham7/golang-system-design/pkg/jwtutils"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/gorm"
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
	dbClient    *gorm.DB
	jwtGen      jwtutils.JWTGenerator
	jwtVal      jwtutils.JWTValidator
}

type EngineOpts struct {
	App         *gin.Engine
	Cfg         *Config
	RedisClient *redis.Client
	DbClient    *gorm.DB
	JwtGen      jwtutils.JWTGenerator
	JwtVal      jwtutils.JWTValidator
}

// NewEngine creates a new engine
func NewEngine(opts *EngineOpts) Engine {
	app := &engine{
		app:         opts.App,
		cfg:         opts.Cfg,
		redisClient: opts.RedisClient,
		dbClient:    opts.DbClient,
		jwtGen:      opts.JwtGen,
		jwtVal:      opts.JwtVal,
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

type handlers struct {
	genPassHandler     handler.GenPass
	urlHandler         handler.ShortenUrl
	healthCheckHandler handler.HealthCheck
	userHandler        userHdl.Handler
	bookmarkHandler    bookmarkHdl.Handler
}

func (e *engine) initHandlers() *handlers {
	genPassService := utils.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	urlStorage := repository.NewUrlStorage(e.redisClient)
	urlService := service.NewShortenUrl(urlStorage, genPassService)
	urlHandler := handler.NewShortenUrl(urlService)

	pingRepo := repository.NewPing(e.redisClient)
	healthCheckService := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID, pingRepo)
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)

	userRepository := userRepo.NewSqlRepository(e.dbClient)
	hasher := utils.NewHasher()
	userService := userSvc.NewService(userRepository, hasher, e.jwtGen)
	userHandler := userHdl.NewHandler(userService)

	bookmarkRepository := bookmarkRepo.NewRepository(e.dbClient)
	bookmarkService := bookmarkSvc.NewService(bookmarkRepository, genPassService)
	bookmarkHandler := bookmarkHdl.NewHandler(bookmarkService)

	return &handlers{
		genPassHandler:     genPassHandler,
		urlHandler:         urlHandler,
		healthCheckHandler: healthCheckHandler,
		userHandler:        userHandler,
		bookmarkHandler:    bookmarkHandler,
	}
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	allHandlers := e.initHandlers()

	// init middleware
	jwtAuth := middleware.NewJWTAuth(e.jwtVal)

	// genpass
	e.app.GET("/genpass", allHandlers.genPassHandler.GeneratePassword)

	// health-check
	e.app.GET("/health-check", allHandlers.healthCheckHandler.HealthCheck)

	// Init swagger routes
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	privateRoutes := e.app.Group("")
	privateRoutes.Use(jwtAuth.JWTAuth())
	{
		privateV1Routes := privateRoutes.Group("/v1")
		{
			selfRoutes := privateV1Routes.Group("/self")
			{
				selfRoutes.GET("/info", allHandlers.userHandler.GetSelfInfo)
			}

			bookmarksRoutes := privateV1Routes.Group("/bookmarks")
			{
				bookmarksRoutes.POST("/", allHandlers.bookmarkHandler.CreateBookmark)
				bookmarksRoutes.GET("/", allHandlers.bookmarkHandler.GetBookmarks)
			}
		}
	}

	publicRoutes := e.app.Group("")
	{
		publicV1Routes := publicRoutes.Group("/v1")
		{
			linksRoutes := publicV1Routes.Group("/links")
			{
				linksRoutes.POST("/shorten", allHandlers.urlHandler.ShortenLink)
				linksRoutes.GET("/redirect/:code", allHandlers.urlHandler.Redirect)
			}

			usersRoutes := publicV1Routes.Group("/users")
			{
				usersRoutes.POST("/register", allHandlers.userHandler.Register)
				usersRoutes.POST("/login", allHandlers.userHandler.Login)
			}
		}
	}

}
