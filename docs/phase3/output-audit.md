# Phase 3 Output Pathway Audit

Date: 2026-05-10

Status key: Green = works fully on real user data. Yellow = works partially. Red = claimed or visible but broken/confusing. Gray = not built and not claimed.

| Output pathway | Status before | Evidence | User impact |
| --- | --- | --- | --- |
| Safe markdown export | Yellow | `Safe export` calls `/exports` and renders the body in a `<pre>`. | Export exists but cannot be downloaded or copied from the UI. |
| Download safe export | Gray | Backend writes export file locally; frontend does not download it. | A user cannot take the generated report out of the browser without manual selection. |
| Copy export to clipboard | Gray | No copy handler. | Common reporting workflow is missing. |
| Downloadable state file | Gray | No case state JSON endpoint or frontend download. | Users cannot back up or move an investigation. |
| Imported state round-trip | Gray | No import endpoint/UI. | Exported state would not be recoverable. |
| CSV export | Gray | README does not claim CSV export. | Out of scope for v0.3.0; avoid claiming it. |
| JSON export | Gray | No user-facing JSON state export, though API returns JSON. | State JSON is the useful Phase 3 target. |
| Code/API export | Yellow | `docs/api.md` has curl examples; UI does not expose a command for the current case. | Automators must read docs and manually fill IDs. |
| Copy search/result snippets | Gray | No copy action. | Lower priority than export copy. |
| Share link | Gray | No shareable state URL. | Out of scope for backend evidence; hash state may leak sensitive data if done casually. |
| Print-friendly view | Yellow | Browser can print the app, but no export print action or print CSS. | Safe export should have a print button. |
| Screenshot | Gray | Not claimed. | Out of scope. |
| Embed code | Gray | Not claimed. | Out of scope. |
| Public Pages links | Green | Header links to GitHub repo and PayPal; Pages bundle contains both URLs. | Users can star/support the project from the live page. |

## Output Findings

1. The core markdown export is not actually portable from the UI.
2. No state backup/import path exists, so persistence depends on one backend instance.
3. API automation is documented but not contextualized inside the app.
4. Share links and screenshots are not claimed and should stay out of scope for Phase 3.
5. Print should target the generated safe export, not the whole operational workspace.
