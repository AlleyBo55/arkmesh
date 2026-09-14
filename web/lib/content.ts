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
    title: "Signed content-addressed capsules",
    status: "verified",
    summary: "Pack assets, commit exact digests and chunk roots, derive a stable capsule ID, and verify explicit Ed25519 authorship offline.",
    acceptance: "Changed, missing, unsigned when required, or ambiguously parsed objects are rejected.",
  },
  {
    id: "lineage",
    command: "status lineage",
    title: "Lineage and authority",
    status: "verified",
    summary: "Bind one child to one parent, authorize exact key rotations, apply local revocations, pin accepted heads, and recover authority through distinct threshold approvals.",
    acceptance: "Wrong parents, reused transitions, revoked signers, rollback, and insufficient approvals fail closed.",
  },
  {
    id: "repair",
    command: "status repair",
    title: "Authenticated repair",
    status: "verified",
    summary: "Prove possession, locate damaged chunks, verify donor bytes before writing, and rebuild a lost object from Reed Solomon shards.",
    acceptance: "The recovered object is installed only after size, digest, and signed chunk-root verification.",
  },
  {
    id: "audit",
    command: "status audit",
    title: "Sampled retrievability",
    status: "verified",
    summary: "Draw unpredictable chunk challenges and report an exact hypergeometric confidence bound instead of claiming unobserved integrity.",
    acceptance: "Observed damage suppresses the bound; exhaustive checks are reported separately.",
  },
];

