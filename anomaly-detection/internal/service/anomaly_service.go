package service

import (
	"anomaly-detection/internal/repository"
	"fmt"
	"net/http"
)

type AnomalyService struct {
	repo *repository.AnomalyRepository
}

func NewAnomalyService(repo *repository.AnomalyRepository) *AnomalyService {
	return &AnomalyService{repo: repo}
}

func (s *AnomalyService) CheckRequest(r *http.Request) (bool, string) {
	ip := r.RemoteAddr
	allowed, err := s.repo.IsAllowed(ip)
	if err != nil {
		return false, fmt.Sprintf("Internal server error: %v", err)
	}
	if !allowed {
		return false, "Too many requests. Your IP has been blocked."
	}
	return true, ""
}
