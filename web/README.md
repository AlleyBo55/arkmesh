# ArkMesh web interface

A mission-first landing page, eight-module learning path, and system wiki for ArkMesh, an open-source research project investigating whether lawfully held AI capabilities can remain obtainable, verifiable, and runnable through outages, withdrawn hosts, and infrastructure failure.

The learning path uses "AI extinction" precisely to mean loss of a lawfully held capability and its required context, not consciousness or biological life. It separates verified experiments from planned research and gives students and researchers a reproducible route from capsule construction through recovery, consent, inference, and open protocol questions.

This is not a SaaS interface. The landing page leads with the public research question, distinguishes the planned peer network from current local evidence, and invites collaboration across distributed systems, cryptography, local AI runtimes, and digital preservation.

## Requirements

- Node.js 20.9 or newer
- `package.json` declares npm 12.0.2 as the intended package manager

The 2026-09-14 fact-check used Node.js 22.22.0 and npm 10.9.4 against the existing installed dependencies; it did not rerun `npm ci`. Package versions are exact in `package.json` and `package-lock.json`.

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
plan transfer
plan execution
wiki security
wiki roadmap
run demo
```

Browser commands reveal documentation and navigate the site. They never execute host shell commands.

## Truth boundary

The interface marks implemented cryptographic and local recovery mechanisms as `VERIFIED`. It marks fresh-device bootstrap, authenticated peer transfer, failure-domain placement, longitudinal retention, attested time, runtime health checks, partition reconciliation, and the independent continuity experiment as `PLANNED`.

A visible roadmap card is not an implementation claim. The status comes from the shared content model in `lib/content.ts` and must remain consistent with the root project README and threat model.

## Stack

- Next.js 16 app router
- React 19
- TypeScript 6
- Tailwind CSS 4
- Redux Toolkit and React Redux
- ESLint 9

The visual system uses local fonts, procedural Three.js geometry, shaders, and no remote runtime assets. The landing page maps native reversible scroll to a deep-sea camera journey through five content chapters. Three.js loads only after the visitor begins scrolling; reduced-motion and no-WebGL paths retain the same chapter text in ordinary document flow. Monospaced command styling remains reserved for the interactive console and protocol evidence.

## Deployment metadata

Set the real public origin before a production deployment:

```bash
cp .env.example .env.local
```

Replace `https://example.com` with the deployed origin. `NEXT_PUBLIC_SITE_URL` controls canonical metadata and sitemap URLs. Do not deploy with the example value.
