# Phase 3 Input Pathway Audit

Date: 2026-05-10

Scope: Civitas v0.2.0 live/workspace UI, API, README, ADRs, and tests.

Status key: Green = works fully on real user data. Yellow = works partially. Red = claimed or visible but broken/confusing. Gray = not built and not claimed.

| Input pathway | Status before | Evidence | User impact |
| --- | --- | --- | --- |
| Single file upload | Yellow | `input[type=file]` uploads one file to `/api/v1/cases/{case_id}/documents`; no reset, progress detail, or per-file result. | A stranger can upload one file after creating a case, but the flow is easy to miss and gives little feedback. |
| Drag and drop | Gray | No drop target or drag handler. | Common desktop evidence workflow is absent. |
| Paste text/HTML | Gray | No paste box or clipboard handler. | Users with copied FOIA text, article HTML, or email excerpts must create files manually. |
| Paste image | Gray | No clipboard image handler. | Screenshot/photo evidence from clipboard cannot enter the system. |
| URL input | Gray | No URL input. README does not claim direct fetching. | Acceptable if documented as out of scope because browser CORS would make it unreliable. |
| Clipboard read button | Gray | No `navigator.clipboard` usage. | A paste box is safer and more predictable for v0.3.0 than a permission-heavy button. |
| Mobile picker | Yellow | Native file input exists but lacks `multiple`, `accept`, or camera hints. | iOS/Android Files may work; multi-select and camera/photo workflows are not intentional. |
| Multi-file upload | Gray | File input accepts one file and mutation accepts one `File`. | Real evidence dumps require repeated manual uploads. |
| Folder upload | Gray | No `webkitdirectory`, no archive-as-folder UX. | Out of scope for v0.3.0; ZIP upload is the safer documented route. |
| Sample/demo input | Red | Demo graph exists, but no one-click loadable evidence sample through the same upload flow. | The app looks demo-able, but the demo is not a real ingestion path. |
| Deep links | Gray | No case/document route, hash state, or query handling. | Fine for v0.3.0 if not claimed; shareable app state is more important. |
| Imported state | Gray | Markdown export exists, but no state import path. | Users cannot move work between backend instances or recover from a downloaded state file. |
| Restored autosave/session | Yellow | API endpoint persists in `localStorage`; selected case and search term do not persist. Backend storage persists cases/documents. | A reload loses context and defaults to the newest case/search. |
| Start fresh / clear state | Gray | No clear current session or delete/reset control. | A stranger has no obvious way to recover from a mistaken endpoint or old local UI state. |

## Input Findings

1. Single-file upload exists but is the only real input path.
2. Multi-file and paste are the highest-impact missing paths because they do not require engine changes.
3. URL input should remain out of scope unless a backend fetch/proxy is designed; direct browser fetch would fail unpredictably on real sites.
4. Demo data should be loaded through the same backend upload path, not remain only an SVG/graph placeholder.
5. Session restore needs to persist selected case and search term in addition to the API endpoint.
