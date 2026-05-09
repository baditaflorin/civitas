# 0049 - Inspectability and Debug Surface

## Status

Accepted

## Context

Power users need to know why Civitas inferred something.

## Decision

Add a debug API endpoint that returns case-level document shapes, states, confidence, field inferences, anomalies, parse durations, and processor decisions. The endpoint uses existing data and does not add a new workflow.

## Consequences

- Support/debugging is possible without reading storage JSON by hand.
- Debug output is excluded from public Pages unless explicitly requested by a user or developer.

## Alternatives Considered

- Frontend-only debug overlay first: deferred to keep Phase 2 focused on engine substance.
