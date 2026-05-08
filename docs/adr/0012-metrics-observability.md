# 0012 - Metrics and Observability

## Status

Accepted

## Context

The backend must expose scrape-ready metrics while the static frontend should avoid tracking users by default.

## Decision

Expose `/metrics` with Prometheus metrics:

- HTTP request count and duration.
- Go runtime metrics.
- Ingestion jobs started and completed.
- Documents processed.
- Exports generated.

Block public `/metrics` access at nginx. The frontend ships with no analytics.

## Consequences

- Operators can monitor ingestion and API health.
- Public visitors are not tracked by default.
- Prometheus is optional and profile-gated in Docker Compose.

## Alternatives Considered

- Client analytics: rejected for v1 privacy posture.
- No metrics: rejected because Mode C requires operational visibility.
