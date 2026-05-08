# 0004 - API Contract

## Status

Accepted

## Context

Mode C requires a runtime API for uploads, ingestion status, search, graph, timeline, redaction preview, export, health, readiness, and metrics.

## Decision

Publish an OpenAPI 3.1 spec at `api/openapi.yaml`. Generate the frontend client types from that spec with `openapi-typescript`; API calls use `openapi-fetch`.

The v1 API is JSON over REST:

- `GET /healthz`
- `GET /readyz`
- `GET /api/v1/version`
- `GET /api/v1/cases`
- `POST /api/v1/cases`
- `POST /api/v1/cases/{case_id}/documents`
- `GET /api/v1/cases/{case_id}/documents`
- `GET /api/v1/cases/{case_id}/search`
- `GET /api/v1/cases/{case_id}/graph`
- `GET /api/v1/cases/{case_id}/timeline`
- `POST /api/v1/cases/{case_id}/exports`
- `GET /api/v1/cases/{case_id}/exports/{export_id}`

## Consequences

- Frontend and backend evolve against a stable contract.
- The generated client avoids hand-written response shapes.
- Breaking API changes require an ADR and versioned path.

## Alternatives Considered

- GraphQL: rejected because v1 workflows are command/resource oriented and REST is simpler to operate.
- Hand-written client types: rejected to avoid drift.
