package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics struct {
	TotalRequests *prometheus.CounterVec
	Latency       *prometheus.Histogram
}

func NewHTTPMetrics(reg prometheus.Registerer) *HTTPMetrics {
	metrics := &HTTPMetrics{
		TotalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests by code and method",
			},
			[]string{"code", "method"},
		),
	}
	reg.MustRegister(metrics.TotalRequests)
	return metrics
}
