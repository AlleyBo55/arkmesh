export type CapabilityStatus = "verified" | "planned";

export type Capability = {
  id: string;
  command: string;
  title: string;
  status: CapabilityStatus;
  summary: string;
  acceptance: string;
};

export const implementedCapabilities: Capability[] = [
  {
    id: "capsules",
    command: "status capsules",
    title: "Content-addressed capability capsules",
    status: "verified",
    summary: "Pack models, runtime material, knowledge, licenses, and recovery instructions as named assets whose object paths come only from validated SHA-256 digests.",
    acceptance: "The manifest records each asset's role, name, size, digest, and optional chunk commitment; the same canonical manifest derives the same capsule ID.",
  },
  {
    id: "inspection",
    command: "status inspection",
    title: "Manifest inspection and strict parsing",
    status: "verified",
    summary: "Inspect capsule metadata without executing its contents, and reject malformed manifests before object paths or trust decisions are used.",
    acceptance: "Unknown fields, duplicate fields, invalid hashes, unsafe object references, and inconsistent metadata fail closed.",
  },
  {
    id: "identity",
    command: "status identity",
    title: "Offline identity and signatures",
    status: "verified",
    summary: "Create local Ed25519 author identities, sign a capsule without copying private keys into it, and verify signatures without a network service.",
    acceptance: "Verification distinguishes trusted, unknown, invalid, and unsigned states without treating a valid key as proof of legal identity or safety.",
  },
  {
    id: "verification",
    command: "status verification",
    title: "Deterministic integrity verification",
    status: "verified",
    summary: "Verify independent objects concurrently while preserving manifest-order error reporting, so performance does not make failures nondeterministic.",
    acceptance: "Missing files, changed bytes, wrong sizes, invalid signatures, and malformed metadata identify the same failing asset on repeated runs.",
  },
  {
    id: "chunks",
    command: "status chunks",
    title: "Merkle chunk commitments",
    status: "verified",
    summary: "Commit a Merkle root for every chunked asset, prove possession of one exact chunk, and locate damaged ranges against an authenticated tree.",
    acceptance: "A proof is bound to the chunk bytes, index, tree size, asset, and signed capsule; forged proofs and forged trees are rejected.",
  },
  {
    id: "repair",
    command: "status repair",
    title: "Authenticated targeted repair",
    status: "verified",
    summary: "Request only damaged chunks and verify every donor chunk against the signed tree before writing it into a replica.",
    acceptance: "Unverified donor bytes are never written, and the complete repaired capsule must pass size, digest, chunk-root, and signature verification.",
  },
  {
    id: "erasure",
    command: "status erasure",
    title: "Donorless erasure reconstruction",
    status: "verified",
    summary: "Protect an object with streaming Reed Solomon shards and reconstruct it after losses even when no surviving holder has the complete bytes.",
    acceptance: "Four of six configured shards recover the exact object; insufficient, forged, or inconsistent shards fail without publishing output.",
  },
  {
    id: "audit",
    command: "status audit",
    title: "Sampled retrievability and history",
    status: "verified",
    summary: "Use unpredictable chunk samples to bound the largest plausible missing fraction, then retain an audit history of checks that actually occurred.",
    acceptance: "Reports state sample size and confidence, suppress clean bounds after observed damage, and never turn local timestamps into claims of attested time.",
  },
  {
    id: "lineage",
    command: "status lineage",
    title: "Signed lineage and key authority",
    status: "verified",
    summary: "Bind one child to one parent, verify same-author updates, authorize one exact replacement key, and apply explicit local revocation policy.",
    acceptance: "Wrong parents, reused transitions, unauthorized descendants, and revoked signers are rejected; valid signatures do not bypass local policy.",
  },
  {
    id: "checkpoints",
    command: "status checkpoints",
    title: "Local checkpoints and threshold recovery",
    status: "verified",
    summary: "Pin one accepted lineage head, advance through one verified direct edge, and recover from a lost author key through distinct approvals.",
    acceptance: "Rollback against the retained checkpoint, checkpoint mismatch, duplicate approvals, revoked custodians, and recovery below the configured threshold fail closed.",
  },
];

