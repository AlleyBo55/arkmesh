# Benchmarks

Status: initial local measurements, not a peer-network comparison.

ArkMesh publishes raw observations and the command that generated them. Results should be treated as measurements of one machine and commit, not universal performance claims.

## Reproduce

```bash
./scripts/benchmark-local.sh --size-mib 64 --runs 3 --out ./benchmark-results/local
```

The script builds ArkMesh, creates a deterministic fixture, and records wall-clock timings for:

- Ordinary full copy followed by a complete SHA-256 read
- Local `rsync --checksum` followed by a complete SHA-256 read, when rsync is installed
- Signed capsule packing with 64 KiB authenticated chunks
- Trusted capsule verification
- Four data plus two parity shard protection
- Recovery after losing two shards and the complete object
- A sampled audit configured for 99 percent confidence and one percent tolerated loss

It writes `raw.csv` and a generated `results.md`. The working data is removed after the report is written.

## Published run: Apple M5, 64 MiB

Raw observations and the complete environment record are in [`benchmarks/2026-09-14-apple-m5-64mib`](../benchmarks/2026-09-14-apple-m5-64mib/results.md).

Commit tested: `057ee8d66f52264f2e5615afb883e907d78b4f85`

| Operation | Mean seconds | Minimum | Maximum |
|---|---:|---:|---:|
| Ordinary full copy plus SHA-256 | 0.140 | 0.136 | 0.144 |
| rsync plus SHA-256 | 0.358 | 0.356 | 0.362 |
| ArkMesh signed pack | 0.125 | 0.125 | 0.126 |
| ArkMesh trusted verify | 0.047 | 0.045 | 0.048 |
| ArkMesh four plus two protection | 0.211 | 0.205 | 0.217 |
| ArkMesh recovery after two shard losses | 0.186 | 0.184 | 0.191 |
| ArkMesh 99 percent sampled audit | 0.013 | 0.012 | 0.014 |

The encoded shard set used 100,663,296 bytes for a 67,108,864 byte object, exactly 50 percent parity overhead. The audit read 349 of 1,024 authenticated chunks.

## Interpretation

The rows perform different work and are not direct feature-equivalent winners. A full copy preserves one additional complete object. ArkMesh pack additionally creates authenticated chunk commitments and a signed manifest. Erasure protection creates recoverability without any surviving full copy. Sampled audit provides a statistical retrievability statement rather than complete verification.

The current run is too small and too local to support broad performance conclusions. It establishes a reproducible measurement floor and exposes raw data for correction.

## Not measured yet

BitTorrent and IPFS transfer are not reported because ArkMesh peer transport is not implemented. Comparing their network transfer to ArkMesh local filesystem operations would be misleading. Future peer experiments must use the same topology, bytes, failure schedule, cache state, and verification endpoint for every system.

Also not measured yet:

- Multiple object sizes and hardware classes
- Peak memory and CPU utilization
- Physical disks or independently administered devices
- WAN and LAN transfer rates
- Model execution and inference health
- Repeated failure curves with confidence intervals
