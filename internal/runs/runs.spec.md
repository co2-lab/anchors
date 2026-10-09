<!-- @anchors
  code: PRCRN
  layer: infra
-->
# ProcessRuns — the record of the project's long processes, and the operating system's view of them

> **Code**: `PRCRN`

## Overview

A long process — a test suite, a mutation run, a build, the commit hook — is recorded in
`.anchors/runs/` by whoever launches it: what it runs, where, its process, its output file and,
once it ended, its exit code and summary. Processes nobody recorded are still seen, from the
operating system's process table: each process with its parent, the CPU time it used, its
command line and, where the system gives it, its working directory. A process belongs to the
project by its working directory, or — where the system gives none — by naming the project in
its command line (DESIGN-process-monitor.md).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | an absolute folder | a relative path | the caller (the commands resolve `--root`) |
| a record file | JSON written by this unit | a file half-written, or another JSON | this unit: such a file is skipped |
| the process table | `ps -axo pid=,ppid=,time=,command=` (macOS, Linux), `Win32_Process` (Windows) | — | the operating system |

## Effects

| Effect | Description |
| --- | --- |
| `PRCRN-B01` | A record is written whole under `.anchors/runs/<id>.json`, by an id made of when the run started and its process, and read back as written; the records list oldest first, skipping any file that does not read as a record and the monitor's own state. (`Save`, `Load`, `List`, `NewID`) |
| `PRCRN-B02` | Finishing a run records when it ended, its exit code when known and its summary; a run still running or stalled becomes finished with an exit, and ended without one. (`Finish`) |
| `PRCRN-B03` | Pruning keeps every run still going and, of the ended ones, the latest `keep` that ended within the maximum age; it removes the rest and says how many. (`Prune`) |
| `PRCRN-B04` | The CPU time `ps` prints — `[dd-][hh:]mm:ss[.ss]` — reads as seconds; each `ps` line gives its process, parent, CPU time and the command line with its own spacing; and `lsof`'s working-directory listing gives each process its folder. (`parseCPUTime`, `parsePS`, `parseLsofCwd`) |
| `PRCRN-B05` | A process belongs to the project when its working directory is the root or under it — not a sibling folder that shares its prefix —, or, with no working directory, when its command line names the root. (`InProject`) |
| `PRCRN-B06` | The process tree gives each process's descendants at any depth, its ancestors nearest first, and the CPU time of a process with everything under it. (`NewTree`, `Descendants`, `Ancestors`, `TreeCPU`) |
| `PRCRN-B07` | Listing the system's processes includes the process asking. (`ListProcs`) |
| `PRCRN-B08` | The system says, of a process, its working directory — none on Windows, which gives it to no other process —, whether it is alive — never for a process id below one —, and the one-minute load average where it has one. (`Cwds`, `Alive`, `LoadAverage`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRCRN-I01` | A reader never sees half a record: a record is written to a temporary file and renamed. | saves a record over an existing one and reads it back whole, with no temporary file left |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRCRN-E01` | Finishing a run that has no record. | The error is returned and nothing is written. | A finish with nothing to finish is a caller's mistake; writing a record of nothing would invent a run. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
