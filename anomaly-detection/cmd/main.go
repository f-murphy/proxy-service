package main

import (
	"anomaly-detection/internal/handler"
	"anomaly-detection/internal/repository"
	"anomaly-detection/internal/service"
	logger "anomaly-detection/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		logrus.WithError(err).Fatal("Error loading .env file")
	}

	logFile, err := logger.InitLogger()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading logrus")
	}
	logrus.Info("logFile initialized successfully")
	defer logFile.Close()

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize db")
	}
	logrus.Info("Database connected successfully")
	defer conn.Close(context.Background())
	
    repo := repository.NewPostgreSQLAnomalyRepository(conn)
    service := service.NewAnomalyService(repo)
    handler := handler.AnomalyHandler(service)

	logrus.Info("Starting proxy server on :8081")
    if err := http.ListenAndServe(":8081", handler); err != nil {
        log.Fatal(err)
    }
}