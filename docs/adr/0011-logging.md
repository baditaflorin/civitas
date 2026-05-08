# 0011 - Logging Strategy

## Status

Accepted

## Context

Mode C needs operational logs while avoiding accidental evidence leakage.

## Decision

Use Go `slog` JSON logs to stdout with level, timestamp, message, request id, method, path, status, duration, and trace id where available. Do not log uploaded file contents, extracted text, or private evidence values.

Frontend production builds avoid console logging except fatal boot diagnostics.

## Consequences

- Docker and server operators can collect structured logs.
- Evidence leakage risk is reduced.
- Debug detail should go into local artifacts only when explicitly requested.

## Alternatives Considered

- Text logs: rejected because JSON is easier to parse.
- Verbose processor logs by default: rejected because evidence pipelines are sensitive.
