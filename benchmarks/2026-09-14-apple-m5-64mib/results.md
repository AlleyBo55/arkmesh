# ArkMesh local benchmark

Generated: `2026-09-14T06:42:21Z`

## Environment

| Field | Value |
|---|---|
| Commit | `057ee8d66f52264f2e5615afb883e907d78b4f85` |
| Host | Darwin 25.2.0 arm64 |
| CPU | Apple M5 |
| Memory bytes | 17179869184 |
| Go | go version go1.22.0 darwin/arm64 |
| Input | 64 MiB (67108864 bytes) |
| Runs per operation | 3 |

The host name was removed from the published environment because it is irrelevant to reproduction.

## Timing

Wall-clock seconds. Lower is better. Raw observations are in [`raw.csv`](raw.csv).

| Operation | Mean | Minimum | Maximum |
|---|---:|---:|---:|
| arkmesh_erasure_protect_4_2 | 0.211 | 0.205 | 0.217 |
| arkmesh_erasure_recover_2_loss | 0.186 | 0.184 | 0.191 |
| arkmesh_pack_signed | 0.125 | 0.125 | 0.126 |
| arkmesh_sample_audit_99pct | 0.013 | 0.012 | 0.014 |
| arkmesh_verify_trusted | 0.047 | 0.045 | 0.048 |
| full_copy_sha256 | 0.140 | 0.136 | 0.144 |
| rsync_sha256 | 0.358 | 0.356 | 0.362 |

## Storage and audit observations

| Metric | Value |
|---|---:|
| Source bytes | 67108864 |
| Complete capsule bytes | 67109876 |
| Encoded four plus two shard bytes | 100663296 |
| Erasure storage overhead | 50.00% |
| Audit chunks sampled | 349 of 1024 |

## Scope

The ordinary-copy and rsync rows are local byte-preservation baselines with a full SHA-256 read. ArkMesh timings additionally cover signed manifests, authenticated chunks, trust verification, erasure protection, reconstruction, and sampled retrievability. This run does not compare peer transport, BitTorrent, or IPFS because ArkMesh networking is not implemented.