export const plannedCapabilities: Capability[] = [
  {
    id: "authority",
    command: "plan authority",
    title: "Organizational policy signing",
    status: "planned",
    summary: "Distribute revocation and threshold-recovery policy as signed, scoped organizational records rather than unmanaged local files.",
    acceptance: "Independent devices resolve the same policy version, reject rollback, and never expand signer authority implicitly.",
  },
  {
    id: "peers",
    command: "plan peers",
    title: "Peer discovery and encrypted transfer",
    status: "planned",
    summary: "Discover only explicitly invited LAN peers, authenticate them, minimize metadata disclosure, and resume encrypted object and shard transfer.",
    acceptance: "No scanning, ambient trust, central tracker, remote installation, or unapproved resource use.",
  },
  {
    id: "placement",
    command: "plan placement",
    title: "Shard placement and repair scheduling",
    status: "planned",
    summary: "Place shards across independent failure domains, measure recoverability, and schedule bounded repairs when redundancy falls below policy.",
    acceptance: "Operators see replica health, repair cost, and failure-domain concentration before approving writes.",
  },
  {
    id: "retention",
    command: "plan retention",
    title: "Measured retention over time",
    status: "planned",
    summary: "Turn isolated possession checks into a longitudinal record of repeated, unpredictable retrievability observations.",
    acceptance: "Reports distinguish checks that occurred from intervals with no evidence and never infer omitted audits.",
  },
  {
    id: "time",
    command: "plan time",
    title: "Independently attested time",
    status: "planned",
    summary: "Bind retention events to evidence stronger than a local wall clock through a documented timestamp or delay mechanism.",
    acceptance: "A local operator cannot silently rewrite when an audit occurred without invalidating its attestation.",
  },
  {
    id: "inference",
    command: "plan inference",
    title: "Local model inference",
    status: "planned",
    summary: "Run a declared open model and deterministic health prompt through an allowlisted local runtime without cloud authentication.",
    acceptance: "ArkMesh reports verified, executable, dormant, and incompatible states without claiming behavioral identity.",
  },
  {
    id: "replay",
    command: "plan replay",
    title: "Replay agreement and branch reconciliation",
    status: "planned",
    summary: "Exchange retained heads, preserve divergent authorized descendants, and require deliberate human reconciliation after partitions.",
    acceptance: "No silent last-writer-wins replacement, ancestry loss, or automatic private-memory merge.",
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
    lead: "ArkMesh preserves verifiable AI capability bundles when the original host and central network are unavailable.",
    paragraphs: [
      "A capsule can contain model weights, tokenizer data, configuration, runtime requirements, knowledge, licenses, and recovery instructions. The signed manifest commits to every asset and its authenticated chunk tree.",
      "The current alpha proves local integrity, authority, repair, donorless reconstruction, and sampled retrievability. Networking and inference remain planned and are never presented as completed behavior.",
    ],
    commands: ["status", "evidence", "roadmap"],
  },
  {
    id: "quickstart",
    index: "01",
    title: "Quick start",
    lead: "Build and exercise the Go reference CLI without a model download.",
    paragraphs: [
      "Go 1.22 or newer is required. ArkMesh uses only the Go standard library. Use only capsule assets whose licenses permit your intended storage and redistribution.",
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
      "Object paths come only from validated SHA-256 digests. The manifest records roles, names, sizes, chunk sizes, roots, ancestry, and a deterministic capsule ID.",
      "Strict decoders reject duplicate and unknown fields. A signature proves control of one specific key, not legal identity, safety, truth, or license compliance.",
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
      "Normal updates retain same-author authority. Planned key rotation authorizes one exact next signer through a parent-signed edge. Local revocation can reject compromised keys.",
      "Threshold recovery requires distinct detached approvals over one exact parent-to-child transition. A retained checkpoint rejects rollback to an older valid capsule.",
    ],
    commands: ["arkmesh checkpoint", "arkmesh recovery", "plan authority"],
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
      "Merkle proofs bind a chunk, index, tree size, asset, and capsule. Repair verifies every donor chunk before writing only damaged ranges.",
      "Streaming Reed Solomon protection uses a systematic GF(2^8) code. Reconstructed output stays temporary until it matches the signed size, digest, and chunk root.",
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
      "Challenges use cryptographic randomness and sample without replacement. After clean samples, ArkMesh reports the largest still-plausible missing fraction under the exact hypergeometric model.",
      "Retention logs record observed checks using local time. They cannot prove omitted checks did not occur, or that the local clock was honest.",
    ],
    commands: ["arkmesh audit sample", "arkmesh audit history", "plan time"],
    links: [{ label: "Audit specification", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/AUDIT.md" }],
  },
  {
    id: "security",
    index: "06",
    title: "Security boundaries",
    lead: "Authenticity before availability: an available but poisoned capsule is a failure.",
    paragraphs: [
      "Peers remain untrusted even after authentication. Device owners control storage, compute, network access, and removal. ArkMesh must never scan, exploit, conceal, install remotely, or treat discovery as consent.",
      "A compromised local operating system remains outside the protection boundary. Received files are data, not permission to execute arbitrary plugins.",
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
      "The release ceremony deletes the publisher source, complete object, and two of six shards before reconstructing the exact signed bytes. A three-shard attempt must fail without writing an object.",
      "The first benchmark uses a 64 MiB fixture on one Apple M5. It is a local measurement floor, not a peer-network or universal performance claim.",
    ],
    commands: ["evidence", "benchmark", "run demo"],
    links: [
      { label: "Demo procedure", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/DEMO.md" },
      { label: "Benchmark data", href: "https://github.com/AlleyBo55/arkmesh/blob/master/docs/BENCHMARKS.md" },
    ],
  },
  {
    id: "roadmap",
    index: "08",
    title: "Roadmap",
    lead: "Seven missing capabilities separate local cryptographic recovery from the project thesis.",
    paragraphs: [
      "The next milestones are organizational policy distribution, invited encrypted peers, failure-domain-aware shard placement, longitudinal retention, attested time, local inference, and deliberate partition reconciliation.",
      "Each remains marked PLANNED until its acceptance criteria are exercised by reproducible tests. The interface never simulates completion.",
    ],
    commands: plannedCapabilities.map((item) => item.command),
  },
  {
    id: "faq",
    index: "09",
    title: "FAQ and limits",
    lead: "ArkMesh is preservation infrastructure, not a cloud-model clone, autonomous replicator, or consciousness claim.",
    paragraphs: [
      "It does not bundle a model. It does not currently discover peers or run inference. It cannot preserve proprietary weights that operators do not possess or survive without compatible hardware, energy, and reconstructible data.",
      "If ordinary backups perform equally well under the target experiment, the simpler method wins. The research claim is deliberately falsifiable.",
    ],
    commands: ["status", "roadmap", "github"],
  },
];

export const allCapabilities = [...implementedCapabilities, ...plannedCapabilities];
