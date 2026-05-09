# Phase 3 Postmortem

Date: 2026-05-10

Version: v0.3.0

## Audit Grids

| Audit | Before | After |
| --- | --- | --- |
| Input pathways | 1 red, 4 yellow, 8 gray | 9 green, 5 explicitly out of scope, 0 red/yellow among claimed paths |
| Output pathways | 1 green, 3 yellow, 10 gray | 8 green, 6 explicitly out of scope, 0 red/yellow among claimed paths |
| Controls | Several controls stopped short of their labels | No visible production stubs; all visible controls complete their labeled workflow |
| Feature claims | 5 reality mismatches | README/API/privacy/in-app claims aligned with implemented behavior |
| Codebase health | Repeated API casts, export formatting in handlers, weak e2e path | Boundary helpers, exporter package, state schema, and fresh-user smoke coverage |

## Half-Baked Feature Triage

| Feature | Outcome | Rationale |
| --- | --- | --- |
| Safe export preview | Finished | Copy/download/print are required for publishable evidence. |
| Single-file upload | Finished | Batch and paste input are required for real evidence dumps. |
| Demo graph | Finished as honest empty state | Keeps orientation without pretending sample data is loaded. |
| Backend endpoint setting | Finished | Connect/status/reset now make it a real settings panel. |
| URL input, folder upload, share links, screenshots, embeds | Hidden/out of scope | They were not claimed and would create privacy or reliability risks. |

## Codebase Health Metrics

| Metric | Before | After |
| --- | --- | --- |
| DRY violations in core modules | 4 listed | Core export/session/API-result duplication reduced; remaining UI coordinator size accepted in ADR 0065 |
| TODO/FIXME/XXX/HACK | 0 production markers | 0 production markers |
| Frontend unsafe casts | Repeated response casts plus multipart cast in feature code | Response casts removed; multipart cast isolated in API boundary helper |
| Dead code | No abandoned files found | No abandoned files found |
| Real-user path tests | Page-load-only e2e | Create case, paste upload, export, state export, state import |

## Stranger Test

The stranger test found three practical issues: a local port collision in smoke, an ambiguous filename assertion, and missing state-import e2e coverage. All three were fixed before the final run.

Final smoke command:

```sh
CIVITAS_SMOKE_BACKEND_PORT=18090 CIVITAS_SMOKE_PAGES_PORT=4180 make smoke
```

Result: 2/2 Playwright tests passed.

## Documentation Fixes

- README now lists verified features and limitations instead of implying all native processors are complete.
- API docs now include state export/import commands.
- Privacy docs now list all localStorage keys and warn that downloaded state files can contain sensitive evidence.
- Architecture docs now describe the current shape-aware pipeline and optional native tools.

## Surprises

The biggest usability gap was not an engine problem. Phase 2 made the backend smart enough, but the UI still forced users to create files manually, repeat uploads, and manually copy previews. Completing the exits made the same engine feel much more real.

The smoke test also exposed a practical local-dev risk: if Pages preview silently shifts ports, Playwright can accidentally test the wrong app. Making the script pick and pass one known port turned that into a stable gate.

## Still Open

1. Split `CivitasWorkspace.tsx` into focused Settings, Upload, Documents, Export, and Search panels.
2. Add upload cancellation for very large files.
3. Add backend-auth and case-level access controls.
4. Add true folder workflows after designing archive privacy and progress behavior.
5. Add native processor images for OCR, PDF extraction, audio transcription, geocoding, and face blurring.

## Honest Take

A stranger can now use Civitas for their own real work end-to-end if they have a reachable backend: create a case, load real files or pasted evidence, inspect results, export a safe report, and back up/import the state. It is not yet a full forensics lab because native OCR/transcription/geospatial/LLM processors are still adapter work, and it is not production-private without deployment-level access controls. But it no longer feels like a toy on the everyday workflow; the main remaining rough edge is depth of native processing, not basic usability.
