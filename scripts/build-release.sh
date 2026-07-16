#!/usr/bin/env bash
# Build multi-platform zhconv CLI artifacts into ./dist
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${1:-dev}"
VERSION="${VERSION#v}"
DIST_DIR="${ROOT_DIR}/dist"
APP_NAME="zhconv"
MODULE="github.com/xylplm/zhconv-go/cmd/zhconv"

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

LDFLAGS="-s -w -X main.version=${VERSION}"

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
  "windows arm64"
)

echo "Building ${APP_NAME} ${VERSION}"

for target in "${targets[@]}"; do
  # shellcheck disable=SC2086
  set -- ${target}
  GOOS="$1"
  GOARCH="$2"

  ext=""
  archive_ext="tar.gz"
  if [[ "${GOOS}" == "windows" ]]; then
    ext=".exe"
    archive_ext="zip"
  fi

  out_dir="${DIST_DIR}/build_${GOOS}_${GOARCH}"
  mkdir -p "${out_dir}"
  bin_path="${out_dir}/${APP_NAME}${ext}"

  echo "-> ${GOOS}/${GOARCH}"
  env CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
    go build -trimpath -ldflags "${LDFLAGS}" -o "${bin_path}" "${MODULE}"

  asset_base="${APP_NAME}_${VERSION}_${GOOS}_${GOARCH}"
  if [[ "${archive_ext}" == "zip" ]]; then
    if command -v zip >/dev/null 2>&1; then
      (
        cd "${out_dir}"
        zip -q "${DIST_DIR}/${asset_base}.zip" "${APP_NAME}${ext}"
      )
    else
      # Git Bash on Windows may not ship zip(1); fall back to Python.
      python - "${out_dir}/${APP_NAME}${ext}" "${DIST_DIR}/${asset_base}.zip" <<'PY'
import sys, zipfile
from pathlib import Path
src, dst = Path(sys.argv[1]), Path(sys.argv[2])
with zipfile.ZipFile(dst, "w", compression=zipfile.ZIP_DEFLATED) as zf:
    zf.write(src, arcname=src.name)
PY
    fi
  else
    tar -C "${out_dir}" -czf "${DIST_DIR}/${asset_base}.tar.gz" "${APP_NAME}${ext}"
  fi

  rm -rf "${out_dir}"
done

echo "Artifacts:"
ls -lh "${DIST_DIR}"
