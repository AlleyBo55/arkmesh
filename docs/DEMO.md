# Offline recovery demonstration

This ceremony demonstrates the strongest implemented ArkMesh recovery claim without requiring a network or a full-object donor.

## Run it

Requirements: Go 1.22 or newer plus standard POSIX shell utilities.

```bash
./scripts/demo-offline-recovery.sh
```

To retain the generated identity, capsule, shards, plan, and recovered object for inspection:

```bash
ARKMESH_DEMO_DIR="$PWD/demo-evidence" ./scripts/demo-offline-recovery.sh
```

Use a new or empty directory. Identity material created by the demonstration is disposable and must not be used for real capsules.

## What the script proves

1. Builds ArkMesh locally.
2. Creates an 8 MiB fixture and Ed25519 signing identity.
3. Packs and signs the fixture with authenticated 256 KiB chunks.
4. Protects it as four data shards plus two parity shards.
5. Deletes the object and three shards in an isolated copy.
6. Requires recovery to fail because only three usable shards remain, and verifies that no object was written.
7. Deletes the publisher source, the real capsule object, and two shards.
8. Reconstructs the object from the four surviving shards.
9. Requires the recovered SHA-256 digest to equal both the original digest and the digest committed by the signed manifest.
10. Verifies the recovered capsule against the explicitly imported author identity.

A successful run ends with:

```text
PASS: donorless recovery rebuilt the exact signed object.
ArkMesh peer/network operations: not implemented
```

ArkMesh currently implements no peer or network transport, and the script invokes no network client directly. A machine with Go 1.22 installed can run the ceremony without WAN access. Go's own toolchain selection may attempt a download if the required toolchain is missing, so operators should install and cache it before an offline rehearsal. The ArkMesh Go implementation itself currently uses only the standard library.

## What it does not prove

- No peer discovery or transfer occurs. Shards are placed in local directories.
- It does not run a model or inference health check.
- It is not a physical multi-device experiment.
- It does not prove long-term retention.
- It does not establish the legal identity or trustworthiness of the signer.

Those limits are part of the release claim, not omitted failure cases.
