# Capsule Protocol

Status: **v0alpha1, unstable, optional signatures implemented**

This document describes the implemented on disk capsule envelope. It is not yet a stable public protocol.

## Directory layout

```text
<capsule>/
├── manifest.json
├── signature.json
├── authority.json
└── objects/
    ├── <lowercase sha256>
    └── ...
```

`signature.json` is optional so existing unsigned capsules remain valid for local integrity checks. `authority.json` appears only on a child signed by a new key after an exact parent-signed rotation. Object paths are derived only from validated lowercase SHA-256 digests. Human provided file names are display metadata and never become storage paths.

## Manifest

A root capsule omits `parent_capsule_id`. A descendant includes the exact capsule ID of its parent:

```json
{
  "schema_version": "arkmesh.capsule/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "name": "field-assistant",
  "created_at": "2026-09-14T00:00:00Z",
  "parent_capsule_id": "sha256:<parent-digest>",
  "assets": [
    {
      "role": "model",
      "name": "model.gguf",
      "size": 1234,
      "sha256": "<digest>"
    }
  ]
}
```

The capsule ID is SHA-256 over the compact JSON representation of the manifest with `capsule_id` set to an empty string. `parent_capsule_id` is omitted from root JSON, which preserves existing root capsule IDs. For descendants, the parent ID is included in the digest and therefore covered by the capsule signature.

Struct field order is fixed by the reference implementation for v0alpha1. This derivation remains experimental until an interoperability review formalizes canonical serialization.

## Lineage and update authority

A root has no declared parent. A child names one parent capsule ID, and the supplied parent must verify and match that exact ID. Parent and child must both have valid Ed25519 signatures.

Two update paths are implemented:

1. **Same author:** parent and child use the same signing key. No authority file is allowed.
2. **Planned key rotation:** the child uses a new key and includes `authority.json`, signed by the parent key, binding the exact parent ID, child ID, old signer, and new signer.

Verification reports:

- `root` when no parent is declared
- `parent_not_checked` when a parent is declared but not supplied
- `verified_same_author` when parent and child use the same valid signer
- `verified_key_rotation` when the parent key authorized the child's exact new signer

Strict lineage mode rejects a descendant when its parent was not supplied. A mismatched parent, invalid signature, missing authority, reused transition, revoked signer, or unauthorized child key is rejected.

Branches from one parent remain possible. ArkMesh does not select a winning branch or merge descendants. Planned rotation also requires access to the old private key before the child is created.

See [Key Authority and Revocation](AUTHORITY.md) for the exact transition payload, trust inheritance rule, local revocation format, and limitations.

## Asset roles

Roles are lowercase identifiers supplied by the packer. Initial conventions:

- `model`: model weights
- `runtime`: executable runtime or source bundle
- `tokenizer`: tokenizer files
- `config`: model or runtime configuration
- `knowledge`: intentionally shared reference material
- `policy`: human readable or machine readable operating policy
- `license`: license and redistribution records
- `recovery`: offline operating and repair instructions

A role describes purpose. It grants no permission and does not make content safe to execute.

## Author identity

ArkMesh creates an Ed25519 key pair in a new owner controlled directory:

```text
<identity>/
├── identity.json
└── identity.key
```

`identity.json` is public and may be copied into a local trust store. `identity.key` contains the private key, is written with owner only file permissions, and must never enter a capsule or source repository.

The public identity format is:

```json
{
  "schema_version": "arkmesh.identity/v0alpha1",
  "algorithm": "ed25519",
  "key_id": "ed25519:<sha256-of-public-key>",
  "public_key": "<standard-base64>"
}
```

The key ID is a fingerprint, not a person's legal identity. Operators establish trust by obtaining and checking `identity.json` through a channel they consider appropriate.

## Detached signature

A signed capsule adds `signature.json`:

```json
{
  "schema_version": "arkmesh.signature/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "signer": {
    "schema_version": "arkmesh.identity/v0alpha1",
    "algorithm": "ed25519",
    "key_id": "ed25519:<digest>",
    "public_key": "<standard-base64>"
  },
  "signature": "<standard-base64>"
}
```

Ed25519 signs these exact UTF-8 bytes:

```text
"arkmesh.capsule.signature/v0alpha1\n" + capsule_id + "\n"
```

The domain prefix prevents the signature from being reused as another message type. The capsule ID already binds the complete manifest, and each listed object is bound by its digest and size.

## Verification states

After object and manifest integrity checks, verification reports one of three successful states:

- `unsigned`: no signature exists and strict signature checks were not requested
- `valid_unknown_author`: the signature is valid, but the signer is absent from the supplied trust set
- `valid_trusted_author`: the signature is valid and exactly matches a supplied public identity

Verification fails for unknown or trailing manifest fields, malformed envelopes, unsupported schemas or algorithms, changed capsule IDs, invalid key fingerprints, invalid key lengths, symlink or nonregular objects, invalid signatures, or a missing trusted signer when strict trust is required.

A valid trusted signature proves control of the signing private key for that capsule ID. It does not prove content safety, factual accuracy, license compliance, or the human identity behind the key.

## Remaining authenticity work

Peer replication must not ship until later work defines:

- Emergency recovery authority for old-key loss
- Distribution and organizational signing of revocation policies
- Trust group invitation and import flow
- Replay resistant version rules
- License and provenance declarations
- Runtime compatibility and health checks
- Chunking rules for resumable large object transfer

## Planned network behavior

1. A device owner installs and starts a node locally.
2. The owner imports a trust group invitation.
3. Authenticated peers advertise capsule IDs and verified object availability.
4. The receiver requests only missing chunks within configured quotas.
5. The receiver verifies every chunk and complete object before storage.
6. A capsule remains data until its author, policy, license, and runtime are approved.
7. Execution occurs only through an allowlisted local runtime.

No central tracker may be required for LAN operation. WAN discovery and relays, if added, must remain optional.
