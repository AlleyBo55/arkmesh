import Link from "next/link";
import { plannedCapabilities } from "@/lib/content";

const claims = [
  ["C1", "Exact local recovery", "Demonstrated", "Deletion ceremony + digest and signature verification"],
  ["C2", "Authenticated partial repair", "Demonstrated", "Merkle proof and hostile-donor rejection tests"],
  ["C3", "Long-term distributed continuity", "Open", "Requires peer, placement, retention, and physical experiments"],
  ["C4", "Executable capability recovery", "Open", "Requires runtime compatibility and offline health evidence"],
] as const;

export function ThesisFrontMatter() {
  return (
    <section className="thesis-front" aria-labelledby="thesis-abstract-title">
      <div className="thesis-paper-meta">
        <span>ARKMESH RESEARCH MONOGRAPH</span>
        <span>V0ALPHA1 · 2026</span>
        <span>STATUS · OPEN TO FALSIFICATION</span>
      </div>
      <div className="thesis-abstract">
        <div>
          <p className="thesis-label">Abstract</p>
          <h2 id="thesis-abstract-title">Verifiable offline preservation and recovery of executable AI capabilities</h2>
        </div>
        <p>ArkMesh investigates whether an AI capability that an operator may lawfully possess can remain obtainable, authentic, reconstructible, and eventually executable after the original publisher, host, wide-area network, or signing key becomes unavailable. The present prototype establishes a local cryptographic and erasure-recovery baseline. It does not yet establish distributed continuity.</p>
      </div>
      <div className="thesis-question">
        <p className="thesis-label">Primary research question</p>
        <blockquote>Can a fresh device recover and verify a complete authorized AI capability from consenting independent peers after every original infrastructure dependency has been removed?</blockquote>
      </div>
      <div className="thesis-hypotheses">
        <article><span>H1 · WORKING HYPOTHESIS</span><p>A content-addressed capsule, explicit authority model, authenticated erasure recovery, and consent-bounded peer protocol are sufficient to recover an executable capability without one permanent coordinator.</p></article>
        <article><span>H0 · NULL HYPOTHESIS</span><p>Ordinary backups or existing distribution systems provide equivalent continuity with lower complexity, or ArkMesh cannot preserve authenticity, consent, and executability under realistic failures.</p></article>
      </div>
      <div className="thesis-method-grid">
        <article><strong>Planned independent variables</strong><p>Loss of publisher, complete objects, shards, network access, keys, peers, and compatible hardware.</p></article>
        <article><strong>Dependent measures</strong><p>Recovery success, exact-byte verification, authorization outcome, time, overhead, and executable health.</p></article>
        <article><strong>Baselines</strong><p>Ordinary backup, shared folder, content distribution, and local model runner where each is applicable.</p></article>
        <article><strong>Failure criterion</strong><p>Any poisoned, unauthorized, incomplete, irreproducible, or silently overclaimed recovery counts as failure.</p></article>
      </div>
      <div className="thesis-claim-table" role="table" aria-label="ArkMesh claim and evidence status">
        <div className="thesis-claim-head" role="row"><span>Claim</span><span>Scope</span><span>Status</span><span>Evidence required</span></div>
        {claims.map(([id, scope, status, evidence]) => (
          <div key={id} role="row"><span>{id}</span><strong>{scope}</strong><em data-state={status}>{status}</em><p>{evidence}</p></div>
        ))}
      </div>
    </section>
  );
}

function OverviewFigure() {
  return (
    <figure className="wiki-figure system-model">
      <figcaption><span>Figure 1</span> Experimental system model and removed dependencies</figcaption>
      <div className="system-model-grid">
        <article className="system-origin"><small>REMOVED</small><strong>Original publisher</strong><span>source + WAN + hosted API</span></article>
        <div className="system-break" aria-hidden="true"><i />X<i /></div>
        <article className="system-peers"><small>CONSENTING SURVIVORS</small><strong>Peer B · Peer C · Offline media</strong><span>authenticated chunks and shards</span></article>
        <div className="system-arrow" aria-hidden="true">→</div>
        <article className="system-target"><small>RECOVERY SUBJECT</small><strong>Fresh device D</strong><span>verify first · execute later</span></article>
      </div>
      <p className="figure-note">The target claim is not that bytes exist somewhere. It is that a new operator can recover the exact authorized capability without silently importing a poisoned object or relying on the failed origin.</p>
    </figure>
  );
}

