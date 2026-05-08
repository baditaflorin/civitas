package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baditaflorin/civitas/internal/config"
	"github.com/baditaflorin/civitas/internal/httpapi"
	"github.com/baditaflorin/civitas/internal/observability"
	"github.com/baditaflorin/civitas/internal/pipeline"
	"github.com/baditaflorin/civitas/internal/storage"
)

var (
	version = "0.1.0"
	commit  = "dev"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := healthcheck(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load(version, commit)
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	store, err := storage.New(cfg.StorageDir)
	if err != nil {
		logger.Error("open storage", "error", err)
		os.Exit(1)
	}

	metrics := observability.NewMetrics()
	pipe := pipeline.New(pipeline.DefaultRegistry(), metrics)
	router := httpapi.NewRouter(httpapi.Dependencies{
		Config:   cfg,
		Logger:   logger,
		Store:    store,
		Pipeline: pipe,
		Metrics:  metrics,
	})

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("server starting", "addr", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	logger.Info("server shutting down")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
}

func healthcheck() error {
	addr := os.Getenv("CIVITAS_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("parse CIVITAS_ADDR: %w", err)
	}
	if host == "" || host == "::" {
		host = "127.0.0.1"
	}
	client := &http.Client{Timeout: 5 * time.Second}
	// #nosec G107,G704 -- healthcheck dials the configured local service address only.
	resp, err := client.Get("http://" + net.JoinHostPort(host, port) + "/healthz")
	if err != nil {
		return fmt.Errorf("request healthz: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthz returned %s", resp.Status)
	}
	return nil
}
