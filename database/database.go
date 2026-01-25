package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"silic0n-wiki/config"
)

var DB *sql.DB

func Connect() error {
	cfg := config.AppConfig.Database

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
