package main

import (
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/sqldb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Debug().Msg("Hello World")
	log.Info().Msg("Hello World")
	log.Warn().Msg("Hello World")
	log.Error().Msg("Hello World")

	db, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{})
	db.Create(&model.User{
		Username:    "test4",
		Password:    "test4",
		Email:       "test4",
		DisplayName: "test4",
	})
}
