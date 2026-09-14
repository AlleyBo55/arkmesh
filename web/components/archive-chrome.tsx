"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect } from "react";

const navigation = [
  { href: "/", code: "01", label: "Overview" },
  { href: "/wiki", code: "02", label: "Thesis" },
  { href: "/learn", code: "03", label: "Field guide" },
] as const;

function isCurrentPath(pathname: string, href: string) {
  return href === "/" ? pathname === href : pathname.startsWith(href);
}

export function ArchiveChrome({ repositoryUrl }: { repositoryUrl: string }) {
  const pathname = usePathname();

  useEffect(() => {
    const root = document.documentElement;
    const motionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    const sections = Array.from(document.querySelectorAll<HTMLElement>("main > section"));
    let frame = 0;

    const updateProgress = () => {
      frame = 0;
      const distance = document.documentElement.scrollHeight - window.innerHeight;
      const progress = distance > 0 ? Math.min(1, Math.max(0, window.scrollY / distance)) : 0;
      root.style.setProperty("--archive-progress", progress.toString());
    };

    const requestProgressUpdate = () => {
      if (!frame) frame = window.requestAnimationFrame(updateProgress);
    };

    const updatePointer = (event: PointerEvent) => {
      root.style.setProperty("--archive-pointer-x", `${event.clientX}px`);
      root.style.setProperty("--archive-pointer-y", `${event.clientY}px`);
    };

    sections.forEach((section) => {
      section.classList.add("archive-reveal");
      if (motionQuery.matches || section.getBoundingClientRect().top < window.innerHeight * 0.92) {
        section.classList.add("is-visible");
      }
    });

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) entry.target.classList.add("is-visible");
        });
      },
      { rootMargin: "0px 0px -10% 0px", threshold: 0.08 },
    );

    sections.forEach((section) => observer.observe(section));
    updateProgress();
    window.addEventListener("scroll", requestProgressUpdate, { passive: true });
    window.addEventListener("resize", requestProgressUpdate, { passive: true });
    window.addEventListener("pointermove", updatePointer, { passive: true });

    return () => {
      observer.disconnect();
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener("scroll", requestProgressUpdate);
      window.removeEventListener("resize", requestProgressUpdate);
      window.removeEventListener("pointermove", updatePointer);
    };
  }, [pathname]);

  return (
    <>
      <a className="archive-skip-link" href="#main-content">Skip to research content</a>
      <header className="archive-header">
        <div className="archive-read-progress" aria-hidden="true"><span /></div>
        <div className="content-grid archive-header-main">
          <Link href="/" className="archive-brand" aria-label="ArkMesh home">
            <span className="archive-brand-mark" aria-hidden="true"><i /><i /><i /></span>
            <span><strong>ArkMesh</strong><small>Continuity research archive</small></span>
          </Link>

          <nav className="archive-nav" aria-label="Primary navigation">
            {navigation.map((item) => {
              const current = isCurrentPath(pathname, item.href);
              return (
                <Link key={item.href} href={item.href} data-active={current ? "true" : "false"} aria-current={current ? "page" : undefined}>
                  <span>{item.code}</span><b>{item.label}</b>
                </Link>
              );
            })}
            <a href={repositoryUrl} target="_blank" rel="noreferrer"><span>04</span>Source ↗</a>
          </nav>
        </div>

        <div className="archive-evidence-band">
          <div className="content-grid">
            <span><i className="is-demonstrated" />Local cryptographic core: demonstrated</span>
            <span><i className="is-open" />Distributed continuity: open research</span>
            <Link href="/wiki#evidence">Open claim ledger <b aria-hidden="true">→</b></Link>
          </div>
        </div>
      </header>
    </>
  );
}
