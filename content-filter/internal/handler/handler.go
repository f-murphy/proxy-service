package handler

import (
	"content-filter/internal/service"
	"content-filter/models"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	service *service.FilterService
	proxy   *httputil.ReverseProxy
}

func NewProxyHandler(service *service.FilterService) *ProxyHandler {
	return &ProxyHandler{
		service: service,
		proxy:   &httputil.ReverseProxy{Director: func(req *http.Request) {}},
	}
}

func (h *ProxyHandler) CreateBlockUrl(c *gin.Context) {
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

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.String()

	allowed, err := h.service.CheckURL(targetURL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "Access to this URL is blocked", http.StatusForbidden)
		return
	}

	h.proxy.ServeHTTP(w, r)
}
