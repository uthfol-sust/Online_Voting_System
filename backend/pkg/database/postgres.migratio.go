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
		{ //3
			"poll_options",
			`CREATE TABLE IF NOT EXISTS poll_options (
               option_id SERIAL PRIMARY KEY,
			   poll_id   INT NOT NULL,
			   option_image VARCHAR(255),
			   score   INT DEFAULT 0,
			   option_text VARCHAR(155) NOT NULL,

			   CONSTRAINT fk_poll_id FOREIGN KEY(poll_id)
			                        REFERENCES polls(poll_id)
									ON DELETE CASCADE
			);`,
		},
		{ //4
			"votes",
			`CREATE TABLE IF NOT EXISTS votes (
				vote_id SERIAL PRIMARY KEY,
				poll_id INT NOT NULL,
				option_id INT NOT NULL,
				user_id INT NOT NULL,
				voted_at TIMESTAMP DEFAULT NOW(),

				CONSTRAINT fk_vote_poll
					FOREIGN KEY (poll_id)
					REFERENCES polls(poll_id)
					ON DELETE CASCADE,

				CONSTRAINT fk_vote_option
					FOREIGN KEY (option_id)
					REFERENCES poll_options(option_id)
					ON DELETE CASCADE,

				CONSTRAINT fk_vote_user
					FOREIGN KEY (user_id)
					REFERENCES users(user_id)
					ON DELETE CASCADE,

				CONSTRAINT unique_user_poll_vote UNIQUE (poll_id, user_id)
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
