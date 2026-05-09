# Phase 3 Codebase Health Audit

Date: 2026-05-10

Measurement commands:

- `rg -n "TODO|FIXME|XXX|HACK" . -g '!docs/assets/**' -g '!test/fixtures/**' -g '!node_modules/**' -g '!bin/**' -g '!test-results/**'`
- `rg -n "\bany\b|as never| as [A-ZA-Za-z_]|unknown|@ts-ignore" src internal api -g '!src/lib/api/schema.d.ts'`
- `find internal src -type f \( -name '*.go' -o -name '*.ts' -o -name '*.tsx' \) -not -path '*/schema.d.ts' -exec wc -l {} + | sort -nr | head -n 20`

## Measurements

| Metric | Before Phase 3 | After Phase 3 |
| --- | --- | --- |
| TODO/FIXME/XXX/HACK | 0 outside generated/binary/test fixture text | 0 production TODO/FIXME/XXX/HACK markers |
| Type safety holes | 1 frontend unsafe error cast, 1 multipart cast in feature code, 9 frontend response casts, Go `any` boundary usage | Frontend response casts removed; error and GitHub JSON use `unknown` narrowing; multipart cast isolated in `src/lib/api/client.ts`; Go `any` remains only JSON/OpenAPI/storage boundary code |
| Largest files | `CivitasWorkspace.tsx` 464 lines, `internal/pipeline/analyze.go` 316 lines, `internal/httpapi/handlers.go` 278 lines, `internal/storage/store.go` 243 lines | `CivitasWorkspace.tsx` 751 lines, accepted Phase 3 UI coordinator debt; exporter/state moved out of handlers; storage grew to support import |
| Dead production UI controls | 0 dead buttons, several half-finished controls | 0 visible production stubs; all visible controls have handlers or are absent |
| Test coverage holes | No e2e coverage for create/upload/export/state/settings paths | Smoke covers published page, public links, create case, paste upload, safe export, state export, and state import |

## DRY / Boundary Results

| Area | Before | After |
| --- | --- | --- |
| API query boilerplate | Repeated casts and response checks in `CivitasWorkspace.tsx` | `requireData` centralizes result/error handling. |
| Cache invalidation | Manual repeated invalidations inline | `invalidateEvidence` centralizes document/graph/timeline/search invalidation. |
| Session persistence | Endpoint-only helper | `src/lib/session.ts` owns endpoint, selected case, search, and reset. |
| Export formatting | Markdown formatting lived in HTTP handler | `internal/exporter` owns markdown, redaction, stable export IDs, and case-state serialization. |
| State schema | No canonical portable state type | `internal/evidence` and OpenAPI now define `CaseState` / `StateDocument`. |

## Accepted Debt

`CivitasWorkspace.tsx` is now larger because Phase 3 intentionally completed visible workflows before splitting UI composition. ADR 0065 accepts this as short-term coordinator debt; the next clean split is Settings, Upload, Export, and Documents panels after behavior is locked.

Go `any` remains in boundary code where JSON maps, OpenAPI response envelopes, storage serialization, and generic JSON inference require it. It is not used as a shortcut inside core typed evidence models.

## Dead Code

No abandoned production files or unused production controls remain by inspection and lint.

## Inconsistent Patterns Closed

- Frontend API data now goes through `requireData`.
- Browser persistence now goes through session helpers.
- Export/state serialization is behind the exporter package.
- User-visible production controls now either complete the labeled workflow or are not present.
