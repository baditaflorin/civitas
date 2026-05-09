# Phase 3 Stranger Test

Date: 2026-05-10

Tester: me-as-stranger in a fresh Playwright browser context, backed by a fresh smoke backend and a Pages preview build.

Command:

```sh
CIVITAS_SMOKE_BACKEND_PORT=18090 CIVITAS_SMOKE_PAGES_PORT=4180 make smoke
```

## Scenario

1. Open the published Pages build.
2. Confirm the repo and PayPal links are visible.
3. Create a new case named `Stranger workflow`.
4. Paste copied evidence text into the upload panel.
5. Upload the pasted evidence.
6. Generate a safe export.
7. Confirm copy/download state controls are visible.
8. Download state through the API command path and import the state JSON through the UI.

## Findings

| Finding | Severity | Response |
| --- | --- | --- |
| Smoke initially hit an unrelated local Vite server on port 4174. | High | `scripts/smoke.sh` now chooses a free Pages preview port and passes that exact backend origin into CORS. |
| Filename appeared in both map and document card, causing a strict e2e selector failure. | Low | Smoke now scopes the filename assertion to the document article. |
| State export existed but import was not covered by e2e. | Medium | Smoke now fetches current case state and imports it through the file picker. |

## Result

The top three issues found during the stranger test were fixed. The final smoke run passed 2/2 Playwright tests, including the fresh-user case/paste/export/import workflow.
