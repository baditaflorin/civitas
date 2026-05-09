# 0060 - Completeness Audit Findings and Phase 3 Metrics

## Status

Accepted

## Context

Phase 2 made the evidence engine honest, but the app still failed common end-to-end user workflows: loading multiple real files, pasting copied evidence, taking exports out, restoring state, and resuming after reload.

## Decision

Phase 3 success is measured by:

- Multi-file upload and paste evidence work on real backend data.
- Safe export can be copied, downloaded, and printed.
- Case state can be exported and imported with a versioned JSON contract.
- Endpoint, selected case, and search term persist across reload.
- README/API/privacy claims match tested behavior.
- Phase 2 real-data fixtures remain green.

## Consequences

Completeness work may add small controls, but only to finish existing workflows. New inference features stay out of scope.

## Alternatives Considered

Polish-first UI changes were rejected because they would not unblock real user work.
