package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	HTTPRequests           *prometheus.CounterVec
	HTTPRequestDuration    *prometheus.HistogramVec
	IngestionJobsStarted   prometheus.Counter
	IngestionJobsCompleted prometheus.Counter
	DocumentsProcessed     prometheus.Counter
	ExportsGenerated       prometheus.Counter
	registry               *prometheus.Registry
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	m := &Metrics{
		HTTPRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "civitas_http_requests_total",
			Help: "HTTP requests by method, path, and status.",
		}, []string{"method", "path", "status"}),
		HTTPRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "civitas_http_request_duration_seconds",
			Help:    "HTTP request duration by method and path.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"}),
		IngestionJobsStarted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "civitas_ingestion_jobs_started_total",
			Help: "Ingestion jobs started.",
		}),
		IngestionJobsCompleted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "civitas_ingestion_jobs_completed_total",
			Help: "Ingestion jobs completed.",
		}),
		DocumentsProcessed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "civitas_documents_processed_total",
			Help: "Documents processed.",
		}),
		ExportsGenerated: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "civitas_exports_generated_total",
			Help: "Publishing exports generated.",
		}),
		registry: registry,
	}
	registry.MustRegister(
		m.HTTPRequests,
		m.HTTPRequestDuration,
		m.IngestionJobsStarted,
		m.IngestionJobsCompleted,
		m.DocumentsProcessed,
		m.ExportsGenerated,
	)
	registry.MustRegister(collectors.NewGoCollector())
	return m
}

func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		status := strconv.Itoa(recorder.status)
		route := r.URL.Path
		m.HTTPRequests.WithLabelValues(r.Method, route, status).Inc()
		m.HTTPRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
