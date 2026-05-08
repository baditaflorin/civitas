# 0005 - Client-Side Storage Strategy

## Status

Accepted

## Context

The frontend needs to remember endpoint settings, dismissed notices, last selected case, and lightweight UI state without storing evidence.

## Decision

Use `localStorage` for non-sensitive preferences and TanStack Query cache for API responses. Evidence files, extracted text, and private derived artifacts stay on the backend filesystem.

## Consequences

- No secrets or evidence are stored in the public frontend.
- The app can reconnect to a configured backend after refresh.
- Cross-device sync is not a v1 goal.

## Alternatives Considered

- IndexedDB or OPFS for evidence: rejected because Mode C keeps evidence server-side.
- Server-side sessions: rejected because v1 avoids auth and user accounts.
