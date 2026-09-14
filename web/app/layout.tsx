import type { Metadata, Viewport } from "next";
import Link from "next/link";
import type { ReactNode } from "react";
import { AppProvider } from "@/components/app-provider";
import "./globals.css";

const publicOrigin = process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, "");
const repositoryUrl = "https://github.com/AlleyBo55/arkmesh";

export const metadata: Metadata = {
  metadataBase: new URL(publicOrigin ?? "http://localhost:3000"),
  applicationName: "ArkMesh",
  title: {
    default: "ArkMesh",
    template: "%s | ArkMesh",
  },
  description: "An open-source foundation project exploring how useful open AI can survive outages, withdrawn hosts, and infrastructure failure through verifiable, consent-based peer preservation.",
  keywords: [
    "ArkMesh",
    "offline AI",
    "local LLM",
    "model preservation",
    "cryptographic provenance",
    "Reed Solomon erasure coding",
    "content addressed storage",
    "distributed systems",
    "disaster recovery",
  ],
  authors: [{ name: "ArkMesh contributors", url: repositoryUrl }],
  creator: "ArkMesh contributors",
  publisher: "ArkMesh",
  category: "technology",
  alternates: publicOrigin ? { canonical: "/" } : undefined,
  manifest: "/manifest.webmanifest",
  icons: { icon: "/icon.svg" },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-image-preview": "large",
      "max-snippet": -1,
      "max-video-preview": -1,
    },
  },
  openGraph: {
    title: "ArkMesh",
    description: "What if the internet goes down or a useful model disappears from its original host? Explore verifiable, consent-based peer preservation for open AI.",
    siteName: "ArkMesh",
    type: "website",
    locale: "en_US",
    url: publicOrigin ?? repositoryUrl,
  },
  twitter: {
    card: "summary_large_image",
    title: "ArkMesh",
    description: "An open-source foundation project exploring how useful open AI can survive without one provider, network, or machine.",
  },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  colorScheme: "dark",
  themeColor: "#001467",
};

const structuredData = {
  "@context": "https://schema.org",
  "@type": "SoftwareSourceCode",
  name: "ArkMesh",
  programmingLanguage: ["Go", "TypeScript"],
  runtimePlatform: "Linux, macOS",
  version: "0.1.0-alpha.1",
  description: "An open-source foundation project researching verifiable, consent-based peer preservation for useful open AI capabilities.",
  url: publicOrigin ?? repositoryUrl,
  codeRepository: repositoryUrl,
  license: "https://opensource.org/license/mit",
  isAccessibleForFree: true,
};

export default function RootLayout({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData).replace(/</g, "\\u003c") }}
        />
        <AppProvider>
          <div className="screen-shell">
            <header className="sticky top-0 z-50 border-b border-white/10 bg-[#090b12]/80 backdrop-blur-xl">
              <div className="content-grid flex min-h-14 items-center justify-between gap-2 py-1 text-xs font-black uppercase tracking-[0.1em] sm:gap-4 sm:py-2 sm:tracking-[0.12em]">
                <Link href="/" className="flex min-h-11 items-center gap-3 text-white hover:text-cyan-200">
                  <span className="border-2 border-cyan-300 px-2 py-1 text-cyan-200">AM</span>
                  <span className="hidden sm:inline">ArkMesh</span>
                  <span className="sr-only">ArkMesh home</span>
                </Link>
                <nav className="flex items-center gap-1 sm:gap-3 md:gap-5" aria-label="Primary navigation">
                  <Link href="/" className="flex min-h-11 items-center px-1 text-blue-200 hover:text-white sm:px-2">[Home]</Link>
                  <Link href="/learn" className="flex min-h-11 items-center px-1 text-blue-200 hover:text-white sm:px-2">[Learn]</Link>
                  <Link href="/wiki" className="flex min-h-11 items-center px-1 text-blue-200 hover:text-white sm:px-2">[Wiki]</Link>
                  <a href={repositoryUrl} target="_blank" rel="noreferrer" className="flex min-h-11 items-center px-1 text-cyan-200 hover:text-white sm:px-2">[Source]</a>
                </nav>
              </div>
            </header>
            {children}
            <footer className="mt-20 border-t border-white/10 bg-[#080a11] py-8 md:mt-24">
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
