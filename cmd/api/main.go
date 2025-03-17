package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"

	"my-go-playground/internal/database"
	"my-go-playground/internal/server"
)

func main() {
	if err := run(context.Background(), os.Getenv, nil); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	runChan chan struct{},
) error {
	defer handlePanic(recover(), debug.Stack())

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := loadConfiguration(getenv)

	port, err := strconv.Atoi(cfg.WebserverPort)
	if err != nil {
		return fmt.Errorf("atoi webserver port: %w", err)
	}

	db := database.New(cfg.DatabaseDriver, cfg.DatabaseURL)
	server := server.New(db, port)

	go func() {
		log.Printf("started server on port %d\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error starting server: %v\n", err)
		}
		log.Println("server stopped")
	}()

	if runChan != nil {
		close(runChan)
	}

	waitForShutdown(ctx, server)

	return nil
}

func handlePanic(r any, stack []byte) {
	if r == nil {
		return
	}

	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}

	log.Fatalf("panic: %v\n%s", err, stack)
}

func waitForShutdown(ctx context.Context, server *http.Server) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("shutting down application")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("error shutting down application: %v", err)
		}
		log.Println("shutdown complete")
	}()
	wg.Wait()
}
