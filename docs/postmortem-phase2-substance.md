# Phase 2 Substance Postmortem

Date: 2026-05-09

Version: v0.2.0

## Real-Data Pass Rate

Before: 0/10 useful results. V1 uploaded without crashing, but treated raw CSV/JSON/HTML and unsupported media as if the evidence had been understood.

After: 10/10 fixture contracts pass. The command `CGO_ENABLED=0 go test ./internal/pipeline -run 'TestRealDataFixtures' -count=1 -v` passes for all 10 real inputs, and `TestRealDataFixturesAreDeterministic` passes for all 10.

| Fixture | Before | After |
| --- | --- | --- |
| R01 OFAC SDN CSV | Wrong-but-confident raw text and noisy graph | `csv` / `ready`, headers, row count, field inference, capped entity input |
| R02 OCDS tender JSON | Raw text, no procurement model | `ocds_json` / `ready`, OCID, buyer, title, value, currency, date |
| R03 DOJ bid-rigging HTML | Markup-heavy text | `html_article` / `ready`, readable text and date timeline |
| R04 ICIJ database page | Page chrome as evidence | `html_data_source` / `ready`, source-page shape recognized |
| R05 BLM FOIA PDF | Placeholder marked complete | `pdf` / `needs_processor`, document processor next step |
| R06 NARA scan JPG | Placeholder marked complete | `image_scan` / `needs_processor`, OCR next step |
| R07 Council audio OGG | Generic binary upload | `audio` / `needs_processor`, transcription next step |
| R08 Partial ZIP | Generic binary upload | `archive_zip` / `recoverable_error`, archive corruption next step |
| R09 Truncated PDF | Same as valid PDF | `pdf` / `recoverable_error`, re-upload/repair next step |
| R10 Empty PDF | Marked complete | `empty` / `failed`, failed-transfer next step |

## Logic Gaps Closed

1. Structured formats are parsed now: CSV, JSON, OCDS JSON, and HTML have shape-specific analysis.
2. Unsupported, unavailable, corrupt, partial, and empty inputs are separate states.
3. Regex entities no longer dominate CSV graph output; extraction uses field-aware and bounded text.
4. Confidence exists on documents, fields, anomalies, entities, and exports.
5. Heavy/native work is honest: PDF/OCR/audio processors are required states instead of fake completed text.

## Smart Behaviors

- Uploading real data now immediately classifies the evidence shape.
- Civitas produces a useful first guess or useful failure state without manual configuration.
- Every non-ready state includes a domain reason and next step.
- Exports include schema version, app version, commit, source IDs, source SHA-256, confidence, processor parameters, and redaction.
- Debug state is available through `/api/v1/cases/{case_id}/debug`.

## Determinism

All 10 fixtures pass deterministic output tests after removing timestamps and measured parse duration from the comparison. Stable IDs cover documents, fields, anomalies, entities, and timeline events. Export IDs are content-derived and export sections are sorted by document ID.

## Performance

Median fixture duration is below 10 ms. The 70 MB audio fixture enters `needs_processor` in about 70-110 ms. The worst fixture is the 5.5 MB OFAC CSV at about 1.7-2.0 seconds because Civitas reads all rows to count records and infer fields.

The most important performance fix was limiting entity extraction for CSV and oversized text. That improved speed and reduced graph noise at the same time.

## Surprises

- The biggest v1 harm was not crashing. It was confidently calling placeholder work complete.
- CSV entity extraction was doing the wrong work very efficiently-looking: lots of graph nodes, little meaning.
- Corrupt files need better product language than parser language. `recoverable_error` plus a next step is much clearer than a stack-ish parse failure.

## Still Open

1. Move CSV parsing and row counting into a streaming or worker-backed job path for much larger files.
2. Add real Tika/PyMuPDF/Tesseract/Whisper processors behind the existing `needs_processor` states.
3. Add upload cancellation through an abortable frontend upload flow.
4. Add richer domain extraction for people, organizations, addresses, and procurement relationships.
5. Add a user-visible debug overlay for confidence explanations, not only the debug API.

## Honest Take

Civitas no longer feels like a toy on the Phase 2 fixture set. It does not magically solve OCR, PDF extraction, or transcription yet, but it now tells the truth about those gaps, preserves provenance, gives useful first guesses for structured civic data, and refuses to pretend broken evidence is complete. The next toy-like edge is scale: huge CSVs and native processor jobs need true streaming/progress/cancellation in Phase 3.
