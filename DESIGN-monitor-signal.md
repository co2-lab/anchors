<!-- @anchors
  code: MNSIG
  layer: desenho
-->

# Monitor signal — every line the monitor sends is worth the turn it wakes

> APPROVED on 2026-10-13 (the user, on MIF's report: "pro item 1 pode ser os 3"). The three
> changes below. The heartbeat of the `action` level follows the proposal made with them.
>
> IMPLEMENTED on 2026-10-13. W01 proven by the command's tests: a heartbeat with the same
> progress sends no summary, a task done sends it with `was …`, `--summary-every` repeats it.
> W02 proven by its tests over every kind of event; the proof over a recorded stretch of MIF's
> events is still to be made, with MIF's next long session.

## The problem

Measured at MIF on 2026-10-13, in a session of about 3 hours: a Maestro round, with 3 agents in
worktrees and the leading agent watching `anchors monitor --summary` as a stream of events.

- The progress summary came out with every event: each heartbeat, each run that ended. That
  was every 1–2 min.
- Each event wakes a whole turn of the model, which re-reads its context. The line to the agent
  asks it to copy the 9 lines to the user, so it did, every time.
- The progress stood at 50% (82/162) for hours: points move only with a commit and a mark.
- So dozens of turns repeated the block and said "nothing new". The token cost rose, and the
  summary carried no information. The user noticed.
- The project filters the stream with its own `grep` for what calls for action, and the summary
  is built to pass every filter.

The summary going out with everything the monitor says was decided on 2026-10-12
(DESIGN-progress-tracking.md). This design narrows it to when it says something new.

## The design

### 1. The summary goes out when it changes

- **What a change is.** The rows of the summary differ from the last ones sent: a weight done or
  added, a task opened, started, done or cancelled, a stage or demand added. A heartbeat or a run
  that ends with the same rows sends no summary.
- **The first line says by how much.** It carries the progress the last summary showed:
  `▸ Progress — 55% (90/162 points) — was 50% (82/162)`. The agent sees the change in one line.
- **The line to the agent stays.** It comes only with news now, so "copy it to the user" is
  always worth it.
- **A hand edit** is still one warning, said once.

### 2. A reminder, if wanted

`monitor.summary_every` (or `--summary-every`) repeats the summary at that interval even when
nothing changed: `30m` brings it back at most every 30 min. It is off by default, so the summary
goes out only on a change.

### 3. The level: what wakes the agent

`monitor.level` (or `--level`): `all`, the default and what the monitor says today, or `action`.

**`action` says only what calls for a reaction:**

| Says | Leaves out |
| --- | --- |
| a run that ended with a failure (`✗`), or that died (`☠`) | a run that started (`▶`) |
| a run the agent launched, ended either way: it waits for it (`✓`/`✗`) | Anchors' own runs that passed, and processes only the process table saw |
| a stall (`⏸`), and moving again (`▶`), which cancels it | the sparse progress marks of a run (`…`) |
| a test report with failures | a report that passed |
| the progress summary when it changed (part 1) | the summary repeated, unless `summary_every` |
| a held task with no activity (DESIGN-agent-activity.md, once approved) | |
| **the heartbeat, only when it says something:** the project fell idle (`running=[NOTHING]`), once per idle stretch; and a HIGH load, once per stretch of it | the heartbeat that repeats what runs |

`running=[NOTHING]` is the line that catches an agent that stopped (DESIGN-process-monitor.md).
In `action` it is said when the project falls idle, not every 5 minutes. A task held with no
activity is the sharper signal for the case that cost MIF hours.

The level replaces the project's own `grep`: the CONTRIBUTING.md that `init` generates
recommends `--level action` for an agent that leads others, and drops the advice about
filters.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `MNSIG-W01` | **The summary on change, with the delta, and `summary_every`.** | A heartbeat and a run ending with the same progress send no summary; one task done sends it with `was …`; with `summary_every: 1ns` it repeats unchanged. |
| `MNSIG-W02` | **The `action` level.** | Over a recorded stretch of MIF's events, `action` keeps every failure, death, stall, change of progress and idle fall, and drops the rest. |

## What NOT to do

- **Repeat what was said.** A line that carries nothing new still costs a turn.
- **Drop a failure to save a turn.** `action` is narrower, never blind: every terminal state
  that is not a success, and every stall, still goes out.
- **Make `action` the default.** A user who watches a terminal wants to see runs start and end;
  the agent that leads is the one who asks for less.
