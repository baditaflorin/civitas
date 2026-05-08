# 0009 - Configuration and Secrets Management

## Status

Accepted

## Context

The frontend must never hold secrets. The backend needs environment-driven configuration for address, storage, allowed origins, version, and commit.

## Decision

Use environment variables documented in `.env.example`. The Go backend reads config through `internal/config`; local compose uses a gitignored `.env`.

No secrets are required for v1. If optional third-party geocoding or identity providers are added later, their credentials stay backend-only.

## Consequences

- Static frontend can be inspected safely.
- Deployment can use Docker Compose `.env` files without committing secrets.
- Gitleaks is enforced by local hooks.

## Alternatives Considered

- Hardcoded production config: rejected.
- Encrypted frontend secrets: rejected because frontend secrets are not secrets.
