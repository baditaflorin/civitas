#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

./scripts/build-pages.sh >/tmp/civitas-build.log

mkdir -p tmp
npm run preview -- --host 127.0.0.1 --port 4174 >tmp/pages-preview.log 2>&1 &
SERVER_PID=$!
trap 'kill "$SERVER_PID" >/dev/null 2>&1 || true' EXIT

for _ in {1..40}; do
  if curl -fsS http://127.0.0.1:4174/civitas/ >/dev/null 2>&1; then
    break
  fi
  sleep 0.25
done

curl -fsS http://127.0.0.1:4174/civitas/ >/dev/null
PLAYWRIGHT_BASE_URL=http://127.0.0.1:4174/civitas/ npx playwright test
