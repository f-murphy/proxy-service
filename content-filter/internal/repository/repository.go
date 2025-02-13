package repository

import (
	"content-filter/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type FilterRepository interface {
	CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error)
	GetBlockUrls(ctx context.Context) ([]string, error)
	GetBlacklistKeywords(ctx context.Context) ([]string, error)
}

type PostgreSQLFilterRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLFilterRepository(db *pgxpool.Pool) *PostgreSQLFilterRepository {
	return &PostgreSQLFilterRepository{db: db}
}

func (r *PostgreSQLFilterRepository) CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error) {
	query := "INSERT INTO filter_urls (id, url) VALUES ($1, $2) RETURNING id"
	var id string
	err := r.db.QueryRow(ctx, query, filter_url.ID, filter_url.Url).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("unable to insert data: %w", err)
	}
	return id, nil
}

func (r *PostgreSQLFilterRepository) GetBlockUrls(ctx context.Context) ([]string, error) {
	query := "SELECT value FROM blacklist WHERE type = 'url'"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query keywords: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, fmt.Errorf("unable to scan keyword: %w", err)
		}
		urls = append(urls, keyword)
	}
	return urls, nil
}

func (r *PostgreSQLFilterRepository) GetBlacklistKeywords(ctx context.Context) ([]string, error) {
	query := "SELECT value FROM blacklist WHERE type = 'keyword'"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query keywords: %w", err)
	}
	defer rows.Close()

	var keywords []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, fmt.Errorf("unable to scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}
	return keywords, nil
}
