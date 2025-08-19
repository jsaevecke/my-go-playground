package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"my-go-playground/internal/database/sqldb"
	"my-go-playground/internal/logging"

	"github.com/pressly/goose"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := loadConfiguration(os.Getenv)

	logger := logging.Init(cfg.LogLevel)
	logger = logger.With().Str(logging.FieldEnvironment, cfg.Environment).Logger()

	if err := run(context.Background(), cfg, &logger, nil); err != nil {
		log.Fatal().Err(err).Msg("error running application")
	}
}

func run(
	_ context.Context,
	cfg configuration,
	logger *zerolog.Logger,
	_ chan struct{},
) error {
	defer handlePanic(recover(), debug.Stack(), logger)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Errorf("migration: goose set dialect: %w", err))
	}

	sqlDB := sqldb.New(cfg.DatabaseDriver, cfg.DatabaseURL)
	if err := goose.Up(sqlDB.DB(), "./migration"); err != nil {
		panic(fmt.Errorf("migration: goose up: %w", err))
	}

	return nil
}

func handlePanic(r any, stack []byte, logger *zerolog.Logger) {
	if r == nil {
		return
	}

	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}

	logger.Fatal().Bytes(logging.FieldStack, stack).Err(err).Msgf("panic")
}
