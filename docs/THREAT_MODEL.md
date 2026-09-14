# Threat Model

## Security objective

ArkMesh should increase the availability of verified AI capsules without turning replication into a path for malware, poisoned models, private-data leakage, or unauthorized control of another person's device.

## Protected assets

- Integrity and authenticity of model weights, runtimes, knowledge, and policies
- Capsule ancestry and update history
- Device-owner control over storage, compute, networking, and removal
- Private conversations, documents, identity keys, and trust-group membership
- Availability of at least one executable trusted capsule

## Trust boundaries

- A local device owner is trusted to administer their own node.
- A peer is untrusted until explicitly invited and authenticated.
- A peer's bytes remain untrusted even after authentication; hashes and signatures must verify them.
- Capsule authors are trusted only for namespaces and versions explicitly accepted by the owner.
- Knowledge publishers are not automatically authorized to modify models, runtimes, or policies.
- Local inference output is untrusted advice and must retain source/provenance boundaries.

## In-scope adversaries and failures

1. A peer sends corrupted, truncated, or mislabeled chunks.
2. A peer presents an unauthorized capsule as an update.
3. A trusted signing key is stolen.
4. A node attempts excessive storage, bandwidth, CPU, or memory consumption.
5. A capsule contains private or malicious knowledge content.
6. Network partitions create divergent descendants.
7. Hardware failure or bit rot removes stored pieces.
8. A runtime becomes unavailable or incompatible with surviving hardware.
9. A malicious peer lies about possessing complete or executable replicas.
10. Metadata leaks group membership or sensitive capsule names.

## Required mitigations before peer replication

- Content-addressed objects with size and digest verification
- Signed manifests using explicit trust roots
- Replay-resistant version and ancestry rules
- Storage, bandwidth, CPU, and concurrency quotas
- No automatic execution of received files
- Runtime allowlists and sandboxing
- Explicit approval for new authors and privilege changes
- Key revocation and offline recovery procedures
- Independent local inference health checks
- Minimal metadata disclosure before authentication
- Audit logs that contain no private prompts by default

## Propagation boundary

ArkMesh must never:

- Scan networks for installable machines
- Exploit software or credentials
- Install or start itself remotely
- Conceal processes, files, traffic, or resource use
- Disable removal or security controls
- Replicate a capsule without device-owner approval and configured limits
- Treat possession of a peer address as consent

Crossing this boundary would make the system worm-like and invalidate the project's coexistence principles.

## v0alpha1 limitations

The current implementation verifies SHA-256 object integrity, Ed25519 capsule signatures, same-author lineage, exact parent-signed key rotations, local operator revocation policies, and exact local lineage checkpoints. A retained checkpoint rejects another valid capsule as the accepted head and advances only through one verified direct lineage edge. Trust roots, parent capsules, revocation files, and checkpoints are supplied manually.

A checkpoint cannot detect replacement or rollback of the checkpoint file itself. ArkMesh still does not provide emergency recovery after old-key loss, organizational revocation signatures, replay consensus across peers or divergent branches, encrypted transfer, peer authentication, model execution, or legal identity proof. A valid transition proves that the parent key authorized one child key for one lineage edge. It does not prove safety, truth, license compliance, or authorization to execute capsule contents.

## Out of scope

- Protecting hardware from physical destruction or electromagnetic effects
- Guaranteeing correctness of model answers
- Defending a fully compromised operating system
- Preserving proprietary models whose weights are unavailable
- Guaranteeing survival when no compatible hardware, energy source, or reconstructible data remains
