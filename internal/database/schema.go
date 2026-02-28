package database

import "database/sql"

// InitSchema creates the users table and seed data. Uses parameterized
// statements only—no user input is involved here.
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (id, username, email) VALUES
		(1, 'admin', 'admin@example.com'),
		(2, 'alice', 'alice@example.com'),
		(3, 'bob', 'bob@example.com');
	`)
	return err
}
