# Experiment Plan

ArkMesh is useful only if it improves measured executable continuity or provenance over simpler systems. Experiments must publish commands, topology, inputs, raw events, results, and known limitations.

## Metrics

- **Executable availability:** Fraction of observation time with at least one verified node passing inference health checks.
- **Complete replicas:** Independently executable capsule copies.
- **Recoverable state:** Whether available trusted pieces can reconstruct a complete capsule.
- **Reconstruction time:** Time from approved join to verified operational state.
- **Storage overhead:** Replicated and parity bytes relative to one complete capsule.
- **Recovery traffic:** Bytes transferred after node or object loss.
- **Corruption rejection:** Altered chunks and objects detected before acceptance.
- **Unauthorized-update rejection:** Unapproved descendants prevented from replacing trusted state.
- **Lineage correctness:** Accepted descendants retain valid, queryable ancestry.
- **Hardware coverage:** Surviving devices capable of executing at least one preserved capability tier.

## Baselines

1. One local llama.cpp installation
2. Ordinary shared-folder synchronization
3. Unsigned BitTorrent or IPFS-style file distribution
4. Central object storage when the central service is reachable

A more complex ArkMesh mechanism should be rejected when a simpler baseline performs equivalently for the target scenario.

## E0: Local integrity

Status: scaffold target.

1. Pack two small fixture assets.
2. Verify the capsule successfully.
3. Modify one stored object.
4. Confirm verification fails with the affected asset identified.
5. Remove one object and confirm verification fails.

Pass condition: all intact assets verify; every tested modification or deletion is detected.

## E1: Three-node LAN continuity

1. Start three invited nodes on one LAN.
2. Seed one capsule from node A.
3. Confirm B and C obtain verified complete replicas.
4. Remove WAN access while retaining LAN.
5. permanently stop A.
6. Confirm B and C remain operational.
7. Add D and reconstruct solely from B and C.

Pass condition: D passes a local inference health check without A, WAN, a central tracker, or cloud authentication.

## E2: Poisoning rejection

1. Send altered chunks under a valid object identifier.
2. Present a well-formed capsule signed by an unauthorized key.
3. Present a descendant with an invalid parent relationship.

Pass condition: none replace trusted state; each rejection is visible in the audit log without exposing private prompts.

## E3: Partition and divergence

1. Split peers into two disconnected islands.
2. Create separately authorized knowledge descendants.
3. Reconnect the islands.
4. Preserve both ancestries and require human approval for reconciliation.

Pass condition: no silent last-writer-wins replacement and no private-memory exchange.

## E4: Failure curve

Run repeated trials while randomly removing 10%, 25%, 50%, 75%, and 90% of nodes. Measure executable availability and reconstruction cost for full replication and candidate erasure-coding strategies.

Pass condition: results identify a defensible resilience/storage tradeoff and include confidence intervals across repeated trials.

## E5: Heterogeneous recovery

Preserve at least two explicitly related capability tiers, then remove hardware capable of executing the larger model.

Pass condition: the system accurately reports the larger tier as dormant/unexecutable while a smaller authenticated tier remains operational. It must not claim they are behaviorally identical.

## Physical recovery rehearsal

The eventual field test must include printed instructions, a fresh compatible machine, local power, and offline installation media. A participant unfamiliar with the implementation should restore and run a capsule without contacting its authors.
