# 0017 - Dependency Policy

## Status

Accepted

## Context

Civitas integrates many mature native and language-specific tools. Custom implementations would be risky and expensive.

## Decision

Use production-ready libraries for core concerns:

- Go: chi, slog, prometheus/client_golang, validator, rs/cors, testify.
- Frontend: Vite, React, TypeScript, Tailwind, TanStack Query, Zod, Lucide, openapi-fetch.
- Processor integrations: shell adapters around battle-tested tools such as Tika, Tesseract, Pandoc, ExifTool, qpdf, Ghostscript, GDAL, and llama.cpp where installed.

Pin dependencies with lockfiles. Avoid adding libraries for trivial code.

## Consequences

- The app can grow into the requested full stack without replacing public contracts.
- Optional native tools can be missing in local development while fallback processors keep tests fast.
- Vulnerability checks are part of local security hygiene.

## Alternatives Considered

- Hand-built OCR/NLP/indexing implementations: rejected.
- Unpinned dependency ranges: rejected for reproducibility.
