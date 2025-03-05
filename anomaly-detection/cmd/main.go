package main

import (
	"anomaly-detection/internal/handler"
	"anomaly-detection/internal/repository"
	"anomaly-detection/internal/service"
	"anomaly-detection/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		logrus.WithError(err).Fatal("Error loading .env file - ", err)
	}

	logFile, err := utils.InitLogger()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading logrus")
	}
	logrus.Info("logFile initialized successfully")
	defer logFile.Close()

	conn, err := pgxpool.Connect(context.Background(), fmt.Sprintf(
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
	defer conn.Close()

	repo := repository.NewPostgreSQLAnomalyRepository(conn)
	service := service.NewAnomalyService(repo)
	handler := handler.NewAnomalyHandler(service, os.Getenv("TARGET_SERVER"))

	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: handler.Middleware(http.DefaultServeMux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown: ", err)
	}
	logrus.Info("Server exiting")
}
