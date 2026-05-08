# 0014 - Error Handling Conventions

## Status

Accepted

## Context

Evidence processing has many optional tools and failure modes. Errors must be inspectable without panics.

## Decision

Use standard Go `errors` and `fmt.Errorf("%w")`. Handlers return stable JSON errors with `code`, `message`, and optional `details`. Never panic for expected failures. Include `internal/utils.HandleErrorOrLogWithMessages(err, errMsg, successMsg)` for CLI and operational scripts.

Frontend errors are surfaced through an error boundary and toast region.

## Consequences

- Pipeline failures can be reported per document.
- API clients receive predictable error shapes.
- Logs remain useful without exposing private content.

## Alternatives Considered

- Global panic recovery as normal control flow: rejected.
- Raw error strings in handlers: rejected because they can leak internals.
