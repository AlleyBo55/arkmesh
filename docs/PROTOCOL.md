# Capsule Protocol

Status: **v0alpha1, unstable, unsigned**

This document describes the implemented on-disk capsule envelope and records invariants for later peer replication. It is not yet a stable public protocol.

## Directory layout

```text
<capsule>/
├── manifest.json
└── objects/
    ├── <lowercase sha256>
    └── ...
```

Object paths are derived only from validated lowercase SHA-256 digests. Human-provided file names are display metadata and never become storage paths.

## Manifest

```json
{
  "schema_version": "arkmesh.capsule/v0alpha1",
  "capsule_id": "sha256:<digest>",
  "name": "field-assistant",
  "created_at": "2026-09-14T00:00:00Z",
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

The capsule ID is SHA-256 over the compact JSON representation of the manifest with `capsule_id` set to an empty string. Struct field order is fixed by the reference implementation for v0alpha1. This derivation is experimental and will be replaced or formally canonicalized before interoperability is promised.

## Asset roles

Roles are currently free-form lowercase identifiers supplied by the packer. Initial conventions:

- `model`: model weights
- `runtime`: executable runtime or source bundle
- `tokenizer`: tokenizer files
- `config`: model/runtime configuration
- `knowledge`: intentionally shared reference material
- `policy`: human-readable or machine-readable operating policy
- `license`: license and redistribution records
- `recovery`: offline operating and repair instructions

A role describes purpose; it grants no permission and does not make content safe to execute.

## Verification

A verifier must reject a capsule when:

- The schema version is unsupported.
- The capsule ID does not match the manifest body.
- The name or asset list is empty.
- An asset role or display name is empty.
- A digest is not exactly 64 lowercase hexadecimal characters.
- An object is absent, not a regular file, has the wrong size, or has the wrong digest.

Extra objects may exist but do not belong to the capsule unless listed in the manifest.

## Planned signed envelope

Peer replication must not ship until a later format defines:

- Ed25519 author and node identities
- Signature scope and canonical serialization
- Parent capsule IDs and update authority
- Trust-root import and invitation flow
- Revocation and key-rotation semantics
- License and provenance declarations
- Runtime compatibility and health-check declarations
- Chunking rules for resumable large-object transfer

## Planned network behavior

1. Device owner installs and starts a node locally.
2. Owner imports a trust-group invitation.
3. Authenticated peers advertise capsule IDs and verified object availability.
4. Receiver requests only missing chunks within configured quotas.
5. Receiver verifies every chunk and complete object before storage.
6. A capsule remains data until its manifest, author, policy, license, and runtime are approved.
7. Execution occurs only through an allowlisted local runtime.

No central tracker may be required for LAN operation. WAN discovery and relays, if added, must remain optional.
