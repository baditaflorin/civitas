# Civitas

Live site: https://baditaflorin.github.io/civitas/

Repository: https://github.com/baditaflorin/civitas

Support the project: https://www.paypal.com/paypalme/florinbadita

Civitas is a civic investigation OS for turning leaked document dumps into searchable, analyzable, publishable evidence.

## Quickstart

```sh
make install-hooks
make dev
make build
make test
make smoke
```

## Architecture

Civitas uses a GitHub Pages frontend and a self-hosted Docker backend. The static frontend is public and safe to cache; the backend runs private ingestion, OCR, transcription, indexing, and redaction workflows against local evidence.

```mermaid
flowchart LR
  investigator["Investigator browser"] --> pages["GitHub Pages frontend"]
  pages --> api["Docker backend API"]
  api --> pipeline["Evidence pipeline"]
  pipeline --> storage["Local evidence storage"]
  pipeline --> index["Search and graph indexes"]
  api --> export["Safe publishing exports"]
```

## Documentation

- Architecture decisions: docs/adr/
- Architecture overview: docs/architecture.md
- API guide: docs/api.md
- Deployment guide: deploy/README.md
- Runbook: docs/runbook.md
- Postmortem: docs/postmortem.md

## Git Hooks

Run `make install-hooks` once after cloning. Hooks run local checks; this project intentionally does not use GitHub Actions.
