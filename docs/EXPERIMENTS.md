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

## E0: Local integrity, authenticity, and lineage

Status: implemented in the reference CLI and automated tests.

1. Pack two small fixture assets.
2. Verify the intact capsule successfully.
3. Modify one stored object and confirm verification fails.
4. Remove one object and confirm verification fails.
5. Create an Ed25519 author identity and sign a fresh capsule.
6. Verify it with the matching public identity and require trusted status.
7. Verify it without a trust input and confirm valid unknown author status.
8. Replace the signature bytes and confirm verification fails.
9. Require a different trusted identity and confirm verification fails.
10. Verify an existing unsigned capsule without strict signature requirements.
11. Create a signed child that includes the signed parent's capsule ID.
12. Verify the parent and child with strict lineage enabled.
13. Supply a different parent and confirm verification fails.
14. Sign a child with another identity and confirm update authority fails.
15. Remove the parent signature and confirm lineage verification fails.
16. Create a new author key and an exact rotation transition signed by the parent key.
17. Trust only the parent key and confirm the rotated child verifies through inherited lineage trust.
18. Copy the transition to another child and confirm verification fails.
19. Modify the transition signature and confirm verification fails.
20. Add the child key to a local revocation policy and confirm verification fails.
21. Create a local checkpoint from a directly trusted signed parent.
22. Verify the checkpointed parent without importing another trust root.
23. Advance the checkpoint through an exact parent-signed key rotation.
24. Verify the rotated child as the exact accepted head.
25. Present the old valid parent and confirm checkpoint verification rejects the rollback.
26. Create three independent recovery identities and a two-of-three policy.
27. Pack an exact replacement child without using the old private key.
28. Produce detached approvals from two recovery identities.
29. Assemble and advance the checkpoint as `verified_threshold_recovery`.
30. Confirm one approval, duplicate approvals, modified evidence, and revoked approvals below threshold fail.
31. Pack a multi chunk asset and confirm the signed chunk root verifies.
32. Produce a possession proof for one chunk and verify it while holding only the manifest.
33. Alter the chunk bytes, the chunk index, and the proof path, and confirm each is rejected.
34. Corrupt a stored object and confirm both capsule verification and proof generation fail.
35. Export an authenticated chunk tree from a healthy replica.
36. Corrupt one chunk of a copy and confirm the scan names that exact chunk index and offset.
37. Repair only the damaged chunk from a donor and confirm the capsule verifies again.
38. Offer a hostile donor and a forged tree, and confirm both are refused without writing bytes.

Pass condition: intact objects, valid signatures, same-author ancestry, exact planned rotations, retained local checkpoints, threshold recovery, chunk commitments, and verified repair all succeed. Tested modifications, wrong parents, reused transitions, revoked signers, insufficient recovery approvals, unauthorized descendants, forged chunk proofs, forged trees, hostile donors, and rollback to an old capsule are rejected. Trust is never inferred from an embedded public key, and unsigned root compatibility remains explicit.

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
