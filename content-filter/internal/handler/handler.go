package handler

import (
	"content-filter/internal/proxy"
	"content-filter/internal/service"
	"content-filter/models"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type FilterHandler struct {
	service service.FilterService
	proxy   *proxy.ReverseProxy
}

func NewFilterHandler(service service.FilterService, target string) *FilterHandler {
	return &FilterHandler{
		service: service,
		proxy:   proxy.NewReverseProxy(target),
	}
}

// Middleware обрабатывает запросы и блокирует их при необходимости
func (h *FilterHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка запроса
		allowed, reason := h.service.CheckRequest(r)
		if !allowed {
			logrus.Warnf("Blocked request: %s", reason)
			http.Error(w, reason, http.StatusForbidden)
			return
		}
		
		// Проксируем запрос, если он разрешён
		h.proxy.ServeHTTP(w, r)
	})
}

// CreateBlockURLHandler добавляет URL в чёрный список
func (h *FilterHandler) CreateBlockURLHandler(w http.ResponseWriter, r *http.Request) {
	var filterURL models.Filter_urls
	if err := json.NewDecoder(r.Body).Decode(&filterURL); err != nil {
		logrus.Errorf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateBlockURL(r.Context(), &filterURL)
	if err != nil {
		logrus.Errorf("Error creating block URL: %v", err)
		http.Error(w, "Failed to create block URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetBlockURLsHandler возвращает список заблокированных URL
func (h *FilterHandler) GetBlockURLsHandler(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetBlockURLs(r.Context())
	if err != nil {
		logrus.Errorf("Error getting blocked URLs: %v", err)
		http.Error(w, "Failed to get blocked URLs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

// GetBlacklistKeywordsHandler возвращает список запрещённых ключевых слов
func (h *FilterHandler) GetBlacklistKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	keywords, err := h.service.GetBlacklistKeywords(r.Context())
	if err != nil {
		logrus.Errorf("Error getting blacklist keywords: %v", err)
		http.Error(w, "Failed to get blacklist keywords", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keywords)
}