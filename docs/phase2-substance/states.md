# Phase 2 Substance State Taxonomy

Date: 2026-05-09

This file enumerates the intentional states reachable after an upload. The UI, API, export, and tests use these state names directly so a stored file cannot masquerade as understood evidence.

| State | Meaning | User-actionable exit |
| --- | --- | --- |
| `ready` | Civitas produced a useful first guess: preview, fields, confidence, provenance, and any timeline/entities it can infer. | Search it, inspect fields/confidence, export it, or upload more evidence. |
| `needs_processor` | The file is valid evidence, but useful extraction requires a native processor not available in this runtime. | Install or enable the named processor and reprocess, or keep the file as source evidence with the warning preserved. |
| `recoverable_error` | The input appears partial, corrupt, malformed, or otherwise fixable. | Re-upload the original, repair the file with the suggested tool, or preserve it as failed-source evidence. |
| `failed` | The input cannot be processed as provided. In v0.2.0 this covers zero-byte failed transfers. | Re-upload the original file or keep the failed-transfer record for provenance. |
| `unsupported` | Civitas stored the file but cannot infer a useful evidence shape from it. | Convert it to PDF, CSV, JSON, HTML, image, audio, or ZIP, then upload the converted evidence. |

## State Matrix

| Input condition | Shape | State |
| --- | --- | --- |
| CSV with readable rows | `csv` | `ready` |
| JSON or OCDS JSON | `json` / `ocds_json` | `ready` |
| HTML article or data-source page | `html_article` / `html_data_source` | `ready` |
| Valid PDF without extraction processors | `pdf` | `needs_processor` |
| PDF missing header/end marker or implausibly small | `pdf` | `recoverable_error` |
| Scan image without OCR processors | `image_scan` | `needs_processor` |
| Audio evidence without transcription processors | `audio` | `needs_processor` |
| Valid ZIP archive | `archive_zip` | `ready` |
| Truncated/corrupt ZIP archive | `archive_zip` | `recoverable_error` |
| Zero-byte upload | `empty` | `failed` |
| Unknown binary | `unknown_binary` | `unsupported` |

No state is terminal without a user-actionable next step. Every non-`ready` state must carry at least one anomaly with `message`, `why`, and `next_step`.
