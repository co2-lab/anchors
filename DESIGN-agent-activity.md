<!-- @anchors
  code: AGACT
  layer: desenho
-->

# Agent activity — the monitor says when an agent's work stopped moving

> APPROVED on 2026-10-13 (the user: "pensei por que não usar o monitor pra colher status dos
> subagentes também, assim o agente veria que tem um agente que terminou e tem um trabalho parado
> que deveria continuar"). Decisions D01–D03 below; D03 as the user set it: "não é só o agente
> pensando, às vezes o agente fica parado esperando outro processo, então o stale do agente não
> deve contar, mas o do processo que ele espera sim".
>
> IMPLEMENTED on 2026-10-13. Measured:
> - W01: a task started 20 min ago with no run of its holder reads `⏸` once; a run of the
>   holder brings `▶`; a line of the older format, with only the date, still reads.
> - W02: the probe, on real Claude Code sessions, corrected the documentation.
>   - `SubagentStart` carries `agent_id`, `agent_type` and the session's transcript, not the
>     subagent's.
>   - `SubagentStop` carries `agent_transcript_path`, `last_assistant_message` and
>     `background_tasks`, whose `description` names the task. It has no `stop_reason`, so
>     the end reads `⏹`, with no verdict.
>   - A `PreToolUse` inside a subagent carries its `agent_id`.
>   - Both events fire for a background subagent.
>   - Through the hook Anchors installs, a real background subagent read
>     `▶ agent "Run sleep and echo command"` and then `⏹ … ended — <its last message>`, and its
>     command carried its id.
>   - A long-lived teammate could not be probed from a headless session; part 1 covers it,
>     through the task it holds.
>
> CORRECTED on 2026-10-13, from MIF's first session with it: the first tick said "no activity
> for 11h08m" of five tasks whose agents were at work. Three causes, three fixes:
> - The agents work in worktrees of their own, and their runs are recorded there. The monitor
>   now reads every worktree of the repository.
> - An agent that only runs quick commands starts no run. The agent's hook now marks its
>   session active at every command.
> - The silence was counted from the task's start. It is now counted from the last activity,
>   never from before the monitor began to watch, so a first tick says nothing.

## The problem

Measured at MIF on 2026-10-12. The main session ran three long-lived teammates (fix-datas,
cob-front, e2e-mif), each in its own worktree. The progress file showed "2 in progress": D01.T7,
started by fix-datas, and D01.T25, started by cob-front. Only one agent was working:

- fix-datas sat idle, sending "Waiting for the E3 app check". That check had already ended or
  died, and the agent's `until … sleep` loop never fired. Its next step depended on no one.
- The monitor told the truth about processes: `running=[until 1]`. Nothing tied that to the
  task fix-datas held.
- The main session read "in progress" from the progress file and the idle notices as normal.
  The user found it by asking.

Two gaps:

1. **`[~] in progress` says only that someone ran `start`.** It does not say that anyone is
   working on the task now. Nothing compares a held task with the activity of its holder.
2. **The monitor does not see agents.** It sees the processes under the project and the
   commands recorded by Anchors, by `monitor run` and by the agent hook (Bash commands sent to
   the background). A subagent launched by the agent's Agent tool is neither, so its start, its
   end and its silence go unseen.

## The design

### 1. A held task with no activity is an event

- **`start` records the time, not only the date.** The line becomes
  `doing since 2026-10-12 21:31 by fix-datas`. A file written before keeps parsing; a date with
  no time counts from the start of that day.
- **Each run record names who launched it.** The `session` field holds the launcher's
  identity, the same one `anchors progress` writes as a task's holder (`ANCHORS_SESSION`, else
  the agent's identity). Anchors' own commands, `monitor run` and the agent hook all write it.
- **A holder's last activity** is the latest of:
  - a run of theirs starting, ending or growing its output, in any worktree of the repository;
  - a command of theirs, quick ones too, through the agent's hook;
  - a progress change they made;
  - for a subagent, its transcript growing.
- **The event.** A task in progress whose holder has no live run and no activity for
  `monitor.idle` (default 10m) gives one line:
  `⏸ D01.T7 in progress — no activity from fix-datas for 20m`.
  It is said once per stretch of silence. Activity again gives `▶ D01.T7 moving again
  (fix-datas)`.

### 2. Subagents are runs

`anchors install-hooks --agent` also registers the `SubagentStart` and `SubagentStop` hooks.
The hook command records them as runs, the way it records the agent's background commands:

- `SubagentStart` opens a run of kind `agent`:
  - named by the agent's type and, when the launch gave one, its name;
  - id `agent_id`;
  - output: its transcript.
- **The subagent's own processes.** A command the subagent sends to the background is recorded
  by the same hook, and its record carries the subagent's `agent_id`. A subagent waiting on a
  process is not stalled: the process it waits on is watched, with its own stall time.
- **A subagent is stalled** only when its transcript has been quiet for `monitor.idle` (10m)
  **and** no process of its own is alive. While one is, its silence is waiting, not a stall.
- `SubagentStop` closes it with the reason, one of completed, failed, interrupted or stopped,
  and with the first line of its last message:
  - `✓ agent code-reviewer ended (completed, 12m) — "No blocking issues…"`;
  - measured: the hook carries no verdict, so every end reads `⏹`.
- The heartbeat counts agents beside commands: `running=[agent 2, jest 1]`.

**Facts the design rests on, from the Claude Code documentation, to be measured before
implementing (W02 starts with a probe):**

- `SubagentStart` and `SubagentStop` carry `agent_id`, `agent_type`, `session_id` and `cwd`.
- `SubagentStop` adds `last_assistant_message` and `stop_reason`, and it fires for background
  subagents too.
- Where the subagent's transcript path comes from: the docs give a path pattern, but disagree on
  whether a hook field carries it.
- That a `PreToolUse` hook inside a subagent carries its `agent_id`, so its commands can be
  tied to it.
- Whether a long-lived **teammate** fires these hooks at all. MIF's case was a teammate. If it
  fires none, part 1 covers it, through the task its holder took.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `AGACT-W01` | **Held tasks without activity.** The start time in the progress file, the launcher on every run record, a holder's last activity, the `⏸`/`▶` events, `monitor.idle`. | A task started by A, with no run of A for longer than the idle time, reads `⏸` once; a run of A brings `▶`; a file of the old format still reads. |
| `AGACT-W02` | **Subagents as runs.** A probe first: a real session launching a foreground and a background subagent, and a teammate, with each hook's payload saved. Then the two hooks in `install-hooks --agent`, the runs and the end lines, the subagent's commands tied to it, its stall. | The probe's payloads become the tests' fixtures; on a real session, a background subagent's end reads `✓ agent … ended` within one tick; a subagent waiting on a live command of its own is not stalled, and the command is. |

## What NOT to do

- **Guess an agent's state from its words.** "Waiting for the E3 app check" is what fix-datas
  said while nothing ran. Activity is measured, never read from the agent's own words.
- **Kill or restart anything.** The monitor says; the agent watching decides.
- **Count "in progress" as work.** It is a claim, and the monitor confronts it.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `AGACT-D01` | The idle time before a held task reads `⏸`. | 10m, by `monitor.idle` or `--idle` (the user, 2026-10-13). |
| `AGACT-D02` | Is the `⏸` event on by default? | Yes: it is one line, said once per stretch; `monitor.idle: off` disables it. It goes out at the `action` level too (DESIGN-monitor-signal.md) (the user, 2026-10-13). |
| `AGACT-D03` | When a subagent reads as stalled. | After 10m of a quiet transcript with no process of its own alive. A subagent waiting on a process is not stalled; the process is watched, with its own stall time (the user, 2026-10-13). |

## Open Decisions

None.
