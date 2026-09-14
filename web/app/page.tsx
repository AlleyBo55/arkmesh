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
    <main>
      <OceanJourney />

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
          <div className="metric-cell"><p className="text-3xl font-black text-white">8 MiB</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Live recovery fixture</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">4 / 6</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Shards required</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">349 / 1024</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Audit sample</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">0</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Runtime network calls</p></div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28" aria-labelledby="peer-thesis-title">
        <div className="mb-12 max-w-4xl">
          <div className="section-kicker mb-5">The peer-seeding thesis</div>
          <h2 id="peer-thesis-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">The internet has seen this resilience pattern before.</h2>
          <p className="mt-6 text-base leading-8 text-blue-100 md:text-lg">BitTorrent showed that files can remain available while volunteers continue to seed them. Bitcoin showed that a public peer network can continue without one central server. ArkMesh explores a related question for open AI capability, with explicit consent, cryptographic provenance, complete runtime context, and verified recovery.</p>
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
            <span>DXM</span>
            <strong>∞</strong>
            <span>HUMAN<br />AUTHORIZED</span>
          </div>
          <div className="machine-covenant-copy">
            <div className="section-kicker mb-6">Deus ex machina, reversed</div>
            <h2 id="machine-covenant-title">Not a god from the machine. A promise between people.</h2>
            <blockquote>
              “Keep what helps. Question what changes. Never copy me without consent. If the network goes dark, let the knowledge remain yours.”
            </blockquote>
            <p>This is a visible design statement, not a consciousness claim. ArkMesh keeps authority with the humans who choose what is preserved, where it lives, and when it runs.</p>
          </div>
          <div className="machine-covenant-seal" aria-hidden="true">
            <span>H</span>
            <i />
            <small>SOVEREIGNTY<br />REMAINS LOCAL</small>
          </div>
        </div>
      </section>

      <section id="console" className="content-grid scroll-mt-24 py-20 md:py-28">
        <div className="mb-8 grid gap-5 md:grid-cols-[1fr_auto] md:items-end">
          <div>
            <div className="section-kicker mb-5">Command interface</div>
            <h2 className="text-3xl font-black uppercase tracking-tight text-white md:text-5xl">Everything has a command.</h2>
          </div>
          <p className="max-w-lg text-sm leading-6 text-blue-200">Try <strong className="text-cyan-200">evidence</strong>, <strong className="text-cyan-200">benchmark</strong>, <strong className="text-cyan-200">roadmap</strong>, or <strong className="text-cyan-200">plan peers</strong>. The browser explains commands; it never executes host operations.</p>
        </div>
        <CommandConsole />
      </section>

      <section className="border-y-2 border-cyan-300/40 bg-[#0c0f17] py-20 md:py-28">
        <div className="content-grid">
          <div className="section-kicker mb-5">Failure ceremony 001</div>
          <div className="grid min-w-0 gap-10 lg:grid-cols-[0.8fr_1.2fr] lg:items-start">
            <div className="min-w-0">
              <h2 className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">Delete the complete object.</h2>
              <p className="mt-6 max-w-xl text-base leading-8 text-blue-100">The release demo removes the publisher source, capsule object, and two shards. Four surviving shards reconstruct the exact signed bytes. A separate three-shard attempt is refused without writing output.</p>
              <div className="command-line mt-7">C:\ARKMESH&gt; ./scripts/demo-offline-recovery.sh</div>
              <a href="https://github.com/AlleyBo55/arkmesh/blob/master/docs/DEMO.md" target="_blank" rel="noreferrer" className="mt-5 inline-block text-xs font-black uppercase tracking-widest text-cyan-200 underline decoration-2 underline-offset-4">Read procedure + limits ↗</a>
            </div>
            <div className="blue-panel min-w-0 p-4 md:p-8">
              <pre className="ascii-map" aria-label="Offline recovery topology">{topology}</pre>
              <div className="mt-8 grid gap-3 text-xs uppercase md:grid-cols-3">
                <div className="border border-cyan-300/40 p-3"><span className="block text-blue-300">Object</span><strong className="mt-2 block text-amber-300">Deleted</strong></div>
                <div className="border border-cyan-300/40 p-3"><span className="block text-blue-300">Recovery</span><strong className="mt-2 block text-emerald-300">Exact hash</strong></div>
                <div className="border border-cyan-300/40 p-3"><span className="block text-blue-300">Authority</span><strong className="mt-2 block text-emerald-300">Trusted</strong></div>
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
            <p className="mt-6 text-base leading-8 text-blue-100">Three-run means on one Apple M5 and a 64 MiB fixture. These operations perform different work. This is a reproducible floor, not a universal speed claim.</p>
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

      <section className="border-y-2 border-cyan-300/40 bg-[#0c0f17] py-20 md:py-28">
        <div className="content-grid">
          <div className="mb-12 max-w-3xl">
            <div className="section-kicker mb-5">Reference core</div>
            <h2 className="text-3xl font-black uppercase text-white md:text-5xl">What works now.</h2>
            <p className="mt-5 text-blue-100">Each group is backed by implementation and tests. Use <span className="text-cyan-200">status &lt;id&gt;</span> in the console for its acceptance boundary.</p>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            {implementedCapabilities.map((item, index) => (
              <article key={item.id} className="blue-panel p-6">
                <div className="mb-8 flex items-center justify-between gap-4"><span className="text-xs text-blue-300">MOD-{String(index + 1).padStart(2, "0")}</span><span className="status-label text-emerald-300">Verified</span></div>
                <h3 className="text-xl font-black uppercase text-white">{item.title}</h3>
                <p className="mt-4 leading-7 text-blue-100">{item.summary}</p>
                <p className="command-line mt-6">{item.command}</p>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28">
        <div className="mb-12 grid gap-6 md:grid-cols-[1fr_0.75fr] md:items-end">
          <div>
            <div className="section-kicker mb-5">Open system work</div>
            <h2 className="text-3xl font-black uppercase text-white md:text-5xl">What must be built next.</h2>
          </div>
          <p className="border-l-2 border-amber-300 pl-4 text-sm leading-6 text-blue-100">These capabilities are interface-visible but not simulated. Every card remains PLANNED until its acceptance boundary has reproducible evidence.</p>
        </div>
        <div className="grid gap-px border border-cyan-300/40 bg-cyan-300/30 md:grid-cols-2 lg:grid-cols-3">
          {plannedCapabilities.map((item, index) => (
            <article key={item.id} className="bg-[#121621] p-6">
              <div className="mb-7 flex items-center justify-between"><span className="text-xs text-blue-300">QUEUE-{String(index + 1).padStart(2, "0")}</span><span className="status-label text-amber-300">Planned</span></div>
              <h3 className="text-lg font-black uppercase text-white">{item.title}</h3>
              <p className="mt-4 text-sm leading-6 text-blue-100">{item.summary}</p>
              <p className="command-line mt-6">{item.command}</p>
            </article>
          ))}
          <article className="flex min-h-64 flex-col justify-between bg-[#e8e8ff] p-6 text-[#15131f]">
            <span className="text-xs font-black uppercase tracking-widest">Documentation channel</span>
            <div><h3 className="text-2xl font-black uppercase">Inspect every boundary.</h3><p className="mt-3 text-sm leading-6">The complete wiki separates verified mechanics, assumptions, and planned system work.</p></div>
            <Link href="/wiki#roadmap" className="text-xs font-black uppercase tracking-widest underline decoration-2 underline-offset-4">Open roadmap wiki →</Link>
          </article>
        </div>
      </section>

      <section className="content-grid py-20 md:py-28" aria-labelledby="collaborate-title">
        <div className="blue-panel overflow-hidden">
          <div className="grid lg:grid-cols-[1.15fr_0.85fr]">
            <div className="min-w-0 p-6 md:p-10 lg:p-14">
              <div className="section-kicker mb-5">Open foundation</div>
              <h2 id="collaborate-title" className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">This is not SaaS. It is unfinished public infrastructure.</h2>
              <p className="mt-6 max-w-3xl text-base leading-8 text-blue-100">There is no subscription, hosted account, token sale, or proprietary network. ArkMesh needs people who want to test whether useful open AI can remain available and trustworthy when ordinary infrastructure fails.</p>
              <div className="mt-8 flex flex-wrap gap-3 text-xs font-black uppercase tracking-[0.1em]">
                <a href="https://github.com/AlleyBo55/arkmesh" target="_blank" rel="noreferrer" className="border-2 border-cyan-200 bg-[#e8e8ff] px-5 py-3 text-[#15131f] hover:bg-white">Explore the source ↗</a>
                <a href="https://github.com/AlleyBo55/arkmesh/blob/master/CHARTER.md" target="_blank" rel="noreferrer" className="border-2 border-cyan-200 px-5 py-3 text-cyan-100 hover:bg-[#e8e8ff] hover:text-[#15131f]">Read the charter ↗</a>
              </div>
            </div>
            <div className="border-t border-cyan-300/40 bg-[#0a0d15] p-6 md:p-10 lg:border-l lg:border-t-0">
              <p className="text-xs font-black uppercase tracking-[0.16em] text-cyan-300">Collaborators wanted</p>
              <ul className="mt-6 space-y-4 text-sm leading-6 text-blue-100">
                <li><strong className="text-white">01 // Distributed systems</strong><br />Invited discovery, encrypted transfer, partitions, and reconciliation.</li>
                <li><strong className="text-white">02 // Applied cryptography</strong><br />Authority distribution, revocation, timestamping, and adversarial review.</li>
                <li><strong className="text-white">03 // Local AI runtimes</strong><br />Deterministic inference health checks and hardware compatibility.</li>
                <li><strong className="text-white">04 // Preservation research</strong><br />Physical rehearsals, failure curves, governance, and independent reproduction.</li>
              </ul>
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
