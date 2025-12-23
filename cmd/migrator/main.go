package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	gooseadapter "my-go-playground/internal/adapter/goose"
	"my-go-playground/internal/adapter/postgres/sqldb"
	"my-go-playground/internal/config"
	"my-go-playground/internal/domain/migration"
	"my-go-playground/internal/infrastructure/cerr"
	"my-go-playground/internal/infrastructure/logging"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var cfg configuration
	if err := config.Parse(os.Getenv); err != nil {
		log.Fatal().Err(err).Msg("error parsing configuration")
	}

	logger := logging.Init(cfg.AppLogLevel)
	logger = logger.With().Str(logging.FieldEnvironment, cfg.AppEnvironment).Logger()
	if err := run(context.Background(), cfg, &logger, nil); err != nil {
		log.Fatal().Err(err).Msg("error during migration")
	}
}

func run(
	ctx context.Context,
	cfg configuration,
	logger *zerolog.Logger,
	_ chan struct{},
) error {
	defer cerr.HandlePanic(recover(), debug.Stack(), logger)

	sqlDB := sqldb.New(cfg.DatabaseDriver, cfg.DatabasePrimaryHost)
	gooseAdapter := gooseadapter.NewGooseAdapter(sqlDB.DB(), "migrations")

	migratorService := migration.New(gooseAdapter)
	if err := migratorService.Up(ctx); err != nil {
		return fmt.Errorf("up: %w", err)
	}
	if err := migratorService.Status(ctx); err != nil {
		return fmt.Errorf("status: %w", err)
	}
	logger.Info().Msg("migrations applied successfully")
	return nil
}