export const plannedCapabilities: Capability[] = [
  {
    id: "bootstrap",
    command: "plan bootstrap",
    title: "Fresh-device trust and consent",
    status: "planned",
    summary: "Let a new device create or import identity, accept a scoped invitation, choose a capsule, and approve storage, bandwidth, expiry, and recovery limits.",
    acceptance: "A fresh device reaches an explicit ready state without ambient trust, remote installation, hidden defaults, or treating discovery as consent.",
  },
  {
    id: "transfer",
    command: "plan transfer",
    title: "Authenticated resumable peer transfer",
    status: "planned",
    summary: "Discover invited peers with bounded metadata disclosure, authenticate them, and resume encrypted chunk and shard transfer under exact authorization.",
    acceptance: "Two-node tests reject unauthorized, expired, replayed, over-quota, wrong-capsule, and hostile-donor requests before atomic publication.",
  },
  {
    id: "placement",
    command: "plan placement",
    title: "Failure-domain placement and repair",
    status: "planned",
    summary: "Place shards across independent physical risks, expose recovery margin, detect degradation, and request bounded authenticated repair.",
    acceptance: "Operators can see correlated-risk concentration, expected recoverability, repair cost, and every write before approving placement or repair.",
  },
  {
    id: "retention",
    command: "plan retention",
    title: "Longitudinal retention and attested time",
    status: "planned",
    summary: "Turn isolated possession checks into repeated unpredictable observations whose timing cannot be silently rewritten by one local clock.",
    acceptance: "Reports distinguish observed checks from missing evidence, bind events to independently verifiable time, and never infer retention across unaudited intervals.",
  },
  {
    id: "execution",
    command: "plan execution",
    title: "Runtime compatibility and local health",
    status: "planned",
    summary: "Determine whether a recovered capsule is executable on the available hardware and run an allowlisted deterministic health check without cloud authentication.",
    acceptance: "ArkMesh reports verified, executable, dormant, or incompatible states without automatic execution or claims of behavioral identity.",
  },
  {
    id: "reconciliation",
    command: "plan reconciliation",
    title: "Partition authority and reconciliation",
    status: "planned",
    summary: "Distribute scoped authority policy, exchange retained heads after partitions, preserve divergent authorized descendants, and require deliberate reconciliation.",
    acceptance: "Independent devices reject policy rollback and silent last-writer-wins replacement while retaining both valid branches and their ancestry.",
  },
  {
    id: "experiment",
    command: "plan experiment",
    title: "Independent continuity experiment",
    status: "planned",
    summary: "Remove the publisher, disconnect WAN access, introduce a fresh device, recover from consenting peers, verify the capsule, and run its local health check.",
    acceptance: "An outside researcher reproduces the complete physical experiment, publishes raw evidence, and compares ArkMesh against ordinary backups and shared folders.",
  },
];

export const benchmarkRows = [
  ["SIGNED PACK", "0.125 s"],
  ["TRUSTED VERIFY", "0.047 s"],
  ["4 + 2 PROTECT", "0.211 s"],
  ["2-LOSS RECOVER", "0.186 s"],
  ["99% AUDIT", "0.013 s"],
] as const;

export type WikiSection = {
  id: string;
  index: string;
  title: string;
  lead: string;
  paragraphs: string[];
  commands?: string[];
  links?: { label: string; href: string }[];
};

