package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

type Config struct {
	Services map[string]string `json:"services"`
}

var serviceURLs map[string]string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize service URLs from environment variables or use defaults
	serviceURLs = map[string]string{
		"database":    getEnv("DATABASE_URL", "http://database:5432"),
		"cache":       getEnv("CACHE_URL", "http://redis:6379"),
		"api":         getEnv("API_URL", "http://api:8000"),
		"monitoring":  getEnv("MONITORING_URL", "http://prometheus:9090"),
	}

	// Register routes
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/health/detailed", detailedHealthCheckHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("Health Checker starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func detailedHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	// Check each service
	for service, url := range serviceURLs {
		if checkServiceHealth(url) {
			status.Services[service] = "UP"
		} else {
			status.Services[service] = "DOWN"
			status.Status = "DEGRADED"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
		"requests":  requestCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func checkServiceHealth(url string) bool {
	// Simulate service check - in a real app, you'd actually try to connect
	// This is just for demonstration
	time.Sleep(100 * time.Millisecond)
	return true
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var (
	startTime    = time.Now()
	requestCount = 0
)
