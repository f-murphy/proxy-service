package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"traffic-log/internal/metrics"
	"traffic-log/internal/proxy"
	"traffic-log/internal/repository"
	"traffic-log/internal/service"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/proxy-service?sslmode=disable"


	fmt.Println("Connecting to DB:", dbURL)

	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	trafficRepo := repository.NewPostgreSQLTrafficRepository(dbpool)
	trafficService := service.NewTrafficService(trafficRepo)

	target :="http://localhost:3000"
	proxyHandler, err := proxy.NewReverseProxy(target, trafficService)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	metrics.Init()

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", proxyHandler)

	port := ":8080"
	fmt.Printf("Proxy server running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
