# Phase 3 Output Pathway Audit

Date: 2026-05-10

Status key: Green = works fully on real user data. Yellow = works partially. Red = claimed or visible but broken/confusing. Gray = not built and not claimed. Out = explicitly out of scope by ADR.

| Output pathway | Status before | Status after | Evidence | Outcome |
| --- | --- | --- | --- | --- |
| Safe markdown export | Yellow | Green | `/exports` still generates markdown; UI now exposes it with downstream actions. | Export is useful, not preview-only. |
| Download safe export | Gray | Green | `Download markdown` creates a `.md` browser download. | Users can take the publishable artifact out. |
| Copy export to clipboard | Gray | Green | `Copy export` writes export body to `navigator.clipboard` and reports success/failure. | Common paste-into-editor flow works. |
| Downloadable state file | Gray | Green | `/api/v1/cases/{case_id}/state` returns versioned `civitas.case_state.v1` JSON and the UI downloads it. | Users can back up and migrate investigations. |
| Imported state round-trip | Gray | Green | `/api/v1/case-states/import` restores case metadata, document metadata, and source bytes; smoke covers import. | State export is recoverable, not a dead artifact. |
| CSV export | Gray | Out | ADR 0062 keeps CSV out of scope because it was not claimed. | Docs do not promise it. |
| JSON export | Gray | Green | Case state JSON is documented and downloadable. | JSON output exists for full investigation state. |
| Code/API export | Yellow | Green | Export panel shows a current-case curl command for state download; `docs/api.md` has import/export commands. | Automators no longer have to invent IDs and endpoints. |
| Copy search/result snippets | Gray | Out | Not claimed; ADR 0062 prioritizes export copy. | No placeholder control exists. |
| Share link | Gray | Out | ADR 0062 rejects share links containing evidence in URLs. | Sensitive state is not encoded into links. |
| Print-friendly view | Yellow | Green | `Print export` opens a markdown-only print view. | Users can print the safe export instead of the whole workspace. |
| Screenshot | Gray | Out | Not claimed. | No misleading control exists. |
| Embed code | Gray | Out | Not claimed. | No misleading control exists. |
| Public Pages links | Green | Green | Header links to GitHub repo and PayPal; e2e asserts both URLs. | Users can star/support from the live page. |

## Output Summary

Before: 1 green, 3 yellow, 10 gray.

After: 8 green, 6 explicitly out of scope by ADR, 0 red, 0 yellow among claimed controls.
