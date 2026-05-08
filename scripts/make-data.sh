#!/usr/bin/env bash
set -euo pipefail

mkdir -p docs/data/v1
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo dev)"
GENERATED_AT="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

cat >docs/data/v1/sample.json <<'JSON'
{
  "case": "demo-civic-dump",
  "documents": [
    {
      "filename": "contract-note.txt",
      "entities": ["source@example.org", "42 Civic Street"],
      "timeline": ["2026-05-08"]
    }
  ]
}
JSON

cat >docs/data/v1/sample.meta.json <<JSON
{
  "generated_at": "$GENERATED_AT",
  "source_commit": "$COMMIT",
  "schema_version": "v1",
  "input_checksums": []
}
JSON
