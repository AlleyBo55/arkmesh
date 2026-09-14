import type { Metadata } from "next";
import Link from "next/link";
import { CommandConsole } from "@/components/command-console";
import { plannedCapabilities, wikiSections } from "@/lib/content";

export const metadata: Metadata = {
  title: "System Wiki",
  description: "ArkMesh protocol concepts, commands, evidence, security boundaries, and research roadmap.",
};

export default function WikiPage() {
  return (
    <main className="content-grid pt-14 md:pt-20">
      <section className="border-b-2 border-cyan-300/50 pb-14">
        <div className="section-kicker mb-5">Documentation memory map</div>
        <div className="grid gap-8 lg:grid-cols-[1fr_auto] lg:items-end">
          <div>
            <h1 className="text-5xl font-black uppercase leading-none tracking-[-0.06em] text-white md:text-8xl">System<br /><span className="text-cyan-200">Wiki.</span></h1>
            <p className="mt-7 max-w-3xl border-l-4 border-cyan-300 pl-5 text-base leading-8 text-blue-100">One screen for the capsule protocol, authority model, repair path, audit claims, security boundaries, evidence, and the seven capabilities still separating the alpha from its full thesis.</p>
          </div>
          <div className="blue-panel min-w-64 p-4 text-xs uppercase tracking-widest">
            <p className="text-blue-300">Documents indexed</p><p className="mt-2 text-3xl font-black text-white">{wikiSections.length.toString().padStart(2, "0")}</p>
            <p className="mt-5 text-blue-300">Planned systems</p><p className="mt-2 text-3xl font-black text-amber-300">{plannedCapabilities.length.toString().padStart(2, "0")}</p>
          </div>
        </div>
      </section>

      <div className="grid gap-12 lg:grid-cols-[240px_minmax(0,1fr)]">
        <aside className="hidden lg:block">
          <nav className="sticky top-20 py-10" aria-label="Wiki table of contents">
            <p className="mb-4 text-xs font-black uppercase tracking-[0.15em] text-cyan-200">C:\WIKI\INDEX</p>
            {wikiSections.map((section) => (
              <a key={section.id} href={`#${section.id}`} className="wiki-link">{section.index} {section.title}</a>
            ))}
            <Link href="/" className="wiki-link mt-5 border-t border-cyan-300/30 pt-4">← Return home</Link>
          </nav>
        </aside>

        <div className="min-w-0 py-4 lg:py-10">
          <details className="blue-panel mb-8 p-4 lg:hidden">
            <summary className="cursor-pointer text-xs font-black uppercase tracking-widest text-cyan-200">Open wiki index</summary>
            <nav className="mt-4" aria-label="Mobile wiki table of contents">
              {wikiSections.map((section) => <a key={section.id} href={`#${section.id}`} className="wiki-link">{section.index} {section.title}</a>)}
            </nav>
          </details>

          {wikiSections.map((section) => (
            <section key={section.id} id={section.id} className="wiki-section first:border-t-0">
              <div className="mb-5 flex items-center justify-between gap-6">
                <span className="text-xs font-black tracking-[0.18em] text-cyan-300">FILE {section.index}</span>
                <a href={`#${section.id}`} className="text-xs text-blue-300 hover:text-cyan-200" aria-label={`Link to ${section.title}`}>#{section.id}</a>
              </div>
              <h2 className="text-3xl font-black uppercase tracking-tight text-white md:text-5xl">{section.title}</h2>
              <p className="mt-5 max-w-4xl text-lg font-bold leading-8 text-cyan-100">{section.lead}</p>
              <div className="mt-7 max-w-4xl space-y-5 text-base leading-8 text-blue-100">
                {section.paragraphs.map((paragraph) => <p key={paragraph}>{paragraph}</p>)}
              </div>

              {section.id === "roadmap" ? (
                <div className="mt-9 grid gap-3 md:grid-cols-2">
                  {plannedCapabilities.map((item) => (
                    <article key={item.id} className="border border-cyan-300/35 bg-[#121621] p-5">
                      <div className="flex items-center justify-between gap-4"><span className="text-xs uppercase text-blue-300">{item.id}</span><span className="status-label text-amber-300">Planned</span></div>
                      <h3 className="mt-5 text-lg font-black uppercase text-white">{item.title}</h3>
                      <p className="mt-3 text-sm leading-6 text-blue-100">{item.summary}</p>
                      <p className="mt-4 border-l-2 border-amber-300 pl-3 text-xs leading-5 text-amber-100"><strong>Done when:</strong> {item.acceptance}</p>
                      <p className="command-line mt-5">{item.command}</p>
                    </article>
                  ))}
                </div>
              ) : null}

              {section.commands?.length ? (
                <div className="mt-8">
                  <p className="mb-3 text-xs font-black uppercase tracking-[0.14em] text-blue-300">Command references</p>
                  <div className="space-y-2">{section.commands.map((command) => <p key={command} className="command-line">C:\ARKMESH&gt; {command}</p>)}</div>
                </div>
              ) : null}

              {section.links?.length ? (
                <div className="mt-7 flex flex-wrap gap-3">
                  {section.links.map((link) => (
                    <a key={link.href} href={link.href} target="_blank" rel="noreferrer" className="border border-cyan-300 px-3 py-2 text-xs font-black uppercase tracking-wider text-cyan-200 hover:bg-cyan-200 hover:text-[#00104f]">{link.label} ↗</a>
                  ))}
                </div>
              ) : null}
            </section>
          ))}

          <section className="wiki-section" aria-labelledby="wiki-console-title">
            <div className="section-kicker mb-5">Interactive index</div>
            <h2 id="wiki-console-title" className="mb-7 text-3xl font-black uppercase text-white md:text-5xl">Query the system.</h2>
            <CommandConsole compact />
          </section>
        </div>
      </div>
    </main>
  );
}
