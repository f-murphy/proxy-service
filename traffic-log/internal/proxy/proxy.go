package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"traffic-log/internal/metrics"
	"traffic-log/internal/service"
	"traffic-log/models"
	"traffic-log/utils"
)

type ReverseProxy struct {
	targetURL      *url.URL
	proxy          *httputil.ReverseProxy
	trafficService service.TrafficService
}

func NewReverseProxy(target string, trafficService service.TrafficService) (*ReverseProxy, error) {
	urlParsed, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}
	return &ReverseProxy{
		targetURL:      urlParsed,
		proxy:          httputil.NewSingleHostReverseProxy(urlParsed),
		trafficService: trafficService,
	}, nil
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	r.Host = p.targetURL.Host

	p.proxy.ServeHTTP(rec, r)

	bytesSent := int64(rec.bytesWritten)
	bytesReceived := r.ContentLength
	clientIP := utils.GetClientIP(r)

	metrics.TotalRequests.WithLabelValues(r.Method, fmt.Sprintf("%d", rec.statusCode)).Inc()
	metrics.TrafficBytes.WithLabelValues("sent").Add(float64(bytesSent))
	metrics.TrafficBytes.WithLabelValues("received").Add(float64(bytesReceived))

	go func() {
		log := &models.TrafficLog{
			Timestamp:     start,
			ClientIP:      clientIP,
			Method:        r.Method,
			URL:           r.URL.String(),
			ResponseCode:  rec.statusCode,
			BytesSent:     bytesSent,
			BytesReceived: bytesReceived,
		}
		_ = p.trafficService.LogRequest(context.Background(), log)
	}()
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}
