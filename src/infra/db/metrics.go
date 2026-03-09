package db

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DBMetrics struct {
	Latency     *prometheus.HistogramVec
	QueryCount  *prometheus.CounterVec
	QueryErrors *prometheus.CounterVec
	PoolActive  *prometheus.GaugeVec
	PoolIdle    *prometheus.GaugeVec
}

func NewDBMetrics(reg prometheus.Registerer) *DBMetrics {
	metrics := &DBMetrics{
		Latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Postgres query latency",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query"},
		),
		QueryCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_query_total",
				Help: "Total number of database queries",
			},
			[]string{"query"},
		),
		QueryErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_query_errors_total",
				Help: "Total number of failed database queries",
			},
			[]string{"query"},
		),
		PoolActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_pool_active_connections",
				Help: "Number of active connections in the pool",
			},
			[]string{},
		),
		PoolIdle: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_pool_idle_connections",
				Help: "Number of idle connections in the pool",
			},
			[]string{},
		),
	}
	reg.MustRegister(metrics.Latency)
	reg.MustRegister(metrics.QueryCount)
	reg.MustRegister(metrics.QueryErrors)
	reg.MustRegister(metrics.PoolActive)
	reg.MustRegister(metrics.PoolIdle)
	return metrics
}
