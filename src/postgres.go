package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPool holds the PostgreSQL connection pool
var PostgresPool *pgxpool.Pool

// InitPostgres initializes the PostgreSQL connection pool
func InitPostgres() error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	PostgresPool = pool
	return nil
}

// ClosePostgres closes the PostgreSQL connection pool
func ClosePostgres() error {
	if PostgresPool != nil {
		PostgresPool.Close()
	}
	return nil
}

// InitSchema creates the required tables if they don't exist
func InitSchema() error {
	ctx := context.Background()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS entries (
			id SERIAL PRIMARY KEY,
			short_code VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			clicks INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`

	_, err := PostgresPool.Exec(ctx, createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}
