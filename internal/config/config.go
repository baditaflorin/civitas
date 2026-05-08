package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env            string   `json:"env" validate:"required"`
	Addr           string   `json:"addr" validate:"required"`
	StorageDir     string   `json:"storage_dir" validate:"required"`
	AllowedOrigins []string `json:"allowed_origins" validate:"min=1"`
	Version        string   `json:"version" validate:"required"`
	CommitSHA      string   `json:"commit_sha" validate:"required"`
}

type rawConfig struct {
	Env            string `envconfig:"ENV" default:"development"`
	Addr           string `envconfig:"ADDR" default:":8080"`
	StorageDir     string `envconfig:"STORAGE_DIR" default:"./storage"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:5173,https://baditaflorin.github.io"`
	Version        string `envconfig:"VERSION" default:"0.1.0"`
	CommitSHA      string `envconfig:"COMMIT_SHA" default:"dev"`
}

func Load(buildVersion, buildCommit string) (Config, error) {
	var raw rawConfig
	if err := envconfig.Process("civitas", &raw); err != nil {
		return Config{}, fmt.Errorf("process env: %w", err)
	}
	if raw.Version == "" || raw.Version == "0.1.0" {
		raw.Version = buildVersion
	}
	if raw.CommitSHA == "" || raw.CommitSHA == "dev" {
		raw.CommitSHA = buildCommit
	}

	cfg := Config{
		Env:            raw.Env,
		Addr:           raw.Addr,
		StorageDir:     raw.StorageDir,
		AllowedOrigins: splitCSV(raw.AllowedOrigins),
		Version:        raw.Version,
		CommitSHA:      raw.CommitSHA,
	}
	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}
	return cfg, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
