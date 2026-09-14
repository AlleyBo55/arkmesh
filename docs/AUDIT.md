# Sampled Retrievability Audits

Status: **v0alpha1, unstable**

A possession proof shows that one chunk existed at one moment. A sampled audit answers a stronger question with an explicit error bound: how much of this object can still be missing, given that every chunk we checked was correct?

## The statistic

Samples are drawn without replacement, so the exact distribution is hypergeometric rather than binomial:

```text
P(all k sampled chunks clean | d of n chunks bad)
  = C(n-d, k) / C(n, k)
  = product over i in [0, k) of (n-d-i) / (n-i)
```

Two directions follow from that one formula.

**Planning.** To rule out losing a fraction of an object unnoticed, sample until the miss probability falls below `1 - confidence`. For a large object and one percent tolerated damage, that is about 300 samples at 95 percent confidence and about 460 at 99 percent, whether the object is 100 MB or 100 GB.

**Reading a result.** After k clean samples, the largest damage level still statistically plausible is reported as a bound. The honest statement is not "the object is intact" but "at least X percent is present, with Y percent confidence".

The implementation is checked against exact rational arithmetic using `math/big` for every combination of n up to 20, so the floating point path is verified against ground truth rather than against itself.

## Unpredictable challenges

Sample indexes come from `crypto/rand`. A holder that kept only the chunks it was asked about last time would pass a predictable audit, so a seeded generator would quietly destroy the guarantee. A test requires two consecutive audits of a large object to draw different samples.

## What a passing audit does and does not mean

- A partial audit never reports full integrity. A test asserts that.
- An exhaustive audit (`--samples` equal to the chunk count) does prove zero damage, and says so differently.
- If any sampled chunk fails, no bound is claimed at all. Observed damage is reported and the retrievability statement is withheld.
- Chunk comparisons are made against an authenticated chunk tree, which was itself verified against the signed manifest.

## Commands

```bash
arkmesh chunks tree --asset DIGEST --out tree.json capsule.ark
arkmesh audit sample --tree tree.json --tolerance 0.01 --confidence 0.99 --log retention.log capsule.ark
arkmesh audit history --log retention.log
```

## Retention log

Each audit can append one line to a log:

```json
{
  "schema_version": "arkmesh.audit/v0alpha1",
  "recorded_at": "2026-09-14T06:20:00Z",
  "capsule_id": "sha256:<digest>",
  "asset_sha256": "<digest>",
  "chunk_count": 64,
  "sampled": 46,
  "failed": 0,
  "tolerance": 0.05,
  "confidence": 0.99,
  "max_bad_chunks": 2,
  "min_intact_fraction": 0.96875,
  "exhaustive": false
}
```

Lines are parsed strictly: unknown fields, duplicate fields, wrong schema, and malformed timestamps are all rejected.

This turns retention from a promise into a record: a history of checks with dates, sample counts, and outcomes.

## Limits

- The log records observed checks. It cannot prove that inconvenient checks were omitted, and the CLI says so on every history report.
- Timestamps come from the local clock and are not independently attested. Proving that time actually passed needs a verifiable delay function or a trusted timestamp authority.
- An audit describes the replica it read. It says nothing about other replicas.
- Sampling bounds assume the auditor chooses indexes the holder cannot predict. Publishing a challenge in advance voids the guarantee.
- The audit reads whole chunks, so a large tolerance and high confidence cost real bandwidth on very large assets.
- The log is not signed, so it is local operator evidence rather than a claim other operators must accept.
