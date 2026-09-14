import type { Metadata } from "next";
import Link from "next/link";
import type { ReactNode } from "react";
import { AppProvider } from "@/components/app-provider";
import "./globals.css";

export const metadata: Metadata = {
  metadataBase: new URL("https://github.com/AlleyBo55/arkmesh"),
  title: {
    default: "ArkMesh // Verified Offline Recovery",
    template: "%s // ArkMesh",
  },
  description: "A research prototype for preserving and recovering signed open AI capabilities without cloud dependence.",
  keywords: [
    "offline AI",
    "model preservation",
    "cryptographic provenance",
    "erasure coding",
    "distributed systems",
  ],
  openGraph: {
    title: "ArkMesh // Verified Offline Recovery",
    description: "Delete the original object. Recover the exact signed bytes from surviving shards.",
    type: "website",
  },
};

export default function RootLayout({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <AppProvider>
          <div className="screen-shell">
            <header className="sticky top-0 z-50 border-b-2 border-cyan-300/80 bg-[#000a46]/95 backdrop-blur-sm">
              <div className="content-grid flex min-h-14 items-center justify-between gap-4 py-2 text-xs font-black uppercase tracking-[0.12em]">
                <Link href="/" className="flex items-center gap-3 text-white hover:text-cyan-200">
                  <span className="border-2 border-cyan-300 px-2 py-1 text-cyan-200">AM</span>
                  <span className="hidden sm:inline">ArkMesh Recovery System</span>
                </Link>
                <nav className="flex items-center gap-2 sm:gap-5" aria-label="Primary navigation">
                  <Link href="/" className="text-blue-200 hover:text-white">[Home]</Link>
                  <Link href="/wiki" className="text-blue-200 hover:text-white">[Wiki]</Link>
                  <a href="https://github.com/AlleyBo55/arkmesh" target="_blank" rel="noreferrer" className="text-cyan-200 hover:text-white">[Source]</a>
                </nav>
              </div>
            </header>
            {children}
            <footer className="mt-24 border-t-2 border-cyan-300/50 bg-[#000735] py-8">
              <div className="content-grid flex flex-col justify-between gap-4 text-xs uppercase tracking-[0.1em] text-blue-300 md:flex-row">
                <p>ArkMesh v0alpha1 // research prototype // MIT</p>
                <p>Networking and inference are not implemented.</p>
              </div>
            </footer>
          </div>
        </AppProvider>
      </body>
    </html>
  );
}
