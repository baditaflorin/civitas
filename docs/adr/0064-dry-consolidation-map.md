# 0064 - DRY Consolidation Map

## Status

Accepted

## Context

The frontend repeats API result handling and cache invalidation. The backend has export formatting inside HTTP handlers.

## Decision

Consolidate:

- Frontend session storage in `src/lib/session.ts`.
- Frontend API result narrowing in `src/lib/api/client.ts`.
- Document-related query invalidation in one helper inside the workspace feature.
- Backend export/state formatting in an exporter package.

Accept for now:

- Small `map[string]any` HTTP envelopes at API boundaries.
- Generated OpenAPI `unknown` index signatures.

## Consequences

The main workspace remains a large feature component, but the most duplicated logic moves to single-purpose helpers.

## Alternatives Considered

A full frontend module split was rejected for Phase 3 because it risks churn without directly improving user completion.
