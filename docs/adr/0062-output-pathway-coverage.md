# 0062 - Output Pathway Coverage Policy

## Status

Accepted

## Context

The app generated markdown export text but did not let users copy, download, print, or back up a case.

## Decision

Support these v0.3.0 outputs:

- Copy safe markdown export.
- Download safe markdown export.
- Print safe markdown export.
- Download versioned case state JSON.
- Import versioned case state JSON.
- Show current-case curl snippets.

Keep out of scope:

- CSV exports not already claimed by the app.
- Share links that encode evidence in the URL.
- Screenshots and embeds.

## Consequences

Users get portable artifacts without adding a new deployment mode.

## Alternatives Considered

Adding more export formats was rejected; completing the existing markdown and state flows is higher value.
