<!-- @anchors
  code: MNCMD
  layer: comando
-->
# MonitorCommand — `anchors monitor`, `anchors monitor run` and the agent's hook: watching the project's long processes and recording them

> **Code**: `MNCMD`

## Overview

`anchors monitor` is the command an agent watches while a long process runs: on a loop, it reads
the project's processes and the records of its runs, prints one line per event, and keeps its
memory in `.anchors/runs/` (DESIGN-process-monitor.md). The records come from whoever launches a
run: Anchors' own long commands record theirs; `anchors monitor run -- <command>` wraps anything else;
and the agent's hook — which `install-hooks --agent` puts in the project's Claude Code settings —
records the agent's long commands before they run and, after, the file their output goes to or
their exit.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the timing flags | Go durations (`30s`, `5m`) | a duration that does not read, or zero | this unit: the command fails naming the value |
| the hook's input | the JSON Claude Code sends a command hook | anything else | this unit: it is ignored, silently |
| the wrapped command | a program and its arguments | — | the caller |

## Effects

| Effect | Description |
| --- | --- |
| `MNCMD-B01` | `anchors monitor --once` reads once, prints the lines of what it read and writes its memory; `--until-done` exits when nothing runs — saying so — or once every run it watched ended. (`runMonitor`) |
| `MNCMD-B02` | The timing comes from the defaults, then the `monitor:` block, then the flags; a flag that is not a positive duration, or a block that does not read, fails the command naming it. |
| `MNCMD-B03` | `anchors monitor run -- <command>` runs the command with its output passed through, records its process, a copy of its output and its exit, and ends with the command's exit code. (`recordedRun`, `ExitCode`) |
| `MNCMD-B04` | Before a command of the agent, the hook records it by the tool call's id when it runs in the background, has a timeout over two minutes, or is a runner's; a short command, one of Anchors' own long commands, the monitor and a command outside a project record nothing. The hook never prints and never fails. (`agentHook`) |
| `MNCMD-B05` | After a command of the agent, the hook records the file the output of a background command goes to, as the answer names it, and ends a foreground command with its exit code. |
| `MNCMD-B06` | Anchors' long commands — test, mutation, check, verify, ingest, map build, docs build — record their run in a project that has an anchors.yaml, and close it with their exit code; another command, or a folder with no anchors.yaml, records nothing. (`BeginOwnRun`, `EndOwnRun`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MNCMD-I01` | The agent's hook never blocks the agent: whatever its input, it prints nothing and returns no error. | feeds it invalid JSON, another tool, a command outside a project and a valid call, and checks each prints nothing and succeeds |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MNCMD-E01` | A timing flag that does not read, or is not positive. | The command fails before reading anything, naming the value. | A monitor on a tick it guessed would flood the agent or say nothing. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
