#!/usr/bin/env bash
# Rebuild embedded OpenCC-derived dictionaries and print a short summary.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

REF="${OPENCC_REF:-${1:-master}}"
export OPENCC_REF="$REF"

echo "==> regenerating dictionaries from OpenCC@${REF}"
go run ./scripts/gendict -ref "$REF" -out-root "$ROOT_DIR"

echo "==> verifying tables load"
go test ./table ./... -count=1

if command -v git >/dev/null 2>&1; then
  echo "==> git status (dict/table)"
  git status --short -- dict table || true
fi

echo "done"