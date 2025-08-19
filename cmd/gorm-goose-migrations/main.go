package main

import (
	"context"
	"embed"
	"fmt"
	"math"
	"os"
	"runtime/debug"

	"my-go-playground/internal/database/sqldb"
	"my-go-playground/internal/logging"

	"github.com/pressly/goose"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var migrationPath string = "migration"

//go:embed migration/*.sql
var embedMigrations embed.FS

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

	migrations, err := goose.CollectMigrations(migrationPath, 0, math.MaxInt64)
	if err != nil {
		panic(fmt.Errorf("migration: goose collect migrations: %w", err))
	}

	if len(migrations) == 0 {
		logger.Info().Msg("no migrations found")
		return nil
	}

	logger.Info().Msgf("applying migrations in %q...", migrationPath)
	for _, migration := range migrations {
		logger.Info().Int64("version", migration.Version).Str("source", migration.Source).Msg("found migration")
	}

	sqlDB := sqldb.New(cfg.DatabaseDriver, cfg.DatabaseURL)
	if err := goose.Up(sqlDB.DB(), migrationPath); err != nil {
		panic(fmt.Errorf("migration: goose up: %w", err))
	}

	if err := goose.Status(sqlDB.DB(), migrationPath); err != nil {
		panic(fmt.Errorf("migration: goose status: %w", err))
	}
	logger.Info().Msg("migrations applied successfully")

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
