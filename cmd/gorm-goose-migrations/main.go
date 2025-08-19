package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"my-go-playground/internal/database"
	"my-go-playground/internal/logging"

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
	defer handlePanic(recover(), debug.Stack(), logger)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db := database.New(cfg.DatabaseDriver, cfg.DatabaseURL)

	if runChan != nil {
		close(runChan)
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

	logger.Fatal().
		Bytes(logging.FieldStack, stack).
		Err(err).
		Msgf("panic")
}
