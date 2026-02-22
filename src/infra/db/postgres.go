package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/mahopon/SmolEarl/config"
)

// PostgresPool holds the PostgreSQL connection pool
var (
	AppConfig = config.AppConfig
)

type DB struct {
	PostgresPool *pgxpool.Pool
}

// InitPostgres initializes the PostgreSQL connection pool
func InitPostgres() (*DB, error) {
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
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return &DB{PostgresPool: pool}, nil
}

// ClosePostgres closes the PostgreSQL connection pool
func (db *DB) ClosePostgres() error {
	if db.PostgresPool != nil {
		db.PostgresPool.Close()
	}
	return nil
}

// InitSchema creates the required tables if they don't exist
func (db *DB) InitSchema() error {
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

	_, err := db.PostgresPool.Exec(ctx, createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func (db *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := db.PostgresPool.Query(ctx, sql, args...)

	return rows, err
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	result, err := db.PostgresPool.Exec(ctx, sql, args...)

	return result, err
}
