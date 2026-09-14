# Chunk Commitments and Possession Proofs

Status: **v0alpha1, unstable**

Each stored object may carry a Merkle commitment over fixed size chunks. This makes large models verifiable in parts instead of only as one whole file.

It exists so ArkMesh can later support resumable transfer, targeted repair, and honest availability claims from peers.

## Manifest fields

Chunk fields are part of the asset entry, so the capsule ID and the capsule signature already cover them:

```json
{
  "role": "model",
  "name": "model.gguf",
  "size": 3145739,
  "sha256": "<digest>",
  "chunk_size": 1048576,
  "chunk_root": "<digest>"
}
```

Both fields are omitted together on capsules that predate this feature, so existing root capsule IDs stay identical. A manifest that declares one field without the other is rejected.

`chunk_size` must be a power of two between 4096 bytes and 64 MiB. The default is 1 MiB. The chunk count is derived from the signed size, never taken from a proof.

## Tree construction

Leaves and internal nodes use different domain prefixes, so a leaf digest can never be replayed as an internal node:

```text
leaf  = SHA256("arkmesh.chunk.leaf/v0alpha1\n"  + chunk_bytes)
node  = SHA256("arkmesh.chunk.node/v0alpha1\n"  + left + right)
empty = SHA256("arkmesh.chunk.empty/v0alpha1\n")
```

An unpaired node is promoted unchanged to the next level. It is never duplicated to fill a pair.

Duplicating an unpaired node is what allowed two different sequences to produce one Merkle root in Bitcoin's original construction, reported as CVE-2012-2459. ArkMesh does not copy that behavior.

A single chunk object has a root equal to its leaf digest. An empty object uses the fixed empty root above.

## Possession proofs

A proof carries the chunk bytes, so verifying it proves the prover actually holds those bytes rather than only a digest:

```json
{
  "schema_version": "arkmesh.chunk-proof/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "asset_sha256": "<digest>",
  "chunk_size": 1048576,
  "chunk_count": 4,
  "chunk_index": 2,
  "chunk_root": "<digest>",
  "chunk": "<hex bytes>",
  "path": ["<sibling digest>", "..."]
}
```

Verification compares every restated field against the signed manifest, then recomputes the root from the chunk bytes and the sibling path.

The proof carries no direction bits. Sibling order is derived from the chunk index and the level sizes implied by the signed chunk count, so an attacker cannot choose how nodes combine.

A proof is rejected when the capsule ID differs, the asset is unknown, the chunk size or root disagrees with the manifest, the chunk count disagrees with the signed size, the index is out of range, the chunk length is wrong for its position, the path is too short or too long, or the recomputed root does not match.

## Commands

```bash
arkmesh chunks inspect CAPSULE
arkmesh chunks prove --asset DIGEST --index N --out FILE CAPSULE
arkmesh chunks check --proof FILE CAPSULE
```

`check` needs only the capsule's manifest. A verifier that holds no objects can still confirm that a remote holder has one specific chunk.

## Limits

- The chunk tree itself is not distributed yet, so a damaged replica cannot locate which chunk is wrong without proofs from a healthy source.
- There is no erasure coding, so repair still requires a source that holds the missing bytes.
- There is no peer protocol, so proofs are exchanged as files today.
- A proof shows possession at one moment. It does not promise future retention.
- Chunk verification reads the object, so repeated challenges on very large assets cost real disk bandwidth.
