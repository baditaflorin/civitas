# Civitas V1 Postmortem

## Built

- Public GitHub Pages app at https://baditaflorin.github.io/civitas/
- Go API backend with health, readiness, metrics, version, processor registry, cases, uploads, search, graph, timeline, and safe markdown exports.
- OpenAPI contract and generated TypeScript client types.
- Dockerfile, GHCR-oriented Make targets, Docker Compose, nginx, Prometheus, and Grafana starter dashboard.
- Local hooks, smoke tests, unit tests, and Playwright e2e.

## Deployment Mode In Hindsight

Mode C was the correct choice. A static-only app could show a demo and maybe query prebuilt public data, but it could not responsibly handle private leaked dumps, native OCR, transcription, face blurring, geospatial tooling, or local LLM processing. Mode B would be useful later for public release corpora, not for the core private ingestion workflow.

## Worked

- GitHub Pages was enabled from the first commit.
- The frontend stayed static and small enough for the first-load asset budget.
- Backend processors are behind a registry, so native tools can be added without changing the API.
- Build output and ADR documentation now coexist under `docs/`.

## Did Not Work Yet

- Native adapters are discovery-ready but not deeply wired to every external tool yet.
- Authentication is intentionally absent in v1; deployment should rely on network controls until auth is added.
- The first graph and entity extraction are useful fallbacks, not replacements for spaCy, Stanza, sentence-transformers, libpostal, or Tantivy.

## Surprises

- Vite's default `emptyOutDir` behavior can erase `docs/adr` when Pages publishes from `docs/`; the build script now cleans generated assets only.
- `go test ./...` can wander into `node_modules`; Make targets filter Go packages.

## Accepted Tech Debt

- Filesystem storage instead of DuckDB/Tantivy-backed indexes.
- Basic fallback extraction until the full native processor image is expanded.
- No user auth, roles, or collaborative editing in v1.

## Next Improvements

1. Add real Tika, Tesseract, Pandoc, ExifTool, qpdf, Ghostscript, Whisper.cpp, and GDAL adapters in the Docker image.
2. Replace fallback search and graph with DuckDB, Tantivy, and sentence-transformer embeddings.
3. Add authentication, case-level access controls, and audit logs.

## Time

Estimated v1 scaffold: 4 to 6 hours.

Actual implementation pass: about 1.5 hours for the functional scaffold and deployment surface, excluding future native processor integration.
