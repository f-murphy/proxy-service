package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AnomalyRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLAnomalyRepository(db *pgxpool.Pool) *AnomalyRepository {
	return &AnomalyRepository{db: db}
}

func (r *AnomalyRepository) IsAllowed(ip string) (bool, error) {
	ctx := context.Background()
	var requestCount int
	var lastRequest time.Time

	err := r.db.QueryRow(ctx, `
        SELECT request_count, last_request 
        FROM request_logs 
        WHERE ip_address = $1
    `, ip).Scan(&requestCount, &lastRequest)

	if err != nil {
		if err.Error() == "no rows in result set" {
			_, err := r.db.Exec(ctx, `
                INSERT INTO request_logs (ip_address, request_count, last_request) 
                VALUES ($1, $2, $3)
            `, ip, 1, time.Now())
			return true, err
		}
		return false, err
	}

	if time.Since(lastRequest) > time.Minute {
		_, err := r.db.Exec(ctx, `
            UPDATE request_logs 
            SET request_count = $1, last_request = $2 
            WHERE ip_address = $3
        `, 1, time.Now(), ip)
		return true, err
	}

	if requestCount >= 100 {
		return false, nil
	}

	_, err = r.db.Exec(ctx, `
        UPDATE request_logs 
        SET request_count = request_count + 1, last_request = $1 
        WHERE ip_address = $2
    `, time.Now(), ip)

	return err == nil, err
}
