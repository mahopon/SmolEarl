package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"

	config "github.com/mahopon/SmolEarl/config"
)

// PostgresPool holds the PostgreSQL connection pool
var (
	AppConfig = config.AppConfig
)

type DB struct {
	PostgresPool *pgxpool.Pool
	Metrics      *DBMetrics
}

// InitPostgres initializes the PostgreSQL connection pool
func InitPostgres(reg prometheus.Registerer) (*DB, error) {
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

	metrics := NewDBMetrics(reg)

	return &DB{PostgresPool: pool, Metrics: metrics}, nil
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
	start := time.Now()
	rows, err := db.PostgresPool.Query(ctx, sql, args...)
	db.Metrics.Latency.WithLabelValues("query").Observe(time.Since(start).Seconds())

	if err != nil {
		db.Metrics.QueryErrors.WithLabelValues("query").Inc()
	}
	db.Metrics.QueryCount.WithLabelValues("query").Inc()

	db.updatePoolMetrics()

	return rows, err
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	start := time.Now()
	result, err := db.PostgresPool.Exec(ctx, sql, args...)
	db.Metrics.Latency.WithLabelValues("exec").Observe(time.Since(start).Seconds())

	if err != nil {
		db.Metrics.QueryErrors.WithLabelValues("exec").Inc()
	}
	db.Metrics.QueryCount.WithLabelValues("exec").Inc()

	db.updatePoolMetrics()

	return result, err
}

func (db *DB) updatePoolMetrics() {
	stats := db.PostgresPool.Stat()
	db.Metrics.PoolActive.WithLabelValues().Set(float64(stats.AcquiredConns()))
	db.Metrics.PoolIdle.WithLabelValues().Set(float64(stats.TotalConns() - stats.AcquiredConns()))
}
