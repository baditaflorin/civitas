package httpapi

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baditaflorin/civitas/internal/config"
	"github.com/baditaflorin/civitas/internal/observability"
	"github.com/baditaflorin/civitas/internal/pipeline"
	"github.com/baditaflorin/civitas/internal/storage"
)

func TestRouterVersionEndpoint(t *testing.T) {
	store, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	metrics := observability.NewMetrics()
	router := NewRouter(Dependencies{
		Config: config.Config{
			Env:            "test",
			Addr:           ":0",
			StorageDir:     t.TempDir(),
			AllowedOrigins: []string{"http://example.test"},
			Version:        "0.1.0",
			CommitSHA:      "abc123",
		},
		Logger:   slog.Default(),
		Store:    store,
		Pipeline: pipeline.New(pipeline.DefaultRegistry(), metrics),
		Metrics:  metrics,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/version", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); body == "" || !containsAll(body, "0.1.0", "abc123") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func containsAll(body string, values ...string) bool {
	for _, value := range values {
		if !strings.Contains(body, value) {
			return false
		}
	}
	return true
}
