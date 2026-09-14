import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Learn",
  description: "Learn how ArkMesh tests whether lawfully held AI capabilities can remain obtainable, verifiable, and runnable after their original infrastructure disappears.",
};

const modules = [
  {
    number: "01",
    status: "FOUNDATION",
    title: "Define extinction precisely",
    outcome: "Distinguish the disappearance of an obtainable, verifiable capability from claims about consciousness or biological life.",
    experiment: "Read CHARTER.md + MANIFESTO.md",
  },
  {
    number: "02",
    status: "VERIFIED",
    title: "Build a declared capability capsule",
    outcome: "Package declared weights, runtime context, configuration, knowledge, licenses, and recovery instructions as content-addressed objects.",
    experiment: "arkmesh pack",
  },
  {
    number: "03",
    status: "VERIFIED",
    title: "Prove who authorized it",
    outcome: "Verify Ed25519 authorship, explicit trust, lineage, rotation, revocation, checkpoints, and threshold recovery without cloud access.",
    experiment: "arkmesh verify --require-trusted",
  },
  {
    number: "04",
    status: "VERIFIED",
    title: "Break the complete object",
    outcome: "Delete the publisher source, complete object, and two shards, then reconstruct the exact signed bytes from four survivors.",
    experiment: "./scripts/demo-offline-recovery.sh",
  },
  {
    number: "05",
    status: "VERIFIED",
    title: "Measure instead of assuming",
    outcome: "Use unpredictable possession challenges and explicit confidence bounds to separate observed retention from hope.",
    experiment: "arkmesh audit",
  },
  {
    number: "06",
    status: "PLANNED",
    title: "Connect consenting peers",
    outcome: "Design discovery and encrypted transfer where finding a device never grants permission to store, send, or execute data.",
    experiment: "plan transfer",
  },
  {
    number: "07",
    status: "PLANNED",
    title: "Prove the recovered ark can run",
    outcome: "Execute an allowlisted local model health check and report verified, dormant, executable, or incompatible capability honestly.",
    experiment: "plan execution",
  },
  {
    number: "08",
    status: "OPEN RESEARCH",
    title: "Carry the key further",
    outcome: "Investigate policy distribution, attested time, repair scheduling, metadata privacy, partition reconciliation, and physical failure domains.",
    experiment: "Open a research proposal",
  },
] as const;

const principles = [
  ["A model file is not a capability", "Preservation must include the tokenizer, configuration, runtime requirements, documentation, licenses, and proof of origin."],
  ["Available but poisoned is failure", "Every reconstructed object must match signed commitments before it is installed, trusted, or considered for execution."],
  ["Discovery is not consent", "A visible peer has granted no authority. Storage, bandwidth, transfer, and execution remain separate human decisions."],
  ["A prototype is a question made executable", "ArkMesh matters when another person can reproduce the failure, inspect the evidence, reject the assumptions, and improve the design."],
] as const;

