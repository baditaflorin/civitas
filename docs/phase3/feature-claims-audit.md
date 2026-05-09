# Phase 3 Feature Claims Audit

Date: 2026-05-10

| Claim source | Claim | Status before | Reality check |
| --- | --- | --- | --- |
| README | "Turning leaked document dumps into searchable, analyzable, publishable evidence." | Shipped partially | Search/graph/timeline/export exist, but input/output paths are incomplete for dumps. |
| README architecture | "private ingestion, OCR, transcription, indexing, and redaction workflows" | Shipped partially | Ingestion/redaction metadata exist; OCR/transcription are honest `needs_processor` states, not full workflows. |
| README diagram | "Search and graph indexes" | Shipped partially | Search/graph are computed from stored docs, not Tantivy/DuckDB-style indexes. |
| README quickstart | `make install-hooks`, `make dev`, `make build`, `make test`, `make smoke` | Shipped fully | Local targets exist and passed in Phase 2. |
| docs/api.md | Curl examples for cases, uploads, search, export | Shipped fully | API endpoints exist. |
| docs/privacy.md | Frontend stores only API endpoint, not evidence | Shipped fully before Phase 3 | Evidence remains backend-side; Phase 3 client session preferences need docs update. |
| Phase 2 docs | Confidence/provenance/state on uploads | Shipped fully | Fixture suite asserts it. |
| In-app header | Version and commit visible | Shipped fully | Bundle shows v0.2.0 and GitHub main commit fallback. |
| In-app demo map | Demo relationship map | Shipped partially | It is a visual placeholder, not a loadable sample. |

## Claim Findings

1. README overstates OCR/transcription as workflows when they are currently processor-needed states.
2. README implies indexing infrastructure that is not present.
3. "Document dump" usability is not true until multi-file and state export/import are finished.
4. Privacy docs must be updated when more UI preferences are persisted.
5. Public project links are true and tested.
