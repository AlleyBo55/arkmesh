# Changelog

All notable ArkMesh milestones will be recorded here.

## v0.1.0-alpha.1

First reproducible research release.

### Implemented

- Deterministic content-addressed capsules and offline object verification
- Ed25519 authorship with explicit local trust
- Signed ancestry, exact key rotation, local revocation, and anti-rollback checkpoints
- Multi-party threshold emergency recovery
- Merkle chunk commitments, possession proofs, damage localization, and verified repair
- Streaming four data plus two parity Reed Solomon protection and donorless reconstruction
- Sampled retrievability audits with exact hypergeometric bounds and local retention history
- Strict JSON parsing, race tests, bounded fuzzing, and standard-library-only Go implementation

### Reproducible evidence

- Offline recovery ceremony in `scripts/demo-offline-recovery.sh`
- Local benchmark in `scripts/benchmark-local.sh`
- Published raw benchmark observations under `benchmarks/2026-09-14-apple-m5-64mib`
- Reproducible release archives from `scripts/build-release.sh`

### Known limits

- No peer discovery or transfer
- No local model inference or executable health check
- No automated shard placement or repair scheduling
- No independent time attestation for retention logs
- No divergent-branch reconciliation
- No physical multi-device recovery result
- No independent security review or second verifier

Do not use this alpha release for disaster-critical work.
