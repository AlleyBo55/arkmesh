"use client";

import { useEffect, useState } from "react";

const phases = [
  {
    id: "distribute",
    step: "01",
    label: "Distribute",
    headline: "Signed capsule split into six shards",
    detail: "Four data shards and two parity shards are placed across independent holders.",
    output: "CAPSULE COMMITTED // 6 SHARDS ONLINE",
  },
  {
    id: "failure",
    step: "02",
    label: "Delete",
    headline: "Publisher, object, and two shards disappear",
    detail: "The complete object no longer exists. Four shards remain, exactly at the recovery threshold.",
    output: "OBJECT ABSENT // SHARDS 01 + 04 LOST",
  },
  {
    id: "rebuild",
    step: "03",
    label: "Rebuild",
    headline: "Surviving shards reconstruct the bytes",
    detail: "No surviving node holds the complete object. Reconstruction streams into a temporary file.",
    output: "REED SOLOMON MATRIX // 4 OF 4 INPUTS",
  },
  {
    id: "verify",
    step: "04",
    label: "Verify",
    headline: "Exact digest and trusted signature match",
    detail: "Only authenticated output is installed. Size, SHA-256, chunk root, and signer trust all pass.",
    output: "SHA-256 MATCH // SIGNATURE TRUSTED // ACCEPT",
  },
] as const;

const shards = ["00", "01", "02", "03", "04", "05"];
const lostShards = new Set(["01", "04"]);

export function ContinuitySimulator() {
  const [phaseIndex, setPhaseIndex] = useState(0);
  const [paused, setPaused] = useState(false);
  const phase = phases[phaseIndex];

  useEffect(() => {
    const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    if (paused || reducedMotion) return;

    const timer = window.setInterval(() => {
      setPhaseIndex((current) => (current + 1) % phases.length);
    }, 3200);

    return () => window.clearInterval(timer);
  }, [paused]);

  const hasFailed = phaseIndex >= 1;
  const isRebuilding = phaseIndex === 2;
  const isVerified = phaseIndex === 3;

  return (
    <section className="simulator" aria-labelledby="simulator-title">
      <div className="terminal-titlebar">
        <span>ARKMESH://CONTINUITY.VIS</span>
        <span className="hidden sm:inline">DETERMINISTIC WALKTHROUGH</span>
        <button
          type="button"
          className="sim-pause"
          onClick={() => setPaused((value) => !value)}
          aria-pressed={paused}
        >
          {paused ? "PLAY" : "PAUSE"}
        </button>
      </div>

      <div className="sim-body">
        <div className="sim-readout">
          <div>
            <p className="text-xs font-black uppercase tracking-[0.16em] text-cyan-300">Phase {phase.step} / 04</p>
            <h2 id="simulator-title" className="mt-3 text-2xl font-black uppercase leading-tight text-white md:text-4xl">{phase.headline}</h2>
          </div>
          <p className="max-w-xl text-sm leading-6 text-blue-100 md:text-right">{phase.detail}</p>
        </div>

        <div className={`sim-layout phase-${phase.id}`}>
          <article className={`sim-node sim-source ${hasFailed ? "is-lost" : "is-live"}`}>
            <div className="sim-node-head"><span>PUBLISHER</span><span>A</span></div>
            <div className="capsule-glyph" aria-hidden="true">
              <span /><span /><span /><span /><span /><span />
            </div>
            <strong>{hasFailed ? "OFFLINE" : "SIGNED"}</strong>
            <small>{hasFailed ? "source deleted" : "capsule ready"}</small>
          </article>

          <div className={`transfer-bus bus-out ${phaseIndex === 0 ? "is-flowing" : ""}`} aria-hidden="true">
            <span /><span /><span /><span />
          </div>

          <div className="shard-bank" aria-label="Six erasure shards">
            <div className="sim-node-head"><span>SHARD BANK</span><span>4 + 2</span></div>
            <div className="shard-grid">
              {shards.map((shard) => {
                const lost = hasFailed && lostShards.has(shard);
                const active = isRebuilding && !lost;
                return (
                  <div key={shard} className={`shard-cell ${lost ? "is-lost" : ""} ${active ? "is-active" : ""}`}>
                    <span>S{shard}</span>
                    <small>{lost ? "LOST" : active ? "TX" : "READY"}</small>
                  </div>
                );
              })}
            </div>
          </div>

          <div className={`transfer-bus bus-in ${isRebuilding || isVerified ? "is-flowing" : ""}`} aria-hidden="true">
            <span /><span /><span /><span />
          </div>

          <article className={`sim-node sim-target ${isVerified ? "is-verified" : isRebuilding ? "is-building" : ""}`}>
            <div className="sim-node-head"><span>FRESH NODE</span><span>D</span></div>
            <div className="verify-glyph" aria-hidden="true">
              <span className="verify-ring" />
              <span className="verify-mark">{isVerified ? "✓" : isRebuilding ? "4/4" : "…"}</span>
            </div>
            <strong>{isVerified ? "VERIFIED" : isRebuilding ? "REBUILDING" : "WAITING"}</strong>
            <small>{isVerified ? "exact signed bytes" : isRebuilding ? "temporary output" : "no complete replica"}</small>
          </article>
        </div>

        <div className="sim-output" aria-live="polite">
          <span>C:\ARKMESH&gt;</span> {phase.output}
        </div>

        <div className="sim-controls" aria-label="Simulation phases">
          {phases.map((item, index) => (
            <button
              key={item.id}
              type="button"
              className={index === phaseIndex ? "is-current" : ""}
              onClick={() => {
                setPhaseIndex(index);
                setPaused(true);
              }}
              aria-current={index === phaseIndex ? "step" : undefined}
            >
              <span>{item.step}</span>
              {item.label}
            </button>
          ))}
        </div>
      </div>
    </section>
  );
}
