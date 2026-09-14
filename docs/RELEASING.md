# Releasing ArkMesh

ArkMesh releases are research snapshots. A tag does not upgrade an unimplemented property into a guarantee.

## Alpha release procedure

1. Start from the intended public default-branch commit with a clean worktree.
2. Run the security gate:

   ```bash
   ./scripts/security-gate.sh
   ```

3. Re-run the offline recovery ceremony:

   ```bash
   ./scripts/demo-offline-recovery.sh
   ```

4. Build the release archives and checksums:

   ```bash
   ./scripts/build-release.sh v0.1.0-alpha.1 dist
   ```

5. Create an annotated tag that names the tested commit:

   ```bash
   git tag -a v0.1.0-alpha.1 -m "ArkMesh v0.1.0-alpha.1"
   git push origin v0.1.0-alpha.1
   ```

6. Create the GitHub release from the tag. Copy the matching section of `CHANGELOG.md` into the release notes and attach every archive plus `SHA256SUMS` from `dist`.
7. Download one archive from GitHub, verify its checksum, run `arkmesh help`, and record any discrepancy before announcing the release.

## Artifact matrix

The builder produces statically linked archives for:

- macOS arm64
- macOS amd64
- Linux arm64
- Linux amd64

Each archive contains the `arkmesh` binary, README, and MIT license. `SHA256SUMS` authenticates accidental or transport corruption but is not a release signature. Signed release provenance remains future work.
