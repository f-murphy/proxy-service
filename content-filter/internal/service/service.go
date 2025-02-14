package service

import (
	"content-filter/internal/repository"
	"content-filter/models"
	"content-filter/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// Интерфейс сервиса
type FilterService interface {
	CheckRequest(r *http.Request) (bool, string)
	CreateBlockURL(ctx context.Context, filterURL *models.Filter_urls) (string, error)
	GetBlockURLs(ctx context.Context) ([]string, error)
	GetBlacklistKeywords(ctx context.Context) ([]string, error)
	BlockIP(ctx context.Context, ip string) error
	IsIPBlocked(ctx context.Context, ip string) (bool, error)
}

// Реализация сервиса
type filterService struct {
	repo   repository.FilterRepository
}

func NewFilterService(repo repository.FilterRepository) FilterService {
	return &filterService{repo: repo}
}

// CheckRequest проверяет запрос на наличие запрещённых URL и ключевых слов
func (s *filterService) CheckRequest(r *http.Request) (bool, string) {
	ctx := context.Background()
	
	// Получаем IP-адрес клиента
	clientIP := utils.GetClientIP(r)
	logrus.Infof("Checking request from IP: %s", clientIP)
	
	// Проверка IP-адреса на блокировку
	isBlocked, err := s.repo.IsIPBlocked(ctx, clientIP)
	if err != nil {
		logrus.Errorf("Error checking IP block status: %v", err)
	}
	if isBlocked {
		logrus.Warnf("Request from blocked IP: %s", clientIP)
		return false, "Your IP is blocked."
	}
	
	// Проверка URL
	urls, err := s.repo.GetBlockURLs(ctx)
	if err != nil {
		logrus.Errorf("Error getting blocked URLs: %v", err)
	}
	
	for _, url := range urls {
		if strings.Contains(r.URL.String(), url) {
			logrus.Warnf("URL blocked: %s", url)
			return false, "URL blocked: " + url
		}
	}
	
	// Проверка ключевых слов в POST-запросах
	if r.Method == http.MethodPost {
		keywords, err := s.repo.GetBlacklistKeywords(ctx)
		if err != nil {
			logrus.Errorf("Error getting keywords: %v", err)
		}
		
		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Error reading request body: %v", err)
			return true, ""
		}
		defer r.Body.Close()
		
		// Проверка на наличие запрещённых ключевых слов
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(string(body)), strings.ToLower(keyword)) {
				// Блокируем IP-адрес
				if err := s.repo.BlockIP(ctx, clientIP); err != nil {
					logrus.Errorf("Error blocking IP: %v", err)
				} else {
					logrus.Warnf("IP %s blocked for keyword: %s", clientIP, keyword)
				}
				return false, fmt.Sprintf("Keyword blocked: %s (IP: %s)", keyword, clientIP)
			}
		}
	}
	
	// Если всё в порядке, разрешаем запрос
	return true, ""
}

// CreateBlockURL добавляет URL в чёрный список
func (s *filterService) CreateBlockURL(ctx context.Context, filterURL *models.Filter_urls) (string, error) {
	return s.repo.CreateBlockURL(ctx, filterURL)
}

// GetBlockURLs возвращает список заблокированных URL
func (s *filterService) GetBlockURLs(ctx context.Context) ([]string, error) {
	return s.repo.GetBlockURLs(ctx)
}

// GetBlacklistKeywords возвращает список запрещённых ключевых слов
func (s *filterService) GetBlacklistKeywords(ctx context.Context) ([]string, error) {
	return s.repo.GetBlacklistKeywords(ctx)
}

// BlockIP добавляет IP-адрес в список заблокированных
func (s *filterService) BlockIP(ctx context.Context, ip string) error {
	return s.repo.BlockIP(ctx, ip)
}

// IsIPBlocked проверяет, заблокирован ли IP-адрес
func (s *filterService) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	return s.repo.IsIPBlocked(ctx, ip)
}