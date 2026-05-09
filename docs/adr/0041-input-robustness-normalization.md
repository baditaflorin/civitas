# 0041 - Input Robustness and Normalization Policy

## Status

Accepted

## Context

Inputs include CSV, JSON, HTML, PDF, images, audio, archives, partial files, and empty transfers.

## Decision

Normalize text before inference:

- Strip UTF-8 BOM.
- Normalize CRLF to LF.
- Collapse repeated whitespace for preview/export snippets.
- Detect empty files before MIME handling.
- Use extension plus detected content type for shape classification.
- Validate archives/PDFs before marking state.

## Consequences

- Corrupt and partial files get recoverable errors.
- CSV/JSON/HTML are parsed before regex entity extraction.
- Binary/media files get explicit processor-needed states.

## Alternatives Considered

- Trust `http.DetectContentType` alone: rejected because OGG and empty PDFs were misclassified in the audit.
