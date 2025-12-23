package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_histogram",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	TasksCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tasks_count",
			Help: "Current number of todo tasks",
		},
	)
)

func MustRegisterMetrics() {
	prometheus.MustRegister(RequestsTotal, RequestLatencyHistogram, TasksCount)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		RequestsTotal.WithLabelValues(method, path, status).Inc()
		RequestLatencyHistogram.
			WithLabelValues(method, path).
			Observe(time.Since(start).Seconds())
	}
}
