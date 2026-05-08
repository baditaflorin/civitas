# 0002 - Architecture Overview and Module Boundaries

## Status

Accepted

## Context

Civitas needs a safe public UI, a private evidence-processing runtime, and clean boundaries between ingestion, extraction, indexing, graphing, export, and operations.

## Decision

Use these boundaries:

- `frontend/`: Vite React TypeScript application, built into `docs/` for GitHub Pages.
- `cmd/server/`: runtime API entrypoint.
- `internal/config`: environment parsing and validation.
- `internal/httpapi`: handlers, middleware, routing, OpenAPI-aligned DTOs.
- `internal/evidence`: case and document domain model.
- `internal/pipeline`: ingestion orchestration and processor registry.
- `internal/processors`: adapters for native tools and local fallback processors.
- `internal/search`: in-process search index abstraction, with Tantivy/DuckDB-ready boundaries.
- `internal/graph`: relationship graph model and export helpers.
- `internal/storage`: filesystem-backed storage for evidence and derived artifacts.
- `internal/observability`: metrics and logging helpers.
- `api/`: OpenAPI spec.
- `deploy/`: Docker Compose, nginx, Prometheus, and server documentation.

## Consequences

- The frontend remains static and does not serve secrets.
- The backend remains API-only and does not serve the frontend.
- Native processors are isolated behind interfaces.
- V1 can ship useful local fallbacks while documenting production adapters.

## Alternatives Considered

- Monolithic script pipeline: rejected because API, metrics, and future scheduling would become brittle.
- Backend serving the frontend: rejected because GitHub Pages is a first-class deliverable.
