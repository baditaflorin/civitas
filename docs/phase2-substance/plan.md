# Phase 2 Substance Plan

Ranked by user impact on the 10 real-data audit inputs.

## Picklist

1. Classify evidence shape on upload: CSV, JSON/OCDS, HTML article/source page, PDF, image scan, audio, archive, empty, corrupt, unknown. Covers 6, 8, 13, 24.
2. Introduce explicit processing states: ready, needs_processor, recoverable_error, failed, unsupported. Covers 24, 25, 34.
3. Add confidence to document classification, fields, anomalies, and export summaries. Covers 16, 44.
4. Parse CSV with delimiter/header inference, row counts, field type inference, and anomaly reporting. Covers 1, 2, 3, 5, 6, 7, 9, 15, 18.
5. Parse JSON and recognize OCDS release packages. Covers 6, 7, 8, 12, 13, 14.
6. Extract readable HTML article/source text with title, times, links, and page-shape detection. Covers 6, 8, 11, 13, 15.
7. Detect empty files as failed transfers. Covers 4, 25, 32, 34.
8. Validate PDFs enough to distinguish likely valid, truncated, and malformed before processor fallback. Covers 4, 12, 17, 18, 32.
9. Detect images as OCR-needed evidence with explicit processor status and next step. Covers 13, 16, 32.
10. Detect audio/video as transcription-needed evidence with size/duration-ish metadata and job-state language. Covers 3, 13, 24, 28.
11. Validate ZIP archives and distinguish valid archive from partial/corrupt archive. Covers 4, 5, 17, 18.
12. Normalize dates in ISO and human-readable formats into timeline events. Covers 9, 12, 15.
13. Normalize whitespace and strip HTML/markup noise by default. Covers 2, 9, 15.
14. Add domain field inference for emails, URLs, money, addresses, procurement IDs, buyer/supplier/title/value/date. Covers 7, 11, 12.
15. Cap and rank entity extraction so graph output is useful instead of noisy. Covers 16, 18, 35.
16. Add anomaly suggestions with what/why/now-what language. Covers 17, 18, 19, 32, 47.
17. Add deterministic stable IDs for documents, entities, fields, anomalies, timeline events. Covers 22, 35.
18. Make exports provenance-rich: schema version, app version, source IDs, processor decisions, confidence summary, parameters. Covers 14, 35, 38.
19. Add debug API surface for internal state and inference explanations. Covers 19, 37.
20. Add real fixture suite asserting the 10 audit inputs and worst edge cases do not regress. Covers 1, 5, 35.
21. Add performance measurement hooks for parse duration and size buckets. Covers 28, 31.
22. Add boundary validation for upload size, filename, and content-type mismatch. Covers 32, 33, 34.

## Implementation Order

1. Fixtures and expected contracts.
2. Evidence model additions.
3. Classifier and normalizer.
4. Structured parsers by shape.
5. State/confidence/anomaly/provenance surfacing.
6. Deterministic exports and debug endpoint.
7. Fixture suite, postmortem, version bump.
