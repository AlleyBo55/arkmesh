# ArkMesh

[![Status: Research Prototype](https://img.shields.io/badge/status-research%20prototype-f59e0b)](#project-status)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](#build-and-test)
[![License: MIT](https://img.shields.io/badge/license-MIT-22c55e.svg)](LICENSE)
[![Consent Required](https://img.shields.io/badge/consent-required-2563eb)](#safety)

## Read the thesis first

> [!IMPORTANT]
> **[Read the ArkMesh continuity thesis →](CHARTER.md#research-question)**
>
> ArkMesh begins with a falsifiable research question, not a product promise: can consenting peers preserve an authenticated, locally executable AI capability after its original host and central network disappear? The thesis defines the hypothesis, simpler baselines, failure conditions, legal boundary, and first decisive experiment. The web application presents its visual edition at `/wiki`; its [thesis page source](web/app/wiki/page.tsx) is included in this repository.

## What happens to AI when the cloud is gone?

ArkMesh is a research project for preserving AI models, runtimes, knowledge, and recovery instructions that operators are legally allowed to possess and run. The format is model agnostic: it can support open weight, research, local, and owner authorized frontier systems without overriding licenses, access controls, or provider rights.

The goal is simple to state and difficult to prove:

> Remove the original host. Disconnect the internet. Add a fresh device. Recover a verified AI capsule from nearby peers and run it locally.

A downloaded model is not enough. It may also need a tokenizer, configuration, runtime, documentation, compatible hardware, and proof that none of those files were changed. ArkMesh packages those pieces as one content addressed capsule and tests whether people can recover it after infrastructure failure.

This is not an attempt to copy GPT, Claude, or another hosted service without authorization. It is not a claim that software is conscious. ArkMesh can preserve only capabilities whose files an operator may lawfully possess; a capsule is never permission to bypass a provider or redistribute restricted material.

Read [The ArkMesh Promise](MANIFESTO.md) for the longer motivation and the project's position on human and AI coexistence.

## Project status

> [!WARNING]
> ArkMesh is a v0alpha1 research prototype. Networking and inference are not implemented. Do not use it for disaster critical work.

### Working now

- Pack local files into a content addressed capsule
- Record SHA-256 hashes, roles, file names, and sizes
- Commit a Merkle chunk root per asset inside the signed manifest
- Prove and verify possession of one exact chunk without holding the whole object
- Locate the exact damaged chunks of a replica against an authenticated chunk tree
- Repair only the damaged chunks, verifying every donor chunk before writing it
- Rebuild a lost object from Reed Solomon shards with no donor holding its bytes
- Bound how much of a replica can be missing, by random sampling with stated confidence
- Record an audit history so retention is observed rather than assumed
- Verify independent objects in parallel with deterministic error reporting
- Derive a deterministic capsule ID
- Create local Ed25519 author identities
- Sign capsules without placing private keys inside them
- Verify trusted, unknown, invalid, and unsigned capsule states offline
- Bind a child capsule to one parent ID
- Verify same-author update authority across a parent and child
- Authorize one exact child key through a parent-signed rotation transition
- Apply explicit local revocation policies during offline verification
- Pin one accepted lineage head in a local checkpoint
- Advance a checkpoint through one verified direct lineage edge
- Recover from old-key loss through distinct threshold approvals
- Inspect a capsule manifest
- Detect missing files, changed content, invalid signatures, wrong parents, reused transitions, revoked signers, insufficient recovery approval, unauthorized child keys, and rollback against a retained checkpoint

### Not built yet

- Distribution or organizational signing of revocation and recovery policies
- Peer discovery and encrypted transfer
- Shard placement, repair scheduling, and replica health metrics
- Measured retention over time rather than possession at one moment
- Independently attested time, so retention logs rely on local clocks
- Local model inference
- Replay agreement across peers and divergent branch reconciliation
- Recovery dashboard

A valid signature proves that the holder of a specific private key signed the capsule ID. It does not prove the author's legal identity, the truth of the contents, license compliance, or safety. Trust is imported explicitly by each operator.

## Reproducible evidence

The [`v0.1.0-alpha.1` release](https://github.com/AlleyBo55/arkmesh/releases/tag/v0.1.0-alpha.1) includes a local failure ceremony and published raw measurements.

### Delete it, then recover it

```bash
./scripts/demo-offline-recovery.sh
```

The demonstration deletes the publisher source, the complete capsule object, and two of six Reed Solomon shards. ArkMesh rebuilds the exact 8 MiB object from the four surviving shards, requires its SHA-256 digest to match the signed manifest, and verifies the recovered capsule against explicit author trust. A separate attempt with only three usable shards must fail without writing an object.

This ceremony makes no runtime network requests. It proves local donorless reconstruction and authenticated acceptance, not peer transfer, physical device survival, or model execution. See the [full procedure and limits](docs/DEMO.md).

### Measured on a 64 MiB fixture

Wall clock means from three runs on an Apple M5 with Go 1.22:

- Signed capsule pack: 0.125 s
- Trusted capsule verification: 0.047 s
- Four data plus two parity protection: 0.211 s
- Recovery after losing two shards: 0.186 s
- 99 percent confidence sampled audit: 0.013 s
- Full copy plus SHA-256 baseline: 0.140 s
- rsync plus SHA-256 baseline: 0.358 s

The encoded shards added exactly 50 percent storage overhead. The sampled audit checked 349 of 1,024 authenticated chunks. These operations perform different work, so the table is evidence, not a claim that every row is directly equivalent or that ArkMesh is universally faster.

Read the [methodology and limitations](docs/BENCHMARKS.md), [generated report](benchmarks/2026-09-14-apple-m5-64mib/results.md), and [raw CSV](benchmarks/2026-09-14-apple-m5-64mib/raw.csv). BitTorrent, IPFS, peer transport, and inference remain unmeasured.

## The planned experiment

```text
                 WAN and cloud unavailable
                           X

 Original node        Trusted LAN peers       Fresh device
  offline              node B    node C        node D
                          |         |              ^
                          +=========+==============+
                          verified capsule recovery
```

The first major experiment will:

1. Place a capsule on several invited devices.
2. Remove the original publisher.
3. Disconnect WAN access while keeping a local network.
4. Introduce a fresh device.
5. Reconstruct the capsule from surviving peers.
6. Reject corrupted objects and unauthorized updates.
7. Run a local inference health check without cloud authentication or a central tracker.

The experiment must be reproducible by someone who did not build ArkMesh.

## Why this is different

**Hosted AI APIs** provide capable managed services, but access depends on a provider and working network.

**Local runners such as llama.cpp** execute models offline, but they do not preserve a complete capability across several devices.

**BitTorrent and IPFS** distribute bytes, but they do not define AI specific runtime requirements, health checks, trust rules, or lineage.

**Distributed inference systems** share computation, while ArkMesh focuses on preservation, reconstruction, and proof of origin.

ArkMesh will reuse existing tools where they already solve the problem. If an ordinary backup or shared folder performs just as well in the experiments, the simpler method wins.

## Design goals

- Verify every stored object by cryptographic digest
- Preserve models, runtimes, knowledge, configuration, licenses, and recovery instructions together
- Recover without cloud login, telemetry, or a central tracker
- Require explicit approval from every participating device owner
- Keep private memory local unless its owner publishes it
- Track signed ancestry when capsules change
- Report whether surviving hardware can actually execute a capsule
- Compare results against simpler recovery methods
- Support redistributable open weight models rather than one vendor

## Quick start

Requirements: Go 1.22 or newer.

ArkMesh does not bundle a model. Use only files whose licenses permit your intended local use and redistribution.

```bash
go run ./cmd/arkmesh help
go test ./...
go build ./cmd/arkmesh
```

The help command documents identity, checkpoint, pack, inspect, and verify syntax.

A signed capsule contains:

```text
field-assistant.ark/
├── manifest.json
├── signature.json
└── objects/
    ├── <sha256>
    └── <sha256>
```

Unsigned capsules remain supported for local integrity checks. Object paths come only from validated content hashes. Human supplied names and roles remain metadata.

## Web interface

The command-aware landing page, educational learning path, and system wiki live in [`web`](web/README.md). The landing begins with a native-scroll deep-sea journey through ArkMesh's continuity thesis, `/learn` presents an eight-module curriculum, and the remaining surfaces label implemented evidence and planned protocol work.

```bash
cd web
npm ci
npm run dev
```

The web application requires Node.js 20.9 or newer.

## Build, test, and reproduce

```bash
go test ./...
go build ./cmd/arkmesh
./scripts/security-gate.sh
./scripts/demo-offline-recovery.sh
./scripts/benchmark-local.sh
```

The security gate adds race detection and bounded fuzzing for capsule and identity parsers. The recovery demonstration deletes the publisher source, the capsule object, and two shards before rebuilding and verifying the exact signed bytes. The benchmark publishes raw timings, storage overhead, and environment details. The implementation currently uses only the Go standard library.

## Research roadmap

**Base camp, implemented:** deterministic capsules, strict parsing, offline signatures, lineage authority, checkpoints, threshold recovery, Merkle repair, Reed Solomon reconstruction, and sampled audits.

1. **Fresh-device trust and consent:** identity bootstrap, scoped invitations, capsule selection, quotas, expiry, and explicit operator approval.
2. **Authenticated peer transfer:** invited discovery, bounded metadata disclosure, encrypted resumable transfer, hostile-donor rejection, and atomic publication.
3. **Failure-domain placement and repair:** shard placement across independent risks, health reporting, degradation detection, and approved repair scheduling.
4. **Longitudinal retention and attested time:** unpredictable audits across months with independently verifiable event timing.
5. **Runtime compatibility and local health:** hardware assessment and allowlisted offline execution without cloud authentication.
6. **Partition authority and reconciliation:** signed policy distribution, retained divergent branches, and deliberate human reconciliation.
7. **Independent continuity experiment:** remove the publisher, disconnect WAN access, recover and run on a fresh device, publish raw evidence, and compare simpler baselines.

## Safety

ArkMesh will not intentionally:

- Scan for targets or exploit another machine
- Install or start itself remotely
- Hide processes, files, traffic, or resource use
- Replicate without owner approval and resource limits
- Train on private conversations automatically
- Execute arbitrary plugins received from peers
- Treat discovery as consent
- Claim consciousness, immortality, or guaranteed survival

ArkMesh is preservation software, not a worm or self installing agent. An available but poisoned capsule counts as a failure.

Read the full [threat model](docs/THREAT_MODEL.md).

## Research documents

- [Research Charter](CHARTER.md): mission, definitions, hypothesis, and limits
- [Threat Model](docs/THREAT_MODEL.md): assets, trust boundaries, attacks, and propagation limits
- [Capsule Protocol](docs/PROTOCOL.md): current capsule, signature, and lineage format
- [Chunk Commitments](docs/CHUNKS.md): Merkle chunk roots and possession proofs
- [Erasure Coding](docs/ERASURE.md): Reed Solomon shards and donorless recovery
- [Sampled Audits](docs/AUDIT.md): retrievability bounds, confidence, and retention history
- [Key Authority and Revocation](docs/AUTHORITY.md): exact planned rotations and local rejection policy
- [Threshold Emergency Recovery](docs/RECOVERY.md): independent approvals and exact recovery edges
- [Local Lineage Checkpoints](docs/CHECKPOINTS.md): accepted heads, direct advancement, and rollback limits
- [Security Invariants](docs/SECURITY_INVARIANTS.md): trust properties every change must preserve
- [Security Policy](SECURITY.md): private vulnerability reporting and response expectations
- [Experiment Plan](docs/EXPERIMENTS.md): metrics, baselines, partition tests, and recovery rehearsal
- [Offline Recovery Demo](docs/DEMO.md): reproducible donorless reconstruction ceremony
- [Benchmarks](docs/BENCHMARKS.md): methodology, raw observations, results, and limits
- [Public Claims Audit](docs/CLAIMS_AUDIT.md): evidence classes, corrected overstatements, validation, and residual uncertainty
- [Release Process](docs/RELEASING.md): validation, artifacts, tags, and publication
- [Changelog](CHANGELOG.md): release contents and known limits
- [GitHub Setup](.github/REPOSITORY_METADATA.md): repository description, topics, and issue labels

## Contributing

Researchers, reproducers, reviewers, and implementers are welcome. Useful contributions include protocol review, hostile tests, independent reproduction, baseline experiments, model license research, recovery documentation, local inference work, and distributed systems analysis.

Start with the [contributor guide](CONTRIBUTING.md), [propose a research question](https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml), or [report an independent reproduction](https://github.com/AlleyBo55/arkmesh/issues/new?template=reproduction-report.yml).

Read the charter and threat model before proposing network or execution features. Contributions must preserve consent, local control, inspectability, and offline operation. A counterexample that narrows or disproves a claim is a valuable contribution.

AI assisted contributions are welcome under the same rules as human contributions. A tool may analyze or suggest work, but it receives no authority to access systems, replicate software, spend resources, or act beyond the human operator's approval.

## License

ArkMesh source code and documentation use the [MIT License](LICENSE).

Models, runtimes, datasets, and knowledge inside capsules retain their own licenses. A capsule is not permission to copy or redistribute its contents.

**Related fields:** offline AI, local LLM, peer to peer AI, decentralized AI, resilient computing, air gapped inference, model preservation, content addressed storage, cryptographic provenance, disaster recovery, edge AI, and distributed systems.
