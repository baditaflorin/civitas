# 0015 - Deployment Topology

## Status

Accepted

## Context

Mode C requires GitHub Pages frontend plus Docker backend. The backend must be pullable from GHCR and deployable behind nginx.

## Decision

Use:

- GitHub Pages for frontend.
- GHCR image `ghcr.io/baditaflorin/civitas:<tag>` for backend.
- `deploy/docker-compose.yml` with `app`, optional `prometheus`, and `nginx`.
- nginx listens on 80 and 443, proxies to internal `app:8080`, and publishes host port `25342`.
- `/metrics` is blocked publicly.

## Consequences

- The backend is independently deployable from the frontend.
- Server operators can update with `docker compose pull && docker compose up -d`.
- TLS and CORS live at nginx.

## Alternatives Considered

- Backend on Pages: impossible.
- App serving frontend: rejected because Pages remains the public frontend.
