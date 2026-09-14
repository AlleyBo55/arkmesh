# Key Authority and Revocation

Status: **v0alpha1, unstable**

ArkMesh separates planned key rotation from local revocation because they answer different questions.

- A key rotation transition proves that the current parent key authorized one exact child key for one exact descendant.
- A revocation policy records an operator's local decision to reject one or more keys.

Neither mechanism proves a person's legal identity.

## Planned key rotation

A rotated child contains three authenticated relationships:

1. The child manifest includes the exact parent capsule ID.
2. The child capsule signature is made by the new key.
3. `authority.json` is made by the old parent key and binds the parent ID, child ID, old key ID, and new key ID.

```json
{
  "schema_version": "arkmesh.authority/v0alpha1",
  "action": "rotate",
  "parent_capsule_id": "sha256:<parent-digest>",
  "child_capsule_id": "sha256:<child-digest>",
  "previous_signer_id": "ed25519:<old-key-digest>",
  "next_signer": {
    "schema_version": "arkmesh.identity/v0alpha1",
    "algorithm": "ed25519",
    "key_id": "ed25519:<new-key-digest>",
    "public_key": "<standard-base64>"
  },
  "signature": "<standard-base64>"
}
```

The old key signs this exact UTF-8 construction:

```text
"arkmesh.capsule.authority/v0alpha1\n" +
action + "\n" +
parent_capsule_id + "\n" +
child_capsule_id + "\n" +
previous_signer_id + "\n" +
next_signer_key_id + "\n"
```

The transition cannot be copied to another child because the child capsule ID is part of the signed payload. It cannot be attached to another parent because the parent capsule ID and previous signer are also part of the payload.

## Trust inheritance

Strict verification may begin with only the old parent identity in the operator's trust set.

1. Verify the parent and require it to be trusted.
2. Verify the child signature without assuming the new key is trusted.
3. Verify the parent and child IDs in `authority.json`.
4. Verify that the old key signed the transition to the child's exact new key.
5. Report `verified_key_rotation`.

After this edge verifies, the child key is accepted for that lineage edge. ArkMesh does not silently add it to a global trust store.

## Same-author descendants

When parent and child use the same key, no authority file is allowed. The lineage result remains `verified_same_author`. This prevents ambiguous or unnecessary transition metadata.

## Local revocation policy

A revocation policy is intentionally simple:

```json
{
  "schema_version": "arkmesh.revocations/v0alpha1",
  "revoked_key_ids": [
    "ed25519:<digest>"
  ]
}
```

Key IDs must be lowercase, unique, and sorted. Multiple local policy files may be merged during verification.

Revocation overrides cryptographic validity and explicit trust. If a capsule signer appears in the active local policy, verification fails before lineage authority is considered.

This file is not signed by the revoked key. A compromised key cannot be trusted to decide whether it is compromised. The operator controls which revocation policies are active and how those files are distributed.

## What this does not solve

- Recovery after the old private key is lost before a transition is signed
- Proving that a revocation policy came from a particular organization
- Consensus about revocation across disconnected communities
- Expiration dates or sequence numbers
- Choosing between divergent descendants
- Undoing content already accepted before a revocation policy arrived

Those problems require a recovery authority and replay policy that have not been designed yet.
