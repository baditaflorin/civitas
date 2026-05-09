# Phase 3 Findings Synthesis

Date: 2026-05-10

## Top 5 Usability Gaps

| Gap | Resolution |
| --- | --- |
| Users could upload only one file at a time. | Multi-file picker and drag/drop share one upload path with progress/result feedback. |
| Users could not paste copied text/HTML evidence. | Paste box uploads copied text/HTML as evidence files. |
| Safe export was preview-only. | Export can be copied, downloaded, and printed. |
| No state backup/import existed. | Versioned case-state JSON can be downloaded and imported. |
| Reload lost selected case/search context. | Endpoint, selected case id, and search term persist locally with a reset control. |

## Half-Baked Feature Outcomes

| Feature | Decision | Outcome |
| --- | --- | --- |
| Safe export preview | Finish | Copy/download/print added. |
| Single-file upload | Finish | Batch, drag/drop, and paste paths added around the same endpoint. |
| Demo graph | Finish as honest empty state | Placeholder copy no longer implies loaded sample evidence. |
| Backend endpoint setting | Finish | Settings panel now includes connect, status, and start-fresh reset. |
| README OCR/transcription/indexing claims | Edit | Claims now match processor-needed/native-tool limitations. |

## Codebase Pain Points

| Pain point | Resolution |
| --- | --- |
| `CivitasWorkspace.tsx` owned many workflows. | Accepted as coordinator debt for Phase 3; core repeated logic was extracted. |
| API response casts repeated in frontend. | `requireData` and JSON narrowing removed repeated unsafe casts. |
| Export formatting embedded in handlers. | Moved to `internal/exporter`. |
| Session persistence was one-off. | Moved to `src/lib/session.ts`. |
| E2E only checked page load. | Smoke now covers the real fresh-user workflow. |

## Documentation/Reality Mismatches

| Mismatch | Resolution |
| --- | --- |
| README implied native OCR/transcription workflows were complete. | Rewritten as processor-needed/future-adapter limitation. |
| README implied heavier indexing infrastructure than exists. | Rewritten around current search/graph/timeline behavior. |
| "Document dumps" was not true with one-file upload. | Batch upload and paste added. |
| Privacy docs mentioned only endpoint storage. | Updated to include selected case id and search term. |
| API docs lacked state backup/import. | State export/import curl examples added. |

## Fully Usable Means

1. A stranger can create a case, load several real files at once, and see per-file outcomes.
2. A stranger can paste copied evidence without first creating a file.
3. A stranger can generate, copy, download, and print a safe export.
4. A stranger can download a state file and import it into a fresh backend.
5. A stranger can reload the app and continue from the same endpoint, selected case, and search context.

## Phase 3 Success Metrics

- Multi-file upload succeeds with partial-success feedback: shipped through `uploadFiles`.
- Paste-to-evidence creates a document without manual file creation: shipped and smoke-tested.
- Safe export supports copy, markdown download, and print: shipped.
- State export/import round-trip restores case metadata and documents: shipped and smoke-tested.
- E2E covers create case, paste/upload, export, state export, and state import: shipped.
- Type-safety audit has no unsafe frontend casts outside API boundary helpers: shipped.
- README claims match implemented/tested behavior: shipped.

## Out of Scope

- No direct URL scraping/fetching.
- No folder upload beyond ZIP-as-archive input.
- No hash/share URLs containing evidence.
- No new inference behavior beyond Phase 2.
- No visual polish work beyond making controls truthful and reachable.
