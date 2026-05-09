# 0063 - Half-Baked Feature Triage Decisions

## Status

Accepted

## Context

Phase 3 must finish, hide, or delete confusing partial features.

## Decision

| Feature | Decision | Rationale |
| --- | --- | --- |
| Safe export preview | Finish | Core publishing promise requires portable output. |
| Single-file upload | Finish | Real document dumps need batch handling. |
| Demo graph | Finish as honest empty state | Keep orientation, remove implication that sample data is loaded. |
| Backend endpoint setting | Finish | Persisted setting needs reset/start-fresh behavior. |
| README OCR/transcription/indexing claims | Edit | Claims must match processor-needed/current search behavior. |
| URL input/folder upload/share links | Keep hidden/out of scope | Not claimed and too risky for v0.3.0. |

## Consequences

Production UI will not expose placeholders or buttons that stop short of their label.

## Alternatives Considered

Deleting the demo map entirely was rejected because it helps orient first-time users before any backend case exists.
