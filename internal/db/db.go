package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(path string) {
	var err error

	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			vote INTEGER NOT NULL,
			PRIMARY KEY (post_id, user_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS comment_likes (
			comment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			vote INTEGER NOT NULL,
			PRIMARY KEY (comment_id, user_id),
			FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range tables {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("Failed to execute DDL query: %s, error: %s", query, err.Error())
		}
	}

	_, _ = DB.Exec("ALTER TABLE users ADD COLUMN favorite_instrument TEXT DEFAULT ''")
	_, _ = DB.Exec("ALTER TABLE users ADD COLUMN preferred_genres TEXT DEFAULT ''")
	_, _ = DB.Exec("ALTER TABLE users ADD COLUMN profile_theme TEXT DEFAULT 'default'")
	_, _ = DB.Exec("ALTER TABLE users ADD COLUMN custom_status TEXT DEFAULT ''")
	_, _ = DB.Exec("ALTER TABLE comments ADD COLUMN parent_id INTEGER DEFAULT NULL")
}