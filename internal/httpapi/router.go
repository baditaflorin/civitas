package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/baditaflorin/civitas/internal/config"
	"github.com/baditaflorin/civitas/internal/observability"
	"github.com/baditaflorin/civitas/internal/pipeline"
	"github.com/baditaflorin/civitas/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Dependencies struct {
	Config   config.Config
	Logger   *slog.Logger
	Store    *storage.Store
	Pipeline *pipeline.Pipeline
	Metrics  *observability.Metrics
}

func NewRouter(deps Dependencies) http.Handler {
	api := &API{deps: deps}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(requestLog(deps.Logger))
	router.Use(deps.Metrics.Middleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.Config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/healthz", api.healthz)
	router.Get("/readyz", api.readyz)
	router.Handle("/metrics", promhttp.HandlerFor(deps.Metrics.Registry(), promhttp.HandlerOpts{}))

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/version", api.version)
		r.Get("/processors", api.processors)
		r.Get("/cases", api.listCases)
		r.Post("/cases", api.createCase)
		r.Post("/case-states/import", api.importCaseState)
		r.Route("/cases/{case_id}", func(r chi.Router) {
			r.Get("/documents", api.listDocuments)
			r.Post("/documents", api.uploadDocument)
			r.Get("/search", api.search)
			r.Get("/graph", api.graph)
			r.Get("/timeline", api.timeline)
			r.Get("/debug", api.debugCase)
			r.Get("/state", api.getCaseState)
			r.Post("/exports", api.createExport)
			r.Get("/exports/{export_id}", api.getExport)
		})
	})

	return router
}
