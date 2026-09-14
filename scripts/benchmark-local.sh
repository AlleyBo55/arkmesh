#!/usr/bin/env bash
set -euo pipefail

ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
SIZE_MIB=64
RUNS=3
OUT=""

usage() {
  echo "usage: $0 [--size-mib N] [--runs N] [--out DIR]" >&2
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --size-mib)
      SIZE_MIB=$2
      shift 2
      ;;
    --runs)
      RUNS=$2
      shift 2
      ;;
    --out)
      OUT=$2
      shift 2
      ;;
    *)
      usage
      exit 2
      ;;
  esac
done

if ! [[ "$SIZE_MIB" =~ ^[1-9][0-9]*$ && "$RUNS" =~ ^[1-9][0-9]*$ ]]; then
  usage
  exit 2
fi

if [[ -z "$OUT" ]]; then
  OUT="$ROOT/benchmark-results/local-$(date -u +%Y%m%dT%H%M%SZ)"
fi
mkdir -p "$OUT"
OUT=$(CDPATH= cd -- "$OUT" && pwd)
WORK=$(mktemp -d "${TMPDIR:-/tmp}/arkmesh-benchmark.XXXXXX")
trap 'rm -rf "$WORK"' EXIT INT TERM

BIN="$WORK/arkmesh"
ASSET="$WORK/model.bin"
IDENTITY="$WORK/author"
CAPSULE="$WORK/benchmark.ark"
TREE="$WORK/tree.json"
SHARDS="$WORK/shards"
PLAN="$WORK/plan.json"
RAW="$OUT/raw.csv"
COMMAND_OUT="$WORK/command.out"
TIME_OUT="$WORK/time.out"

hash_file() {
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    sha256sum "$1" | awk '{print $1}'
  fi
}

directory_bytes() {
  find "$1" -type f -exec wc -c {} \; | awk '{sum += $1} END {print sum + 0}'
}

