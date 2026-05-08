# 0001 - Deployment Mode

## Status

Accepted

## Context

Civitas must process leaked document dumps that can contain PDFs, scans, office files, images, audio, video, geospatial files, and multilingual text. The requested toolchain includes Tika, Tesseract, Pandoc, PyMuPDF, Camelot, pdfminer, qpdf, Ghostscript, spaCy, Stanza, sentence-transformers, DuckDB, Polars, Tantivy, ExifTool, libpostal, libmagic, ssdeep, GraphViz, Whisper.cpp, pyannote, NLLB-200, ImageMagick, dlib, MediaPipe, libosmscout, GDAL, llama.cpp, and Pandoc for press export.

GitHub Pages is still the desired public surface, but browser-only execution cannot reliably run this native, heavyweight, long-running, and privacy-sensitive evidence pipeline.

## Decision

Use Mode C: GitHub Pages frontend plus Docker backend.

The frontend is a static Vite app published from `main` `/docs` to `https://baditaflorin.github.io/civitas/`. The backend is an API-only Go service deployed with Docker Compose behind nginx on public host port `25342`. Heavy processors run behind backend interfaces so deployments can add or remove optional native tools without changing the frontend contract.

## Consequences

- Private evidence never needs to be uploaded to a public static host.
- The public Pages app can link to the repository, PayPal support, version, and commit.
- Docker is required for real ingestion, indexing, redaction, and export workflows.
- The backend must expose health, readiness, metrics, and a stable OpenAPI contract.
- ADRs and deployment docs must explicitly justify the runtime backend.

## Alternatives Considered

- Mode A: rejected because browser-only OCR, transcription, NLP, media processing, geocoding, and local LLM workflows would exceed asset budget and capability limits.
- Mode B: rejected because prebuilt data artifacts do not satisfy the core v1 workflow of dropping private leaked dumps for processing.
