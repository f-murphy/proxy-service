package handler

import (
	"anomaly-detection/internal/proxy"
	"anomaly-detection/internal/service"
	"net/http"

	"github.com/sirupsen/logrus"
)

type AnomalyHandler struct {
	service *service.AnomalyService
	proxy   *proxy.ReverseProxy
}

func NewAnomalyHandler(service *service.AnomalyService, target string) *AnomalyHandler {
	return &AnomalyHandler{
		service: service,
		proxy:   proxy.NewReverseProxy(target),
	}
}

func (h *AnomalyHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, reason := h.service.CheckRequest(r)
		if !allowed {
			logrus.Warnf("Blocked request: %s", reason)
			http.Error(w, reason, http.StatusForbidden)
			return
		}

		h.proxy.ServeHTTP(w, r)
	})
}
