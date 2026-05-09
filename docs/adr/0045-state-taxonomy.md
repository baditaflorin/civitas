# 0045 - State Taxonomy and State Machine

## Status

Accepted

## Context

V1 used `completed` for too many incompatible outcomes.

## Decision

Document states:

- `ready`: useful first guess exists.
- `needs_processor`: file is valid but requires an unavailable native processor.
- `recoverable_error`: input appears partial/corrupt or needs user action.
- `failed`: input cannot be processed as provided.
- `unsupported`: recognized but not currently processable.

Every state has at least one next step through anomalies.

## Consequences

- UI and exports can distinguish "understood" from "stored."
- Empty and corrupt inputs cannot masquerade as complete.

## Alternatives Considered

- Keep `completed` plus messages: rejected because users read status first.
