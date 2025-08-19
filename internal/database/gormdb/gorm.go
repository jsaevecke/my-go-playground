package gormdb

import (
	"fmt"
	"my-go-playground/internal/database/sqldb"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDatabase struct {
	sqlDB  *sqldb.SqlDatabase
	gormDB *gorm.DB

	gormMigrator gorm.Migrator
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

		gormMigrator: gormDB.Migrator(),
	}
}

func (g *GormDatabase) DB() *gorm.DB {
	return g.gormDB
}

func (g *GormDatabase) Migrator() gorm.Migrator {
	return g.gormMigrator
}

func (g *GormDatabase) Close() error {
	if err := g.sqlDB.Close(); err != nil {
		return fmt.Errorf("sql db close: %w", err)
	}

	return nil
}
