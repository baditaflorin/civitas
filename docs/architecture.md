# Civitas Architecture

Live site:

https://baditaflorin.github.io/civitas/

Repository:

https://github.com/baditaflorin/civitas

## Context

```mermaid
C4Context
  title Civitas context
  Person(investigator, "Investigator", "Journalist, researcher, civic watchdog")
  System_Boundary(pages, "GitHub Pages") {
    System(frontend, "Civitas frontend", "Static Vite React app")
  }
  System_Boundary(server, "Self-hosted Docker server") {
    System(api, "Civitas API", "Go REST API")
    System(pipeline, "Evidence pipeline", "Native processor adapters")
    SystemDb(storage, "Local storage", "Evidence and derived artifacts")
  }
  System_Ext(github, "GitHub", "Repository and Pages hosting")
  System_Ext(ghcr, "GHCR", "Backend container images")
  Rel(investigator, frontend, "Uses", "HTTPS")
  Rel(frontend, api, "Calls configured backend", "HTTPS/JSON")
  Rel(api, pipeline, "Runs ingestion and export jobs")
  Rel(pipeline, storage, "Reads and writes")
  Rel(github, frontend, "Publishes")
  Rel(ghcr, api, "Provides image")
```

## Containers

```mermaid
flowchart LR
  browser["Investigator browser"] --> pages["GitHub Pages static app"]
  pages --> client["Generated OpenAPI client"]
  client --> nginx["nginx on host port 25342"]
  nginx --> api["Go API :8080"]
  api --> cases["Case service"]
  api --> processors["Processor registry"]
  cases --> fs["Filesystem storage"]
  processors --> native["Tika, Tesseract, Pandoc, qpdf, ExifTool, Whisper.cpp, GDAL, llama.cpp adapters"]
  api --> metrics["Prometheus /metrics"]
```

## Module Boundaries

- `frontend`: `src/` plus Vite config, built into `docs/`.
- `api`: OpenAPI contract consumed by generated frontend types.
- `cmd/server`: backend entrypoint and graceful shutdown.
- `internal/httpapi`: routing, handlers, JSON responses, CORS, metrics.
- `internal/pipeline`: processor registry and v1 extraction orchestration.
- `internal/storage`: filesystem case, document, graph, timeline, and export storage.
- `deploy`: production Docker Compose, nginx, Prometheus, and run instructions.
