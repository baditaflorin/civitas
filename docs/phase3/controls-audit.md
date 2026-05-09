# Phase 3 Controls Audit

Date: 2026-05-10

| Control | Status before | Status after | Handler/result |
| --- | --- | --- | --- |
| Star on GitHub | Green | Green | Opens https://github.com/baditaflorin/civitas. |
| PayPal | Green | Green | Opens https://www.paypal.com/paypalme/florinbadita. |
| API endpoint input | Yellow | Green | Valid endpoint is persisted; settings show backend status. |
| Connect | Yellow | Green | Applies endpoint, persists it, refetches data, and reports success. |
| Start fresh | Gray | Green | Clears local endpoint/selected-case/search state and leaves backend evidence intact. |
| New case | Yellow | Green | Creates a case, selects it, and smoke verifies the fresh-user path. |
| Case row | Green | Green | Selects case and persists selected id across reload. |
| Search input | Yellow | Green | Runs search and persists term across reload. |
| File picker | Yellow | Green | Accepts multiple files and shares batch progress/result handling. |
| Drag/drop area | Gray | Green | Uploads dropped files. |
| Paste box / Upload paste | Gray | Green | Creates text/HTML evidence and uploads it. |
| Load sample | Red | Green | Loads bundled sample through the same backend upload path. |
| Import state | Gray | Green | Validates and imports `civitas.case_state.v1` JSON. |
| Safe export | Yellow | Green | Generates markdown and unlocks copy/download/print. |
| Download state | Gray | Green | Downloads portable case state JSON. |
| Evidence map | Yellow | Green | Shows real graph when evidence exists and honest placeholder copy otherwise. |
| Processors list | Yellow | Green | Lists processor availability; README explains native-heavy limitations. |
| Timeline list | Yellow | Green | Shows timeline candidates and no longer claims more than the backend provides. |

## Control Summary

Production UI now has no visible stub controls. Controls either perform the labeled end-to-end action on real data or are absent and documented out of scope.
