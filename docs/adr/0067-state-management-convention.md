# 0067 - State-Management Convention

## Status

Accepted

## Context

Endpoint persists, but selected case and search term do not.

## Decision

- Server state stays in TanStack Query.
- UI preferences/session state use typed `localStorage` helpers.
- Persist endpoint, selected case ID, and search term.
- Reset/start-fresh clears local UI state only; backend evidence remains intact.

## Consequences

Reloads become coherent without adding auth or cross-device sync.

## Alternatives Considered

Persisting evidence in browser storage was rejected because evidence belongs on the private backend.
