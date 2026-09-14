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
