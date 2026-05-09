# 0046 - Performance Budgets

## Status

Accepted

## Context

Large CSV/audio/archive uploads must not feel frozen or dishonest.

## Decision

Budgets:

- Classify any file within 300 ms after upload bytes are available.
- Produce text/CSV/JSON/HTML preview within 1 second median on audit fixtures.
- For files over 5 MB, include size bucket and parse duration metadata.
- Native media transcription remains processor-needed until a real job runner is added.

## Consequences

- Phase 2 records parse duration and size bucket.
- Long native work is not faked.

## Alternatives Considered

- Synchronous full processing for every file: rejected for media and large archives.
