import type { Metadata, Viewport } from "next";
import Link from "next/link";
import type { ReactNode } from "react";
import { AppProvider } from "@/components/app-provider";
import { ArchiveChrome } from "@/components/archive-chrome";
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
  description: "A research prototype for content-addressed AI capsules, offline trust, authenticated repair, and donorless local reconstruction. Distributed peer continuity remains planned.",
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
    description: "ArkMesh demonstrates signed local capsules, authenticated repair, and donorless reconstruction. It asks whether consenting peers can extend that evidence beyond one original host.",
    siteName: "ArkMesh",
    type: "website",
    locale: "en_US",
    url: publicOrigin ?? repositoryUrl,
  },
  twitter: {
    card: "summary_large_image",
    title: "ArkMesh",
    description: "Local signed-capsule integrity, authenticated repair, and donorless reconstruction are implemented; distributed continuity is the open research question.",
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
  description: "Research prototype for content-addressed AI capsules, offline trust, authenticated repair, and donorless local reconstruction.",
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
            <ArchiveChrome repositoryUrl={repositoryUrl} />
            {children}
            <footer className="archive-footer">
              <div className="content-grid archive-footer-grid">
                <div>
                  <p className="archive-footer-label">ArkMesh // Open continuity research</p>
                  <h2>Preserve the evidence.<br />Challenge the thesis.</h2>
                </div>
                <div className="archive-footer-links">
                  <Link href="/wiki">Read the thesis</Link>
                  <Link href="/wiki#evidence">Inspect evidence</Link>
                  <Link href="/learn">Run the field guide</Link>
                  <a href={repositoryUrl} target="_blank" rel="noreferrer">View source ↗</a>
                </div>
              </div>
              <div className="content-grid archive-footer-status">
                <p>v0alpha1 · MIT · research prototype</p>
                <p><span />Peer networking and inference remain unimplemented.</p>
              </div>
            </footer>
          </div>
        </AppProvider>
      </body>
    </html>
  );
}
