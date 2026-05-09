# Phase 3 Findings Synthesis

Date: 2026-05-10

## Top 5 Usability Gaps

1. Users can upload only one file at a time; document dumps require repetitive manual work.
2. Users cannot paste copied text/HTML evidence directly.
3. Safe export is preview-only in the UI; no copy, download, or print action exists.
4. There is no state backup/import, so users cannot move or restore an investigation themselves.
5. Reload loses selected case and search context.

## Top 5 Half-Baked Features

| Feature | Decision | Rationale |
| --- | --- | --- |
| Safe export preview | Finish | Core promise says publishable evidence; preview-only is not usable. |
| Single-file upload | Finish | Phase 2 engine is useful only if users can load their actual dump. |
| Demo graph | Finish as honest empty state | Keep visual affordance, but stop implying loaded sample evidence. |
| Backend endpoint setting | Finish | It persists but has no reset/clear-state story. |
| README OCR/transcription/indexing claims | Edit | They overstate current processor-needed behavior. |

## Top 5 Codebase Pain Points

1. `CivitasWorkspace.tsx` is the primary god module.
2. API response casts are repeated and unsafe in the frontend.
3. Export formatting is embedded in HTTP handlers.
4. Session persistence is a one-off endpoint helper.
5. E2E tests do not exercise the actual user workflow.

## Top 5 Documentation/Reality Mismatches

1. README says OCR/transcription workflows; reality is explicit `needs_processor`.
2. README says indexing; reality is simple stored-document search/graph.
3. README says document dumps; UI supports one file at a time.
4. Privacy docs mention only endpoint storage; Phase 3 will persist more UI session state.
5. API docs do not mention state backup/import because it does not exist yet.

## Fully Usable Means

1. A stranger can create a case, load several real files at once, and see per-file outcomes.
2. A stranger can paste copied evidence without first creating a file.
3. A stranger can generate, copy, download, and print a safe export.
4. A stranger can download a state file and import it into a fresh backend.
5. A stranger can reload the app and continue from the same endpoint, selected case, and search context.

## Phase 3 Success Metrics

- Multi-file upload succeeds with partial-success feedback for at least three real fixtures.
- Paste-to-evidence creates a document without manual file creation.
- Safe export supports copy, markdown download, and print from the UI.
- State export/import round-trip restores case metadata and documents with deterministic IDs.
- E2E test covers create case, upload/paste, export, and visible result.
- Type-safety audit has no unsafe frontend casts outside API boundary helpers.
- README claims match implemented/tested behavior.

## Out of Scope

- No new engine inference behavior.
- No direct URL scraping/fetching.
- No folder upload; ZIP remains the supported dump container.
- No share link that encodes sensitive evidence in the URL.
- No UI polish beyond controls needed for completeness.
