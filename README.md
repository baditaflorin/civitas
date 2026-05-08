# Civitas

![GitHub Repo stars](https://img.shields.io/github/stars/baditaflorin/civitas?style=social)
![GitHub Pages](https://img.shields.io/badge/GitHub%20Pages-live-0b5cad)
![License](https://img.shields.io/badge/license-MIT-2f6f4e)

Live site: https://baditaflorin.github.io/civitas/

Repository: https://github.com/baditaflorin/civitas

Support the project: https://www.paypal.com/paypalme/florinbadita

Civitas is a civic investigation OS for turning leaked document dumps into searchable, analyzable, publishable evidence.

![Civitas dashboard preview](docs/demo.svg)

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

- Architecture decisions: https://github.com/baditaflorin/civitas/tree/main/docs/adr
- Architecture overview: https://github.com/baditaflorin/civitas/blob/main/docs/architecture.md
- API guide: https://github.com/baditaflorin/civitas/blob/main/docs/api.md
- Deployment guide: https://github.com/baditaflorin/civitas/blob/main/deploy/README.md
- Runbook: https://github.com/baditaflorin/civitas/blob/main/docs/runbook.md
- Postmortem: https://github.com/baditaflorin/civitas/blob/main/docs/postmortem.md

## Git Hooks

Run `make install-hooks` once after cloning. Hooks run local checks; this project intentionally does not use GitHub Actions.
