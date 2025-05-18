package repository

import (
	"context"
	"fmt"
	"traffic-log/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TrafficRepository interface {
	LogTraffic(ctx context.Context, log *models.TrafficLog) error
}

type PostgreSQLTrafficRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLTrafficRepository(db *pgxpool.Pool) *PostgreSQLTrafficRepository {
	return &PostgreSQLTrafficRepository{db: db}
}

func (r *PostgreSQLTrafficRepository) LogTraffic(ctx context.Context, log *models.TrafficLog) error {
	query := `
		INSERT INTO traffic_logs 
			(timestamp, client_ip, method, url, response_code, bytes_sent, bytes_received)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		log.Timestamp,
		log.ClientIP,
		log.Method,
		log.URL,
		log.ResponseCode,
		log.BytesSent,
		log.BytesReceived,
	)
	if err != nil {
		return fmt.Errorf("failed to log traffic: %w", err)
	}
	return nil
}