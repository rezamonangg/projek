package common

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "timestamp"

	logger := zerolog.New(os.Stdout).
		With().
		Str("appname", "projek").
		Timestamp().
		Logger()

	logger.Level(zerolog.InfoLevel)
	return logger
}
