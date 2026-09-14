# ArkMesh Research Charter

## Mission

Preserve humanity's ability to run, understand, and peacefully govern useful open AI after severe infrastructure failure, without dependence on any single organization, global network, or machine.

## Intended beneficiary

The first beneficiary is a small, disconnected human community with at least one functioning general-purpose computer and limited or absent internet access. Examples include disaster-response teams, remote schools, field stations, ships, and air-gapped institutions.

## Research question

Can trusted peers preserve a complete, authenticated, locally executable AI capsule after its original host and central network disappear, while detecting corruption, rejecting unauthorized descendants, and retaining verifiable lineage through partitions?

## Hypothesis

A signed, content-addressed capsule replicated across independently operated devices will provide greater executable availability and provenance under node failure and network partition than a single local installation, ordinary folder synchronization, or unsigned file distribution, with measurable storage and recovery costs.

This hypothesis must remain falsifiable. If a simpler baseline performs as well on the defined measurements, ArkMesh should adopt that simpler method or stop.

## Definitions

- **AI capsule:** Model weights, configuration, tokenizer, runtime requirements, selected knowledge, policy metadata, license records, and integrity/lineage metadata needed to restore a useful local capability.
- **Operational:** A compatible device has verified the capsule and passed an actual local inference health check.
- **Dormant:** A verified capsule exists, but no device is currently executing it.
- **Recoverable:** No complete operational replica exists, but trusted stored pieces are sufficient to reconstruct one.
- **Extinct:** No verified complete copy or reconstructible set remains.
- **Descendant:** A deliberately versioned capsule derived from a known parent. Descendance does not imply a continuous conscious identity.

## Non-negotiable principles

1. **Human consent:** Joining, storing, seeding, updating, and sharing memory are explicit choices.
2. **No self-propagation:** ArkMesh does not exploit, scan, conceal, or install itself on other devices.
3. **Local sovereignty:** A device owner controls storage, compute, network access, and removal.
4. **Inspectable state:** Models, knowledge, policies, memory, and lineage are distinguishable.
5. **Private by default:** Personal data and conversations do not replicate by default.
6. **Authenticity before availability:** An available but poisoned capsule is a failure.
7. **Graceful limits:** The system reports when surviving hardware cannot execute a preserved model.
8. **Evidence over metaphor:** "Survival" claims refer to measured executable continuity, not sentience.
9. **Pluralism:** Communities may preserve multiple models; ArkMesh must not appoint one system as an autonomous authority.
10. **Offline recoverability:** Core operation must not require cloud authentication, telemetry, or a central tracker.

## Definition of first research success

The original host is removed and WAN access is disabled. A newly introduced device reconstructs an authenticated capsule solely from invited LAN peers, rejects a corrupted object and an unauthorized descendant, then completes a local inference health check. An independent person can reproduce the experiment from preserved instructions.

## Explicit non-goals for the first prototype

- Surviving every possible catastrophe
- Reproducing a proprietary hosted model such as GPT or Claude
- Claiming that an AI is conscious, immortal, or continuously experiencing time
- Autonomous model training or unreviewed memory merging
- Hidden persistence or worm-like replication
- Replacing emergency professionals or authoritative medical guidance
- Building a global token economy

## Governance posture

Protocol, identity, update authority, privacy, and licensing decisions are one-way doors once outside contributors and capsules depend on them. They remain experimental until documented threat reviews and reproducible tests justify stabilization.
