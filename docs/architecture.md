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
    System(pipeline, "Evidence pipeline", "Shape-aware analysis and processor states")
    SystemDb(storage, "Local storage", "Evidence, derived artifacts, and exports")
  }
  System_Ext(github, "GitHub", "Repository and Pages hosting")
  System_Ext(ghcr, "GHCR", "Backend container images")
  Rel(investigator, frontend, "Uses", "HTTPS")
  Rel(frontend, api, "Calls configured backend", "HTTPS/JSON")
  Rel(api, pipeline, "Runs ingestion, analysis, and export")
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
  api --> exporter["Safe export and case state serializer"]
  cases --> fs["Filesystem storage"]
  processors --> native["Optional native tools: Tika, Tesseract, Pandoc, qpdf, ExifTool, Whisper.cpp, GDAL, llama.cpp"]
  api --> metrics["Prometheus /metrics"]
```

## Module Boundaries

- `src/`: static React frontend, built into `docs/` for GitHub Pages.
- `src/lib/api`: generated OpenAPI client types plus small boundary helpers.
- `src/lib/session`: browser session preference storage and migration surface.
- `api/openapi.yaml`: REST contract consumed by the frontend.
- `cmd/server`: backend entrypoint and graceful shutdown.
- `internal/httpapi`: routing, handlers, JSON responses, CORS, metrics.
- `internal/pipeline`: evidence shape classification, inference, confidence, and processor-needed states.
- `internal/exporter`: safe markdown export and portable case-state serialization.
- `internal/storage`: filesystem case, document, export, and state import storage.
- `deploy`: production Docker Compose, nginx, Prometheus, and run instructions.
