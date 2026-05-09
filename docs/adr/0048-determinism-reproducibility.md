# 0048 - Determinism and Reproducibility

## Status

Accepted

## Context

Investigative output must be reproducible.

## Decision

Use stable IDs derived from content hashes and normalized values. Sort entities, fields, anomalies, timeline events, graph nodes, and export sections. Export provenance includes schema version, app version, source SHA-256, processor decisions, and parameters.

Timestamps may exist as metadata but are excluded from deterministic normalized-output tests.

## Consequences

- Re-running fixture analysis produces stable IDs and ordering.
- Graph output stops changing due to map iteration.

## Alternatives Considered

- Random IDs for documents/entities: rejected because exports and tests need stable references.
