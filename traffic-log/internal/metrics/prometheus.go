package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_total_requests",
			Help: "Total number of HTTP requests through proxy",
		},
		[]string{"method", "status"},
	)

	TrafficBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_traffic_bytes",
			Help: "Bytes transferred through proxy",
		},
		[]string{"direction"},
	)
)

func Init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(TrafficBytes)
}
