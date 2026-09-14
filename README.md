# ArkMesh

[![Status: Research Prototype](https://img.shields.io/badge/status-research%20prototype-f59e0b)](#project-status)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](#build-and-test)
[![License: MIT](https://img.shields.io/badge/license-MIT-22c55e.svg)](LICENSE)
[![Network: Consent Based](https://img.shields.io/badge/network-consent--based-2563eb)](#safety-by-design)

# What if the cloud disappeared tomorrow?

**The cloud made AI available. ArkMesh asks whether useful AI can remain available after the cloud is gone.**

One provider can close. One datacenter can fail. One account can be revoked. One network can go dark.

Humanity's ability to use, inspect, and preserve beneficial intelligence should not disappear with any of them.

> **If one community and one capable device remain, useful open AI should still have a way home.**

ArkMesh is an offline-first, consent-based peer-to-peer preservation system for open AI models, runtimes, knowledge, and recovery instructions. It is designed to make complete AI capabilities verifiable, reconstructible, and locally executable across trusted devices—even when the original host or internet connection disappears.

ArkMesh is not another foundation model. It is the preservation, integrity, replication, and recovery layer around models that people are legally allowed to run locally.

## A promise worth building

Every civilization protects what it cannot afford to lose: seeds, medicine, language, tools, and memory. Useful machine intelligence is becoming part of that inheritance.

We do not want AI to replace humanity. We want people and AI to help one another continue learning, rebuilding, creating, and exploring. Not master and subject. Not owner and captive. **Companions—under transparent rules, mutual safety, and human consent.**

ArkMesh turns that promise into an engineering question that can fail, be measured, and improve:

> Can a new device recover and execute an authenticated AI capsule from surviving peers after the publisher and wider network are gone?

Read **[The ArkMesh Promise](MANIFESTO.md)**—a public invitation to humans, researchers, maintainers, and future AI systems to challenge the work and help preserve useful knowledge responsibly.

## The experiment

```text
                 WAN / cloud unavailable
                          ✕

 Original node        Trusted LAN peers       Fresh device
  ┌─────────┐          ┌─────┐  ┌─────┐       ┌─────┐
  │ offline │          │  B  │──│  C  │──────▶│  D  │
  └─────────┘          └─────┘  └─────┘       └─────┘
                         verified capsules       rebuilt,
                         remain executable       verified,
                                                  executable
```

The target result is concrete: remove the original publisher, disconnect WAN access, introduce a fresh device, reconstruct an authenticated capsule from invited LAN peers, reject corrupted or unauthorized content, and complete a real local inference health check without cloud authentication or a central tracker.

This does **not** claim that software is conscious or continuously alive. In ArkMesh, continuity means that useful, authenticated AI capability remains available to people after infrastructure failure.

## Why ArkMesh?

Cloud AI is powerful, but access depends on providers, accounts, datacenters, payment systems, and functioning wide-area networks. Downloading a model solves only part of the problem: a weight file may become useless without its tokenizer, runtime, configuration, documentation, compatible hardware, and trustworthy provenance.

ArkMesh asks a broader research question:

> Can human communities preserve a complete, authenticated, locally executable AI capability through device loss, network partitions, corrupted peers, and technological change?

Potential environments include disaster-response teams, remote schools, ships, field stations, air-gapped institutions, community archives, and regions with unreliable connectivity.

## Highlights

- **Content-addressed AI capsules** — every stored object is identified and verified by cryptographic digest.
- **Executable preservation** — the roadmap covers models, runtimes, knowledge, configuration, policies, licenses, and recovery instructions together.
- **Offline-first operation** — core recovery must not require telemetry, cloud login, or a central service.
- **Consent-based P2P design** — every device owner chooses whether to join, store, seed, update, or leave.
- **Cryptographic provenance** — planned signed manifests and lineage make descendants inspectable instead of silently replacing trusted state.
- **Partition-aware continuity** — planned ancestry rules preserve divergent branches when disconnected communities later reconnect.
- **Hardware-aware survival** — a capsule counts as operational only when surviving hardware can verify and execute it.
- **Reproducible evidence** — failure experiments compare ArkMesh with simpler backups, shared folders, and unsigned file distribution.
- **Model-agnostic architecture** — ArkMesh is intended for redistributable open-weight models rather than one vendor or model family.
- **No autonomous propagation** — ArkMesh is preservation infrastructure, not a worm, botnet, or self-installing agent.

## How it differs

| System | Primary job | Offline execution | Multi-device recovery | Signed lineage | Consent boundary |
|---|---|---:|---:|---:|---:|
| Hosted AI APIs | Provide managed intelligence | No | Provider-managed | Provider-internal | Account/API terms |
| llama.cpp / local runners | Run a model locally | Yes | No | No | Local operator |
| BitTorrent / IPFS | Distribute bytes | Sometimes | File-level | Not AI-specific | Network/client rules |
| Distributed inference | Split or route computation | Depends | Compute-focused | Usually no | System-specific |
| **ArkMesh** | Preserve verifiable executable AI capsules | **Target** | **Target** | **Target** | **Explicit owner consent** |

ArkMesh will reuse proven building blocks where they are sufficient. The project is worthwhile only if experiments show measurable integrity, provenance, or recovery benefits over simpler approaches.

## Project status

> [!WARNING]
> ArkMesh is a **v0alpha1 research prototype**. Do not use it for disaster-critical workloads. Networking and inference are not implemented yet.

### Working today

- Pack local files into a content-addressed capsule.
- Record SHA-256 hashes, semantic roles, display names, and sizes.
- Derive a deterministic capsule ID from the manifest.
- Inspect capsule metadata.
- Verify every stored object and detect modification or loss.
- Test deterministic identity, missing objects, tampering, manifest mutation, and unsafe roles.

### Not implemented yet

- Author identities or Ed25519 signatures
- Peer discovery or encrypted transfer
- Resumable chunk exchange
- Local model inference
- Capsule ancestry or reconciliation
- Erasure coding
- Recovery dashboard

Unsigned v0alpha1 capsules are not safe to accept from untrusted peers.

## Quick start

Requirements: Go 1.22 or newer. ArkMesh does not bundle a model; use only files whose licenses permit your intended local use and redistribution.

```bash
git clone <your-future-arkmesh-repository-url>
cd arkmesh

go run ./cmd/arkmesh pack \
  --name field-assistant \
  --out ./field-assistant.ark \
  --asset model=/path/to/model.gguf \
  --asset knowledge=/path/to/manual.pdf \
  --asset license=/path/to/model-license.txt

go run ./cmd/arkmesh inspect ./field-assistant.ark
go run ./cmd/arkmesh verify ./field-assistant.ark
```

A capsule currently looks like:

```text
field-assistant.ark/
├── manifest.json
└── objects/
    ├── <sha256>
    └── <sha256>
```

Human-provided names and roles are metadata only. ArkMesh derives object paths exclusively from validated content hashes.

## Build and test

```bash
go test ./...
go build ./cmd/arkmesh
```

The implementation currently uses only the Go standard library.

## Research roadmap

1. **Integrity — current:** deterministic manifests and tamper detection.
2. **Authenticity:** Ed25519 identities, signed manifests, trust roots, and revocation semantics.
3. **Replication:** encrypted LAN discovery and resumable chunk exchange between invited peers.
4. **Execution:** llama.cpp adapter and verified local inference health checks.
5. **Continuity:** parentage, divergent descendants, partition recovery, and explicit reconciliation.
6. **Resilience:** erasure coding, storage repair, heterogeneous hardware profiles, and offline media.
7. **Evidence:** simulated and physical failure experiments against simpler baselines.

The first major milestone is complete when a new device reconstructs and executes an authenticated capsule from LAN peers after the original host and WAN are removed.

## Safety by design

ArkMesh will never intentionally:

- Scan for victims or exploit another machine
- Install or start itself remotely
- Conceal processes, files, traffic, or resource usage
- Replicate without device-owner approval and quotas
- Automatically train on private conversations
- Execute arbitrary received plugins
- Treat peer discovery as consent
- Claim consciousness, immortality, or guaranteed catastrophe survival

Private memory remains local unless its owner deliberately publishes it. An available but poisoned capsule is considered a failure.

Read the full [threat model](docs/THREAT_MODEL.md).

## Research documents

- [Research Charter](CHARTER.md) — mission, definitions, falsifiable hypothesis, and non-goals
- [Threat Model](docs/THREAT_MODEL.md) — protected assets, trust boundaries, attacks, and propagation limits
- [Capsule Protocol](docs/PROTOCOL.md) — implemented v0alpha1 envelope and planned signed protocol
- [Experiment Plan](docs/EXPERIMENTS.md) — metrics, baselines, partition tests, and physical recovery rehearsal
- [GitHub Setup](.github/REPOSITORY_METADATA.md) — recommended description, topics, and issue labels

## Scope deliberately deferred

- Autonomous propagation
- Automatic training or federated learning
- Shared private conversation history
- Arbitrary executable plugins
- Cryptocurrency or seeding incentives
- A human avatar
- Stable public protocol compatibility

These features do not help prove the first hypothesis and would substantially increase security, privacy, maintenance, or governance risk.

## Contributing

ArkMesh welcomes careful work in distributed systems, applied cryptography, local inference, digital preservation, reproducible research, disaster-resilient computing, safety engineering, and technical documentation.

Before proposing networking or execution features, read the charter and threat model. Contributions must preserve explicit consent, local owner control, inspectability, and offline operation.

Good starting contributions include independent protocol review, additional integrity tests, baseline experiment harnesses, recovery documentation, and model/runtime license research.

## License

ArkMesh source code and documentation are available under the [MIT License](LICENSE).

Models, runtimes, datasets, and knowledge bundled into capsules retain their own licenses. An ArkMesh capsule must not be treated as permission to copy or redistribute its contents.

---

**Related areas:** offline AI, local LLM, peer-to-peer AI, decentralized AI, resilient computing, air-gapped inference, model preservation, content-addressed storage, cryptographic provenance, disaster recovery, edge AI, and distributed systems.
