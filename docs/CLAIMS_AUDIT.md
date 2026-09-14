# Public Claims Audit

Date: 2026-09-14

## Conclusion

ArkMesh cannot be certified as "100% factual" in an absolute sense. Software can contain undiscovered defects, measurements are environment-specific, legal authorization depends on facts outside the repository, and the distributed continuity thesis has not been tested.

The narrower result of this audit is:

- Current local implementation claims are supported by source, tests, and reproducible commands.
- Published benchmark numbers match the recorded raw observations.
- Networking, peer transfer, placement, attested time, inference, reconciliation, and the independent field experiment are labeled planned or unproven.
- Mission language is presented as normative or hypothetical rather than evidence.
- Concrete overstatements found during the audit were corrected.

This is an internal evidence audit, not independent peer review, security certification, legal review, or proof that no misleading statement remains.

## Claim classes

### Implemented

An `implemented` or `verified` label means the named acceptance behavior exists in the current Go reference implementation and has direct automated evidence. It does not mean safe to execute, independently audited, production-ready, legally authorized, or disaster-ready.

| Claim area | Evidence |
|---|---|
| Deterministic capsule identity | `internal/capsule/manifest_test.go` verifies the same canonical manifest derives the same ID |
| Strict parsing | parser rejection tests plus bounded fuzz targets for capsule, identity, and duplicate-key decoding |
| Offline Ed25519 signatures and trust states | `internal/capsule/signature_test.go` and CLI workflow tests |
| Deterministic parallel verification | `TestParallelVerifyReportsFirstAssetInManifestOrder` |
| Merkle commitments and possession proofs | chunk proof, malformed proof, tree, and damage-localization tests |
| Authenticated repair | hostile-donor and forged-tree rejection tests |
| Reed Solomon reconstruction | every recoverable loss pattern plus below-threshold and forged-shard rejection tests |
| Sampled retrievability | exact rational cross-checks, confidence-bound tests, partial/exhaustive distinction, and damaged-sample suppression |
| Lineage and exact rotation | same-author, wrong-author, modified transition, wrong key, and transition-reuse tests |
| Revocation, checkpoints, and threshold recovery | rollback, revoked signer, exact-edge, insufficient approval, duplicate approval, and recovery workflow tests |

### Measured

Benchmark values refer only to the published run in `benchmarks/2026-09-14-apple-m5-64mib`:

- Source revision: `057ee8d66f52264f2e5615afb883e907d78b4f85`
- Apple M5, Go 1.22.0, 64 MiB fixture, three runs per operation
- Raw timing observations are stored in `raw.csv`
- Displayed means match those observations after rounding to three decimals
- The audit sample was 349 of 1024 chunks at the recorded policy
- Four data plus two parity shards produced 50.00% encoded storage overhead

These are local observations, not universal performance bounds or network comparisons. This audit did not rerun the benchmark.

### Demonstrated locally

`./scripts/demo-offline-recovery.sh` passed during this audit. It:

1. Created an 8 MiB generated fixture.
2. Packed and signed it with authenticated chunks.
3. Created four data and two parity shard files in a local temporary directory.
4. Confirmed that three surviving shards cannot recover and write no object.
5. Deleted the source, complete capsule object, and two shard files.
6. Reconstructed the exact object from four local shard files.
7. Matched the original and recovered SHA-256 digests.
8. Verified the capsule using an explicitly imported author identity.

It did not use peers, transfer over a network, place shards on independent devices, retain data over time, or run a model.

### Planned or unproven

The following must not be described as working behavior:

- Fresh-device trust bootstrap and consent workflow
- Peer discovery, authentication, encrypted transfer, and resume
- Failure-domain-aware shard placement and repair scheduling
- Longitudinal retention with independently verifiable time
- Runtime compatibility assessment and local inference health checks
- Organizational policy distribution and partition reconciliation
- Independent multi-device continuity experiment
- Superiority over ordinary backup, shared folders, BitTorrent, IPFS, or another baseline

### Normative or aspirational

Statements about preserving human progress, enduring coexistence, consent, pluralism, and invitations to future intelligences are values or research motivations. They are not empirical findings. The site must keep them separate from `verified`, `measured`, and `demonstrated` labels.

## Corrections made by this audit

1. Replaced multi-holder and fresh-node language in the simulator with the implemented local shard-file ceremony.
2. Replaced an uninstrumented "zero runtime network requests" metric with the fact that peer transport is not implemented.
3. Changed "live recovery fixture" to "published demo fixture."
4. Changed "reproducible floor" to an environment-specific published baseline.
5. Replaced fake geographic coordinates on the conceptual roadmap with abstract grid labels.
6. Changed deterministic identity wording from "same content" to "same canonical manifest."
7. Removed semantic-completeness claims that the capsule format cannot automatically prove.
8. Replaced a browser-console "live pass" claim with a documented procedure and pass conditions.
9. Bound displayed benchmark values to their recorded source commit.
10. Replaced achieved cross-provider preservation language in repository metadata with implemented local scope.
11. Removed unsourced "newest supported" dependency assertions and documented the actual validation runtime.
12. Replaced self-assessed "academic standards" language with the factual absence of institutional authorship, peer review, and formal publication.
13. Replaced uninstrumented network-request wording in the demo with the absence of ArkMesh peer/network operations.
14. Corrected the Merkle proceedings citation year from 1987 to 1988.

## Validation performed

The following completed successfully on the current working tree:

```text
go test ./...
go vet ./...
go build ./cmd/arkmesh
./scripts/demo-offline-recovery.sh
./scripts/security-gate.sh
npm run lint
npm run typecheck
npm run build
```

The security gate included race detection and three bounded fuzz targets. The web checks used Node.js 22.22.0 and npm 10.9.4 against existing installed dependencies; `npm ci` was not rerun.

## Residual uncertainty

- No independent security review has been completed.
- No second implementation verifies protocol interoperability.
- No legal review establishes rights for any particular model or capsule contents.
- No physical multi-device or long-duration retention experiment has been completed.
- No model inference or runtime compatibility evidence exists.
- The browser panel was unavailable, so visual composition was not independently inspected during this audit.
- Passing tests do not establish absence of defects.
- External scholarly references were spot-checked, not subjected to a systematic literature review.

Future changes that add or strengthen a public claim should update this audit or add a direct evidence pointer.
