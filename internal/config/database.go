package config

import (
	"fmt"
	"time"
)

type Database struct {
	DatabasePrimaryHost   string `env:"DATABASE_PRIMARY_HOST" envDefault:"localhost"`
	DatabaseSecondaryHost string `env:"DATABASE_SECONDARY_HOST"`
	DatabaseName          string `env:"DATABASE_NAME" envDefault:"app_db"`
	DatabaseUser          string `env:"DATABASE_USER" envDefault:"app_user"`
	DatabasePassword      string `env:"DATABASE_PASSWORD" envDefault:"password"`
	DatabaseDriver        string `env:"DATABASE_DRIVER" envDefault:"postgres"`
	DatabasePort          int    `env:"DATABASE_PORT" envDefault:"5432"`

	DatabaseStatementTimeout time.Duration
}

func (cfg Database) ToPrimaryDSN() string {
	return fmt.Sprintf("%s dbname=%s", cfg.ToDSNNoDatabase(cfg.DatabasePrimaryHost), cfg.DatabaseName)
}
func (cfg Database) ToSecondaryDSN() string {
	return fmt.Sprintf("%s dbname=%s", cfg.ToDSNNoDatabase(cfg.DatabaseSecondaryHost), cfg.DatabaseName)
}
func (cfg Database) ToDSNNoDatabase(host string) string {
	statementTimeout := cfg.DatabaseStatementTimeout.Milliseconds()
	if statementTimeout == 0 {
		statementTimeout = (2 * time.Minute).Milliseconds()
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=disable application_name=%s statement_timeout=%d",
		host,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		"app",
		statementTimeout,
	)
}