function MethodFigure() {
  return (
    <figure className="wiki-figure method-figure">
      <figcaption><span>Figure 2</span> Reproducible experimental method</figcaption>
      <ol>
        <li><span>01</span><strong>Construct</strong><p>Pack, sign, protect, and record the initial evidence.</p></li>
        <li><span>02</span><strong>Destroy</strong><p>Remove the source, complete object, and selected shards.</p></li>
        <li><span>03</span><strong>Recover</strong><p>Reconstruct into temporary storage from surviving evidence.</p></li>
        <li><span>04</span><strong>Adjudicate</strong><p>Accept only after integrity, authority, and policy checks pass.</p></li>
      </ol>
    </figure>
  );
}

function CapsuleFigure() {
  const layers = ["Recovery instructions", "Licenses + provenance", "Runtime requirements", "Tokenizer + configuration", "Knowledge + documentation", "Model or capability assets"];
  return (
    <figure className="wiki-figure capsule-anatomy">
      <figcaption><span>Figure 3</span> A capability is a dependency closure, not a weight file</figcaption>
      <div className="capsule-stack">
        {layers.map((layer, index) => <div key={layer}><span>{String(layers.length - index).padStart(2, "0")}</span><strong>{layer}</strong></div>)}
      </div>
      <aside><strong>Signed envelope</strong><p>Manifest, deterministic capsule ID, chunk commitments, signer evidence, lineage, and authority transitions bind the layers without claiming that their contents are safe or legally redistributable.</p></aside>
    </figure>
  );
}

function AuthorityFigure() {
  return (
    <figure className="wiki-figure authority-graph">
      <figcaption><span>Figure 4</span> Mutually exclusive authority paths for one direct lineage edge</figcaption>
      <div className="authority-parent"><span>CHECKPOINTED HEAD</span><strong>Parent P · signer K₀</strong></div>
      <div className="authority-paths">
        <article><span>A</span><strong>Same author</strong><p>K₀ signs child C.</p></article>
        <article><span>B</span><strong>Exact rotation</strong><p>K₀ authorizes exactly K₁ and C.</p></article>
        <article><span>C</span><strong>Threshold recovery</strong><p>Distinct approved custodians authorize P → C.</p></article>
      </div>
      <div className="authority-child"><strong>Accepted child C</strong><span>one valid path · local revocation still applies</span></div>
    </figure>
  );
}

function RepairFigure() {
  return (
    <figure className="wiki-figure repair-figure">
      <figcaption><span>Figure 5</span> Damage localization and donorless reconstruction</figcaption>
      <div className="shard-visual">
        {["S0", "S1", "S2", "S3", "S4", "S5"].map((shard, index) => <span key={shard} data-lost={index === 1 || index === 4}>{shard}<small>{index === 1 || index === 4 ? "LOST" : "VALID"}</small></span>)}
        <i>4 / 6</i><strong>temporary reconstruction</strong><b>SHA-256 + chunk root + signature</b><em>ACCEPT</em>
      </div>
      <p className="figure-note">Erasure coding restores availability, not trust. Reconstructed bytes remain unaccepted until they match the commitments already covered by the signed manifest.</p>
    </figure>
  );
}

function AuditFigure() {
  return (
    <figure className="wiki-figure audit-figure">
      <figcaption><span>Figure 6</span> Sampling supports a bounded statement, not certainty</figcaption>
      <div className="audit-population" aria-label="Illustrative chunk population with sampled positions">
        {Array.from({ length: 48 }, (_, index) => <i key={index} data-sampled={[2, 7, 11, 18, 24, 31, 38, 44].includes(index)} />)}
      </div>
      <div className="audit-inference"><span>Observed</span><strong>8 unpredictable clean samples</strong><i>⇒</i><span>Permitted conclusion</span><strong>A stated confidence bound on missing fraction</strong><b>≠ proof of future retention</b></div>
    </figure>
  );
}

function SecurityFigure() {
  return (
    <figure className="wiki-figure boundary-figure">
      <figcaption><span>Figure 7</span> Trust boundaries remain explicit</figcaption>
      <div className="boundary-outer"><span>UNTRUSTED ENVIRONMENT · peers · network · capsule contents · clocks</span><div className="boundary-device"><span>OPERATOR-CONTROLLED DEVICE</span><div className="boundary-core"><strong>ARKMESH VERIFIER</strong><p>strict parse → integrity → authority → local policy → publication</p></div></div></div>
      <p className="figure-note">A compromised local operating system remains inside the trusted computing boundary and can replace policy or verifier state. ArkMesh does not claim protection from that attacker.</p>
    </figure>
  );
}

