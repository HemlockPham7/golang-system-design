package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Debug().Msg("Hello World")
	log.Info().Msg("Hello World")
	log.Warn().Msg("Hello World")
	log.Error().Msg("Hello World")
}
