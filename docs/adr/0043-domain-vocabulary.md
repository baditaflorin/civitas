# 0043 - Domain Vocabulary and UI Language

## Status

Accepted

## Context

V1 exposed implementation words like binary, placeholder, and completed.

## Decision

Use civic investigation vocabulary:

- Evidence shape, source document, scan, transcript, archive, transfer failed.
- Processing state, not generic status.
- Next step language on every anomaly.
- Domain field names such as buyer, supplier, tender title, source URL, amount, date.

## Consequences

- API anomaly messages are user-facing.
- Export text uses the same language as the UI.

## Alternatives Considered

- Keep messages terse and technical: rejected because real users need recovery guidance.
