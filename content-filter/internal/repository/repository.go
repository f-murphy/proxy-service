package repository

import (
	"content-filter/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type FilterRepository interface {
	CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error)
	GetBlockUrls(ctx context.Context) ([]models.Filter_urls, error)
}

type PostgreSQLFilterRepository struct {
	conn *pgx.Conn
}

func NewPostgreSQLFilterRepository(conn *pgx.Conn) *PostgreSQLFilterRepository {
	return &PostgreSQLFilterRepository{conn: conn}
}

func (r *PostgreSQLFilterRepository) CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error) {
	query := "INSERT INTO filter_urls (id, url) VALUES ($1) RETURNING id"
	args := pgx.NamedArgs{
		"id":  filter_url.ID,
		"url": filter_url.Url,
	}

	var id string
	err := r.conn.QueryRow(context.Background(), query, args).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("unable to insert data: %w", err)
	}
	return id, nil
}

func (r *PostgreSQLFilterRepository) GetBlockUrls(ctx context.Context) ([]models.Filter_urls, error) {
	query := "SELECT id, url, action FROM filter_urls"
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Filter_urls])
}
