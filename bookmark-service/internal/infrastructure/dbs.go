package infrastructure

import (
	"github.com/HemlockPham7/golang-system-design/internal/model"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/HemlockPham7/golang-system-design/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateDB(envPrefix string) *gorm.DB {
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

func CreateRedisClient(envPrefix string) *redis.Client {
	redisClient, err := redisPkg.NewClient(envPrefix)
	if err != nil {
		panic(err)
	}
	return redisClient
}
