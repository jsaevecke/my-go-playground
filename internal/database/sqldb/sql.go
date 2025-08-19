package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	driver string
	url    string

	db *sql.DB
}

func New(driver, url string) *database {
	log.Printf("open database connection to %s with driver %s", url, driver)
	db, err := sql.Open(driver, url)
	if err != nil {
		panic(fmt.Errorf("sql open: %w", err))
	}
	log.Println("connected to database")

	return &database{
		driver: driver,
		url:    url,

		db: db,
	}
}

func (s *database) DB() *sql.DB {
	return s.db
}

func (s *database) Close() error {
	log.Printf("closing db connection to %s with driver %s", s.url, s.driver)
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	log.Println("db connection closed")

	return nil
}
