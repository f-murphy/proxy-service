package service

import (
	"context"
	"traffic-log/internal/repository"
	"traffic-log/models"
)

type TrafficService interface {
	LogRequest(ctx context.Context, log *models.TrafficLog) error
}

type trafficService struct {
	repo repository.TrafficRepository
}

func NewTrafficService(repo repository.TrafficRepository) TrafficService {
	return &trafficService{repo: repo}
}

func (s *trafficService) LogRequest(ctx context.Context, log *models.TrafficLog) error {
	return s.repo.LogTraffic(ctx, log)
}
