"use client";

import { FormEvent, KeyboardEvent, useEffect, useRef, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { useDispatch, useSelector } from "react-redux";
import {
  allCapabilities,
  benchmarkRows,
  implementedCapabilities,
  plannedCapabilities,
  wikiSections,
} from "@/lib/content";
import {
  appendEntry,
  clearEntries,
  type AppDispatch,
  type RootState,
  type TerminalEntry,
  type TerminalTone,
} from "@/lib/store";

const quickCommands = ["help", "status", "evidence", "benchmark", "roadmap", "wiki"];

const toneClass: Record<TerminalTone, string> = {
  system: "text-cyan-200",
  success: "text-emerald-300",
  warning: "text-amber-300",
  muted: "text-blue-300",
};

function entry(command: string, lines: string[], tone: TerminalTone = "system"): TerminalEntry {
  return {
    id: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    command,
    lines,
    tone,
  };
}

export function CommandConsole({ compact = false }: { compact?: boolean }) {
  const dispatch = useDispatch<AppDispatch>();
  const entries = useSelector((state: RootState) => state.terminal.entries);
  const history = useSelector((state: RootState) => state.terminal.history);
  const [value, setValue] = useState("");
  const [historyIndex, setHistoryIndex] = useState(-1);
  const outputRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    outputRef.current?.scrollTo({ top: outputRef.current.scrollHeight, behavior: "smooth" });
  }, [entries]);

  const runCommand = (raw: string) => {
    const command = raw.trim();
    const normalized = command.toLowerCase().replace(/\s+/g, " ");
    if (!normalized) return;

    if (normalized === "clear") {
      dispatch(clearEntries());
      setHistoryIndex(-1);
      return;
    }

    const [verb, ...rest] = normalized.split(" ");
    const target = rest.join(" ");

    if (verb === "help") {
      dispatch(
        appendEntry(
          entry(command, [
            "COMMAND INDEX",
            "STATUS [CAPABILITY]   verified implementation state",
            "EVIDENCE              deletion ceremony and release proof",
            "BENCHMARK             measured local observations",
            "ROADMAP               seven planned system capabilities",
            "PLAN <ID>             design goal and acceptance boundary",
            "WIKI [SECTION]        open the complete protocol wiki",
            "RUN DEMO              show the reproducible recovery command",
            "RELEASE               open v0.1.0-alpha.1",
            "GITHUB                open source repository",
            "CLEAR                 reset this console",
          ]),
        ),
      );
      return;
    }

    if (verb === "status") {
      if (target) {
        const capability = allCapabilities.find((item) => item.id === target);
        dispatch(
          appendEntry(
            capability
              ? entry(
                  command,
                  [
                    `${capability.status.toUpperCase()} :: ${capability.title}`,
                    capability.summary,
                    `ACCEPTANCE :: ${capability.acceptance}`,
                  ],
                  capability.status === "verified" ? "success" : "warning",
                )
              : entry(command, [`UNKNOWN CAPABILITY :: ${target}`, "Run ROADMAP or STATUS for valid identifiers."], "warning"),
          ),
        );
        return;
      }
      dispatch(
        appendEntry(
          entry(
            command,
            [
              `REFERENCE CORE :: ${implementedCapabilities.length} VERIFIED CAPABILITY GROUPS`,
              `SYSTEM THESIS :: ${plannedCapabilities.length} PLANNED CAPABILITY GROUPS`,
              "NETWORKING :: NOT IMPLEMENTED",
              "LOCAL INFERENCE :: NOT IMPLEMENTED",
              "DISASTER-CRITICAL USE :: NOT APPROVED",
            ],
            "success",
          ),
        ),
      );
      return;
    }

    if (verb === "evidence") {
      dispatch(
        appendEntry(
          entry(command, [
            "LIVE CEREMONY :: PASS",
            "Deleted source + complete object + 2 of 6 shards",
            "Recovered exact 8 MiB object from 4 surviving shards",
            "SHA-256 matched signed manifest; trusted author verified",
            "3-shard attempt refused without writing an object",
            "COMMAND :: ./scripts/demo-offline-recovery.sh",
          ], "success"),
        ),
      );
      return;
    }

    if (verb === "benchmark") {
      dispatch(
        appendEntry(
          entry(command, [
            "APPLE M5 / GO 1.22 / 64 MiB / 3 RUN MEANS",
            ...benchmarkRows.map(([label, result]) => `${label.padEnd(20, ".")} ${result}`),
            "ERASURE OVERHEAD .... 50.00%",
            "AUDIT SAMPLE ......... 349 / 1024 chunks",
            "SCOPE :: local measurements, not a peer-network comparison",
          ]),
        ),
      );
      return;
    }

    if (verb === "roadmap") {
      dispatch(
        appendEntry(
          entry(command, [
            "PLANNED SYSTEM WORK",
            ...plannedCapabilities.map((item) => `${item.id.padEnd(12, ".")} ${item.title}`),
            "Use PLAN <ID> to inspect one acceptance boundary.",
          ], "warning"),
        ),
      );
      return;
    }

    if (verb === "plan") {
      const capability = plannedCapabilities.find((item) => item.id === target);
      dispatch(
        appendEntry(
          capability
            ? entry(
                command,
                [
                  `PLANNED :: ${capability.title}`,
                  capability.summary,
                  `DONE WHEN :: ${capability.acceptance}`,
                ],
                "warning",
              )
            : entry(command, [`UNKNOWN PLAN :: ${target}`, "Run ROADMAP for valid plan identifiers."], "warning"),
        ),
      );
      return;
    }

    if (verb === "wiki") {
      const section = target ? wikiSections.find((item) => item.id === target) : undefined;
      if (target && !section) {
        dispatch(appendEntry(entry(command, [`UNKNOWN WIKI SECTION :: ${target}`, `VALID :: ${wikiSections.map((item) => item.id).join(", ")}`], "warning")));
        return;
      }
      const destination = section ? `/wiki#${section.id}` : "/wiki";
      dispatch(appendEntry(entry(command, [`OPENING ${destination.toUpperCase()}`], "success")));
      router.push(destination);
      return;
    }

    if (normalized === "run demo" || normalized === "demo") {
      dispatch(
        appendEntry(
          entry(command, [
            "LOCAL SHELL COMMAND",
            "./scripts/demo-offline-recovery.sh",
            "This web console does not execute host commands. Run it in a checked-out release.",
          ], "warning"),
        ),
      );
      return;
    }

    if (verb === "release") {
      dispatch(appendEntry(entry(command, ["OPENING TAG :: v0.1.0-alpha.1"], "success")));
      window.open("https://github.com/AlleyBo55/arkmesh/releases/tag/v0.1.0-alpha.1", "_blank", "noopener,noreferrer");
      return;
    }

    if (verb === "github") {
      dispatch(appendEntry(entry(command, ["OPENING SOURCE :: github.com/AlleyBo55/arkmesh"], "success")));
      window.open("https://github.com/AlleyBo55/arkmesh", "_blank", "noopener,noreferrer");
      return;
    }

    if (verb === "home") {
      dispatch(appendEntry(entry(command, ["OPENING /"], "success")));
      router.push("/");
      return;
    }

    dispatch(appendEntry(entry(command, [`COMMAND NOT FOUND :: ${command}`, "Type HELP for the command index."], "warning")));
  };

  const submit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    runCommand(value);
    setValue("");
    setHistoryIndex(-1);
  };

  const handleKeys = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key !== "ArrowUp" && event.key !== "ArrowDown") return;
    event.preventDefault();
    if (history.length === 0) return;

    const nextIndex = event.key === "ArrowUp"
      ? Math.min(historyIndex + 1, history.length - 1)
      : Math.max(historyIndex - 1, -1);
    setHistoryIndex(nextIndex);
    setValue(nextIndex === -1 ? "" : history[history.length - 1 - nextIndex]);
  };

  return (
    <section className="terminal-frame" aria-label="ArkMesh command console">
      <div className="terminal-titlebar">
        <span>ARKMESH://{pathname === "/wiki" ? "WIKI" : "ROOT"}</span>
        <span className="hidden text-blue-200 sm:inline">TTY-01 · LOCAL ONLY</span>
        <span className="status-pulse">READY</span>
      </div>
      <div
        ref={outputRef}
        className={`${compact ? "h-64" : "h-80 md:h-96"} terminal-output`}
        aria-live="polite"
      >
        {entries.map((item) => (
          <div key={item.id} className="mb-4 last:mb-0">
            {item.command ? <p className="text-white"><span className="text-cyan-300">C:\ARKMESH&gt;</span> {item.command}</p> : null}
            {item.lines.map((line, index) => (
              <p key={`${item.id}-${index}`} className={`${toneClass[item.tone]} leading-6`}>
                {line}
              </p>
            ))}
          </div>
        ))}
      </div>
      <div className="border-t-2 border-blue-300/70 p-3 md:p-4">
        <div className="mb-3 flex flex-wrap gap-2">
          {quickCommands.map((command) => (
            <button key={command} type="button" className="command-key" onClick={() => runCommand(command)}>
              {command}
            </button>
          ))}
        </div>
        <form onSubmit={submit} className="flex items-center gap-2">
          <label htmlFor={`command-${compact ? "compact" : "full"}`} className="shrink-0 text-cyan-300">
            C:\ARKMESH&gt;
          </label>
          <input
            id={`command-${compact ? "compact" : "full"}`}
            value={value}
            onChange={(event) => setValue(event.target.value)}
            onKeyDown={handleKeys}
            autoComplete="off"
            spellCheck={false}
            className="min-w-0 flex-1 border-0 bg-transparent p-1 text-white caret-cyan-200 outline-none placeholder:text-blue-400"
            placeholder="type help"
          />
          <button type="submit" className="command-enter">ENTER</button>
        </form>
      </div>
    </section>
  );
}
