package service

import (
	"content-filter/internal/repository"
	"content-filter/models"
	"context"
	"strings"
)

type ICpntentFilterService interface {
	CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error)
}

type FilterService struct {
	repo repository.FilterRepository
}

func NewFilterService(repo repository.FilterRepository) *FilterService {
	return &FilterService{repo: repo}
}

func (s *FilterService) CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error) {
	return s.repo.CreateBlockUrl(ctx, filter_url)
}

func (s *FilterService) CheckURL(url string) (bool, error) {
	rules, err := s.repo.GetBlockUrls(context.Background())
	if err != nil {
		return false, err
	}

	for _, rule := range rules {
		if strings.Contains(url, rule.Url) {
			return false, nil
		}
	}

	return true, nil
}
