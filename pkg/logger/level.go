package logger

import "github.com/rs/zerolog"

func SetLogLevel(levelStr string) {
	level := zerolog.NoLevel
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.NoLevel
	}
	zerolog.SetGlobalLevel(level)
}
