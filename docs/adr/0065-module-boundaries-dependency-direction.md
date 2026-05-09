# 0065 - Module Boundaries and Dependency Direction

## Status

Accepted

## Context

Completeness adds state export/import and more UI controls. Without boundaries, handlers and components will grow harder to change.

## Decision

Use this dependency direction:

- Frontend UI -> API client/session helpers -> generated API types.
- HTTP handlers -> storage/pipeline/exporter -> evidence domain.
- Storage owns filesystem details.
- Exporter owns markdown and state artifact construction.

## Consequences

New state/export behavior avoids living directly in HTTP handlers.

## Alternatives Considered

Introducing broad interfaces everywhere was rejected; only storage/exporter seams need clear boundaries now.
