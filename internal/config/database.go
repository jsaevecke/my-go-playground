package config

import (
	"fmt"
	"time"
)

type Database struct {
	PrimaryHost   string `env:"DATABASE_PRIMARY_HOST" envDefault:"localhost"`
	SecondaryHost string `env:"DATABASE_SECONDARY_HOST"`
	DatabaseName  string `env:"DATABASE_NAME" envDefault:"app_db"`
	User          string `env:"DATABASE_USER" envDefault:"app_user"`
	Password      string `env:"DATABASE_PASSWORD" envDefault:"password"`
	Driver        string `env:"DATABASE_DRIVER" envDefault:"postgres"`
	Port          int    `env:"DATABASE_PORT" envDefault:"5432"`

	StatementTimeout time.Duration
}

func (cfg Database) ToPrimaryDSN() string {
	return fmt.Sprintf("%s dbname=%s", cfg.ToDSNNoDatabase(cfg.PrimaryHost), cfg.DatabaseName)
}
func (cfg Database) ToSecondaryDSN() string {
	return fmt.Sprintf("%s dbname=%s", cfg.ToDSNNoDatabase(cfg.SecondaryHost), cfg.DatabaseName)
}
func (cfg Database) ToDSNNoDatabase(host string) string {
	statementTimeout := cfg.StatementTimeout.Milliseconds()
	if statementTimeout == 0 {
		statementTimeout = (2 * time.Minute).Milliseconds()
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=disable application_name=%s statement_timeout=%d",
		host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		"app",
		statementTimeout,
	)
}
