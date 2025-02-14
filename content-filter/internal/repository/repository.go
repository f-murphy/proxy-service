package repository

import (
	"content-filter/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Интерфейс репозитория
type FilterRepository interface {
	CreateBlockURL(ctx context.Context, filterURL *models.Filter_urls) (string, error)
	GetBlockURLs(ctx context.Context) ([]string, error)
	GetBlacklistKeywords(ctx context.Context) ([]string, error)
	BlockIP(ctx context.Context, ip string) error
	IsIPBlocked(ctx context.Context, ip string) (bool, error)
}

// Реализация репозитория
type PostgreSQLFilterRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLFilterRepository(db *pgxpool.Pool) *PostgreSQLFilterRepository {
	return &PostgreSQLFilterRepository{db: db}
}

// CreateBlockURL добавляет URL в чёрный список
func (r *PostgreSQLFilterRepository) CreateBlockURL(ctx context.Context, filterURL *models.Filter_urls) (string, error) {
	query := "INSERT INTO filter_urls (id, url, action) VALUES ($1, $2, $3) RETURNING id"
	var id string
	err := r.db.QueryRow(ctx, query, filterURL.ID, filterURL.Url, filterURL.Action).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("unable to insert data: %w", err)
	}
	return id, nil
}

// GetBlockURLs возвращает список заблокированных URL
func (r *PostgreSQLFilterRepository) GetBlockURLs(ctx context.Context) ([]string, error) {
	query := "SELECT value FROM blacklist WHERE type = 'url'"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query blocks url: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, fmt.Errorf("unable to scan url: %w", err)
		}
		urls = append(urls, keyword)
	}
	return urls, nil
}

// GetBlacklistKeywords возвращает список запрещённых ключевых слов
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

// BlockIP добавляет IP-адрес в список заблокированных
func (r *PostgreSQLFilterRepository) BlockIP(ctx context.Context, ip string) error {
	query := "INSERT INTO blocked_ips (ip_address) VALUES ($1) ON CONFLICT (ip_address) DO NOTHING"
	_, err := r.db.Exec(ctx, query, ip)
	if err != nil {
		return fmt.Errorf("unable to block IP: %w", err)
	}
	return nil
}

// IsIPBlocked проверяет, заблокирован ли IP-адрес
func (r *PostgreSQLFilterRepository) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM blocked_ips WHERE ip_address = $1)"
	var exists bool
	err := r.db.QueryRow(ctx, query, ip).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("unable to check IP block status: %w", err)
	}
	return exists, nil
}