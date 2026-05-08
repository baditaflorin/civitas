# 0007 - Data Generation Pipeline

## Status

Accepted

## Context

Mode B would require a static data generation pipeline. Civitas is Mode C, but it still benefits from deterministic sample fixtures and optional export generation.

## Decision

Do not build a Mode B artifact pipeline in v1. Provide `make data` as a lightweight deterministic sample fixture target that writes small public demo data only. Private evidence generation runs through the backend ingestion pipeline and writes to local backend storage.

## Consequences

- The public Pages app can show a useful demo without private data.
- Release-hosted data artifacts are out of scope for v1.
- Backend export artifacts are not committed.

## Alternatives Considered

- Full Mode B static corpus generation: rejected because private leaked dumps cannot be safely committed or released.
