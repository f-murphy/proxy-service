package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AnomalyRepository interface {
	UpdateRequestMetrics(ctx context.Context, ip string) error
	GetRequestMetrics(ctx context.Context, ip string, interval string) (int, error)
	BlockIP(ctx context.Context, ip string) error
	IsIPBlocked(ctx context.Context, ip string) (bool, error)
}

type PostgreSQLAnomalyRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLAnomalyRepository(db *pgxpool.Pool) *PostgreSQLAnomalyRepository {
	return &PostgreSQLAnomalyRepository{db: db}
}

func (r *PostgreSQLAnomalyRepository) UpdateRequestMetrics(ctx context.Context, ip string) error {
	query := `
		INSERT INTO request_metrics (ip_address, request_count)
		VALUES ($1, 1)
		ON CONFLICT (ip_address)
		DO UPDATE SET request_count = request_metrics.request_count + 1, last_request = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query, ip)
	if err != nil {
		return fmt.Errorf("unable to update request metrics: %w", err)
	}
	return nil
}

func (r *PostgreSQLAnomalyRepository) GetRequestMetrics(ctx context.Context, ip string, interval string) (int, error) {
	query := `
		SELECT request_count
		FROM request_metrics
		WHERE ip_address = $1 AND last_request >= NOW() - INTERVAL '1 ' || $2
	`
	var count int
	err := r.db.QueryRow(ctx, query, ip, interval).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unable to get request metrics: %w", err)
	}
	return count, nil
}

func (r *PostgreSQLAnomalyRepository) BlockIP(ctx context.Context, ip string) error {
	query := "INSERT INTO blocked_ips (ip_address) VALUES ($1) ON CONFLICT (ip_address) DO NOTHING"
	_, err := r.db.Exec(ctx, query, ip)
	if err != nil {
		return fmt.Errorf("unable to block IP: %w", err)
	}
	return nil
}

func (r *PostgreSQLAnomalyRepository) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM blocked_ips WHERE ip_address = $1)"
	var exists bool
	err := r.db.QueryRow(ctx, query, ip).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("unable to check IP block status: %w", err)
	}
	return exists, nil
}
