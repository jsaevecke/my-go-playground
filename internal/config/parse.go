package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Parse(cfg any) error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load env file: %w", err)
	} // load .env if present
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse env into cfg: %w", err)
	}
	return nil
}
