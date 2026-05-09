# 0066 - Error-Handling Convention

## Status

Accepted

## Context

Backend handlers use `writeError`; frontend uses `apiError` and some query errors vanish from the UI.

## Decision

- Backend HTTP handlers keep `writeError` for unexpected errors and explicit `errorResponse` for validation errors.
- Frontend API errors are converted through one `apiError` helper.
- User-facing mutations show an actionable inline status or result message.
- Go CLI/server startup keeps explicit logged errors instead of panics.

## Consequences

No new error framework is needed. Completeness work focuses on making existing errors visible where users act.

## Alternatives Considered

Replacing all Go handler errors with `internal/utils.HandleErrorOrLogWithMessages` was rejected because HTTP response mapping still needs status-specific handling.
