# Phase 3 Codebase Health Audit

Date: 2026-05-10

Measurement command highlights:

- `rg -n "TODO|FIXME|XXX|HACK" . -g '!docs/assets/**' -g '!test/fixtures/**' -g '!node_modules/**' -g '!bin/**'`
- `rg -n "\\bany\\b|as never| as [A-ZA-Za-z_]|unknown" src internal api -g '!src/lib/api/schema.d.ts'`
- `find internal src -type f \( -name '*.go' -o -name '*.ts' -o -name '*.tsx' \) -not -path '*/schema.d.ts' -exec wc -l {} + | sort -nr | head -n 15`

## Measurements Before Phase 3

| Metric | Count / finding |
| --- | --- |
| TODO/FIXME/XXX/HACK | 0 outside generated/binary/test fixture text. |
| Type safety holes | 1 frontend unsafe error cast, 1 `as never` multipart cast, 9 frontend response casts, Go `any` boundary usage in JSON/OpenAPI/pipeline JSON parsing. |
| Largest files | `CivitasWorkspace.tsx` 464 lines, `internal/pipeline/analyze.go` 316 lines, `internal/httpapi/handlers.go` 278 lines, `internal/storage/store.go` 243 lines. |
| Dead production UI controls | 0 dead buttons, but several half-finished controls. |
| Test coverage holes | No e2e coverage for creating a case, uploading real data, multi-file, export copy/download, state backup/import, or settings persistence. |

## DRY Violations

| Area | Files | Finding |
| --- | --- | --- |
| API query boilerplate | `src/features/investigation/CivitasWorkspace.tsx:62-149` | Repeated `{ data, error }` handling and unsafe casts for every query. |
| Cache invalidation | `src/features/investigation/CivitasWorkspace.tsx:184-197` | Upload success invalidates four related queries manually. |
| JSON response envelopes | `internal/httpapi/handlers.go` | Repeated `map[string]any` envelopes for list responses. Acceptable boundary duplication until typed envelopes are introduced. |
| Export/state serialization | `internal/httpapi/handlers.go` | Markdown export lives in handler file; adding state export will worsen responsibility unless moved. |

## SOLID / Module Boundary Findings

| Area | Finding |
| --- | --- |
| Frontend workspace | `CivitasWorkspace.tsx` owns endpoint settings, cases, upload, search, export, documents, processors, and layout. It has too many reasons to change. |
| HTTP handlers | `handlers.go` owns request handling, export formatting, redaction, and ID generation. Export formatting belongs behind an exporter boundary. |
| Storage | `store.go` owns case/document/export persistence. It is cohesive enough for filesystem storage, but state import/export will need explicit helpers. |
| Generated OpenAPI types | Generated `schema.d.ts` contains expected `unknown` indexes; exclude from type-safety audit. |

## Dead Code

No unreferenced production files were found by inspection. The public demo graph is used as an empty-state visual, but it reads like sample data and needs clearer language.

## Inconsistent Patterns

1. Frontend API data is trusted through casts instead of boundary helpers.
2. Only the API endpoint persists; other UI state is volatile.
3. Error messages appear inline in upload only; other mutation/query errors mostly vanish.
4. Case/document state can be stored on the backend, but no export/import contract exists for a user-controlled backup.

## Codebase Health Targets

- Extract frontend persistence/session helpers.
- Extract export/state serialization out of HTTP handler flow.
- Replace unsafe frontend API casts with typed helpers where practical.
- Add tests for real-user paths, not only static page load.
