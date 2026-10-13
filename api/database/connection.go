package database

import (
	"database/sql"
	"fmt"
	"time"

	"bitpos.com/m/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func UnganishaDatabase() (*sql.DB, error) {
	// Environment variables
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbUser := config.GetEnv("DB_USER", "postgres")
	dbPassword := config.GetEnv("DB_PASSWORD", "postgres")
	dbName := config.GetEnv("DB_NAME", "countbox_db")
	sslMode := config.GetEnv("DB_SSLMODE", "disable")

	// PostgreSQL DSN
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		sslMode,
	)

	// Open connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Ping test
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	return db, nil
}
