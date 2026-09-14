import Link from "next/link";
import { CommandConsole } from "@/components/command-console";
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
      <section className="content-grid pb-20 pt-14 md:pb-28 md:pt-24">
        <div className="section-kicker mb-7">Offline capability preservation</div>
        <div className="grid items-end gap-12 lg:grid-cols-[1.35fr_0.65fr]">
          <div>
            <h1 className="big-type">
              Recover.
              <br />
              <span className="text-cyan-200">Verify.</span>
              <br />
              <span className="outline-type">Continue.</span>
            </h1>
            <p className="mt-10 max-w-3xl border-l-4 border-cyan-300 pl-5 text-base leading-8 text-blue-100 md:text-xl">
              Preserve open models, runtimes, knowledge, and recovery instructions as one signed capsule. Delete the original object. Rebuild the exact bytes from surviving shards without cloud authentication.
            </p>
            <div className="mt-9 flex flex-wrap gap-3 text-xs font-black uppercase tracking-[0.12em]">
              <Link href="#console" className="border-2 border-cyan-200 bg-cyan-200 px-5 py-3 text-[#00104f] hover:bg-white">
                Open console
              </Link>
              <Link href="/wiki" className="border-2 border-cyan-200 px-5 py-3 text-cyan-100 hover:bg-cyan-200 hover:text-[#00104f]">
                Read system wiki
              </Link>
            </div>
          </div>
          <aside className="blue-panel p-5 text-xs uppercase tracking-[0.1em]" aria-label="Release status">
            <div className="mb-5 flex items-center justify-between border-b border-cyan-300/40 pb-3">
              <span className="text-blue-300">System status</span>
              <span className="status-label text-amber-300">Alpha</span>
            </div>
            <dl className="space-y-4">
              <div className="flex justify-between gap-4"><dt className="text-blue-300">Release</dt><dd className="text-white">0.1.0-A1</dd></div>
              <div className="flex justify-between gap-4"><dt className="text-blue-300">Core</dt><dd className="text-emerald-300">Verified</dd></div>
              <div className="flex justify-between gap-4"><dt className="text-blue-300">Transport</dt><dd className="text-amber-300">Planned</dd></div>
              <div className="flex justify-between gap-4"><dt className="text-blue-300">Inference</dt><dd className="text-amber-300">Planned</dd></div>
              <div className="flex justify-between gap-4"><dt className="text-blue-300">Dependencies</dt><dd className="text-white">Go stdlib</dd></div>
            </dl>
            <a
              href="https://github.com/AlleyBo55/arkmesh/releases/tag/v0.1.0-alpha.1"
              target="_blank"
              rel="noreferrer"
              className="mt-6 block border border-cyan-300 px-3 py-2 text-center text-cyan-200 hover:bg-cyan-200 hover:text-[#00104f]"
            >
              Inspect tagged release ↗
            </a>
          </aside>
        </div>
      </section>

      <section className="border-y-2 border-cyan-300/50 bg-[#00125f]">
        <div className="content-grid grid md:grid-cols-4">
          <div className="metric-cell"><p className="text-3xl font-black text-white">8 MiB</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Live recovery fixture</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">4 / 6</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Shards required</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">349 / 1024</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Audit sample</p></div>
          <div className="metric-cell"><p className="text-3xl font-black text-white">0</p><p className="mt-2 text-xs uppercase tracking-widest text-blue-300">Runtime network calls</p></div>
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

      <section className="border-y-2 border-cyan-300/40 bg-[#000a43] py-20 md:py-28">
        <div className="content-grid">
          <div className="section-kicker mb-5">Failure ceremony 001</div>
          <div className="grid gap-10 lg:grid-cols-[0.8fr_1.2fr] lg:items-start">
            <div>
              <h2 className="text-3xl font-black uppercase leading-tight text-white md:text-5xl">Delete the complete object.</h2>
              <p className="mt-6 max-w-xl text-base leading-8 text-blue-100">The release demo removes the publisher source, capsule object, and two shards. Four surviving shards reconstruct the exact signed bytes. A separate three-shard attempt is refused without writing output.</p>
              <div className="command-line mt-7">C:\ARKMESH&gt; ./scripts/demo-offline-recovery.sh</div>
              <a href="https://github.com/AlleyBo55/arkmesh/blob/master/docs/DEMO.md" target="_blank" rel="noreferrer" className="mt-5 inline-block text-xs font-black uppercase tracking-widest text-cyan-200 underline decoration-2 underline-offset-4">Read procedure + limits ↗</a>
            </div>
            <div className="blue-panel p-4 md:p-8">
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
            <a href="https://github.com/AlleyBo55/arkmesh/blob/master/benchmarks/2026-09-14-apple-m5-64mib/raw.csv" target="_blank" rel="noreferrer" className="mt-6 inline-block border border-cyan-300 px-4 py-3 text-xs font-black uppercase tracking-widest text-cyan-200 hover:bg-cyan-200 hover:text-[#00104f]">Download raw CSV ↗</a>
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

      <section className="border-y-2 border-cyan-300/40 bg-[#00105a] py-20 md:py-28">
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
            <article key={item.id} className="bg-[#000d50] p-6">
              <div className="mb-7 flex items-center justify-between"><span className="text-xs text-blue-300">QUEUE-{String(index + 1).padStart(2, "0")}</span><span className="status-label text-amber-300">Planned</span></div>
              <h3 className="text-lg font-black uppercase text-white">{item.title}</h3>
              <p className="mt-4 text-sm leading-6 text-blue-100">{item.summary}</p>
              <p className="command-line mt-6">{item.command}</p>
            </article>
          ))}
          <article className="flex min-h-64 flex-col justify-between bg-cyan-200 p-6 text-[#00104f]">
            <span className="text-xs font-black uppercase tracking-widest">Documentation channel</span>
            <div><h3 className="text-2xl font-black uppercase">Inspect every boundary.</h3><p className="mt-3 text-sm leading-6">The complete wiki separates verified mechanics, assumptions, and planned system work.</p></div>
            <Link href="/wiki#roadmap" className="text-xs font-black uppercase tracking-widest underline decoration-2 underline-offset-4">Open roadmap wiki →</Link>
          </article>
        </div>
      </section>

      <section className="border-y-2 border-cyan-300/40 bg-[#000832] py-20">
        <div className="content-grid grid gap-10 md:grid-cols-[1fr_auto] md:items-center">
          <div><div className="section-kicker mb-5">Research posture</div><h2 className="text-3xl font-black uppercase text-white md:text-5xl">Evidence before mythology.</h2><p className="mt-5 max-w-3xl leading-8 text-blue-100">A signature is not legal identity. A sample is not whole-object proof. A reconstructed file is not executable continuity. ArkMesh names each boundary so the thesis can be falsified.</p></div>
          <Link href="/wiki#security" className="border-2 border-cyan-200 px-6 py-4 text-center text-xs font-black uppercase tracking-widest text-cyan-100 hover:bg-cyan-200 hover:text-[#00104f]">Read security boundaries</Link>
        </div>
      </section>
    </main>
  );
}
