package gormdb

import (
	"fmt"
	"my-go-playground/internal/adapter/postgres/sqldb"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDatabase struct {
	sqlDB  *sqldb.SqlDatabase
	gormDB *gorm.DB
}

func New(sqlDB *sqldb.SqlDatabase, gormCfg *gorm.Config) *GormDatabase {
	if gormCfg == nil {
		gormCfg = new(gorm.Config)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB.DB()}), gormCfg)
	if err != nil {
		panic(fmt.Errorf("gorm open: %w", err))
	}

	return &GormDatabase{
		sqlDB:  sqlDB,
		gormDB: gormDB,
	}
}

func (g *GormDatabase) Close() error {
	if err := g.sqlDB.Close(); err != nil {
		return fmt.Errorf("sql db close: %w", err)
	}
	return nil
}
