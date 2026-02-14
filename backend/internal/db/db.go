package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Need side effect only
)

var DB *sql.DB

func Init(dbPath string) error {
	var finalPath string

	if dbPath != "" {
		finalPath = dbPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		tickerDir := filepath.Join(homeDir, ".ticker")
		if err := os.MkdirAll(tickerDir, 0755); err != nil {
			return fmt.Errorf("failed to create ticker directory: %w", err)
		}

		finalPath = filepath.Join(tickerDir, "ticker.db")
	}

	dir := filepath.Dir(finalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", finalPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// sql.Open() doesn't connect, so we try pinging
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Sqlite uses 1 writer
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	DB = db
	fmt.Printf("Database initialized: %s\n", finalPath)
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
