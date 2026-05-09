# 0047 - Error Taxonomy and Messaging

## Status

Accepted

## Context

Errors must say what failed, why, and what to do next.

## Decision

Use anomaly objects for recoverable domain issues:

- `code`
- `severity`
- `message`
- `why`
- `next_step`
- `confidence`

Fatal API errors retain JSON error responses, but document-level problems should stay attached to the document whenever possible.

## Consequences

- Partial work remains visible.
- Exports can carry warnings.

## Alternatives Considered

- Fail uploads on every processor miss: rejected because storing evidence with an honest state is better.
