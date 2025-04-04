package main

import (
	_ "github.com/joho/godotenv/autoload"
)

type configuration struct {
	Environment    string
	LogLevel       string
	WebserverPort  string
	DatabaseURL    string
	DatabaseDriver string
}

func loadConfiguration(getenv func(string) string) configuration {
	return configuration{
		Environment:    getenv("ENVIRONMENT"),
		LogLevel:       getenv("LOG_LEVEL"),
		WebserverPort:  getenv("WEBSERVER_PORT"),
		DatabaseURL:    getenv("DATABASE_URL"),
		DatabaseDriver: getenv("DATABASE_DRIVER"),
	}
}
