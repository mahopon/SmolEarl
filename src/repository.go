package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EntryRepository handles database operations for entries
type EntryRepository struct {
	pool *pgxpool.Pool
}

// NewEntryRepository creates a new EntryRepository
func NewEntryRepository(pool *pgxpool.Pool) *EntryRepository {
	return &EntryRepository{
		pool: pool,
	}
}

// Entry represents a shortened URL entry
type Entry struct {
	ShortCode  string
	OriginalURL string
	Clicks     int
	CreatedAt  time.Time
}

// Create inserts a new entry into the database
func (r *EntryRepository) Create(ctx context.Context, shortCode, originalURL string, clicks int, createdAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO entries (short_code, original_url, clicks, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (short_code) DO NOTHING",
		shortCode, originalURL, clicks, createdAt)
	return err
}

// GetByShortCode retrieves an entry by its short code
func (r *EntryRepository) GetByShortCode(ctx context.Context, shortCode string) (*Entry, error) {
	var entry Entry
	err := r.pool.QueryRow(ctx,
		"SELECT original_url, clicks, created_at FROM entries WHERE short_code = $1", shortCode).
		Scan(&entry.OriginalURL, &entry.Clicks, &entry.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	entry.ShortCode = shortCode
	return &entry, nil
}

// GetStats retrieves stats for an entry by its short code
func (r *EntryRepository) GetStats(ctx context.Context, shortCode string) (int, time.Time, error) {
	var clicks int
	var createdAt time.Time
	err := r.pool.QueryRow(ctx,
		"SELECT clicks, created_at FROM entries WHERE short_code = $1", shortCode).
		Scan(&clicks, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, time.Time{}, nil
		}
		return 0, time.Time{}, err
	}
	return clicks, createdAt, nil
}
