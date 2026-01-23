package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init initializes the database connection
func Init() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://galleryblue:galleryblue@localhost:5432/galleryblue?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
