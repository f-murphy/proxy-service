package service

import (
	"anomaly-detection/internal/repository"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type AnomalyService interface {
	CheckRequest(r *http.Request) (bool, string)
}

type anomalyService struct {
	repo repository.AnomalyRepository
}

func NewAnomalyService(repo repository.AnomalyRepository) AnomalyService {
	return &anomalyService{repo: repo}
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	ip = strings.Trim(ip, "[]")
	return ip
}

func (s *anomalyService) CheckForAnomalies(ctx context.Context, ip string) (bool, error) {
	threshold := 100     // Пороговое значение (например, 100 запросов в минуту)
	interval := "minute"

	count, err := s.repo.GetRequestMetrics(ctx, ip, interval)
	if err != nil {
		return false, fmt.Errorf("unable to check request metrics: %w", err)
	}

	if count > threshold {
		return true, nil
	}

	return false, nil
}

// CheckRequest проверяет запрос на наличие аномалий и блокирует IP при необходимости
func (s *anomalyService) CheckRequest(r *http.Request) (bool, string) {
	ctx := context.Background()

	clientIP := getClientIP(r)
	logrus.Infof("Checking request from IP: %s", clientIP)

	// Обновляем метрики запросов
	if err := s.repo.UpdateRequestMetrics(ctx, clientIP); err != nil {
		logrus.Errorf("Error updating request metrics: %v", err)
	}

	// Проверка на аномалии (DDoS)
	isAnomaly, err := s.CheckForAnomalies(ctx, clientIP)
	if err != nil {
		logrus.Errorf("Error checking for anomalies: %v", err)
	}
	if isAnomaly {
		// Блокируем IP-адрес
		if err := s.repo.BlockIP(ctx, clientIP); err != nil {
			logrus.Errorf("Error blocking IP: %v", err)
		} else {
			logrus.Warnf("IP %s blocked for DDoS attack", clientIP)
		}
		return false, "Your IP is blocked due to suspicious activity."
	}

	// Если всё в порядке, разрешаем запрос
	return true, ""
}
