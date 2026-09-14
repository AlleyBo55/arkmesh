# ArkMesh web interface

A command-first landing page and system wiki for ArkMesh, an open-source foundation project researching how useful open AI can survive outages, withdrawn hosts, and infrastructure failure through verifiable, consent-based peer preservation.

This is not a SaaS interface. The landing page leads with the public research question, distinguishes the planned peer network from current local evidence, and invites collaboration across distributed systems, cryptography, local AI runtimes, and digital preservation.

## Requirements

- Node.js 20.9 or newer
- npm 12.0.2 when reproducing the checked lockfile workflow

The application was validated with Node.js 26.8.2. Package versions are exact in `package.json` and `package-lock.json`.

## Run

```bash
npm ci
npm run dev
```

Open `http://localhost:3000`.

Production validation:

```bash
npm run lint
npm run typecheck
npm run build
npm run start
```

## Interface

The landing page and wiki share a Redux-backed command console. Useful commands include:

```text
help
status
status repair
evidence
benchmark
roadmap
plan peers
plan inference
wiki security
wiki roadmap
run demo
```

Browser commands reveal documentation and navigate the site. They never execute host shell commands.

## Truth boundary

The interface marks implemented cryptographic and recovery mechanisms as `VERIFIED`. It marks organizational policy distribution, peer transfer, shard scheduling, longitudinal retention, attested time, local inference, and partition reconciliation as `PLANNED`.

A visible roadmap card is not an implementation claim. The status comes from the shared content model in `lib/content.ts` and must remain consistent with the root project README and threat model.

## Stack

- Next.js 16 app router
- React 19
- TypeScript 6, the newest release currently supported by the Next.js ESLint parser
- Tailwind CSS 4
- Redux Toolkit and React Redux
- ESLint 9, the newest release currently supported by every plugin bundled with the Next.js configuration

The visual system uses local fonts, procedural Three.js geometry, shaders, and no remote runtime assets. The landing page maps native reversible scroll to a deep-sea camera journey through five truthful mission beats. Three.js loads only after the visitor begins scrolling; reduced-motion and no-WebGL paths retain the complete story in ordinary document flow. Monospaced command styling remains reserved for the interactive console and protocol evidence.

## Deployment metadata

Set the real public origin before a production deployment:

```bash
cp .env.example .env.local
```

Replace `https://example.com` with the deployed origin. `NEXT_PUBLIC_SITE_URL` controls canonical metadata and sitemap URLs. Do not deploy with the example value.
