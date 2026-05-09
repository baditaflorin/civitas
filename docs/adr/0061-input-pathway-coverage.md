# 0061 - Input Pathway Coverage Policy

## Status

Accepted

## Context

Users bring civic evidence as files, batches, copied text, copied HTML, and occasionally archives. Browser URL fetching and folder upload are unreliable or privacy-sensitive.

## Decision

Support these v0.3.0 inputs:

- Single and multi-file upload.
- Drag and drop onto the upload surface.
- Paste box for text/HTML.
- One-click sample evidence through the same backend upload path.
- Mobile picker through standard file input.

Explicitly out of scope:

- Direct URL scraping/fetching.
- Folder upload; users should ZIP folders.
- Clipboard image read; image files remain supported through file upload.

## Consequences

The frontend remains static and Mode C stays unchanged. All ingestion continues through the backend upload endpoint.

## Alternatives Considered

Browser-side URL fetch was rejected because CORS would fail on many civic sites and produce confusing behavior.
