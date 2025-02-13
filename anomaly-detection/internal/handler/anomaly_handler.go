package handler

import (
	"anomaly-detection/internal/service"
)

type AnomalyHandler struct {
	service service.AnomalyService
}

func NewAnomalyHandler(service service.AnomalyService) *AnomalyHandler {
	return &AnomalyHandler{service: service}
}
