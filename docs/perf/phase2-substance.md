# Phase 2 Substance Performance Notes

Date: 2026-05-09

Command: `CGO_ENABLED=0 go test ./internal/pipeline -run 'TestRealDataFixtures' -count=1 -v`

## Fixture Timing

| Fixture | Duration |
| --- | --- |
| BLM FOIA sample PDF | <10 ms |
| Truncated BLM FOIA PDF | <10 ms |
| Ann Arbor council audio OGG, 70 MB | about 70-110 ms |
| DOJ bid-rigging HTML | <10 ms |
| Empty failed-transfer PDF | <10 ms |
| ICIJ database HTML | <10 ms |
| ICIJ partial ZIP | <10 ms |
| NARA scan JPG | <10 ms |
| OCDS tender JSON | <10 ms |
| OFAC SDN CSV, 5.5 MB | about 1.7-2.0 seconds |

Median fixture duration is below 10 ms. The p95 and worst case are the OFAC CSV, which dominates because Civitas parses all records to count rows and infer a stable field model.

## Hot Paths Fixed

1. Raw CSV entity regex scanning was removed from the main path. Entity extraction now uses structured CSV preview and inferred field values, which keeps the graph from filling with regex noise.
2. Large evidence text passed to entity extraction is capped at a deterministic 200 KB window.
3. Media, image, valid PDF, and corrupt archive inputs enter explicit states quickly instead of doing fake synchronous processing.

## Remaining Cliff

CSV row counting is still synchronous. It is acceptable for the current 5.5 MB audit fixture, but a future Phase 3 should move CSV parsing into a worker/job path and stream row counts for files much larger than this fixture.
