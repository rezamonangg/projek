package common

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *Config) zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("appname", "projek").Logger()
	logger.Level(zerolog.InfoLevel)
	return logger
}
