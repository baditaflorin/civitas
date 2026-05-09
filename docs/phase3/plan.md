# Phase 3 Completeness Plan

Date: 2026-05-10

Ranked by real-user impact and constrained to completeness, not polish or engine changes.

## Picklist

1. Multi-file upload with per-file progress/result feedback. Covers catalog 1, 4, 5.
2. Drag-and-drop upload onto the existing Upload surface. Covers 1.
3. Paste text/HTML evidence through a first-class paste box. Covers 1, 6.
4. Clipboard permission fallback: use paste box rather than a fragile clipboard-read button. Covers 6.
5. Sample evidence loader that uses the same backend upload path as user data. Covers 7.
6. Persist selected case across reload. Covers 8, 38.
7. Persist search term across reload. Covers 8, 38.
8. Add one-click local UI reset/start fresh. Covers 8, 40.
9. Copy safe export to clipboard with visible confirmation. Covers 10.
10. Download safe markdown export from the browser. Covers 9.
11. Print the safe export instead of printing the whole workspace. Covers 13.
12. Add downloadable case state JSON with versioned schema. Covers 11, 38.
13. Add state import endpoint and UI. Covers 9, 41.
14. Add state export/import round-trip tests. Covers 9, 39, 41.
15. Add current-case API/curl snippet in the app. Covers 14.
16. Hide/directly mark URL input, folder upload, share links, screenshots, embeds as out of scope in ADR/docs. Covers 15, 17, 19.
17. Replace misleading demo graph copy with honest empty-state language. Covers 16, 17.
18. Add settings panel with endpoint, reset, and session persistence behavior; every setting must do something. Covers 18, 38, 40.
19. Extract frontend session persistence helpers into one module. Covers 20, 21, 32.
20. Extract frontend API result helpers to remove repeated unsafe casts. Covers 20, 21, 35, 36.
21. Extract backend markdown/state export logic out of handlers. Covers 24, 27.
22. Add canonical state schema types in the domain model. Covers 22, 23, 39.
23. Consolidate query invalidation after document/state changes. Covers 20, 32.
24. Audit and reduce `any`/unsafe cast usage outside boundary code. Covers 35, 36.
25. Update README/API/privacy/docs so every claim matches reality and has a test or audit row. Covers 42, 43, 45.
26. Add e2e happy path for create case, paste/upload, export copy/download visibility, and state controls. Covers 43, 46.
27. Run stranger test in a fresh browser context and fix top three findings. Covers 46, 47.

## Batches

1. ADRs 0060-0071. Done.
2. Backend state export/import and exporter extraction. Done.
3. Frontend input/output/session/settings completion. Done.
4. Tests and documentation alignment. Done.
5. Stranger test, postmortem, version bump, Pages build, tag, push. Done after final verification/tag.

## Explicit Non-Goals

- No direct URL scraping in the browser.
- No folder upload beyond ZIP-as-archive input.
- No hash/share URLs containing evidence.
- No new inference behavior beyond surfacing existing Phase 2 output.
- No visual polish work beyond making controls truthful and reachable.
