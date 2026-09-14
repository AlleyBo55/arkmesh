# Erasure Coding

Status: **v0alpha1, unstable**

Chunk repair needs a donor that still holds the missing bytes. Erasure coding removes that requirement: any `data` shards out of `data + parity` rebuild the object, even when no surviving copy of it exists anywhere.

## Codec

Systematic Reed Solomon over GF(2^8), implemented with the Go standard library only.

- Reducing polynomial `x^8 + x^4 + x^3 + x^2 + 1` (0x11d), generator 2
- Encoding matrix is a Vandermonde matrix multiplied by the inverse of its top square block, so the first `data` rows are the identity and data shards are stored unchanged
- Reconstruction inverts the submatrix of surviving rows by Gauss-Jordan elimination
- Up to 255 shards total

Two properties are proven by test rather than assumed:

1. The generator is primitive. Every nonzero field element appears exactly once in the exponent table. A non primitive generator produces a codec that encodes happily and then fails to reconstruct.
2. Recovery works for every recoverable loss pattern. For a 4 + 3 layout, the test removes all combinations of one, two, and three shards and requires exact recovery each time.

## Plan format

```json
{
  "schema_version": "arkmesh.erasure/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "asset_sha256": "<digest>",
  "asset_size": 200000,
  "data_shards": 4,
  "parity_shards": 2,
  "shard_size": 65536,
  "block_size": 65536,
  "shard_digests": ["<digest>", "..."]
}
```

Shards are plain files named `shard-000`, `shard-001`, and so on, so they can be copied to separate disks, machines, or offline media.

## Trust chain

The plan carries no signature and needs none.

Shard digests exclude corrupt shards cheaply, but they are not the trust anchor. Recovery reconstructs into a temporary file and then verifies it against the digest, size, and chunk root in the **signed manifest** before renaming it into place. A hostile plan with self consistent digests still fails, and a test proves exactly that: forged shard bytes plus a rewritten plan digest are rejected with `recovered bytes rejected`, and no object is written.

Recovery also refuses to write anything when fewer than `data` usable shards survive.

## Memory behavior

Encoding and recovery stream one stripe at a time, so memory stays near `(data + parity) * block_size` regardless of asset size. A 4 + 2 layout with 64 KiB blocks holds under 400 KiB of shard buffers while protecting an arbitrarily large model.

## Commands

```bash
arkmesh erasure protect --asset DIGEST --data 4 --parity 2 \
  --shards ./shards --plan plan.json capsule.ark

arkmesh erasure recover --plan plan.json --shards ./shards capsule.ark
```

`recover` reads the manifest without requiring the object, so it works when the object is missing entirely.

## Limits

- Parity costs storage. A 4 + 2 layout adds 50 percent overhead and tolerates two lost shards.
- Shards only add resilience if they are stored separately. Six shards on one disk survive nothing that disk does not.
- There is no peer protocol, so shards are placed and gathered by the operator today.
- Padding means shard files are slightly larger than `asset_size / data`.
- Shard placement, repair scheduling, and replica health metrics are not implemented.
- Erasure coding protects bytes. It does not protect the signing key, the checkpoint, or the recovery policy.
