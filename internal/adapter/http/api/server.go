package api

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"my-go-playground/internal/infrastructure/database"
)

type Server struct {
	db   database.Database
	port int
}

func New(db database.Database, port int) *http.Server {
	// TODO: validation
	srv := &Server{
		port: port,
		db:   db,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
