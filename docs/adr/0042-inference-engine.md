# 0042 - Inference Engine

## Status

Accepted

## Context

Users should correct Civitas, not configure it from zero.

## Decision

The inference engine produces:

- Evidence shape.
- Processing state.
- Preview text.
- Field inferences with type, value, confidence, and reason.
- Anomalies with severity, what, why, and next step.
- Timeline events from ISO and common human date formats.

Confidence is heuristic in Phase 2 but deterministic.

## Consequences

- The backend owns first-guess intelligence.
- Frontend can display better states without introducing new feature areas.
- Fixture tests assert shape/state/field/anomaly contracts.

## Alternatives Considered

- Send raw text to a local LLM first: deferred because deterministic parsers should handle obvious structure before probabilistic summarization.
