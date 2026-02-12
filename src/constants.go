package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port            string
	Host            string
	RedisHost       string
	RedisPort       string
	RedisPass       string
	RedisDB         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	AppName         string
	AppVersion      string
	KafkaBroker     string
	KafkaClickTopic string
}

// LoadConfig loads configuration from .env file
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration with defaults
	config := &Config{
		Port:            getEnv("PORT", "8000"),
		Host:            getEnv("HOST", "localhost"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPass:       getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnv("REDIS_DB", "0"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "smolearl"),
		AppName:         getEnv("APP_NAME", "SmolEarl"),
		AppVersion:      getEnv("APP_VERSION", "1.0.0"),
		KafkaBroker:     getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaClickTopic: getEnv("KAFKA_CLICK_TOPIC", "click-events"),
	}

	return config
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Global configuration instance
var AppConfig = LoadConfig()
