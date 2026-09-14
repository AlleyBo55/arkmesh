# Security Invariants

These are release gates, not aspirations.

## Cryptographic identity

1. A key ID is derived from the exact Ed25519 public key.
2. Embedded public keys never become trusted merely by appearing in a capsule.
3. Private keys never enter capsules, policies, approvals, logs, or source control.
4. Every signature payload has a versioned domain separator.

## Capsule integrity

1. A capsule ID binds the exact manifest representation defined by v0alpha1.
2. Every asset path is derived from a validated lowercase SHA-256 digest.
3. Unknown fields, duplicate fields, trailing JSON values, symlinks, and nonregular objects fail closed.
4. Unsigned compatibility never satisfies a requirement for signatures or trust.

## Chunk commitments

1. A chunk root and chunk size live in the manifest, so the capsule signature covers them.
2. Leaves and internal nodes use distinct domain prefixes.
3. An unpaired node is promoted, never duplicated.
4. Chunk count is derived from the signed asset size, never from a proof.
5. Proof verification compares every restated field against the signed manifest.
6. Sibling direction is derived from the index and level sizes, so proofs carry no direction bits.
7. A possession proof carries real chunk bytes and proves possession only at that moment.
8. A chunk tree file is authenticated by recomputing the signed root, never trusted because it was supplied.
9. Repair verifies each donor chunk against the authenticated tree before writing it, and never writes unverified bytes.
10. Parallel object verification reports the same asset as sequential verification.
11. An erasure plan is not a trust anchor: reconstructed bytes are verified against the signed manifest before they are written.
12. Recovery writes through a temporary file and refuses to replace an object when verification fails.
13. Recovery refuses to proceed when fewer than the required number of usable shards survive.

## Authority

1. Same-author lineage, planned rotation, and threshold recovery are mutually exclusive evidence paths.
2. Rotation binds one exact parent, child, previous signer, and next signer.
3. Recovery binds one exact policy and the same exact edge fields.
4. Recovery requires at least two distinct policy members and two valid approvals.
5. Revoked capsule signers fail verification. Revoked recovery approvers do not count.
6. No authority transition silently modifies a global trust store.

## Rollback and concurrency

1. A retained checkpoint matches one exact capsule ID and signer ID.
2. Advancement starts from the checkpointed parent and verifies one direct edge.
3. Concurrent local advancement is refused instead of selecting a silent last writer.
4. Local policy-file compromise remains inside the trusted device boundary and must never be described as solved.

## Execution and replication

1. Verified never means safe to execute.
2. Received content remains data until an explicit local execution policy approves it.
3. Discovery is not consent.
4. Replication must require invitation, authentication, quotas, and owner control.
5. No scanning, exploitation, concealed persistence, or autonomous propagation is permitted.

## Required evidence

Every security-sensitive change needs targeted rejection tests, compatibility tests, formatting checks, all unit tests, `go vet`, a build, and a live smoke test when a public workflow changes. Parser changes also require fuzz seeds or a bounded fuzz run.
