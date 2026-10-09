<!-- @anchors
  code: PRMNP
  layer: desenho
-->

# Process monitor — an agent waiting on a long process hears when it moves, stalls, dies or ends

> IMPLEMENTED on 2026-10-09 (W01–W03, one delivery). An agent that launches a long process — a test suite, a mutation run,
> a build, a commit hook — waits on a watch that speaks only at the end. When the process dies
> or hangs, the end never comes: the agent waits for nothing and the user stares at a screen
> that says nothing. `anchors monitor` is the command the agent watches: it reads the project's
> processes on a loop, records what it sees for the user, and prints one line for each event
> that calls for a reaction.

## The problem

Measured at MIF, which built three monitors of its own as inline shell loops (reported on
2026-10-09):

- **A watch that speaks only at the end misses the costly cases:** a process gone, and the
  next step never started. Agents finished and sat idle for **40 min to 4 h** before anyone
  noticed. MIF's working answer was a heartbeat every 5 min that says `rodando=[NADA]`
  ("nothing running") with work pending.
- **Death looks like work.** A Jest run that ended with no artifact looked the same as one still
  running, until `ps` showed it gone. The exit code of a process the watcher did not launch —
  an agent's run, a commit hook — is invisible.
- **Broad matching lies.** Patterns over process names and file times matched other sessions'
  processes (a Stryker run in another project, another session's watch script), a Metro server
  kept alive between runs, and unrelated writes in the scratchpad.
- **The watch expires** (30 min in Claude Code). A re-arm by hand lost the cursor, and either
  repeated events or skipped them.
- **Load explained most flakes** (load average 45–94: Maestro waits timed out), and the monitor
  did not show it.
- **The commit hook is a long run** (over 10 min at MIF), and nobody saw it run.

The agent does not use Anchors commands alone: it runs `npx jest`, `maestro test`, shell
scripts and binaries directly. A record kept only by Anchors commands would miss the common
case.

## The design

- **`anchors monitor` is the sentinel, run by the agent.** Anchors cannot wake an agent; the
  agent's watch tool can, one notification per output line. The agent runs `anchors monitor` in
  it. Each tick the monitor reads the project's processes and the outputs it knows of, writes
  the state, and prints a line only for an event. Nothing runs in the background when nobody
  watches.
- **The project's processes, from the operating system.** Every process whose working directory
  is under the project root — or whose command line names it — with its children. Whoever
  launched it does not matter: an Anchors command, the agent's shell, a script, a hook. MIF's
  false readings came from matching by name; filtering by the project's directory removes them.
  macOS and Linux read the process table (`ps`/`/proc`); Windows reads `Win32_Process`.
- **A record per run, written by whoever launches it, when it can.** The record holds the
  command, the directory, the pid, the start, the output file, the expected artifacts (a JUnit
  report), and at the end the exit code and the summary. Three writers:
  1. Anchors' own long commands (`test`, `mutation`, `map build`, `check --all`, the commit
     hook's run) write it themselves.
  2. **An agent hook** (`anchors install-hooks --agent`) registers each long or background
     command the agent runs — its command line and output file — and, when the harness reports
     it, its exit.
  3. `anchors monitor run -- <command>` wraps anything else and writes the same record (under
     `monitor`: a root `run` would read as the watcher's `watch run`).

  A process with no record is still seen, by the operating system: the monitor says it ended,
  and reads pass or fail from its output's final summary when it knows the output file.
- **The states, and what each one means:**
  - *running*: output grows, a report lands, or CPU time advances;
  - *stalled*: no output growth, no report, **and** near-zero CPU for X minutes — X per kind of
    run, since a Maestro flow takes 2–3 min and a Jest run longer;
  - *died*: the pid is gone with no exit record and no final summary;
  - *finished*: the exit record or the final summary says pass or fail.
- **One line per event, and a heartbeat:**

  | Event | Line |
  | --- | --- |
  | started | `▶ jest (pid 812) apps/mobile — by the agent` |
  | report landed | `✗ maestro AuthScreen-A01: 1 failure — <failure text>` |
  | progress, sparse | `… jest 340/1673, 2 failed, 6 min` |
  | stalled | `⏸ maestro: no output or CPU for 8 min (pid 901)` |
  | died | `☠ jest gone without a summary — last lines: …` |
  | finished | `✓ jest 1673/1673 (exit 0)` / `✗ jest 3 failed (exit 1)` |
  | heartbeat | `14:05 running=[jest 1, maestro 1] load=12.4` / `running=[NOTHING] — 2 runs ended since 13:40` |

  The heartbeat is MIF's lesson: "nothing running" is the line that catches an agent that
  stopped. The load average goes on the heartbeat, flagged when high. Every terminal state is
  a line, the monitor's own failure to read included: silence means only that nothing changed.
- **The state survives the watch.** `.anchors/runs/` holds each run's record and the monitor's
  cursor. A re-armed monitor resumes where the last one stopped: what happened in between —
  a death, a report — is its first line, and nothing is said twice.
- **Two modes:** by default it streams until stopped; `--until-done` exits when the runs it
  watches have all ended, with their outcome as the last line.
- **Configurable timing.** The tick, the heartbeat, the progress spacing and the stall
  thresholds per kind have defaults in the code, a `monitor:` block in `anchors.yaml`, and a
  flag on the command (`--every 30s`, `--heartbeat 5m`). Proposed defaults, from MIF's
  experience: tick 30 s; heartbeat 5 min; stall 6 min for an e2e flow, 10 min for a unit
  suite, 5 min for a map build or anything else.
- **For the user:** `anchors status` gains a "running" section read from `.anchors/runs/` —
  what runs, its progress and state, and what ended since. The user knows where things stand
  without asking the agent.
- **For the agent:** `anchors guide` teaches the habit: launch the long command in the
  background, then watch `anchors monitor`, react to each line, and re-arm on expiry.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `PRMNP-W01` | **The run record and the monitor's core.** `.anchors/runs/` records written by Anchors' long commands and the commit hook; `anchors monitor` reads the operating system's processes under the project, with their children; it prints started, stalled, died and finished, plus the heartbeat with the load; it keeps a cursor; it has the timing configuration. | On a clone of MIF: a Jest run killed with `kill -9` reads *died* within one tick; a process stopped with `SIGSTOP` reads *stalled* after its threshold; another project's process is not seen; a re-armed monitor reports the death that happened while it was down, once. On macOS, Linux and Windows in CI. |
| `PRMNP-W02` | **Outputs and reports.** A JUnit report landing is an event, with its failures; the progress line; pass or fail read from a known output file's final summary, by patterns per runner in `anchors.yaml` with defaults for the common ones. | On a clone of MIF: a Maestro flow failing mid-run is a line with its failure text before the run ends; a Jest run's final summary gives its pass/fail without an exit record. |
| `PRMNP-W03` | **Who launches.** The agent hook (`install-hooks --agent`) that records each background command, its output file and its exit; `anchors monitor run --` for anything else; `--until-done`; the "running" section of `anchors status`; the guide. | An agent session on a clone: a background `npx jest` the agent launches carries its command line, output file and exit code into the monitor's line and into `anchors status`. |

All three are one delivery, at the user's request (2026-10-09): the feature is complete only
with the hook.

## What NOT to do

- Match processes by name across the machine. The project's directory is the filter; a name
  pattern saw other sessions' runs as this one's.
- Print the log. A line is an event, a sparse progress mark or the heartbeat — a watch that
  floods is stopped by the harness.
- Stay silent on a failure. Every terminal state, and the monitor's own trouble reading, is a
  line.
- Run a daemon nobody asked for. The monitor runs while the agent watches it; the records
  written by the commands that launch keep the history.
- Tailor it to one runner. Runner kinds, thresholds and summary patterns are configuration with
  defaults, not code per tool.

## What was proven

The phases' proofs ran on a scratch project with real processes, not on a clone of MIF: a clone
has no `node_modules`, and linking the real ones is how a copy once wrote into the real tree
(see the lessons of 0.1.238). A stand-in `jest` — a script printing a line a second — ran
under the monitor:

- killed with `kill -9` under `anchors monitor run`: `☠ jest gone without a summary`, with its
  last lines;
- stopped with `SIGSTOP`: `⏸ jest: no output or CPU for 17s`; continued: `▶ jest moving again`;
- the same script running in another folder: not seen;
- killed while no monitor ran: the re-armed monitor's first line, once, and nothing the next
  time;
- the agent hook fed the JSON Claude Code sends: a background command recorded with its output
  file and ended by its summary; a foreground one ended with its exit code; a short one not
  recorded.

## Lessons

- **A child that ends is no sign of life.** The first reading counted any change in the
  processes under a run as activity: a script stopped with `SIGSTOP` looked alive because its
  last `sleep` ended. Only a new process under the run counts.
- **CPU alone misses quiet scripts.** A script that waits on short commands uses almost no CPU
  of its own; the new processes it starts are what shows it moving.
- **The wrapper names what it wraps.** `anchors monitor run -- ./bin/jest` matched the jest
  runner and was found as a second run; Anchors' own wrapping and watching processes are
  excluded.
- **The hook's answer is read in words.** Claude Code documents no field for a background
  command's output file; the answer says it ("Output is being written to: …"), and the hook
  reads it there. No hook fires when a background command ends: the monitor reads its end
  from the process table and its output.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `PRMNP-D01` | What is the command called? | `anchors monitor`; the user's view stays in the existing `anchors status`, which gains a "running" section (the user, 2026-10-09). |
| `PRMNP-D02` | Does the agent hook ship in the first delivery? | Yes — the feature ships complete (the user, 2026-10-09). |
| `PRMNP-D03` | Where does the timing come from? | Defaults in the code, a block in `anchors.yaml`, and an override on the command (the user, 2026-10-09). |
| `PRMNP-D04` | Should `.anchors/runs/` keep ended runs, and for how long? | The last 50 runs, or 7 days, whichever keeps less: enough for "what ended since" and a post-mortem, without growing forever (the user, 2026-10-09). |
| `PRMNP-D05` | Which agents does the hook support first? | Claude Code (its hooks run before and after each command); other agents are still covered by the operating system's view (the user, 2026-10-09). |

## Open Decisions

None.
