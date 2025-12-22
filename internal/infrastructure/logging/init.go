package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func Init(level string) zerolog.Logger {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "timestamp"

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if level == "" {
		level = "info"
	}

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	return logger.Level(logLevel).With().Logger()
}
