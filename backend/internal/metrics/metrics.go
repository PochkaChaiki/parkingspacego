package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// === HTTP метрики (по методологии RED) ===
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "parking",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "handler", "status"}, // handler — имя обработчика, не полный путь!
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "parking",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "handler", "status"},
	)
)

// RegisterAll регистрирует все метрики в дефолтном регистре Prometheus
func RegisterAll() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
	)
}
