package goose

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Adapter struct {
	db   *sql.DB
	path string
}

func NewGooseAdapter(db *sql.DB, path string) *Adapter {
	return &Adapter{db: db, path: path}
}

func (a *Adapter) Up(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.UpContext(ctx, a.db, a.path); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func (a *Adapter) Status(ctx context.Context) error {
	if err := goose.StatusContext(ctx, a.db, a.path); err != nil {
		return fmt.Errorf("goose status: %w", err)
	}
	return nil
}
