# 0008 - Go Backend Project Layout

## Status

Accepted

## Context

The backend must be maintainable, testable, and close to common Go project layout conventions.

## Decision

Use a layout inspired by `golang-standards/project-layout`:

- `cmd/server`
- `internal/...`
- `pkg/...` only for reusable public packages, initially empty
- `api/openapi.yaml`
- `configs/`
- `scripts/`
- `test/integration/`

Use `chi`, `slog`, `prometheus/client_golang`, `validator`, and standard Go errors with wrapping. Include `internal/utils.HandleErrorOrLogWithMessages(err, errMsg, successMsg)`.

## Consequences

- Runtime concerns are separated from domain and pipeline code.
- Handlers can be unit tested without starting Docker.
- Processor adapters can evolve independently.

## Alternatives Considered

- Flat package layout: rejected because the app has multiple clear concerns.
- Framework-heavy Go backend: rejected because standard Go plus focused libraries is easier to audit.