function EvidenceFigure() {
  const rows = [
    ["Integrity", "Changed bytes are rejected", "Unit + hostile fixture", "Content is safe"],
    ["Authority", "An accepted key authorized this edge", "Signature + policy tests", "Signer legal identity"],
    ["Recovery", "Four valid shards reconstruct exact bytes", "Destructive local ceremony", "Distributed survival"],
    ["Audit", "Observed samples bound plausible loss", "Exact statistical calculation", "Future availability"],
  ];
  return (
    <figure className="wiki-figure evidence-matrix">
      <figcaption><span>Figure 8</span> Claim discipline</figcaption>
      <div><p>Area</p><p>What evidence supports</p><p>Current method</p><p>What it does not support</p></div>
      {rows.map((row) => <div key={row[0]}>{row.map((cell, index) => index === 0 ? <strong key={cell}>{cell}</strong> : <p key={cell}>{cell}</p>)}</div>)}
    </figure>
  );
}

function RoadmapFigure() {
  return (
    <figure className="wiki-figure thesis-roadmap">
      <figcaption><span>Figure 9</span> Dependency-ordered research route</figcaption>
      <ol>
        <li className="is-complete"><span>BASE</span><strong>Verified local core</strong><p>Implemented evidence</p></li>
        {plannedCapabilities.map((item, index) => <li key={item.id}><span>{String(index + 1).padStart(2, "0")}</span><strong>{item.title}</strong><p>{index === plannedCapabilities.length - 1 ? "Independent summit" : "Acceptance evidence required"}</p></li>)}
      </ol>
    </figure>
  );
}

function ScopeFigure() {
  return (
    <figure className="wiki-figure scope-figure">
      <figcaption><span>Figure 10</span> Model scope is broad; authorization remains narrow</figcaption>
      <div><article><strong>Technically in scope</strong><p>Open-weight, research, local, and owner-authorized frontier systems with obtainable files.</p></article><article><strong>Not granted by ArkMesh</strong><p>Provider access, redistribution rights, legal identity, model safety, compatible hardware, or permission to execute.</p></article></div>
    </figure>
  );
}

export function WikiFigure({ section }: { section: string }) {
  switch (section) {
    case "overview": return <OverviewFigure />;
    case "quickstart": return <MethodFigure />;
    case "capsules": return <CapsuleFigure />;
    case "authority": return <AuthorityFigure />;
    case "repair": return <RepairFigure />;
    case "audits": return <AuditFigure />;
    case "security": return <SecurityFigure />;
    case "evidence": return <EvidenceFigure />;
    case "roadmap": return <RoadmapFigure />;
    case "faq": return <ScopeFigure />;
    default: return null;
  }
}

export function ResearchReferences() {
  return (
    <section className="wiki-section research-references" aria-labelledby="references-title">
      <div className="section-kicker mb-5">Selected foundations</div>
      <h2 id="references-title">Research lineage</h2>
      <p>ArkMesh combines established primitives but does not inherit their guarantees automatically. Each adaptation requires its own threat model and experiment.</p>
      <ol>
        <li><span>[1]</span><p>Reed, I. S. and Solomon, G. (1960). “Polynomial Codes over Certain Finite Fields.” <a href="https://doi.org/10.1137/0108018">SIAM Journal</a>.</p></li>
        <li><span>[2]</span><p>Merkle, R. C. (1988). “A Digital Signature Based on a Conventional Encryption Function.” <a href="https://doi.org/10.1007/3-540-48184-2_32">CRYPTO</a>.</p></li>
        <li><span>[3]</span><p>Juels, A. and Kaliski, B. S. (2007). “PORs: Proofs of Retrievability for Large Files.” <a href="https://doi.org/10.1145/1315245.1315317">ACM CCS</a>.</p></li>
        <li><span>[4]</span><p>ArkMesh protocol, threat model, experiments, and security invariants are maintained with the reference implementation.</p></li>
      </ol>
      <div className="mt-7 flex flex-wrap gap-3"><Link href="/wiki#evidence">Inspect claim discipline</Link><a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml">Challenge the thesis ↗</a></div>
    </section>
  );
}
