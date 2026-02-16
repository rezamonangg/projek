package common

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *Config) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Str("appname", "projek").Logger()
	logger.Level(zerolog.InfoLevel)
	return logger
}
