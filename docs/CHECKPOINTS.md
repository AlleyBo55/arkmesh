# Local Lineage Checkpoints

Status: **v0alpha1, unstable**

A local checkpoint records the exact capsule and signer that one operator currently accepts as a lineage head. It prevents another valid capsule, including an older release, from silently satisfying that local policy.

A checkpoint is not a global version number, publication claim, or consensus record.

## Format

```json
{
  "schema_version": "arkmesh.checkpoint/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "signer_id": "ed25519:<digest>"
}
```

The file is stored outside capsules. It contains no private key and grants no authority to modify a capsule.

## Creating a checkpoint

Initial creation requires a valid capsule signature and an explicitly trusted public identity. This prevents an embedded public key from becoming trusted merely because it appears in a capsule.

```bash
arkmesh checkpoint create \
  --out checkpoint.json \
  --trust author/identity.json \
  capsule.ark
```

Creation refuses to replace an existing file.

## Verifying the accepted head

```bash
arkmesh verify --checkpoint checkpoint.json capsule.ark
```

Verification checks object integrity, the capsule signature, exact capsule ID equality, and exact signer ID equality. A matching checkpoint may satisfy strict local trust without adding that signer to a global trust store.

A different capsule is rejected even when it has a valid signature from the same signer. This is the anti-rollback property.

## Advancing a checkpoint

Advancement requires the currently checkpointed capsule as the parent and one direct child:

```bash
arkmesh checkpoint advance \
  --checkpoint checkpoint.json \
  --parent current.ark \
  next.ark
```

ArkMesh verifies both signatures and then accepts one of two lineage edges:

1. Parent and child have the same signer.
2. The parent signer authorized the child's exact new signer through `authority.json`.

The checkpoint file is replaced only after those checks pass. The new head records the child capsule ID and child signer ID. The old capsule remains valid historical data, but it no longer matches the local checkpoint.

Local revocation policies apply during creation, verification, and advancement. A revoked parent cannot authorize advancement, and a revoked current signer cannot satisfy the checkpoint.

## Security boundary

The checkpoint file is trusted local policy. ArkMesh assumes the device owner can protect it with operating system permissions, backups, and physical controls.

This mechanism does not defend against an attacker who can replace or roll back the checkpoint file itself. It also does not provide:

- A global latest version
- Consensus between operators
- Trusted timestamps
- Sequence numbers across branches
- Automatic branch selection
- Recovery when the current parent capsule is unavailable
- Recovery after the current signing key is lost

Stronger protection can place checkpoint copies on offline media or independently administered devices. Future threshold recovery must not silently weaken this local acceptance boundary.
