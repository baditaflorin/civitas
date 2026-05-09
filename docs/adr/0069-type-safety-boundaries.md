# 0069 - Type-Safety Policy at Boundaries

## Status

Accepted

## Context

Generated OpenAPI types and JSON parsing require `unknown`/`any`, but UI feature code should not scatter casts.

## Decision

- Generated API types may contain `unknown`.
- Go JSON helpers may use `any` at serialization boundaries.
- Pipeline JSON parsing may use `any` with explicit narrowing.
- Frontend feature code should use typed API helpers instead of repeated casts.
- Multipart `FormData` remains the one explicitly documented openapi-fetch boundary cast if required.

## Consequences

The audit can distinguish safe boundary code from unsafe feature code.

## Alternatives Considered

Removing all `any` in Go JSON code was rejected because it would make generic JSON helpers worse.