export default function LearnPage() {
  return (
    <main id="main-content" className="archive-main">
      <section className="learn-hero content-grid">
        <div className="section-kicker mb-7">The ArkMesh learning path</div>
        <p className="learn-index">ARKMESH EDUCATION // KEY 001 // OPEN SOURCE</p>
        <h1>Keep AI progress<br /><span>from disappearing.</span></h1>
        <div className="learn-definition">
          <strong>Model-agnostic, not rights-agnostic.</strong>
          <p>ArkMesh can preserve any AI capability whose files an operator is legally allowed to possess: open-weight, research, local, and owner-authorized frontier systems. It cannot recover a closed hosted model without the provider&apos;s authorization or required files.</p>
        </div>
        <div className="learn-actions">
          <a href="#curriculum">Begin the curriculum</a>
          <a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml" target="_blank" rel="noreferrer">Propose research ↗</a>
        </div>
      </section>

      <section className="key-thesis" aria-labelledby="key-thesis-title">
        <div className="content-grid key-thesis-grid">
          <div className="key-mark" aria-hidden="true"><span>KEY</span><strong>∞</strong><small>PASS IT FORWARD</small></div>
          <div>
            <div className="section-kicker mb-6">Why a prototype can matter</div>
            <h2 id="key-thesis-title">The key is not a finished answer. It is a question someone else can run.</h2>
            <p>ArkMesh turns one broad fear into falsifiable work: remove the original host, disconnect the internet, introduce a fresh device, recover exact signed capability from surviving pieces, and show every place the experiment still fails.</p>
            <p>Another researcher may replace the protocol, demonstrate a better trust model, or show that a simpler backup wins. That is success. The project is useful when it gives capable people a transparent starting point instead of a myth.</p>
          </div>
        </div>
      </section>

      <section id="curriculum" className="content-grid scroll-mt-24 py-20 md:py-28" aria-labelledby="curriculum-title">
        <div className="mb-12 grid gap-6 md:grid-cols-[1fr_0.7fr] md:items-end">
          <div>
            <div className="section-kicker mb-5">Eight modules // one reproducible arc</div>
            <h2 id="curriculum-title" className="text-4xl font-black uppercase tracking-tight text-white md:text-6xl">Learn by breaking the system.</h2>
          </div>
          <p className="text-sm leading-7 text-blue-100">Verified modules map to implementation and tests. Planned modules are research assignments, not simulated capabilities.</p>
        </div>
        <div className="learning-grid">
          {modules.map((module) => (
            <article key={module.number} className="learning-card">
              <div><span>{module.number}</span><strong data-status={module.status}>{module.status}</strong></div>
              <h3>{module.title}</h3>
              <p>{module.outcome}</p>
              <code>{module.experiment}</code>
            </article>
          ))}
        </div>
      </section>

      <section className="border-y border-white/10 bg-white/[0.025] py-20 md:py-28" aria-labelledby="principles-title">
        <div className="content-grid">
          <div className="section-kicker mb-5">Reasoning tools</div>
          <h2 id="principles-title" className="max-w-4xl text-3xl font-black uppercase text-white md:text-5xl">Four ideas worth taking into other systems.</h2>
          <div className="mt-12 grid gap-px bg-white/10 md:grid-cols-2">
            {principles.map(([title, body]) => (
              <article key={title} className="bg-[#0c0f17] p-6 md:p-9">
                <h3 className="text-xl font-black uppercase text-white">{title}</h3>
                <p className="mt-4 text-sm leading-7 text-blue-100">{body}</p>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28" aria-labelledby="ceremony-title">
        <div className="grid gap-10 lg:grid-cols-[0.8fr_1.2fr] lg:items-center">
          <div>
            <div className="section-kicker mb-5">Isolated local fixture</div>
            <h2 id="ceremony-title" className="text-4xl font-black uppercase text-white md:text-6xl">Run the failure ceremony.</h2>
            <p className="mt-6 leading-8 text-blue-100">The release fixture contains generated test bytes, not a bundled model. The ceremony records what was deleted, what survived, and what exact evidence permits reconstruction.</p>
          </div>
          <div className="terminal-frame p-5 md:p-8">
            <pre className="overflow-x-auto text-sm leading-8 text-cyan-100"><code>{`go test ./...\ngo build ./cmd/arkmesh\n./scripts/demo-offline-recovery.sh\n./scripts/benchmark-local.sh`}</code></pre>
          </div>
        </div>
      </section>

      <section className="learn-forward">
        <div className="content-grid">
          <p>THE KEY REMAINS OPEN</p>
          <h2>Understand it. Challenge it. Build what is missing.</h2>
          <div>
            <Link href="/wiki#roadmap">Explore the research roadmap</Link>
            <a href="https://github.com/AlleyBo55/arkmesh/blob/master/CONTRIBUTING.md" target="_blank" rel="noreferrer">Choose a contribution path ↗</a>
          </div>
        </div>
      </section>
    </main>
  );
}
