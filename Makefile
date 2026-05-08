SHELL := /bin/bash

APP := civitas
OWNER := baditaflorin
VERSION := $(shell node -p "require('./package.json').version" 2>/dev/null || echo 0.1.0)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
GO_PACKAGES := $(shell go list ./... 2>/dev/null | grep -v /node_modules/ || true)
GO_LINT_PACKAGES := ./cmd/server ./internal/...

.PHONY: help install-hooks dev build data test test-integration smoke lint fmt pages-preview docker-build docker-push release compose-up compose-down clean hooks-pre-commit hooks-commit-msg hooks-pre-push

help: ## list all targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-22s %s\n", $$1, $$2}'

install-hooks: ## wire .githooks
	git config core.hooksPath .githooks
	chmod +x .githooks/*

dev: ## run locally
	@mkdir -p tmp
	@echo "Starting backend on http://localhost:8080"
	@CIVITAS_ENV=development CIVITAS_ADDR=:8080 CIVITAS_STORAGE_DIR=./storage go run ./cmd/server & echo $$! > tmp/backend.pid
	@echo "Starting frontend on http://localhost:5173/civitas/"
	npm run dev

build: ## build frontend into docs and backend binary
	./scripts/build-pages.sh
	@mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)" -o bin/civitas ./cmd/server

data: ## regenerate static demo data
	./scripts/make-data.sh

test: ## unit tests
	@if [ -n "$(GO_PACKAGES)" ]; then CGO_ENABLED=0 go test $(GO_PACKAGES); fi
	npm run test

test-integration: ## integration tests
	@if find test/integration -name '*_test.go' -print -quit 2>/dev/null | grep -q .; then CGO_ENABLED=0 go test -tags=integration ./test/integration/...; else echo "No integration tests yet"; fi

smoke: ## smoke tests
	./scripts/smoke.sh

lint: ## all linters
	@if [ -n "$(GO_PACKAGES)" ]; then CGO_ENABLED=0 go vet $(GO_PACKAGES); fi
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run $(GO_LINT_PACKAGES); else echo "golangci-lint not installed; skipping"; fi
	@if command -v govulncheck >/dev/null 2>&1; then govulncheck $(GO_PACKAGES); else echo "govulncheck not installed; skipping"; fi
	npm run lint
	npm run typecheck
	npm audit --audit-level=high

fmt: ## autoformat
	gofmt -w cmd internal
	npx prettier --write "api/**/*.yaml" "src/**/*.{ts,tsx,css}" "*.{js,ts,json}" "public/*.{html,webmanifest}"

pages-preview: ## serve docs locally as Pages would
	npm run preview -- --host 127.0.0.1 --port 4174

docker-build: ## buildx amd64
	docker buildx build --platform linux/amd64 --load \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT) \
		--build-arg CREATED=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		-t ghcr.io/$(OWNER)/$(APP):$(COMMIT) .

docker-push: ## push to ghcr
	docker buildx build --platform linux/amd64 --push \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT) \
		--build-arg CREATED=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		-t ghcr.io/$(OWNER)/$(APP):latest \
		-t ghcr.io/$(OWNER)/$(APP):v$(VERSION) \
		-t ghcr.io/$(OWNER)/$(APP):$(COMMIT) .

release: ## tag and prepare release artifacts
	$(MAKE) test build smoke
	git tag v$(VERSION)
	@echo "Run: git push origin main v$(VERSION) && make docker-push"

compose-up: ## local stack
	docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.dev.yml up -d --build

compose-down: ## stop local stack
	docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.dev.yml down

clean: ## remove generated local outputs
	rm -rf bin coverage tmp dist-data

hooks-pre-commit:
	$(MAKE) fmt lint
	@if command -v gitleaks >/dev/null 2>&1; then gitleaks protect --staged --redact; else echo "gitleaks not installed; skipping"; fi

hooks-commit-msg:
	./.githooks/commit-msg .git/COMMIT_EDITMSG

hooks-pre-push:
	$(MAKE) test build
	@if ! git diff --quiet -- docs; then echo "Pages build changed docs/. Commit the build output before pushing."; git status --short docs; exit 1; fi
	$(MAKE) smoke
	@if ! git diff --quiet -- docs; then echo "Smoke test changed docs/. Commit the build output before pushing."; git status --short docs; exit 1; fi
