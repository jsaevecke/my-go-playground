package gormdb

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

func New(sqlDB *sql.DB, gormCfg *gorm.Config) *database {
	if gormCfg == nil {
		gormCfg = new(gorm.Config)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormCfg)
	if err != nil {
		panic(fmt.Errorf("sql open: %w", err))
	}

	return &database{db: db}
}
