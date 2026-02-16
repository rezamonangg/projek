package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "projek_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "projek_http_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	DBQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "projek_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"query"},
	)

	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "projek_active_sessions",
			Help: "Number of active sessions",
		},
	)

	ProjectsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "projek_projects_total",
			Help: "Total number of projects",
		},
	)

	TasksTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "projek_tasks_total",
			Help: "Total number of tasks",
		},
	)
)

func Init() {}
