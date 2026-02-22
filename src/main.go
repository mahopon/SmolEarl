package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mahopon/SmolEarl/config"
	"github.com/mahopon/SmolEarl/infra/db"
	"github.com/mahopon/SmolEarl/infra/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewRegistry()

	redisClient, err := redis.InitRedis()
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	dbClient, err := db.InitPostgres()
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	if err := dbClient.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	repo := NewEntryRepository(dbClient.PostgresPool)
	service := NewService()
	service.SetRepository(repo)
	service.SetRedis(redisClient)
	controller := NewController(service)
	router := NewRouter(controller)

	mux := http.NewServeMux()
	router.Init(mux)

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	addr := ":" + config.AppConfig.Port
	fmt.Printf("Starting server on http://%s:%s\n", config.AppConfig.Host, config.AppConfig.Port)
	fmt.Printf("Prometheus metrics available at http://%s:%s/metrics\n", config.AppConfig.Host, config.AppConfig.Port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
