#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TMP="$(mktemp -d)"

cleanup() {
  trash "$TMP" >/dev/null 2>&1 || true
}
trap cleanup EXIT

require_tool() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "SKIP: $1 not installed"
    exit 0
  fi
}

require_tool typst
require_tool pdf-to-png

cd "$ROOT"

go run ./cmd/folio-manuscript --style british examples/dummy-manuscript.md "$TMP/british.pdf"
go run ./cmd/folio-manuscript --style us examples/dummy-manuscript.md "$TMP/us.pdf"

(cd "$TMP" && pdf-to-png british.pdf 120 >/dev/null)
(cd "$TMP" && pdf-to-png us.pdf 120 >/dev/null)

for image in "$TMP"/british-*.png "$TMP"/us-*.png; do
  if [[ ! -s "$image" ]]; then
    echo "FAIL: missing raster output $image"
    exit 1
  fi
done

echo "PASS: issue #9 manuscript PDFs rasterize with pdf-to-png"
