#!/bin/sh
set -eu

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
  echo "usage: $0 VERSION [OUT_DIR]" >&2
  exit 2
fi

VERSION=$1
OUT=${2:-dist}
ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)

case "$VERSION" in
  v[0-9]*.[0-9]*.[0-9]*-alpha.[0-9]*) ;;
  *)
    echo "version must look like v0.1.0-alpha.1" >&2
    exit 2
    ;;
esac

case "$OUT" in
  /*) ;;
  *) OUT="$ROOT/$OUT" ;;
esac

if [ -d "$OUT" ] && [ -n "$(find "$OUT" -mindepth 1 -maxdepth 1 -print -quit)" ]; then
  echo "output directory must be empty: $OUT" >&2
  exit 1
fi
if [ -e "$OUT" ] && [ ! -d "$OUT" ]; then
  echo "output path is not a directory: $OUT" >&2
  exit 1
fi
mkdir -p "$OUT"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/arkmesh-release.XXXXXX")
trap 'rm -rf "$WORK"' EXIT INT TERM

for target in darwin/arm64 darwin/amd64 linux/arm64 linux/amd64; do
  GOOS=${target%/*}
  GOARCH=${target#*/}
  NAME="arkmesh_${VERSION}_${GOOS}_${GOARCH}"
  STAGE="$WORK/$NAME"
  mkdir -p "$STAGE"
  (
    cd "$ROOT"
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
      go build -trimpath -ldflags "-s -w" -o "$STAGE/arkmesh" ./cmd/arkmesh
  )
  cp "$ROOT/README.md" "$ROOT/LICENSE" "$STAGE/"
  COPYFILE_DISABLE=1 tar -C "$WORK" -czf "$OUT/$NAME.tar.gz" "$NAME"
done

(
  cd "$OUT"
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 ./*.tar.gz >SHA256SUMS
  else
    sha256sum ./*.tar.gz >SHA256SUMS
  fi
)

printf 'Release artifacts created in %s\n' "$OUT"
cat "$OUT/SHA256SUMS"
