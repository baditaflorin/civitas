# 0044 - Confidence Model

## Status

Accepted

## Context

No silent wrongness is a Phase 2 bar.

## Decision

Use confidence values from 0.0 to 1.0:

- 0.9-1.0: deterministic structural signal, such as valid JSON parse or ZIP validation.
- 0.7-0.89: strong heuristic, such as delimiter/header inference or HTML article title.
- 0.4-0.69: weak heuristic, such as regex-only entities.
- Below 0.4: do not present as a confident inference.

Document confidence is the weighted summary of shape detection, parse success, and anomaly severity.

## Consequences

- Confidence appears in documents, field inferences, anomalies, debug output, and exports.
- Placeholder-only output has low confidence and non-ready state.

## Alternatives Considered

- Binary high/low labels only: rejected because fixture tests need measurable thresholds.
