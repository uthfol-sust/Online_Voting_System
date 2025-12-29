package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func execMigration(db *sql.DB, name, query string) error {
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("migration failed (%s): %w", name, err)
	}
	return nil
}

func AutoMigrate(db *sql.DB) error {
	if db == nil {
		return errors.New("migrations: nil db provided")
	}

	migration := []struct {
		name  string
		query string
	}{
		{ //1
			"users",
			`CREATE TABLE IF NOT EXISTS users (
				user_id SERIAL PRIMARY KEY,
				user_name VARCHAR(100),
				password VARCHAR(250),
				email VARCHAR(50) UNIQUE,
				role VARCHAR(20) DEFAULT 'user',
				image VARCHAR(255),
				created_at TIMESTAMP DEFAULT NOW()
         );`,
		},
		{ //2
			"polls",
			`CREATE TABLE IF NOT EXISTS polls (
				poll_id SERIAL PRIMARY KEY,
				title VARCHAR(150) NOT NULL,
				description VARCHAR(250),
				created_by INT NOT NULL,
				is_active BOOLEAN DEFAULT TRUE,
				expires_at TIMESTAMP,
				created_at TIMESTAMP DEFAULT NOW(),

				CONSTRAINT fk_created_by
					FOREIGN KEY (created_by)
					REFERENCES users(user_id)
					ON DELETE CASCADE
			);`,
		},
	}

	for _, m := range migration {
		if err := execMigration(db, m.name, m.query); err != nil {
			return err
		}
	}

	log.Println("Database migration completed successfully!")
	return nil
}
