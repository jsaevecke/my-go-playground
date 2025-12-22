package gormdb

import (
	"context"
	"database/sql"
	"my-go-playground/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	// Actual SQL logic here
	return nil
}
