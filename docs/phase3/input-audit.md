# Phase 3 Input Pathway Audit

Date: 2026-05-10

Scope: Civitas v0.2.0 before Phase 3 and Civitas v0.3.0 after Phase 3.

Status key: Green = works fully on real user data. Yellow = works partially. Red = claimed or visible but broken/confusing. Gray = not built and not claimed. Out = explicitly out of scope by ADR.

| Input pathway | Status before | Status after | Evidence | Outcome |
| --- | --- | --- | --- | --- |
| Single file upload | Yellow | Green | Existing file input now shares the same `uploadFiles` path and resets after selection. | A user can upload a single file and see upload feedback. |
| Drag and drop | Gray | Green | Upload panel handles `dragover` and `drop`, then uploads dropped files. | Desktop evidence dumps can enter without the file picker. |
| Paste text/HTML | Gray | Green | Paste box creates `pasted-evidence.txt` or `pasted-evidence.html` and sends it to the backend. | Copied articles, email excerpts, and FOIA text no longer require manual files. |
| Paste image | Gray | Out | ADR 0061 keeps clipboard image read out of scope; image files are supported through upload/drop. | No misleading image-paste control is exposed. |
| URL input | Gray | Out | ADR 0061 rejects direct browser URL fetching because CORS failures would be unpredictable. | Users are guided to paste rendered HTML/text or upload saved files. |
| Clipboard read button | Gray | Out | ADR 0061 chooses the paste box instead of permission-heavy `navigator.clipboard` reads. | No fragile browser permission flow is exposed. |
| Mobile picker | Yellow | Green | Standard `input[type=file][multiple]` remains available and works with mobile Files pickers. | Mobile users can select available files through the OS picker. |
| Multi-file upload | Gray | Green | File input has `multiple`; smoke and unit paths use the batch upload helper with per-file progress text. | Evidence batches no longer require repeated manual single uploads. |
| Folder upload | Gray | Out | ADR 0061 keeps folder upload out of scope; ZIP archives are the documented route. | No partial directory API support is implied. |
| Sample/demo input | Red | Green | `Load sample` fetches `data/v1/sample.json` and uploads it through the same backend endpoint. | Demo and real input now exercise the same ingestion path. |
| Deep links | Gray | Out | ADR 0061/0062 avoid URL-encoded case state because evidence can be sensitive. | No sensitive data is leaked through shareable URLs. |
| Imported state | Gray | Green | Frontend validates `civitas.case_state.v1` JSON and posts it to `/api/v1/case-states/import`. | A downloaded investigation can be restored into a backend. |
| Restored autosave/session | Yellow | Green | Endpoint, selected case id, and search term persist in `localStorage`. | Reload returns users to their working context. |
| Start fresh / clear state | Gray | Green | Settings panel clears local endpoint/case/search state while leaving backend evidence intact. | Users can recover from stale local UI state. |

## Input Summary

Before: 0 red-to-green paths, 4 yellow, 8 gray, 1 red.

After: 9 green, 5 explicitly out of scope by ADR, 0 red, 0 yellow among claimed controls.
