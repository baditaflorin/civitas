# Phase 2 Substance Real-Data Audit

Date: 2026-05-08

Mode: Mode C, same as Phase 1.

V1 path tested: create case -> upload document -> list documents -> search `contract` -> graph -> timeline -> markdown export.

Technical result: 10/10 uploads completed without crashing.

Useful-result pass rate: 0/10. Every input either produced raw unstructured text, a generic placeholder, or a completed status that hid important uncertainty.

## Inputs

| ID | Input | Source | Shape | Why it belongs |
| --- | --- | --- | --- | --- |
| R01 | OFAC SDN CSV | https://www.treasury.gov/ofac/downloads/sdn.csv | 5.3 MB CSV | Real sanctions list, large-ish structured civic data. |
| R02 | OCDS tender JSON | https://raw.githubusercontent.com/open-contracting/sample-data/main/fictional-example/1.1/ocds-213czf-000-00001-02-tender.json | Nested JSON | Real domain standard for procurement data. |
| R03 | DOJ bid-rigging press release | https://www.justice.gov/opa/pr/executive-pleads-guilty-multi-million-dollar-bid-rigging-conspiracy | HTML article | Real investigation context with people, money, contracts, dates. |
| R04 | ICIJ Offshore Leaks download page | https://offshoreleaks.icij.org/pages/database | HTML page | Real investigative data source with disclaimers, address, graph vocabulary. |
| R05 | BLM sample FOIA PDF | https://www.blm.gov/sites/blm.gov/files/foia_samplefoiarequestletter.pdf | PDF | Common public-records document format. |
| R06 | NARA scanned document image | https://commons.wikimedia.org/wiki/Special:Redirect/file/Py1_50.jpg | JPEG scan | Real scanned government record needing OCR. |
| R07 | Ann Arbor council audio | https://commons.wikimedia.org/wiki/Special:Redirect/file/A2Council-Caucus-2020-10-04.ogg | 70 MB OGG audio | Real public-meeting media, large enough to expose performance honesty gaps. |
| R08 | ICIJ Offshore Leaks ZIP partial | https://offshoreleaks-data.icij.org/offshoreleaks/csv/full-oldb.LATEST.zip | First 1 MB of a ZIP | Huge archive / interrupted transfer case. |
| R09 | Truncated BLM FOIA PDF | Derived from R05, first 4 KB | Broken PDF | Partial/corrupt upload case. |
| R10 | Empty failed transfer PDF | Zero-byte file named `failed_transfer_empty.pdf` | Empty file | Failed transfer edge case users really produce. |

## Walkthrough Results

| ID | What v1 did | What it should have done | Why it failed | Failure type | Manual work forced onto user |
| --- | --- | --- | --- | --- | --- |
| R01 | Uploaded as `text/plain`, stored 5.5 MB of raw CSV text, extracted 12,118 regex entities, generated a huge graph, no timeline. | Detect CSV dialect, infer columns, identify names/countries/program fields, warn about large list, avoid turning every regex hit into a graph node. | No structured parser, no schema inference, no field typing, no row-level confidence. | Wrong-but-confident. | User must know it is CSV, understand the columns, decide which fields matter, and mentally ignore graph noise. |
| R02 | Uploaded as `text/plain`, preserved raw JSON text, extracted 9 URL-like entities, no procurement model. | Detect OCDS JSON, summarize tender title/value/buyer/suppliers, preserve OCID, expose linked documents and dates. | No JSON parsing, no domain-shape detection, no nested-field extraction. | Silent under-reading. | User must manually read JSON and translate it into investigation concepts. |
| R03 | Uploaded raw HTML, summary starts with markup, extracted 18 regex entities, no dates despite visible human dates. | Strip boilerplate, extract article title/date/people/orgs/money/locations, infer procurement-fraud article shape. | No HTML readability extraction, no human-date parsing, no money phrase normalization. | Wrong-but-confident. | User must mentally remove markup and spot the real article facts. |
| R04 | Uploaded raw HTML, extracted 4 entities, treated disclaimer/page chrome as evidence text. | Identify this as a data-source page, capture dataset facts, download/archive links, disclaimer, address, and graph vocabulary. | No page-type recognition or semantic extraction. | Silent under-reading. | User must identify the actual download links and data caveats. |
| R05 | Uploaded PDF and marked `completed`, but text is only a generic "configure Tika/PyMuPDF/pdfminer" placeholder. | Extract PDF text, detect it as FOIA request template, identify agency/contact/fee-waiver language, cite extraction failure if tools missing. | PDF adapter is a placeholder; missing processor availability is not promoted to document status. | Obvious but still marked complete. | User must install tools and rerun, with no guided next step. |
| R06 | Uploaded JPEG and marked `completed`, but text is only a generic OCR placeholder. | OCR scan, report OCR confidence, preserve image metadata, flag if OCR unavailable. | OCR adapter is a placeholder; image evidence is not treated as needing OCR. | Obvious but still marked complete. | User must know OCR is required and which tool to install. |
| R07 | Uploaded 70 MB OGG as `application/ogg`, returned "Binary evidence uploaded", no media classification, no progress or cancellation. | Detect audio, estimate duration, queue transcription/diarization, show progress, make job cancellable, mark transcription unavailable if no Whisper/pyannote. | Media MIME handling is too generic; upload is synchronous; no job state machine. | Silent under-reading plus performance dishonesty. | User must know audio was not transcribed and wait without progress. |
| R08 | Uploaded partial ZIP and marked `completed`, said binary evidence uploaded. | Detect ZIP archive, validate central directory, report truncated archive, suggest re-upload or partial recovery. | No archive detection/validation; corrupt archive is not distinguished from valid binary. | Wrong-but-confident. | User must discover later that the archive is unusable. |
| R09 | Uploaded truncated PDF and gave the same message as a normal PDF. | Detect malformed/truncated PDF, try qpdf repair, report recoverability, avoid same status as valid PDF. | No PDF validation; no distinction between unsupported and corrupt. | Wrong-but-confident. | User must know the source file is broken and re-upload. |
| R10 | Uploaded zero-byte PDF as `text/plain`, marked complete with "contains no extracted text yet." | Treat zero bytes as failed transfer, block "completed", ask user to re-upload or keep as failed evidence with provenance. | Empty input is not a first-class state. | Wrong-but-confident. | User must infer that the upload itself failed. |

