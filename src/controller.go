package main

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

// CreateHandler handles POST /create requests
func (c *Controller) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if customAlias is provided in the request body
	var customAlias string
	if alias, exists := data["customAlias"]; exists {
		if aliasStr, ok := alias.(string); ok {
			customAlias = aliasStr
		}
	}

	// Call service to create entry with custom alias if provided
	id, err := c.service.Create(data, customAlias)
	if err != nil {
		http.Error(w, "Failed to create entry", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Entry created successfully",
		"id":      id,
	})
}

// GetHandler handles GET /{id} requests
func (c *Controller) GetHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := r.PathValue("path")

	// If path is empty (root path), return a default response or redirect
	if path == "" {
		http.Error(w, "Missing input", http.StatusBadRequest)
		return
	}

	// Call service to get entry
	entry, err := c.service.Get(path)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Return entry data
	w.Header().Set("Content-Type", "application/json")
	entryUrl := entry["url"].(string)
	http.Redirect(w, r, entryUrl, http.StatusFound)
}

// StatsHandler handles GET /stats/{id} requests
func (c *Controller) StatsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	stats, err := c.service.GetStats(id)
	if err != nil {
		http.Error(w, "Stats not found", http.StatusNotFound)
		return
	}

	// Return stats data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// StatusHandler handles GET /status requests
func (c *Controller) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Welcome to SmolEarl API",
		"version": "1.0",
	})
}
