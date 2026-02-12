package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize Redis
	if err := InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer CloseRedis()

	// Initialize Kafka
	if err := InitKafka(); err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	defer CloseKafka()

	// Initialize PostgreSQL
	if err := InitPostgres(); err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer ClosePostgres()

	// Initialize schema
	if err := InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	repo := NewEntryRepository(PostgresPool)
	service := NewService()
	service.SetRepository(repo)
	controller := NewController(service)
	router := NewRouter(controller)

	mux := http.NewServeMux()
	router.Init(mux)

	// Middleware
	loggedMux := LoggingMiddleware(mux)
	corsMux := CORSMiddleware(loggedMux)

	// Server config
	addr := ":" + AppConfig.Port
	fmt.Printf("Starting server on http://%s:%s\n", AppConfig.Host, AppConfig.Port)

	if err := http.ListenAndServe(addr, corsMux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