export const wikiSections: WikiSection[] = [
  {
    id: "overview",
    index: "00",
    title: "System overview",
    lead: "What happens to a lawfully held AI capability when the internet, original publisher, or hosting organization disappears?",
    paragraphs: [
      "ArkMesh treats an AI capability as a dependency closure: model or program assets plus tokenizer data, configuration, runtime requirements, knowledge, licenses, operating instructions, provenance, and recovery evidence. Preserving only weights may preserve bytes while losing the ability to interpret, authorize, or execute them.",
      "The thesis decomposes continuity into four independently falsifiable properties: obtainability, integrity, authority, and executability. The current alpha demonstrates local integrity, authority, repair, donorless reconstruction, and sampled retrievability. It does not yet demonstrate peer distribution or executable recovery.",
      "The target system removes one permanent coordinator but does not remove human governance. Every device owner retains control of identity, storage, bandwidth, trust policy, and execution. Discovery provides reachability information, never consent.",
      "The decisive experiment removes the original publisher, disconnects WAN access, introduces a fresh device, and asks whether independent consenting holders can recover and run one exact authorized capsule. Until that physical experiment is independently reproduced, distributed continuity remains a hypothesis.",
    ],
    commands: ["status", "evidence", "roadmap"],
  },
  {
    id: "quickstart",
    index: "01",
    title: "Quick start",
    lead: "Build and exercise the Go reference CLI without a model download.",
    paragraphs: [
      "Go 1.22 or newer is required. The reference core uses only the Go standard library. Use generated fixtures or assets whose licenses permit the intended storage, testing, and redistribution; the repository intentionally bundles no model.",
      "A valid reproduction records the exact revision, environment, commands, raw output, and discrepancy. Running the happy path alone is insufficient: the ceremony also requires recovery below threshold to fail without publishing an object.",
      "The quick start establishes a local baseline, not the project thesis. It removes complete local objects but does not exercise peer discovery, encrypted transport, physical failure domains, long-term retention, or inference hardware.",
    ],
    commands: ["go test ./...", "go build ./cmd/arkmesh", "./scripts/demo-offline-recovery.sh"],
    links: [{ label: "README", href: "https://github.com/AlleyBo55/arkmesh#quick-start" }],
  },
  {
    id: "capsules",
    index: "02",
    title: "Capsules and manifests",
    lead: "Content addressing separates human metadata from the bytes that determine identity.",
    paragraphs: [
      "Object paths come only from validated SHA-256 digests. Human names and roles remain metadata, preventing an untrusted manifest from selecting arbitrary filesystem paths. The manifest commits names, roles, sizes, digests, optional chunk roots, ancestry, and a deterministic capsule ID.",
      "Canonical encoding is a security boundary because signatures are meaningful only when every verifier derives the same identity from the same logical manifest. Strict decoders therefore reject duplicate fields, unknown fields, malformed hashes, and ambiguous structures instead of normalizing them silently.",
      "Completeness is partly semantic. A valid capsule can still omit a required driver, undocumented prompt template, hardware assumption, or license condition. ArkMesh verifies declared contents; execution experiments and independent review must test whether the declaration is sufficient.",
      "A signature proves control of one specific private key over one capsule ID. It does not prove legal identity, factual truth, model safety, license compliance, or permission to execute.",
    ],
    commands: ["arkmesh pack", "arkmesh inspect", "arkmesh verify"],
    links: [{ label: "Capsule protocol", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/PROTOCOL.md" }],
  },
  {
    id: "authority",
    index: "03",
    title: "Authority and continuity",
    lead: "Trust is imported explicitly and every exceptional authority path is scoped to one exact transition.",
    paragraphs: [
      "Trust is imported by each operator; it is not discovered from the network. Normal updates retain same-author authority. Planned rotation authorizes one exact child key and transition, while local revocation can reject a cryptographically valid but no longer accepted signer.",
      "Threshold recovery is intentionally distinct from rotation. It requires a retained checkpointed parent, one explicit local policy, distinct detached approvals, and the exact parent-to-child edge. Duplicate or revoked approvals do not count.",
      "A checkpoint provides local rollback resistance by pinning one accepted capsule and signer. It is not global consensus, and an attacker who can replace the checkpoint file remains inside the trusted device boundary.",
      "Disconnected communities can still produce two valid descendants. The future reconciliation problem is therefore governance as well as cryptography: preserve both histories, expose authority evidence, and require a deliberate human choice.",
    ],
    commands: ["arkmesh checkpoint", "arkmesh recovery", "plan reconciliation"],
    links: [
      { label: "Authority", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/AUTHORITY.md" },
      { label: "Recovery", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/RECOVERY.md" },
    ],
  },
  {
    id: "repair",
    index: "04",
    title: "Chunks, repair, and erasure",
    lead: "Authenticated chunks localize damage; erasure shards reconstruct bytes when no full donor survives.",
    paragraphs: [
      "Merkle proofs bind a chunk's bytes, index, tree size, asset, and signed capsule commitment. This allows a verifier to authenticate one range without reading the entire object and to localize damage against an independently retained tree.",
      "Targeted repair treats donors as hostile. Each candidate chunk is verified before any write, only damaged ranges are replaced, the result is truncated to the signed size, and the complete capsule is verified again.",
      "Streaming Reed Solomon protection uses a systematic GF(2^8) code to trade storage overhead for tolerance of missing shards. Coding restores enough information to reconstruct bytes; it does not establish their identity or authority.",
      "Reconstructed output therefore remains temporary until size, SHA-256 digest, Merkle root, signature, and active local policy agree. Loss beyond parity, forged shards, and inconsistent shard sets must fail without publication.",
    ],
    commands: ["arkmesh chunks", "arkmesh erasure", "run demo"],
    links: [
      { label: "Chunk commitments", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/CHUNKS.md" },
      { label: "Erasure coding", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/ERASURE.md" },
    ],
  },
  {
    id: "audits",
    index: "05",
    title: "Retrievability audits",
    lead: "A sample yields a confidence bound, not an unsupported claim that an entire replica is intact.",
    paragraphs: [
      "Challenges use cryptographic randomness and sample chunks without replacement. After a clean sample, ArkMesh inverts the exact hypergeometric model to report the largest missing fraction still compatible with the observation at a stated confidence level.",
      "The conclusion is deliberately asymmetric: a sampled failure is direct evidence of damage, while clean samples only reduce the range of plausible unseen loss. A sample is never presented as whole-object verification.",
      "Audit history records checks that actually occurred using local timestamps. It cannot prove omitted checks did not occur, that the local clock was honest, or that bytes remained available between observations.",
      "Longitudinal retention requires unpredictable scheduling across independent holders plus time evidence that one operator cannot silently rewrite. That extension remains planned.",
    ],
    commands: ["arkmesh audit sample", "arkmesh audit history", "plan retention"],
    links: [{ label: "Audit specification", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/AUDIT.md" }],
  },
  {
    id: "security",
    index: "06",
    title: "Security boundaries",
    lead: "Authenticity before availability: an available but poisoned capsule is a failure.",
    paragraphs: [
      "The adversary may control peers, transport messages, capsule contents, donor bytes, stale lineage, and local wall-clock claims. Authentication narrows who sent a message; it does not make the peer honest or grant consent to consume resources.",
      "The verifier applies ordered gates: strict parsing, path safety, size and digest integrity, chunk commitments, signature validity, lineage authority, revocation, checkpoint policy, and only then atomic publication. Availability never bypasses authenticity.",
      "Device owners control identity, storage, compute, bandwidth, trust policy, removal, and execution. ArkMesh must never scan, exploit, conceal, install remotely, treat discovery as consent, or execute arbitrary received plugins.",
      "A compromised local operating system remains outside the protection claim. Such an attacker can replace verifier binaries, policy files, checkpoints, or displayed results. Hardware roots and independently administered policy are future research, not implied guarantees.",
    ],
    commands: ["security invariants", "threat model"],
    links: [
      { label: "Threat model", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/THREAT_MODEL.md" },
      { label: "Security invariants", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/SECURITY_INVARIANTS.md" },
    ],
  },
  {
    id: "evidence",
    index: "07",
    title: "Evidence and benchmarks",
    lead: "Claims are tied to reproducible commands, raw observations, and named limits.",
    paragraphs: [
      "The release ceremony is a destructive controlled experiment: it creates a signed fixture, protects it with four data and two parity shards, deletes the publisher source, complete object, and two shards, then reconstructs and verifies the exact bytes. A three-shard attempt is the negative control and must leave no object.",
      "Benchmarks publish fixture size, environment, commands, storage overhead, and raw observations. The current Apple M5 measurements are a local implementation baseline, not a peer-network comparison or universal performance estimate.",
      "Evidence is versioned with the implementation because protocol claims can regress. A credible result identifies the exact revision and distinguishes cryptographic validity, policy acceptance, recoverability, and executability.",
      "The decisive external-validity test remains undone: independent operators on varied physical devices must compare ArkMesh with ordinary backups, shared folders, and existing distribution tools under the same failure schedule.",
    ],
    commands: ["evidence", "benchmark", "run demo"],
    links: [
      { label: "Demo procedure", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/DEMO.md" },
      { label: "Benchmark data", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/BENCHMARKS.md" },
      { label: "Claims audit", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/CLAIMS_AUDIT.md" },
    ],
  },
  {
    id: "roadmap",
    index: "08",
    title: "Roadmap",
    lead: "Seven dependency-ordered experiments separate verified local recovery from independently reproduced executable continuity.",
    paragraphs: [
      "The roadmap begins from the implemented local core, not from organizational scale. A fresh device must first establish explicit trust and consent; only then can authenticated transfer, failure-domain placement, longitudinal evidence, and execution be evaluated.",
      "Authority distribution and branch reconciliation move later because they become necessary when independent communities remain disconnected long enough to evolve divergent but valid histories. They do not block the first two-node recovery experiment.",
      "The route ends with an integrated physical experiment rather than another protocol feature: remove the publisher, disconnect WAN access, recover on a fresh device, verify authority, run a local health check, and publish raw results against simpler baselines.",
      "Support for different model families is a cross-cutting constraint. Every milestone must remain format-agnostic while refusing artifacts the operator is not legally allowed to possess or redistribute.",
    ],
    commands: plannedCapabilities.map((item) => item.command),
  },
  {
    id: "faq",
    index: "09",
    title: "FAQ and limits",
    lead: "ArkMesh is preservation infrastructure, not a cloud-model clone, autonomous replicator, or consciousness claim.",
    paragraphs: [
      "ArkMesh does not bundle a model, grant access to proprietary weights, or convert a hosted API into a preservable system. It can package only files an operator may lawfully possess. Model-agnostic format support is not rights-agnostic acquisition.",
      "The prototype does not currently discover peers, transfer over a network, schedule placement, attest time independently, or run inference. It cannot survive without sufficient reconstructible data, compatible hardware, energy, and operators who continue to consent.",
      "Cryptographic provenance does not establish truth, safety, legal identity, license compliance, or behavioral continuity. A recovered model may be authentic and still be harmful, obsolete, incompatible, or inappropriate to execute.",
      "If ordinary backups or existing distribution systems perform equally well under the same experiment with lower cost and risk, the simpler method wins. The research program is designed to permit that conclusion.",
    ],
    commands: ["status", "roadmap", "github"],
  },
];

export const allCapabilities = [...implementedCapabilities, ...plannedCapabilities];
