package database

import (
	"database/sql"
	"fmt"
	"pollvoting/pkg/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connect() error {
	dbURL := config.LocalConfig.PostgresqlURL

	d, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	if err := d.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db = d
	fmt.Println("Connected to Postgres successfully")
	return nil
}

func ConnectDB() (*sql.DB, error) {
	if db == nil {
		if err := connect(); err != nil {
			return nil, err
		}
	}

	return db, nil
}
