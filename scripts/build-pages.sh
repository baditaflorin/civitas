#!/usr/bin/env bash
set -euo pipefail

rm -rf docs/assets
rm -f docs/index.html docs/404.html docs/manifest.webmanifest

export VITE_COMMIT_SHA="${VITE_COMMIT_SHA:-static}"
export VITE_APP_VERSION="${VITE_APP_VERSION:-$(node -p "require('./package.json').version")}"

npx vite build
