package main

import (
	"github.com/HemlockPham7/golang-system-design/internal/api"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/jwtutils"
	"github.com/HemlockPham7/golang-system-design/pkg/logger"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/HemlockPham7/golang-system-design/pkg/sqldb"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// @title Bookmark Management API
// @version 1.0.7
// @description API for managing bookmarks
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// create app config
	cfg := createAPIConfig()

	// set log level
	logger.SetLogLevel(cfg.LogLevel)

	// create redis client
	redisClient := createRedisClient("")

	// Init db
	db := createDB("")

	app := createAPIApp(cfg, redisClient, db)
	err := app.Start()
	if err != nil {
		panic(err)
	}
}

func createRedisClient(envPrefix string) *redis.Client {
	redisClient, err := redisPkg.NewClient(envPrefix)
	if err != nil {
		panic(err)
	}
	return redisClient
}

func createDB(envPrefix string) *gorm.DB {
	db, err := sqldb.NewClient(envPrefix)
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	return db
}

func createAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func createAPIApp(cfg *api.Config, redis *redis.Client, db *gorm.DB) api.Engine {
	app := gin.Default()

	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	if err != nil {
		panic(err)
	}

	jwtVal, err := jwtutils.NewJWTValidator("./public.pem")
	if err != nil {
		panic(err)
	}

	a := api.NewEngine(&api.EngineOpts{
		App:         app,
		Cfg:         cfg,
		RedisClient: redis,
		DbClient:    db,
		JwtGen:      jwtGen,
		JwtVal:      jwtVal,
	})

	return a
}
