# Phase 3 Feature Claims Audit

Date: 2026-05-10

| Claim source | Claim | Status before | Status after | Reality check |
| --- | --- | --- | --- | --- |
| README | "Turning leaked document dumps into searchable, analyzable, publishable evidence." | Shipped partially | Shipped with documented limits | Batch upload, paste, search, map, timeline, export, and state backup now work; native processors remain limitations. |
| README architecture | "private ingestion, OCR, transcription, indexing, and redaction workflows" | Overclaimed | Fixed | README now says OCR/transcription/native processing are processor-needed states or future adapters unless the backend image is extended. |
| README diagram | "Search and graph indexes" | Overclaimed | Fixed | Diagram now says shape-aware analysis and safe/state exports. |
| README quickstart | `make install-hooks`, `make dev`, `make build`, `make test`, `make smoke` | Shipped fully | Shipped fully | Targets exist; smoke passed on v0.3.0. |
| docs/api.md | Curl examples for cases, uploads, search, export | Shipped fully | Shipped fully | API endpoints exist and now include state export/import examples. |
| docs/privacy.md | Frontend stores only endpoint, not evidence | Drift risk | Fixed | Privacy docs list endpoint, selected case id, and search term; evidence remains backend-side. |
| Phase 2 docs | Confidence/provenance/state on uploads | Shipped fully | Shipped fully | Fixture suite still passes. |
| In-app header | Version and commit visible | Shipped fully | Shipped fully | Header shows v0.3.0 and GitHub commit fallback. |
| In-app demo map | Demo relationship map | Shipped partially | Fixed | Placeholder copy now says example layout until evidence is uploaded. |
| Phase 3 docs | "No stubs left in production UI" | Not applicable | Shipped | Controls audit and smoke test cover visible production controls. |

## Claim Summary

All README/API/privacy/in-app claims now either match implemented behavior or explicitly name the limitation. Claims without tests or audit evidence were removed.