TIMEFORMAT='%R'
measure() {
  local operation=$1
  local iteration=$2
  shift 2
  : >"$TIME_OUT"
  { time "$@" >"$COMMAND_OUT"; } 2>"$TIME_OUT"
  local elapsed
  elapsed=$(tail -n 1 "$TIME_OUT")
  if ! [[ "$elapsed" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
    cat "$TIME_OUT" >&2
    echo "could not parse elapsed time for $operation" >&2
    exit 1
  fi
  printf '%s,%d,%s,%d\n' "$operation" "$iteration" "$elapsed" "$ASSET_BYTES" >>"$RAW"
}

pack_once() {
  rm -rf "$CAPSULE"
  "$BIN" pack --name benchmark --out "$CAPSULE" \
    --asset "model=$ASSET" \
    --signing-key "$IDENTITY/identity.key" \
    --chunk-size 65536
}

verify_once() {
  "$BIN" verify --trust "$IDENTITY/identity.json" --require-trusted "$CAPSULE"
}

copy_hash_once() {
  rm -f "$WORK/full-copy.bin"
  cp "$ASSET" "$WORK/full-copy.bin"
  hash_file "$WORK/full-copy.bin" >/dev/null
}

rsync_hash_once() {
  rm -f "$WORK/rsync-copy.bin"
  rsync -a --checksum "$ASSET" "$WORK/rsync-copy.bin"
  hash_file "$WORK/rsync-copy.bin" >/dev/null
}

protect_once() {
  rm -rf "$SHARDS"
  rm -f "$PLAN"
  "$BIN" erasure protect --asset "$DIGEST" --data 4 --parity 2 \
    --shards "$SHARDS" --plan "$PLAN" "$CAPSULE"
}

audit_once() {
  "$BIN" audit sample --tree "$TREE" --tolerance 0.01 --confidence 0.99 "$CAPSULE"
}

recover_once() {
  rm -f "$CAPSULE/objects/$DIGEST"
  "$BIN" erasure recover --plan "$PLAN" --shards "$SHARDS" "$CAPSULE"
}

printf 'Building ArkMesh and creating a %s MiB fixture...\n' "$SIZE_MIB"
(cd "$ROOT" && go build -o "$BIN" ./cmd/arkmesh)
dd if=/dev/zero of="$ASSET" bs=1048576 count="$SIZE_MIB" 2>/dev/null
printf 'ArkMesh benchmark fixture\n' | dd of="$ASSET" conv=notrunc 2>/dev/null
ASSET_BYTES=$(wc -c <"$ASSET" | tr -d ' ')
"$BIN" identity create --out "$IDENTITY" >/dev/null
printf 'operation,run,seconds,input_bytes\n' >"$RAW"

for ((run = 1; run <= RUNS; run++)); do
  measure full_copy_sha256 "$run" copy_hash_once
  if command -v rsync >/dev/null 2>&1; then
    measure rsync_sha256 "$run" rsync_hash_once
  fi
  measure arkmesh_pack_signed "$run" pack_once
done

DIGEST=$(sed -n 's/^[[:space:]]*"sha256": "\([0-9a-f]*\)",*$/\1/p' "$CAPSULE/manifest.json" | head -n 1)
if [[ -z "$DIGEST" ]]; then
  echo "could not read asset digest from manifest" >&2
  exit 1
fi
"$BIN" chunks tree --asset "$DIGEST" --out "$TREE" "$CAPSULE" >/dev/null

for ((run = 1; run <= RUNS; run++)); do
  measure arkmesh_verify_trusted "$run" verify_once
  measure arkmesh_erasure_protect_4_2 "$run" protect_once
  measure arkmesh_sample_audit_99pct "$run" audit_once
done

# Retain the final protected set, remove two shards, and repeatedly rebuild the object.
rm "$SHARDS/shard-001" "$SHARDS/shard-004"
for ((run = 1; run <= RUNS; run++)); do
  measure arkmesh_erasure_recover_2_loss "$run" recover_once
done

AUDIT_OUTPUT=$(audit_once)
SAMPLED=$(printf '%s\n' "$AUDIT_OUTPUT" | sed -n 's/^sampled: //p')
CHUNKS=$(printf '%s\n' "$AUDIT_OUTPUT" | sed -n 's/^chunks: //p')
CAPSULE_BYTES=$(directory_bytes "$CAPSULE")
SHARD_BYTES=$(directory_bytes "$SHARDS")
# Two shard files were removed for recovery. Derive encoded total from the plan's equal shard size.
SURVIVING_SHARDS=$(find "$SHARDS" -type f | wc -l | tr -d ' ')
ENCODED_BYTES=$(awk -v bytes="$SHARD_BYTES" -v surviving="$SURVIVING_SHARDS" 'BEGIN {printf "%.0f", bytes * 6 / surviving}')

UNAME_TEXT=$(uname -srm | tr '|' '/')
GO_TEXT=$(go version | tr '|' '/')
CPU_TEXT=$(sysctl -n machdep.cpu.brand_string 2>/dev/null || true)
if [[ -z "$CPU_TEXT" ]] && command -v lscpu >/dev/null 2>&1; then
  CPU_TEXT=$(lscpu | sed -n 's/^Model name:[[:space:]]*//p' | head -n 1)
fi
MEMORY_BYTES=$(sysctl -n hw.memsize 2>/dev/null || true)
COMMIT=$(git -C "$ROOT" rev-parse HEAD)
GENERATED=$(date -u +%Y-%m-%dT%H:%M:%SZ)

{
  echo "# ArkMesh local benchmark"
  echo
  echo "Generated: \`$GENERATED\`"
  echo
  echo "## Environment"
  echo
  echo "| Field | Value |"
  echo "|---|---|"
  echo "| Commit | \`$COMMIT\` |"
  echo "| Host | $UNAME_TEXT |"
  echo "| CPU | ${CPU_TEXT:-not reported} |"
  echo "| Memory bytes | ${MEMORY_BYTES:-not reported} |"
  echo "| Go | $GO_TEXT |"
  echo "| Input | $SIZE_MIB MiB ($ASSET_BYTES bytes) |"
  echo "| Runs per operation | $RUNS |"
  echo
  echo "## Timing"
  echo
  echo "Wall-clock seconds. Lower is better. Raw observations are in \`raw.csv\`."
  echo
  echo "| Operation | Mean | Minimum | Maximum |"
  echo "|---|---:|---:|---:|"
  awk -F, 'NR > 1 {
    sum[$1] += $3; count[$1]++
    if (!($1 in min) || $3 < min[$1]) min[$1] = $3
    if (!($1 in max) || $3 > max[$1]) max[$1] = $3
  }
  END {
    for (name in sum) printf "| %s | %.3f | %.3f | %.3f |\n", name, sum[name] / count[name], min[name], max[name]
  }' "$RAW" | sort
  echo
  echo "## Storage and audit observations"
  echo
  echo "| Metric | Value |"
  echo "|---|---:|"
  echo "| Source bytes | $ASSET_BYTES |"
  echo "| Complete capsule bytes | $CAPSULE_BYTES |"
  echo "| Encoded 4 + 2 shard bytes | $ENCODED_BYTES |"
  awk -v encoded="$ENCODED_BYTES" -v source="$ASSET_BYTES" 'BEGIN {printf "| Erasure storage overhead | %.2f%% |\n", (encoded / source - 1) * 100}'
  echo "| Audit chunks sampled | $SAMPLED of $CHUNKS |"
  echo
  echo "## Scope"
  echo
  echo "The ordinary-copy and rsync rows are local byte-preservation baselines with a full SHA-256 read. ArkMesh timings additionally cover signed manifests, authenticated chunks, trust verification, erasure protection, reconstruction, and sampled retrievability. This run does not compare peer transport, BitTorrent, or IPFS because ArkMesh networking is not implemented."
} >"$OUT/results.md"

printf 'Benchmark complete.\nRaw data: %s\nReport: %s\n' "$RAW" "$OUT/results.md"
cat "$OUT/results.md"
