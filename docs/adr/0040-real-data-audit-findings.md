# 0040 - Real-Data Audit Findings and Substance Metrics

## Status

Accepted

## Context

The Phase 2 audit showed that v1 uploaded 10/10 real inputs without crashing but produced 0/10 useful first investigative guesses. The dominant failure mode was wrong confidence: placeholder or raw text output was marked `completed`.

## Decision

Phase 2 success is defined by useful first guesses, explicit processing states, confidence, and deterministic provenance on the same Mode C architecture.

Target metrics:

- At least 7/10 real fixtures produce useful first guesses.
- 10/10 fixtures expose evidence shape, state, confidence, and anomalies.
- Same input produces deterministic normalized output.
- Broken/empty/corrupt evidence is never marked as cleanly complete.

## Consequences

- Tests are fixture-driven.
- Exports must include confidence and provenance.
- "No crash" is no longer considered success.

## Alternatives Considered

- Add more native tools first: rejected because the app must still behave honestly when tools are absent.
