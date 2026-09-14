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

`chunk_size` must be a power of two between 4096 bytes and 64 MiB. The default is 64 KiB. The chunk count is derived from the signed size, never taken from a proof.

A smaller default keeps possession challenges and repairs cheap, since both move one chunk rather than one file. The cost is more leaves per object, which matters only while a tree is held in memory: a 70 GB asset has about 1.1 million leaves at 64 KiB.

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

## Tree distribution

A holder can export every leaf digest of one object:

```json
{
  "schema_version": "arkmesh.chunk-tree/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "asset_sha256": "<digest>",
  "chunk_size": 65536,
  "chunk_count": 6,
  "chunk_root": "<digest>",
  "leaves": ["<digest>", "..."]
}
```

This file carries no signature and needs none. Verification recomputes the root from the leaves and compares it with the chunk root inside the signed manifest, so a forged or truncated tree is rejected. A tree is also rejected when the capsule ID, asset, chunk size, root, or leaf count disagrees with the signed manifest.

## Damage localization and repair

With an authenticated tree, a damaged replica reports exactly which chunks are wrong instead of only failing:

```bash
arkmesh chunks tree --asset DIGEST --out tree.json healthy.ark
arkmesh chunks scan --tree tree.json damaged.ark
arkmesh chunks repair --tree tree.json --source donor.bin damaged.ark
```

`scan` prints each damaged chunk index, byte offset, length, and whether the bytes were missing or mismatched, then exits nonzero.

`repair` verifies every donor chunk against the authenticated tree BEFORE writing it, so a hostile donor cannot use repair as an injection path. After writing, the object is truncated to its signed size and the whole capsule is verified again. A donor that supplies a wrong chunk is refused and no bytes are written for it.

## Parallel verification

Objects are verified concurrently with a bounded worker pool, capped at four workers. Results are collected per asset and reported in manifest order, so the surfaced error is identical to sequential verification. A regression test repeats verification on a capsule with two damaged assets and requires the same asset to be reported every time.

## Limits

- The tree file lists every leaf, so it grows with object size. Range limited tree exchange is future work.
- There is no erasure coding, so repair still requires a donor that holds the missing bytes.
- There is no peer protocol, so trees and proofs are exchanged as files today.
- A possession proof shows possession at one moment. It does not promise future retention, and nothing here measures retention over time.
- Chunk verification reads the object, so repeated challenges on very large assets cost real disk bandwidth.
