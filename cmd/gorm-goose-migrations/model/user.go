package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.Null[time.Time]

	Name  string
	Email string
}
