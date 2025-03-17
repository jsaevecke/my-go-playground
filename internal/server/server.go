package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"my-go-playground/internal/service"
)

type server struct {
	db   service.Database
	port int
}

func New(serviceDatabase service.Database, port int) *http.Server {
	// TODO: validation
	srv := &server{
		port: port,
		db:   serviceDatabase,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
