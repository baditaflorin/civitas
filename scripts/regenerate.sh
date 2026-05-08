#!/usr/bin/env bash
set -euo pipefail

if [ -d node_modules ]; then
  npm run gen:api
else
  echo "node_modules missing; skipping API type generation"
fi
