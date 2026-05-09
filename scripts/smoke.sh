#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

./scripts/build-pages.sh >/tmp/civitas-build.log

mkdir -p tmp
BACKEND_PORT="${CIVITAS_SMOKE_BACKEND_PORT:-18089}"
BACKEND_URL="http://127.0.0.1:${BACKEND_PORT}"
PAGES_PORT="${CIVITAS_SMOKE_PAGES_PORT:-4174}"
while lsof -iTCP:"$PAGES_PORT" -sTCP:LISTEN >/dev/null 2>&1; do
  PAGES_PORT=$((PAGES_PORT + 1))
done
PAGES_URL="http://127.0.0.1:${PAGES_PORT}/civitas/"
STORAGE_DIR="$(mktemp -d)"
CGO_ENABLED=0 CIVITAS_ENV=smoke CIVITAS_ADDR=":${BACKEND_PORT}" CIVITAS_STORAGE_DIR="$STORAGE_DIR" CIVITAS_ALLOWED_ORIGINS="http://127.0.0.1:${PAGES_PORT}" go run ./cmd/server >tmp/backend-smoke.log 2>&1 &
BACKEND_PID=$!
npm run preview -- --host 127.0.0.1 --port "$PAGES_PORT" >tmp/pages-preview.log 2>&1 &
SERVER_PID=$!
trap 'kill "$SERVER_PID" "$BACKEND_PID" >/dev/null 2>&1 || true; rm -rf "$STORAGE_DIR"' EXIT

for _ in {1..40}; do
  if curl -fsS "$PAGES_URL" >/dev/null 2>&1; then
    break
  fi
  sleep 0.25
done

curl -fsS "$PAGES_URL" >/dev/null
for _ in {1..40}; do
  if curl -fsS "$BACKEND_URL/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 0.25
done

curl -fsS "$BACKEND_URL/healthz" >/dev/null
PLAYWRIGHT_BASE_URL="$PAGES_URL" PLAYWRIGHT_API_BASE_URL="$BACKEND_URL" npx playwright test
