package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

	if err = runMigrations("internal/db/migrations"); err != nil {
		log.Fatal(err)
	}

	createTables()
}

func runMigrations(migrationsDir string) error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file)

		applied, err := isMigrationApplied(version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		upSQL := extractUpMigration(string(content))
		if strings.TrimSpace(upSQL) == "" {
			continue
		}

		tx, err := DB.Begin()
		if err != nil {
			return err
		}

		if _, err = tx.Exec(upSQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s: %w", version, err)
		}

		if _, err = tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func isMigrationApplied(version string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func extractUpMigration(content string) string {
	upMarker := "-- +goose Up"
	downMarker := "-- +goose Down"

	upIndex := strings.Index(content, upMarker)
	if upIndex >= 0 {
		content = content[upIndex+len(upMarker):]
	}

	downIndex := strings.Index(content, downMarker)
	if downIndex >= 0 {
		content = content[:downIndex]
	}

	return strings.TrimSpace(content)
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
