package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"

	"my-go-playground/internal/adapter/http/api"
	"my-go-playground/internal/adapter/postgres/sqldb"
	"my-go-playground/internal/infrastructure/cerr"
	"my-go-playground/internal/infrastructure/logging"

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
	ctx context.Context,
	cfg configuration,
	logger *zerolog.Logger,
	runChan chan struct{},
) error {
	defer cerr.HandlePanic(recover(), debug.Stack(), logger)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port, err := strconv.Atoi(cfg.WebserverPort)
	if err != nil {
		return fmt.Errorf("atoi webserver port: %w", err)
	}

	db := sqldb.New(cfg.DatabaseDriver, cfg.DatabaseURL)
	server := api.New(db, port)

	go func() {
		logger.Info().Msgf("started server on port %q", cfg.WebserverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(fmt.Errorf("listen and serve %w", err)).Msg("error starting server")
		}
		logger.Info().Msg("server stopped")
	}()

	if runChan != nil {
		close(runChan)
	}

	waitForShutdown(ctx, server, logger)

	return nil
}

func waitForShutdown(ctx context.Context, server *http.Server, logger *zerolog.Logger) {
	var wg sync.WaitGroup
	wg.Go(func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Info().Msg("shutting down application")
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(fmt.Errorf("shutdown: %w", err)).Msg("error shutting down application")
		}
		logger.Info().Msg("application shut down")
	})
	wg.Wait()
}
