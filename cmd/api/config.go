package main

import (
	"my-go-playground/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

type configuration struct {
	config.App
	config.Database
	config.HttpServer
}
