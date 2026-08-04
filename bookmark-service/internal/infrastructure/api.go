package infrastructure

import (
	"github.com/HemlockPham7/golang-system-design/internal/api"
	"github.com/HemlockPham7/golang-system-design/pkg/logger"
	"github.com/gin-gonic/gin"
)

func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func CreateAPI() api.Engine {
	// create app config
	cfg := CreateAPIConfig()

	// set log level
	logger.SetLogLevel(cfg.LogLevel)

	// create redis client
	redisClient := CreateRedisClient("")

	// Init db
	db := CreateDB("")

	jwtGen, jwtVal := CreateJWTProvider()

	app := gin.Default()

	return api.NewEngine(&api.EngineOpts{
		App:         app,
		Cfg:         cfg,
		RedisClient: redisClient,
		DbClient:    db,
		JwtGen:      jwtGen,
		JwtVal:      jwtVal,
	})
}
