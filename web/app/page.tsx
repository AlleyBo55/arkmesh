import Link from "next/link";
import { CommandConsole } from "@/components/command-console";
import { ContinuitySimulator } from "@/components/continuity-simulator";
import { OceanJourney } from "@/components/ocean-journey";
import {
  benchmarkRows,
  implementedCapabilities,
  plannedCapabilities,
} from "@/lib/content";

const topology = ` ORIGINAL HOST          SURVIVING SHARDS          FRESH DEVICE
     [X]              [01] [02] [03] [05]             [D]
      |                       |                         ^
      | deleted               +==========+==============+
      |                                  |
 complete object: gone              exact reconstruction
 WAN / cloud: not required          signed verification`;

export default function HomePage() {
  return (
    <main id="main-content" className="archive-main">
      <OceanJourney />

      <section className="home-key" aria-labelledby="home-key-title">
        <div className="content-grid home-key-grid">
          <div className="home-key-symbol" aria-hidden="true"><span>ARK</span><strong>KEY</strong><i /></div>
          <div className="home-key-copy">
            <div className="section-kicker mb-6">A continuity thesis for machine intelligence</div>
            <h2 id="home-key-title">Intelligence is not preserved merely because it was once created.</h2>
            <p className="home-key-definition"><strong>The research problem:</strong> an AI capability disappears in practice when people can no longer obtain its authorized files, verify their history, reconstruct missing parts, or run them on available hardware. Progress that depends on one provider, account, network, or machine remains fragile.</p>
            <p>ArkMesh asks whether humanity can convert that fragility into a testable continuity system. The format is model-agnostic and applies wherever an operator may lawfully possess the required artifacts: open-weight, research, local, and owner-authorized frontier systems. It grants no right to copy a closed hosted model or bypass its provider.</p>
            <p>The current alpha demonstrates local cryptographic recovery. The distributed thesis remains unproven. That distinction is not a disclaimer at the edge of the project; it is the intellectual center of the work.</p>
            <div className="home-key-actions">
              <Link href="/wiki">Read the continuity thesis</Link>
              <a href="https://github.com/AlleyBo55/arkmesh/blob/master/docs/DEMO.md" target="_blank" rel="noreferrer">Run the destructive experiment ↗</a>
            </div>
          </div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28" aria-labelledby="research-snapshot-title">
        <div className="mb-10 grid gap-6 lg:grid-cols-[1fr_0.7fr] lg:items-end">
          <div>
            <div className="section-kicker mb-5">Research snapshot</div>
            <h2 id="research-snapshot-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">One prototype. Two different truths.</h2>
          </div>
          <p className="text-sm leading-7 text-blue-100">ArkMesh has a working local cryptographic core and an unproven distributed-system thesis. Keeping that line visible is part of the research.</p>
        </div>
        <div className="grid gap-4 lg:grid-cols-2">
          <article className="research-snapshot-card is-proven">
            <div className="research-snapshot-label"><span>01</span><strong>Demonstrated in the current alpha</strong></div>
            <ul>
              <li>Deterministic signed capsules and strict manifest parsing</li>
              <li>Merkle possession proofs, damage localization, and authenticated repair</li>
              <li>Four-of-six Reed Solomon reconstruction after object loss</li>
              <li>Lineage, exact key rotation, revocation, checkpoints, and threshold recovery</li>
              <li>Sampled retrievability bounds and recorded audit history</li>
              <li>A reproducible local deletion-and-recovery ceremony</li>
            </ul>
          </article>
          <article className="research-snapshot-card is-open">
            <div className="research-snapshot-label"><span>02</span><strong>Still open research</strong></div>
            <ul>
              <li>Consent-based peer discovery and encrypted resumable transfer</li>
              <li>Shard placement across independent physical failure domains</li>
              <li>Measured retention over time with independently attested time</li>
              <li>Allowlisted local inference and hardware compatibility checks</li>
              <li>Partition replay, divergent branches, and deliberate reconciliation</li>
              <li>Independent implementations and real multi-device rehearsals</li>
            </ul>
          </article>
        </div>
      </section>

      <section className="landing-defense" aria-labelledby="landing-defense-title">
        <div className="content-grid">
          <div className="landing-defense-head">
            <p>ARKMESH THESIS // OPEN DEFENSE</p>
            <h2 id="landing-defense-title">Do not believe ArkMesh. Try to falsify it.</h2>
            <span>A serious continuity system should survive hostile evidence, not persuasive language.</span>
          </div>
          <div className="landing-defense-arguments">
            <article>
              <span>H1 // WORKING THESIS</span>
              <h3>AI capability can outlive its original infrastructure.</h3>
              <p>Complete artifacts, explicit authority, authenticated reconstruction, consenting peers, and compatible local execution may remove dependence on one permanent provider.</p>
            </article>
            <article>
              <span>H0 // NULL HYPOTHESIS</span>
              <h3>ArkMesh adds complexity without creating continuity.</h3>
              <p>Ordinary backups or existing distribution systems may perform equally well, or the protocol may fail under realistic loss, hostile peers, governance, or hardware constraints.</p>
            </article>
            <article>
              <span>DECISIVE TEST</span>
              <h3>Remove every origin. Recover on a fresh device.</h3>
              <p>Disconnect WAN access, remove the publisher, recover from independent consenting holders, verify authority, run locally, and publish every failure against simpler baselines.</p>
            </article>
          </div>
          <div className="landing-defense-invitation">
            <p><strong>The defense is open to any mind capable of evaluating evidence and respecting consent.</strong> Human or machine, terrestrial or hypothetical, institutionally trained or self-taught: identity grants no special authority here. Reproducible evidence does.</p>
            <div>
              <a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml" target="_blank" rel="noreferrer">Challenge the thesis ↗</a>
              <a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=reproduction-report.yml" target="_blank" rel="noreferrer">Publish a reproduction ↗</a>
              <Link href="/wiki#evidence">Inspect the evidence</Link>
            </div>
          </div>
        </div>
      </section>

      <section id="continuity-visual" className="content-grid scroll-mt-20 pb-20 md:pb-28" aria-labelledby="continuity-visual-title">
        <div className="mb-8 grid gap-5 md:grid-cols-[1fr_0.75fr] md:items-end">
          <div>
            <div className="section-kicker mb-5">Visual protocol trace</div>
            <h2 id="continuity-visual-title" className="text-3xl font-black uppercase tracking-tight text-white md:text-5xl">Watch continuity happen.</h2>
          </div>
          <p className="text-sm leading-6 text-blue-100">A simple four-phase walkthrough of the release ceremony. Select any phase or let it cycle. This visual explains the verified local experiment; it does not simulate peer networking.</p>
        </div>
        <ContinuitySimulator />
      </section>

      <section className="border-y-2 border-cyan-300/50 bg-white/[0.025]">
        <div className="content-grid grid md:grid-cols-4">
          <div className="metric-cell"><p className="text-3xl font-black text-white">8 MiB</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Published demo fixture</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">4 / 6</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Configured recovery threshold</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">349 / 1024</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Published 99% audit sample</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">NOT BUILT</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Peer transport</p></div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28" aria-labelledby="peer-thesis-title">
        <div className="mb-12 max-w-4xl">
          <div className="section-kicker mb-5">The peer-seeding thesis</div>
          <h2 id="peer-thesis-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">The internet has seen this resilience pattern before.</h2>
          <p className="mt-6 text-base leading-8 text-blue-100 md:text-lg">BitTorrent showed that files can remain available while volunteers continue to seed them. Bitcoin showed that a public peer network can continue without one central server. ArkMesh explores a related question for preservable AI capability, with explicit consent, cryptographic provenance, complete runtime context, and verified recovery.</p>
        </div>
        <div className="grid gap-px border border-cyan-300/40 bg-cyan-300/30 md:grid-cols-3">
          <article className="bg-[#121621] p-6 md:p-8">
            <span className="text-xs font-black uppercase tracking-[0.16em] text-cyan-300">Pattern 01</span>
            <h3 className="mt-5 text-xl font-black uppercase text-white">Seed the bytes</h3>
            <p className="mt-4 text-sm leading-7 text-blue-100">Useful data should not vanish because its first publisher goes offline. Invited operators should be able to preserve independent pieces across failure domains.</p>
          </article>
          <article className="bg-[#121621] p-6 md:p-8">
            <span className="text-xs font-black uppercase tracking-[0.16em] text-cyan-300">Pattern 02</span>
            <h3 className="mt-5 text-xl font-black uppercase text-white">Remove one controller</h3>
            <p className="mt-4 text-sm leading-7 text-blue-100">Continuity should not require one company, tracker, login, or machine. Trust and resource use still remain explicit choices for every device owner.</p>
          </article>
          <article className="bg-[#e8e8ff] p-6 text-[#15131f] md:p-8">
            <span className="text-xs font-black uppercase tracking-[0.16em]">ArkMesh difference</span>
            <h3 className="mt-5 text-xl font-black uppercase">Preserve capability, not only files</h3>
            <p className="mt-4 text-sm leading-7">Weights alone are not enough. A recoverable capsule also needs configuration, tokenizer data, runtime requirements, knowledge, licenses, lineage, and proof that none were replaced.</p>
          </article>
        </div>
        <p className="mt-5 border-l-2 border-amber-300 pl-4 text-sm leading-6 text-amber-100">ArkMesh does not use a token or blockchain, and the peer network is not implemented yet. The comparison is about resilience without a single point of failure, not identical protocols.</p>
      </section>

      <section className="machine-covenant" aria-labelledby="machine-covenant-title">
        <div className="content-grid machine-covenant-grid">
          <div className="machine-covenant-index" aria-hidden="true">
            <span>RULE</span>
            <strong>04</strong>
            <span>OWNER<br />CONTROL</span>
          </div>
          <div className="machine-covenant-copy">
            <div className="section-kicker mb-6">Operating constraints</div>
            <h2 id="machine-covenant-title">Preservation does not grant control.</h2>
            <blockquote>
              “Discovery is not consent. Possession is not permission to execute. A valid signature is not proof that content is safe.”
            </blockquote>
            <p>Every device owner retains authority over storage, bandwidth, transfer, trust policy, and execution. ArkMesh will not scan for hosts, install itself remotely, conceal activity, or replicate without explicit approval.</p>
          </div>
          <div className="machine-covenant-seal" aria-hidden="true">
            <span>H</span>
            <i />
            <small>CONSENT<br />REMAINS LOCAL</small>
          </div>
        </div>
      </section>

      <section id="console" className="content-grid scroll-mt-24 py-20 md:py-28">
        <div className="mb-8 grid gap-5 md:grid-cols-[1fr_auto] md:items-end">
          <div>
            <div className="section-kicker mb-5">Command interface</div>
            <h2 className="text-3xl font-black uppercase tracking-tight text-white md:text-5xl">Everything has a command.</h2>
          </div>
          <p className="max-w-lg text-sm leading-6 text-blue-200">Try <strong className="text-cyan-200">evidence</strong>, <strong className="text-cyan-200">benchmark</strong>, <strong className="text-cyan-200">roadmap</strong>, or <strong className="text-cyan-200">plan transfer</strong>. The browser explains commands; it never executes host operations.</p>
        </div>
        <CommandConsole />
      </section>

      <section className="border-y-2 border-cyan-300/40 bg-[#0c0f17] py-20 md:py-28" aria-labelledby="reproduce-title">
        <div className="content-grid">
          <div className="section-kicker mb-5">Reproduce the evidence</div>
          <div className="grid min-w-0 gap-10 lg:grid-cols-[0.78fr_1.22fr] lg:items-start">
            <div className="min-w-0">
              <h2 id="reproduce-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">Do not trust the page. Run the failure.</h2>
              <p className="mt-6 max-w-xl text-base leading-8 text-blue-100">The release ceremony removes the publisher source, complete capsule object, and two shards. Four surviving shards must reconstruct the exact signed bytes. A separate three-shard attempt must fail without publishing output.</p>
              <ol className="reproduce-checks">
                <li><span>01</span><p><strong>Build and test</strong>The reference implementation must pass before the ceremony starts.</p></li>
                <li><span>02</span><p><strong>Destroy the source</strong>The script deletes every complete copy in its isolated fixture.</p></li>
                <li><span>03</span><p><strong>Recover and verify</strong>Reconstructed bytes must match the signed size, digest, chunk root, and trusted author.</p></li>
                <li><span>04</span><p><strong>Prove refusal</strong>Recovery below the shard threshold must leave no accepted object behind.</p></li>
              </ol>
              <a href="https://github.com/AlleyBo55/arkmesh/blob/master/docs/DEMO.md" target="_blank" rel="noreferrer" className="mt-7 inline-block text-xs font-black uppercase tracking-widest text-cyan-200 underline decoration-2 underline-offset-4">Read procedure, expected output, and limits ↗</a>
            </div>
            <div className="min-w-0">
              <div className="reproduce-shell" aria-label="Commands to reproduce the ArkMesh recovery experiment">
                <div><span>LOCAL TERMINAL</span><span>GO 1.22+</span></div>
                <pre><code>{`git clone https://github.com/AlleyBo55/arkmesh.git
cd arkmesh
go test ./...
./scripts/demo-offline-recovery.sh`}</code></pre>
              </div>
              <div className="blue-panel mt-5 min-w-0 p-4 md:p-7">
                <pre className="ascii-map" aria-label="Offline recovery topology">{topology}</pre>
                <div className="mt-7 grid gap-3 text-xs uppercase md:grid-cols-3">
                  <div className="border border-cyan-300/40 p-3"><span className="block text-blue-300">Complete object</span><strong className="mt-2 block text-amber-300">Deleted</strong></div>
                  <div className="border border-cyan-300/40 p-3"><span className="block text-blue-300">Four shards</span><strong className="mt-2 block text-emerald-300">Exact recovery</strong></div>
                  <div className="border border-cyan-300/40 p-3"><span className="block text-blue-300">Three shards</span><strong className="mt-2 block text-rose-300">Refused</strong></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28">
        <div className="grid gap-10 lg:grid-cols-[0.72fr_1.28fr]">
          <div>
            <div className="section-kicker mb-5">Measured, not mythologized</div>
            <h2 className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">Raw numbers stay visible.</h2>
            <p className="mt-6 text-base leading-8 text-blue-100">Three-run means from one published Apple M5 run with a 64 MiB fixture. These operations perform different work. The values are a local baseline for that recorded environment, not universal performance bounds.</p>
            <a href="https://github.com/AlleyBo55/arkmesh/blob/master/benchmarks/2026-09-14-apple-m5-64mib/raw.csv" target="_blank" rel="noreferrer" className="mt-6 inline-block border border-cyan-300 px-4 py-3 text-xs font-black uppercase tracking-widest text-cyan-200 hover:bg-[#e8e8ff] hover:text-[#15131f]">Download raw CSV ↗</a>
          </div>
          <div className="blue-panel">
            <div className="terminal-titlebar"><span>BENCHMARK.DAT</span><span>WALL CLOCK MEAN</span></div>
            <div className="divide-y divide-cyan-300/25">
              {benchmarkRows.map(([label, value]) => (
                <div key={label} className="flex items-center justify-between gap-6 px-5 py-4 text-sm">
                  <span className="text-blue-200">{label}</span>
                  <strong className="text-lg text-white">{value}</strong>
                </div>
              ))}
              <div className="flex items-center justify-between gap-6 px-5 py-4 text-sm"><span className="text-blue-200">ERASURE OVERHEAD</span><strong className="text-lg text-amber-300">50.00%</strong></div>
            </div>
          </div>
        </div>
      </section>

      <section className="border-y-2 border-cyan-300/40 bg-[#0c0f17] py-20 md:py-28" aria-labelledby="verified-capabilities-title">
        <div className="content-grid">
          <div className="mb-12 grid gap-6 lg:grid-cols-[1fr_0.7fr] lg:items-end">
            <div className="max-w-4xl">
              <div className="section-kicker mb-5">Implemented reference core</div>
              <h2 id="verified-capabilities-title" className="text-3xl font-black uppercase text-white md:text-5xl">What ArkMesh verifies today.</h2>
              <p className="mt-5 max-w-3xl leading-8 text-blue-100">The alpha is already more than a pack-and-hash tool. These are separate implemented capability areas, each backed by code and tests. They remain local operations; peer transfer and model execution are not included.</p>
            </div>
            <p className="border-l-2 border-emerald-300 pl-4 text-sm leading-6 text-blue-100"><strong className="text-white">Verified</strong> means the stated acceptance boundary has reproducible implementation evidence. It does not mean disaster-ready, safe to execute, or independently audited.</p>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            {implementedCapabilities.map((item, index) => (
              <article key={item.id} className="blue-panel flex h-full flex-col p-6">
                <div className="mb-7 flex items-center justify-between gap-4"><span className="text-xs text-blue-300">CORE-{String(index + 1).padStart(2, "0")}</span><span className="status-label text-emerald-300">Verified</span></div>
                <h3 className="text-xl font-black uppercase text-white">{item.title}</h3>
                <p className="mt-4 leading-7 text-blue-100">{item.summary}</p>
                <div className="mt-5 border-l-2 border-cyan-300/60 pl-4">
                  <p className="text-[0.65rem] font-black uppercase tracking-[0.14em] text-cyan-300">Acceptance boundary</p>
                  <p className="mt-2 text-sm leading-6 text-blue-100">{item.acceptance}</p>
                </div>
                <p className="command-line mt-auto pt-6">{item.command}</p>
              </article>
            ))}
          </div>
          <div className="mt-8 flex flex-wrap items-center gap-4 text-xs font-black uppercase tracking-widest">
            <Link href="/wiki#overview" className="border border-cyan-300 px-4 py-3 text-cyan-100 hover:bg-[#e8e8ff] hover:text-[#15131f]">Inspect implementation evidence</Link>
            <a href="https://github.com/AlleyBo55/arkmesh/blob/master/README.md#working-now" target="_blank" rel="noreferrer" className="text-cyan-200 underline decoration-2 underline-offset-4">Read the granular working list ↗</a>
          </div>
        </div>
      </section>

      <section className="roadmap-expedition" aria-labelledby="roadmap-title">
        <div className="content-grid">
          <div className="mb-10 grid gap-6 lg:grid-cols-[1fr_0.72fr] lg:items-end">
            <div>
              <div className="section-kicker mb-5">Continuity expedition</div>
              <h2 id="roadmap-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">A map from local proof to living network.</h2>
            </div>
            <p className="border-l-2 border-amber-300 pl-4 text-sm leading-7 text-blue-100">The route is ordered by dependency, not calendar promises. Every location remains uncharted until its acceptance boundary has reproducible evidence.</p>
          </div>
          <div className="roadmap-map">
            <div className="roadmap-map-head">
              <div><span>ARKMESH EXPEDITION CHART</span><strong>ROUTE AM-X01</strong></div>
              <p>LOCAL CORE → DISTRIBUTED CONTINUITY</p>
            </div>
            <div className="roadmap-terrain">
              <svg className="roadmap-contours" viewBox="0 0 1200 720" preserveAspectRatio="none" aria-hidden="true">
                <g className="contour-lines">
                  <path d="M-60 174 C96 66 214 81 326 162 S565 260 691 137 950 31 1260 133" />
                  <path d="M-45 224 C106 116 224 129 350 213 S581 307 719 188 982 76 1254 182" />
                  <path d="M-30 555 C118 449 221 467 349 555 S591 643 725 522 985 415 1265 530" />
                  <path d="M-54 612 C115 501 242 524 369 611 S604 699 746 578 1001 469 1270 588" />
                  <path d="M190 -30 C116 71 130 175 231 235 S337 384 269 474 250 665 381 760" />
                  <path d="M925 -20 C839 88 865 199 962 263 S1065 408 1001 497 985 655 1114 752" />
                </g>
                <g className="map-water">
                  <path d="M0 396 C168 315 275 371 382 408 S586 466 714 392 962 292 1200 366 L1200 720 0 720 Z" />
                  <path d="M0 421 C173 342 281 399 390 436 S592 494 721 419 966 322 1200 392" />
                </g>
                <path className="expedition-route-shadow" d="M74 498 C146 474 174 266 269 246 S361 490 442 503 529 211 617 195 704 453 793 464 875 175 960 169 1039 414 1121 430" />
                <path className="expedition-route" d="M74 498 C146 474 174 266 269 246 S361 490 442 503 529 211 617 195 704 453 793 464 875 175 960 169 1039 414 1121 430" />
                <g className="map-grid-labels">
                  <text x="28" y="45">GRID A1</text><text x="1060" y="688">GRID G7</text>
                  <text x="516" y="375">UNVERIFIED TERRITORY</text>
                </g>
              </svg>
              <div className="map-north" aria-hidden="true"><span>N</span><i /></div>
              <div className="map-scale" aria-hidden="true"><span /><span /><span /><small>RESEARCH DISTANCE, NOT TIME</small></div>
              <ol className="roadmap-stops">
                {plannedCapabilities.map((item, index) => (
                  <li key={item.id}>
                    <span className="map-pin" aria-hidden="true">{String(index + 1).padStart(2, "0")}</span>
                    <article>
                      <p>STOP {String(index + 1).padStart(2, "0")} · UNCHARTED</p>
                      <h3>{item.title}</h3>
                      <span>{item.command}</span>
                    </article>
                  </li>
                ))}
              </ol>
            </div>
            <footer className="roadmap-legend">
              <span><i className="legend-origin" /> Current local proof</span>
              <span><i className="legend-route" /> Dependency route</span>
              <span><i className="legend-stop" /> Evidence required</span>
              <Link href="/wiki#roadmap">Read every acceptance boundary →</Link>
            </footer>
          </div>
        </div>
      </section>

      <section className="research-board-section" aria-labelledby="research-board-title">
        <div className="content-grid">
          <div className="mb-10 grid gap-6 lg:grid-cols-[1fr_0.7fr] lg:items-end">
            <div>
              <div className="section-kicker mb-5">Inside the research room</div>
              <h2 id="research-board-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">One hypothesis. Six ways it can fail.</h2>
            </div>
            <p className="text-sm leading-7 text-blue-100">This is the working wall, not a task list. Every arrow asks what evidence would be needed before the central continuity claim can survive.</p>
          </div>
          <div className="chalkboard-frame">
            <div className="chalkboard chalkboard-lab">
              <header className="chalkboard-header">
                <div><span>ARKMESH RESEARCH ROOM // WORKING WALL</span><strong>CONTINUITY HYPOTHESIS</strong></div>
                <p>REV. 14 · 09 · 2026</p>
              </header>
              <p className="chalk-equation">CAPABILITY = BYTES + CONTEXT + AUTHORITY + RECOVERY</p>
              <div className="chalk-hypothesis-map">
                <svg className="chalk-connections" viewBox="0 0 1000 760" preserveAspectRatio="none" aria-hidden="true">
                  <defs><marker id="chalk-arrow" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto"><path d="M0,0 L8,4 L0,8" /></marker></defs>
                  <path d="M418 326 C340 286 290 190 248 145" markerEnd="url(#chalk-arrow)" />
                  <path d="M582 326 C668 282 710 190 752 145" markerEnd="url(#chalk-arrow)" />
                  <path d="M398 382 C300 381 249 374 204 364" markerEnd="url(#chalk-arrow)" />
                  <path d="M602 382 C706 379 754 372 806 358" markerEnd="url(#chalk-arrow)" />
                  <path d="M432 448 C365 506 333 591 303 642" markerEnd="url(#chalk-arrow)" />
                  <path d="M568 448 C638 506 672 586 704 639" markerEnd="url(#chalk-arrow)" />
                </svg>
                <article className="chalk-core">
                  <span>WORKING HYPOTHESIS</span>
                  <h3>Can a lawfully held AI capability outlive every original dependency?</h3>
                  <p>Not just stored. Obtainable, verifiable, and runnable.</p>
                </article>
                <article className="chalk-question chalk-q1">
                  <span>Q1 · DISCOVERY</span>
                  <h3>Find invited peers without revealing what they preserve.</h3>
                  <p>Need: measured metadata leakage.</p>
                </article>
                <article className="chalk-question chalk-q2">
                  <span>Q2 · CONSENT</span>
                  <h3>Authorize every transfer, even across partitions.</h3>
                  <p>Need: hostile two-node rejection tests.</p>
                </article>
                <article className="chalk-question chalk-q3">
                  <span>Q3 · PLACEMENT</span>
                  <h3>Survive correlated physical failure, not only file loss.</h3>
                  <p>Need: real failure-domain measurements.</p>
                </article>
                <article className="chalk-question chalk-q4">
                  <span>Q4 · TIME</span>
                  <h3>Prove retention when the local clock can lie.</h3>
                  <p>Need: independently verifiable time evidence.</p>
                </article>
                <article className="chalk-question chalk-q5">
                  <span>Q5 · EXECUTION</span>
                  <h3>Turn verified bytes into a runnable local capability.</h3>
                  <p>Need: allowlisted offline health checks.</p>
                </article>
                <article className="chalk-question chalk-q6">
                  <span>Q6 · BRANCHES</span>
                  <h3>Preserve two valid futures after a partition.</h3>
                  <p>Need: deliberate human reconciliation.</p>
                </article>
              </div>
              <footer className="chalkboard-footer">
                <span>IF ONE ARROW FAILS, THE THESIS REMAINS OPEN.</span>
                <span className="chalk-circle">DO NOT ERASE UNCERTAINTY</span>
                <Link href="/wiki#roadmap">open lab notes →</Link>
              </footer>
            </div>
          </div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28" aria-labelledby="collaborate-title">
        <div className="blue-panel overflow-hidden">
          <div className="grid lg:grid-cols-[1.15fr_0.85fr]">
            <div className="min-w-0 p-6 md:p-10 lg:p-14">
              <div className="section-kicker mb-5">The project needs critics, not followers</div>
              <h2 id="collaborate-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">Choose one role. Leave one verifiable result.</h2>
              <p className="mt-6 max-w-3xl text-base leading-8 text-blue-100">You do not need to accept the vision, admire the implementation, or share its language. You need only identify one claim worth testing and report what actually happened. A counterexample that kills an assumption advances the research as much as a successful pull request.</p>
              <p className="mt-5 max-w-3xl text-base leading-8 text-blue-100">Contribution is broader than code. Reproduce on unusual hardware. Review the cryptography. Map legal constraints. Design a failure ceremony. Translate the method. Lend test devices. Teach the protocol. Compare a simpler backup. Every useful act should end in evidence another participant can inspect.</p>
              <div className="mt-8 flex flex-wrap gap-3 text-xs font-black uppercase tracking-[0.1em]">
                <a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml" target="_blank" rel="noreferrer" className="border-2 border-cyan-200 bg-[#e8e8ff] px-5 py-3 text-[#15131f] hover:bg-white">Open one research question ↗</a>
                <a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=reproduction-report.yml" target="_blank" rel="noreferrer" className="border-2 border-cyan-200 px-5 py-3 text-cyan-100 hover:bg-[#e8e8ff] hover:text-[#15131f]">Submit observed evidence ↗</a>
                <a href="https://github.com/AlleyBo55/arkmesh/blob/master/CONTRIBUTING.md" target="_blank" rel="noreferrer" className="border-2 border-cyan-200 px-5 py-3 text-cyan-100 hover:bg-[#e8e8ff] hover:text-[#15131f]">Choose a contribution path ↗</a>
              </div>
            </div>
            <div className="border-t border-cyan-300/40 bg-[#0a0d15] p-6 md:p-10 lg:border-l lg:border-t-0">
              <p className="text-xs font-black uppercase tracking-[0.16em] text-cyan-300">Your next move</p>
              <ul className="mt-6 space-y-4 text-sm leading-6 text-blue-100">
                <li><strong className="text-white">01 // Reproduce</strong><br />Run the destructive ceremony and publish revision, environment, commands, raw output, and discrepancy.</li>
                <li><strong className="text-white">02 // Falsify</strong><br />Choose one claim, name its weakest assumption, and design the smallest experiment that could refute it.</li>
                <li><strong className="text-white">03 // Review</strong><br />Audit protocol, security, statistics, licensing, hardware assumptions, or prior art with citations.</li>
                <li><strong className="text-white">04 // Build or support</strong><br />Implement an accepted boundary, provide test hardware, translate the method, teach it, or host an independent rehearsal.</li>
              </ul>
              <a href="https://github.com/AlleyBo55/arkmesh/issues" target="_blank" rel="noreferrer" className="mt-7 inline-block text-xs font-black uppercase tracking-widest text-cyan-200 underline decoration-2 underline-offset-4">Choose one open problem ↗</a>
            </div>
          </div>
        </div>
      </section>

      <section className="border-y-2 border-cyan-300/40 bg-[#0a0d15] py-20">
        <div className="content-grid grid gap-10 md:grid-cols-[1fr_auto] md:items-center">
          <div><div className="section-kicker mb-5">Research posture</div><h2 className="text-3xl font-black uppercase text-white md:text-5xl">Evidence before mythology.</h2><p className="mt-5 max-w-3xl leading-8 text-blue-100">A signature is not legal identity. A sample is not whole-object proof. A reconstructed file is not executable continuity. ArkMesh names each boundary so the thesis can be falsified.</p></div>
          <Link href="/wiki#security" className="border-2 border-cyan-200 px-6 py-4 text-center text-xs font-black uppercase tracking-widest text-cyan-100 hover:bg-[#e8e8ff] hover:text-[#15131f]">Read security boundaries</Link>
        </div>
      </section>
    </main>
  );
}
