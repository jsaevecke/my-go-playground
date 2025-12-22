package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type SqlDatabase struct {
	driver string
	url    string

	db *sql.DB
}

func New(driver, url string) *SqlDatabase {
	log.Printf("open database connection to %s with driver %s", url, driver)
	db, err := sql.Open(driver, url)
	if err != nil {
		panic(fmt.Errorf("sql open: %w", err))
	}
	log.Println("connected to database")

	return &SqlDatabase{
		driver: driver,
		url:    url,

		db: db,
	}
}

func (s *SqlDatabase) DB() *sql.DB {
	return s.db
}

func (s *SqlDatabase) Close() error {
	log.Printf("closing db connection to %s with driver %s", s.url, s.driver)
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db close: %w", err)
	}
	log.Println("db connection closed")

	return nil
}
