package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics struct {
	TotalRequests    *prometheus.CounterVec
	Latency          *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
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
		Latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		),
		RequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of in-flight HTTP requests",
			},
		),
	}

	reg.MustRegister(metrics.TotalRequests)
	reg.MustRegister(metrics.Latency)
	reg.MustRegister(metrics.RequestsInFlight)
	return metrics
}
