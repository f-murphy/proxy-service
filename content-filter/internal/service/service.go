package service

import (
	"content-filter/internal/repository"
	"content-filter/models"
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type FilterService interface {
	CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error)
	GetBlockUrls(ctx context.Context) ([]string, error)
	CheckRequest(r *http.Request) (bool, string)
}

type filterService struct {
	repo repository.FilterRepository
}

func NewFilterService(repo repository.FilterRepository) FilterService {
	return &filterService{repo: repo}
}

func (s *filterService) CheckRequest(r *http.Request) (bool, string) {
	ctx := context.Background()

	// urls, err := s.repo.GetBlockUrls(ctx)
	// if err != nil {
	// 	logrus.Errorf("Error getting blocked URLs: %v", err)
	// }
	// fmt.Println("urls - ", urls)
	// logrus.Infof("Checking URL: %s", r.URL.String())
	// for _, url := range urls {
		
	// 	if strings.Contains(r.URL.Path, url) {
	// 		logrus.Info("Url blocked")
	// 		return false, "URL blocked: " + url
	// 	}
	// }

	if r.Method == http.MethodPost {
		keywords, err := s.repo.GetBlacklistKeywords(ctx)
		if err != nil {
			logrus.Errorf("Error getting keywords: %v", err)
		}

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		defer r.Body.Close()

		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(string(body)), strings.ToLower(keyword)) {
				return false, "Keyword blocked: " + keyword
			}
		}
	}

	return true, ""
}

func (s *filterService) CreateBlockUrl(ctx context.Context, filter_url *models.Filter_urls) (string, error) {
	return s.repo.CreateBlockUrl(ctx, filter_url)
}

func (s *filterService) GetBlockUrls(ctx context.Context) ([]string, error) {
	return s.repo.GetBlockUrls(ctx)
}