## Top 5 Logic Gaps

1. V1 does not parse structured formats. CSV and JSON are treated as plain text, so the app misses rows, columns, nested procurement fields, IDs, values, buyers, suppliers, and schema cues.
2. V1 does not distinguish unsupported, unavailable, corrupt, partial, and empty inputs. PDFs, truncated PDFs, images, ZIPs, audio, and zero-byte uploads all become `completed` documents.
3. V1 extracts regex hits without domain meaning. A sanctions CSV creates thousands of graph nodes, while procurement-specific facts like award value, buyer, supplier, conspiracy period, and OCID are not recognized.
4. V1 has no confidence model. It cannot tell the user "high-confidence text extraction" versus "placeholder only" versus "file appears broken."
5. V1 has no job/state model for heavy inputs. A 70 MB audio upload is synchronous, has no progress, no cancellation, no duration estimate, and no transcription status.

## Top 3 Intuition Failures

1. `completed` means "the upload did not crash," not "evidence was understood." Users will assume completed documents are processed.
2. The graph looks authoritative even when it is mostly regex noise from raw CSV/HTML.
3. The export exists even when source documents contain only placeholders, so the export can look like a deliverable while carrying almost no evidence.

## Top 3 Feels-Stupid Moments

1. User uploads a CSV and the app does not infer rows, headers, or field types.
2. User uploads a PDF/image/audio file and the app says, in effect, "install tools," instead of making the processing state and next step explicit.
3. User uploads a broken or empty file and the app does not say "this upload is broken."

## What Smart Means For Civitas

Smart Civitas means:

1. On upload, Civitas classifies the evidence shape: CSV, JSON/OCDS, HTML article, PDF, scanned image, audio/video, archive, empty, corrupt, or unknown.
2. Civitas gives a useful first guess: structured rows for CSV, key procurement fields for OCDS, clean article text for HTML, extracted text/OCR/transcript status for document media.
3. Every inference carries confidence and provenance, and low confidence is visible in the UI and export.
4. Broken, partial, empty, or unsupported evidence gets a domain-specific status with a next step, not a fake completion.
5. Large/heavy work becomes a tracked job with progress, cancellation, deterministic output metadata, and no UI freeze.

## Phase 2 Substance Success Metrics

1. Useful-result pass rate: at least 7/10 audit inputs produce a useful first guess with no manual configuration.
2. No silent wrongness: 10/10 inputs expose confidence and processing state; corrupt/empty/unsupported inputs are never marked as cleanly completed.
3. Determinism: running the same 10 inputs twice produces byte-identical normalized outputs, excluding explicitly recorded run timestamps.
4. Format detection: 10/10 inputs are classified into the correct evidence shape.
5. Heavy input honesty: inputs over 5 MB show progress; operations over 5 seconds are cancellable.
6. Export provenance: 10/10 exports include source ID, schema version, app version, processor decisions, confidence summary, and parameters used.
7. Performance: median time from upload completion to first useful preview is under 1 second for text/CSV/JSON/HTML fixtures; heavy media enters a queued state within 300 ms.

## Out Of Scope For Phase 2 Substance

- No new user-facing feature areas.
- No UI polish, redesign, dark mode, command palette, marketing, or visual chrome.
- No architecture mode change; remain Mode C.
- No auth system.
- No multi-user collaboration.
- No full replacement of every native processor; Phase 2 may improve adapters and statuses, but the target is smarter behavior under the existing surface area.
- No hosted SaaS workflow.
- No new export formats beyond making existing exports honest, deterministic, and provenance-rich.
