# Contributing to ArkMesh

ArkMesh is a research prototype for preserving AI capabilities that an operator is legally allowed to possess, store, and run. It is model-agnostic, but it does not override licenses, access controls, or provider rights. Do not upload model weights, private datasets, credentials, or other restricted material to issues or pull requests.

The project values reproducible counterexamples as highly as successful implementations. A result that narrows or disproves a claim is useful research.

## Choose an entry path

### Reproduce an experiment

Run the reference checks on a different operating system or hardware configuration:

```bash
go test ./...
go build ./cmd/arkmesh
./scripts/demo-offline-recovery.sh
./scripts/benchmark-local.sh
```

Open a [reproduction report](https://github.com/AlleyBo55/arkmesh/issues/new?template=reproduction-report.yml) with the exact revision, environment, commands, output, and any discrepancy. Do not summarize a failure when the raw error can be included safely.

### Propose research

Use the [research proposal](https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml) for prior art, threat hypotheses, protocol questions, model-license analysis, experiment design, or a falsifiable alternative to an existing assumption.

A useful proposal names:

1. The question or claim.
2. Why the current evidence is insufficient.
3. The smallest experiment that could change our confidence.
4. The expected observation for success and failure.
5. Safety, consent, privacy, and licensing constraints.

### Review and attack

High-value review areas include strict manifest parsing, chunk commitments, erasure reconstruction, lineage authority, local checkpoints, threshold recovery, sampled audits, and experimental methodology.

For a public correctness issue, include a minimal fixture or reproduction. Report vulnerabilities privately through [GitHub Security Advisories](https://github.com/AlleyBo55/arkmesh/security/advisories/new) instead of opening a public issue.

### Implement a bounded issue

Start from an accepted issue with a named acceptance boundary. Keep the change focused. Networking and execution proposals must preserve explicit consent, local control, inspectability, and offline operation.

Before requesting review:

```bash
gofmt -w <changed-go-files>
go test ./...
go vet ./...
go build ./cmd/arkmesh
```

For web changes, use Node.js 20.9 or newer:

```bash
cd web
npm ci
npm run lint
npm run typecheck
npm run build
```

Do not remove tests, weaken assertions, suppress errors, or present planned behavior as implemented. Add a regression test for every corrected failure when practical.

## Pull requests

A pull request should state:

- The problem and affected user or experiment
- The behavioral change, if any
- The evidence added or updated
- Commands actually run and their results
- Known limits and deliberately excluded work
- Any protocol, security, compatibility, or licensing consequence

Public formats, trust rules, identity behavior, and protocol changes are one-way doors. Discuss them before implementation and update the relevant protocol and threat-model documents with the code.

## AI-assisted contributions

AI-assisted work is welcome under the same standards as human work. The contributor remains responsible for understanding the change, protecting private data, verifying every claim, and ensuring the tool received no authority to access systems, replicate software, spend resources, or act beyond explicit approval.

## Research posture

Read the [research charter](CHARTER.md), [threat model](docs/THREAT_MODEL.md), and [security invariants](docs/SECURITY_INVARIANTS.md) before proposing peer networking, authority, recovery, or model execution work. Authenticity comes before availability. Discovery is not consent. A valid signature is not proof that content is safe.
