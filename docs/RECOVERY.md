# Threshold Emergency Recovery

Status: **v0alpha1, unstable**

Threshold recovery replaces a lost signing key without creating one universal recovery key. It is accepted only from the locally checkpointed parent and requires distinct approvals under an explicit local policy.

## Policy

A policy contains at least two public recovery identities and requires at least two approvals:

```json
{
  "schema_version": "arkmesh.recovery-policy/v0alpha1",
  "policy_id": "sha256:<digest>",
  "threshold": 2,
  "members": ["<sorted public identities>"]
}
```

The policy ID binds the threshold and sorted member key IDs. The policy is local operator policy. It is not global consensus.

## Independent approvals

Each custodian signs on a separate device with its own private identity. Private keys are never combined. Every approval binds:

- Policy ID
- Recovery action
- Exact parent capsule ID
- Exact child capsule ID
- Previous signer ID
- Exact next signer ID

The signed UTF-8 payload is:

```text
"arkmesh.capsule.recovery/v0alpha1\n" +
policy_id + "\n" +
action + "\n" +
parent_capsule_id + "\n" +
child_capsule_id + "\n" +
previous_signer_id + "\n" +
next_signer_id + "\n"
```

ArkMesh assembles distinct valid approvals into the child's `recovery.json`. Revoked recovery identities do not count toward the threshold.

## Workflow

1. Create several offline recovery identities on independently controlled devices.
2. Create a policy with a threshold of at least two.
3. Retain a trusted checkpoint for the current parent.
4. Pack a child signed by the replacement key using the recovery policy.
5. Have custodians produce detached approvals independently.
6. Assemble enough approvals into the child.
7. Advance the checkpoint with the same recovery policy.

Successful advancement reports `verified_threshold_recovery`.

Normal same-author lineage and planned old-key rotation remain separate. A child cannot contain both `authority.json` and `recovery.json`.

## Limits

Recovery fails safely when too few unrevoked custodians approve, evidence targets another edge, the parent does not match the checkpoint, or the policy differs.

This mechanism does not protect a recovery policy or checkpoint from a local attacker who can replace those policy files. It does not establish global identity, timestamps, branch consensus, or legal authority. Custodians must be independent in practice, not merely represented by different key files controlled by one person.
