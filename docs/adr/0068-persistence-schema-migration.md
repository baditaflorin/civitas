# 0068 - Persistence Schema and Migration Policy

## Status

Accepted

## Context

Users need a downloadable state file to move or restore a case.

## Decision

Add `civitas.case_state.v1` JSON:

- Case metadata.
- Document metadata.
- Original uploaded bytes as base64.
- Export timestamp and app version metadata.

Importing v1 replays stored document metadata and bytes into filesystem storage without re-running inference. Future breaking changes use a new schema version and importer branch.

## Consequences

State files may be large and sensitive; docs must warn users to protect them.

## Alternatives Considered

State without source bytes was rejected because it would not truly restore the case.
