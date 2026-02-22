package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/mahopon/SmolEarl/infra/redis"
)

// Service handles business logic for the application
type Service struct {
	repo  *EntryRepository
	redis *redis.Redis
}

// NewService creates a new Service instance
func NewService() *Service {
	return &Service{}
}

// SetRepository sets the repository for the service
func (s *Service) SetRepository(repo *EntryRepository) {
	s.repo = repo
}

// SetRedis sets the Redis client for the service
func (s *Service) SetRedis(r *redis.Redis) {
	s.redis = r
}

// Create creates a new entry with write-through to PostgreSQL
func (s *Service) Create(data map[string]any, customAlias string) (string, error) {
	incomingUrl := data["url"].(string)
	if _, err := url.Parse(incomingUrl); err != nil {
		return "", err
	}

	// Use customAlias if provided, otherwise generate a short code
	var shortCode string
	if customAlias != "" {
		shortCode = customAlias
	} else {
		shortCode = generateShortCode(incomingUrl)
	}

	// Prepare the entry data with timestamp
	entryData := map[string]any{
		"url":       incomingUrl,
		"shortCode": shortCode,
		"createdAt": time.Now().UTC().Format(time.RFC3339),
		"clicks":    0,
	}

	// Store in Redis with 24-hour expiration
	jsonData, err := json.Marshal(entryData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	ctx := context.Background()
	err = s.redis.Set(ctx, shortCode, jsonData, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to store in Redis: %w", err)
	}

	// Write-through to PostgreSQL via repository
	createdAt := time.Now().UTC()
	if err := s.repo.Create(ctx, shortCode, incomingUrl, 0, createdAt); err != nil {
		return "", fmt.Errorf("failed to store in PostgreSQL: %w", err)
	}

	return shortCode, nil
}

// Get retrieves an entry by ID (from Redis first, fallback to PostgreSQL)
func (s *Service) Get(id string) (map[string]any, error) {
	ctx := context.Background()

	// Try Redis first
	data, err := s.redis.Get(ctx, id)
	if err == nil {
		// Cache hit - parse the JSON data
		var result map[string]any
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return nil, fmt.Errorf("invalid data format: %w", err)
		}
		return result, nil
	}

	// Cache miss - try PostgreSQL via repository
	entry, err := s.repo.GetByShortCode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query PostgreSQL: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("entry not found")
	}

	// Rebuild entry data from PostgreSQL
	result := map[string]any{
		"url":       entry.OriginalURL,
		"shortCode": entry.ShortCode,
		"createdAt": entry.CreatedAt.Format(time.RFC3339),
		"clicks":    entry.Clicks,
	}

	// Repopulate Redis cache
	jsonData, _ := json.Marshal(result)
	s.redis.Set(ctx, id, jsonData, 24*time.Hour)

	return result, nil
}

// GetStats retrieves statistics for an entry by ID (from Redis first, fallback to PostgreSQL)
func (s *Service) GetStats(id string) (map[string]any, error) {
	ctx := context.Background()

	// Try Redis first
	data, err := s.redis.Get(ctx, id)
	var clicks int
	var createdAt string
	var size int

	if err == nil {
		// Cache hit - parse the JSON data
		var entryData map[string]any
		if err := json.Unmarshal([]byte(data), &entryData); err != nil {
			return nil, fmt.Errorf("invalid data format: %w", err)
		}

		if clicksVal, ok := entryData["clicks"].(float64); ok {
			clicks = int(clicksVal)
		}
		if createdAtVal, ok := entryData["createdAt"].(string); ok {
			createdAt = createdAtVal
		}
		size = len(data)
	} else {
		// Cache miss - try PostgreSQL via repository
		clicks, dbCreatedAt, err := s.repo.GetStats(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to query PostgreSQL: %w", err)
		}
		if clicks == 0 && dbCreatedAt.IsZero() {
			return nil, fmt.Errorf("entry not found")
		}
		createdAt = dbCreatedAt.Format(time.RFC3339)
		size = len(id) // Approximate size
	}

	stats := map[string]any{
		"entry_id":  id,
		"clicks":    clicks,
		"createdAt": createdAt,
		"size":      size,
	}

	return stats, nil
}
