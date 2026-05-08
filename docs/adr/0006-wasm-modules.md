# 0006 - WASM Modules

## Status

Accepted

## Context

The original architecture bias prefers browser-side WASM when feasible. Civitas, however, needs many native tools that are large, long-running, and privacy-sensitive.

## Decision

Use no required frontend WASM modules in v1. Heavy compute runs in the Docker backend. Future optional browser modules may be lazy-loaded behind explicit user actions and documented with COOP/COEP constraints.

## Consequences

- GitHub Pages remains simple and compatible.
- Initial payload stays within budget.
- Offline browser-only ingestion is not a v1 feature.

## Alternatives Considered

- DuckDB-WASM in v1: deferred because current v1 search and graph data come from the backend API.
- OCR or transcription WASM in v1: rejected because model and binary sizes are too large for the first-load budget.
