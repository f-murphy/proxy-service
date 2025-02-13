package handler

import (
	"content-filter/internal/Proxy"
	"content-filter/internal/service"
	"content-filter/models"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *FilterHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, reason := h.service.CheckRequest(r)
		if !allowed {
			logrus.Warnf("Blocked request to %s: %s", r.URL.String(), reason)
			http.Error(w, reason, http.StatusForbidden)
			return
		}
		
		h.proxy.ServeHTTP(w, r)
	})
}

func (h *FilterHandler) CreateBlockUrl(c *gin.Context) {
	var request models.Filter_urls

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := h.service.CreateBlockUrl(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

