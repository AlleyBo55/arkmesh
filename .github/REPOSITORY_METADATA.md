# GitHub Repository Metadata

Use these values when the repository is created on GitHub. This file is documentation; GitHub does not automatically apply repository settings from it.

## Description

> Research prototype for content-addressed AI capsules, offline signatures, authenticated repair, and donorless local recovery. Peer distribution and inference are planned.

## Suggested topics

GitHub allows up to 20 topics. Start with this focused set:

```text
offline-ai
local-llm
peer-to-peer
decentralized-ai
ai-resilience
model-preservation
disaster-recovery
air-gapped
edge-ai
content-addressed-storage
cryptographic-provenance
distributed-systems
digital-preservation
golang
open-source-ai
```

Do not add implementation topics such as `libp2p`, `llama-cpp`, or `ed25519` until those integrations exist in the repository.

## Social preview copy

**Headline:** Can an AI capability outlive its original host?

**Subheadline:** A local research prototype for signed capsules, authenticated repair, Reed Solomon reconstruction, and measured retrievability. Peer continuity remains unproven.

## Issue labels

| Label | Color | Description |
|---|---|---|
| `type: bug` | `d73a4a` | Confirmed incorrect behavior in implemented functionality |
| `type: feature` | `1d76db` | Bounded product or protocol capability |
| `type: research` | `7057ff` | Research question, prior art, or hypothesis work |
| `type: experiment` | `5319e7` | Reproducible benchmark or failure experiment |
| `type: security` | `b60205` | Threat model, cryptography, trust, or vulnerability work |
| `type: documentation` | `0075ca` | Documentation or recovery instructions |
| `area: capsule` | `0e8a16` | Capsule manifest, storage, packing, or verification |
| `area: identity` | `0052cc` | Signing identities, trust roots, rotation, or revocation |
| `area: network` | `006b75` | Discovery, transport, synchronization, or partitions |
| `area: runtime` | `008672` | Local inference, compatibility, or health checks |
| `area: resilience` | `2cbe4e` | Repair, erasure coding, hardware diversity, or recovery media |
| `area: dashboard` | `54d6eb` | Evidence visualization and local operator interface |
| `status: needs evidence` | `fbca04` | Claim requires a test, measurement, or citation |
| `status: blocked` | `bdbdbd` | Cannot proceed until a named dependency or decision is resolved |
| `good first issue` | `7057ff` | Small, documented task suitable for a first contribution |
| `help wanted` | `008672` | Maintainers welcome outside implementation or research help |
| `breaking change` | `e99695` | Alters a public format, command, trust rule, or compatibility promise |

## Label rules

- Apply one primary `type:` label.
- Apply every relevant `area:` label.
- Use `status: needs evidence` for attractive claims that lack reproducible support.
- Security reports that may expose a vulnerability should not begin as public issues; establish a private reporting channel before launch.
- Do not label speculative ideas as roadmap commitments.

## Suggested About-panel values

- **Website:** Leave blank until a maintained project site or documentation deployment exists.
- **Releases:** Disable expectations until signed artifacts and a versioning policy exist.
- **Discussions:** Enable after the research charter is public; use it for conceptual questions instead of issue noise.
- **Sponsorships:** Defer until governance and funding-use policies are documented.
