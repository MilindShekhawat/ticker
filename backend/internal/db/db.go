package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Need side effect only
)

var DB *sql.DB

func Init() error {
	// TODO: also need to allow the user to set this dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	tickerDir := filepath.Join(homeDir, ".ticker")
	if err := os.MkdirAll(tickerDir, 0755); err != nil {
		return fmt.Errorf("failed to create ticker directory: %w", err)
	}

	dbPath := filepath.Join(tickerDir, "ticker.db")

	db, err := sql.Open("sqlite", dbPath)
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
	fmt.Printf("Database initialized: %s\n", dbPath)
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
