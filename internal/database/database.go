package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"phoenix-alliance-be/internal/config"
)

// DB is the database connection pool
var DB *sql.DB

// Connect establishes a connection to the database
func Connect(cfg *config.DatabaseConfig) error {
	var err error
	DB, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

