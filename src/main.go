package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mahopon/SmolEarl/config"
	"github.com/mahopon/SmolEarl/infra/db"
	infra_prom "github.com/mahopon/SmolEarl/infra/prometheus"
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

	dbClient, err := db.InitPostgres(reg)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	if err := dbClient.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	httpMetrics := infra_prom.NewHTTPMetrics(reg)

	repo := NewEntryRepository(dbClient.PostgresPool)
	service := NewService()
	service.SetRepository(repo)
	service.SetRedis(redisClient)
	controller := NewController(service)
	router := NewRouter(controller).Init()
	linkRouter := NewLinkRouter(controller).Init()

	mux := http.NewServeMux()
	mux.Handle("/link/", http.StripPrefix("/link", linkRouter))
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.Handle("/", router)
	handler := StripTrailingSlashMiddleware(mux)
	handler = LoggingMiddleware(handler)
	handler = CORSMiddleware(handler)
	handler = PrometheusHTTPMiddleware(httpMetrics)(handler)

	addr := ":" + config.AppConfig.Port
	fmt.Printf("Starting server on http://%s:%s\n", config.AppConfig.Host, config.AppConfig.Port)
	fmt.Printf("Prometheus metrics available at http://%s:%s/metrics\n", config.AppConfig.Host, config.AppConfig.Port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
