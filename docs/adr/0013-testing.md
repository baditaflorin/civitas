# 0013 - Testing Strategy

## Status

Accepted

## Context

The project needs fast local checks because it intentionally does not use GitHub Actions.

## Decision

Use:

- Go unit tests colocated with source.
- Frontend unit tests with Vitest and Testing Library.
- Playwright smoke/e2e tests against a locally served Pages build.
- `scripts/smoke.sh` for end-to-end sanity.
- `make test`, `make build`, and `make smoke` in pre-push.

Coverage target is at least 70 percent for important logic packages.

## Consequences

- Contributors get local feedback before pushing.
- Heavy processor integration tests can be build-tagged or skipped if native tools are absent.
- Smoke tests protect Pages base-path behavior.

## Alternatives Considered

- GitHub Actions: rejected by project constraint.
- Manual browser-only QA: rejected because Pages regressions are easy to automate.
