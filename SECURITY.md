# Security Policy

## Reporting a vulnerability

Please do not open a public issue for an unpatched vulnerability.

Send a private report through GitHub Security Advisories for this repository. Include:

- Affected commit or release
- Preconditions and trust boundary
- Minimal reproduction
- Expected and observed behavior
- Security impact
- Any suggested mitigation

Do not include real private keys, private capsule contents, credentials, or personal data.

## Response principles

ArkMesh maintainers should preserve evidence, reproduce the report, identify the violated security invariant, add a failing regression test, and fix the narrow root cause before expanding scope.

A cryptographic format change requires explicit compatibility analysis. Existing capsule IDs and signed payloads must not change silently.

## Current scope

ArkMesh is a v0alpha1 research prototype. Security review currently covers local capsule integrity, signatures, lineage, planned rotation, local revocation, checkpoints, and threshold recovery.

Networking, peer authentication, encrypted transfer, runtime sandboxing, and model execution are not implemented. Reports about those future components are design feedback rather than vulnerabilities in deployed code.

No software is guaranteed free of defects. Responsible reports that improve ArkMesh's evidence and failure handling are welcome.
