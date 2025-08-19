package main

import (
	_ "github.com/joho/godotenv/autoload"
)

type configuration struct {
	Environment           string
	LogLevel              string
	DatabaseURL           string
	DatabaseDriver        string
	DatabaseMigrationPath string
}

func loadConfiguration(getenv func(string) string) configuration {
	return configuration{
		Environment:           getenv("ENVIRONMENT"),
		LogLevel:              getenv("LOG_LEVEL"),
		DatabaseURL:           getenv("DATABASE_URL"),
		DatabaseDriver:        getenv("DATABASE_DRIVER"),
		DatabaseMigrationPath: getenv("DATABASE_MIGRATION_PATH"),
	}
}
