package service

import (
	"anomaly-detection/internal/repository"
)

type IAnomalyService interface {
}

type AnomalyService struct {
	repo repository.IAnomalyRepository
}

func NewAnomalyService(repo repository.IAnomalyRepository) IAnomalyService {
	return &AnomalyService{repo: repo}
}
