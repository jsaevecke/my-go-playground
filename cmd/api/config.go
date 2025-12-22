package main

import (
	"my-go-playground/internal/config"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type configuration struct {
	config.App
	config.Database
	config.HttpServer
}

func loadConfiguration(getenv func(string) string) configuration {
	return configuration{
		App: config.App{
			Environment: getenv("APP_ENVIRONMENT"),
			Name:        getenv("APP_NAME"),
			LogLevel:    getenv("APP_LOG_LEVEL"),
		},
		Database: config.Database{
			PrimaryHost:      getenv("DATABASE_PRIMARY_HOST"),
			SecondaryHost:    getenv("DATABASE_SECONDARY_HOST"),
			DatabaseName:     getenv("DATABASE_NAME"),
			User:             getenv("DATABASE_USER"),
			Password:         getenv("DATABASE_PASSWORD"),
			Driver:           getenv("DATABASE_DRIVER"),
			Port:             getenv("DATABASE_PORT"),
			StatementTimeout: 2 * time.Minute,
		},
		HttpServer: config.HttpServer{
			Port: getenv("HTTP_SERVER_PORT"),
		},
	}
}
