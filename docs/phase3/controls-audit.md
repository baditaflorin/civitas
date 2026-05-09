# Phase 3 Controls Audit

Date: 2026-05-10

| Control | Status before | Handler/result | Decision |
| --- | --- | --- | --- |
| Star on GitHub | Green | Opens https://github.com/baditaflorin/civitas. | Keep. |
| PayPal | Green | Opens https://www.paypal.com/paypalme/florinbadita. | Keep. |
| API endpoint input | Yellow | Saves endpoint to `localStorage`, but no reset/default validation message. | Finish with settings/reset behavior. |
| Connect | Yellow | Applies endpoint and refetches. | Keep, add persistence coherence. |
| New case | Yellow | Creates a case, but form values remain and errors are not surfaced inline. | Finish enough to be usable. |
| Case row | Green | Selects case. | Persist selected case. |
| Search input | Yellow | Runs search, but default search term is arbitrary and reload loses the user's query. | Persist query and allow clear. |
| File picker | Yellow | Uploads only first file. | Finish as multi-file input. |
| Safe export | Yellow | Generates export preview only. | Finish with copy/download/print. |
| Evidence map | Yellow | Shows demo graph when no backend data. | Keep but label as empty state, not sample evidence. |
| Processors list | Yellow | Lists tools but no next step. | Keep; not a stub, but docs should explain local processors. |
| Timeline list | Yellow | Shows five events only; no empty-state guidance. | Keep, not top priority. |

## Control Findings

There are no dead buttons, but several controls stop one step short of what their labels imply. `Safe export` is the clearest half-baked control: it generates text, but does not let the user safely take it with them.
