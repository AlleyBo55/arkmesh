#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
WORK=${ARKMESH_DEMO_DIR:-}
KEEP=0

if [ -z "$WORK" ]; then
  WORK=$(mktemp -d "${TMPDIR:-/tmp}/arkmesh-demo.XXXXXX")
else
  mkdir -p "$WORK"
  if [ -n "$(find "$WORK" -mindepth 1 -maxdepth 1 -print -quit)" ]; then
    echo "ARKMESH_DEMO_DIR must be empty: $WORK" >&2
    exit 1
  fi
  KEEP=1
fi

cleanup() {
  if [ "$KEEP" -eq 0 ]; then
    rm -rf "$WORK"
  fi
}
trap cleanup EXIT INT TERM

BIN="$WORK/arkmesh"
ASSET="$WORK/model.bin"
CAPSULE="$WORK/rescue.ark"
IDENTITY="$WORK/author"
SHARDS="$WORK/shards"
PLAN="$WORK/plan.json"
THRESHOLD_CAPSULE="$WORK/threshold.ark"
THRESHOLD_SHARDS="$WORK/threshold-shards"

hash_file() {
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    sha256sum "$1" | awk '{print $1}'
  fi
}

step() {
  printf '\n==> %s\n' "$1"
}

step "Build the local ArkMesh binary"
(cd "$ROOT" && go build -o "$BIN" ./cmd/arkmesh)

step "Create an 8 MiB fixture and an offline signing identity"
dd if=/dev/zero of="$ASSET" bs=1048576 count=8 2>/dev/null
printf 'ArkMesh offline recovery demonstration\n' | dd of="$ASSET" conv=notrunc 2>/dev/null
ORIGINAL_HASH=$(hash_file "$ASSET")
"$BIN" identity create --out "$IDENTITY"

step "Pack and sign a capsule with authenticated 256 KiB chunks"
"$BIN" pack \
  --name "offline-recovery-demo" \
  --out "$CAPSULE" \
  --asset "model=$ASSET" \
  --signing-key "$IDENTITY/identity.key" \
  --chunk-size 262144
DIGEST=$(sed -n 's/^[[:space:]]*"sha256": "\([0-9a-f]*\)",*$/\1/p' "$CAPSULE/manifest.json" | head -n 1)
if [ -z "$DIGEST" ]; then
  echo "could not read asset digest from manifest" >&2
  exit 1
fi
"$BIN" verify --trust "$IDENTITY/identity.json" --require-trusted "$CAPSULE"

step "Create four data shards and two parity shards"
"$BIN" erasure protect \
  --asset "$DIGEST" \
  --data 4 \
  --parity 2 \
  --shards "$SHARDS" \
  --plan "$PLAN" \
  "$CAPSULE"

step "Prove that losses beyond parity are refused without writing an object"
cp -R "$CAPSULE" "$THRESHOLD_CAPSULE"
cp -R "$SHARDS" "$THRESHOLD_SHARDS"
rm "$THRESHOLD_CAPSULE/objects/$DIGEST"
rm "$THRESHOLD_SHARDS/shard-001" "$THRESHOLD_SHARDS/shard-003" "$THRESHOLD_SHARDS/shard-005"
if "$BIN" erasure recover --plan "$PLAN" --shards "$THRESHOLD_SHARDS" "$THRESHOLD_CAPSULE"; then
  echo "unexpected recovery with only three surviving shards" >&2
  exit 1
fi
if [ -e "$THRESHOLD_CAPSULE/objects/$DIGEST" ]; then
  echo "failed recovery wrote an object" >&2
  exit 1
fi
echo "expected refusal observed; no object was written"

step "Destroy the publisher source, capsule object, and two shards"
rm "$ASSET" "$CAPSULE/objects/$DIGEST"
rm "$SHARDS/shard-001" "$SHARDS/shard-004"
test ! -e "$CAPSULE/objects/$DIGEST"
echo "the original object no longer exists; four of six shards remain"

step "Recover with no full-object donor"
"$BIN" erasure recover --plan "$PLAN" --shards "$SHARDS" "$CAPSULE"

step "Verify exact bytes, signature, and imported author trust offline"
RECOVERED_HASH=$(hash_file "$CAPSULE/objects/$DIGEST")
if [ "$RECOVERED_HASH" != "$ORIGINAL_HASH" ] || [ "$RECOVERED_HASH" != "$DIGEST" ]; then
  echo "recovered hash does not match the signed manifest" >&2
  exit 1
fi
"$BIN" verify --trust "$IDENTITY/identity.json" --require-trusted "$CAPSULE"

printf '\nPASS: donorless recovery rebuilt the exact signed object.\n'
printf 'original sha256:  %s\n' "$ORIGINAL_HASH"
printf 'recovered sha256: %s\n' "$RECOVERED_HASH"
printf 'ArkMesh peer/network operations: not implemented\n'
if [ "$KEEP" -eq 1 ]; then
  printf 'evidence retained at: %s\n' "$WORK"
fi
