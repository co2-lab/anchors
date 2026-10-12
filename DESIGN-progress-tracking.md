<!-- @anchors
  code: PRTRP
  layer: desenho
-->

# Progress tracking — one file of the project's demands, measured, written only by Anchors

> APPROVED on 2026-10-12 (the user: "quero adicionar uma feature de acompanhamento de progresso … salvar
> um arquivo de progress no root, como anchors.progress.md com um padrão bem específico … o progresso
> deve ser alterado unicamente via comando do anchors"; and "vamos migrar o progress que existe hoje
> para esse novo modelo"). Decisions D01–D04 below.
>
> IMPLEMENTED on 2026-10-12. Measured:
> - W01: 30 processes adding tasks and 15 finishing them at once, on one clone, lost nothing; 20
>   goroutines in the unit test likewise; every percentage is the done weight over the counted one.
> - W02: a line changed by hand is refused on read and on change, and `progress-intact` says so; a
>   merge of two branches that each added tasks keeps both, the newer side renumbered.
> - W03: on a clone of jokenpo, `anchors migrate` moved `plans/skill-game-progress.md` into demand
>   D01, titled as its plan, one stage per phase, the TODO placeholders left out, and deleted it;
>   nothing in the project cites a `-progress.md` any more.

## The problem

Progress today is one `-progress.md` beside each plan: phases with checkboxes, written by hand. It says
what is done, not how much is left — every item weighs the same —, it holds only plans, and any edit
can break its shape. Agents working in parallel edit it at the same time, and the last write wins.
Measured: MIF has no progress file at all; jokenpo has one, still the unfilled template.

## The design

- **One file, `anchors.progress.md`, at the root.** It lists the project's DEMANDS — a plan in
  development, or any other piece of work —, each with optional STAGES and its TASKS. Every task
  carries a WEIGHT: an estimate of its effort, a number, in the unit the project names
  (`progress.unit`, default `points`). Progress is the done weight over the weight not cancelled, per
  stage, per demand, and over the whole file.
- **A delimited format, so Anchors edits it deterministically.** Each line is a record with an id:

      <!-- anchors:progress v1 — written by `anchors progress`; never edit by hand -->
      # Progress — 47% (28/60 points)

      ## D01 — Checkout (plan: plans/0012-checkout.md) — 45% (9/20)
      ### D01.S1 — Foundation — 100% (5/5)
      - [x] D01.T1 — cart spec · 3 · done 2026-10-12 by agent-a
      - [x] D01.T2 — cart code · 2 · done 2026-10-12 by agent-b
      ### D01.S2 — Payment — 27% (4/15)
      - [~] D01.T3 — payment provider · 8 · doing since 2026-10-12 by agent-c
      - [ ] D01.T4 — refund rules · 5
      - [-] D01.T5 — legacy gateway · 2 · cancelled: provider dropped

      <!-- anchors:progress sum=3f9c2a1b7e40 -->

  The percentages are written by the command, never typed. The last line is a hash of the content:
  an edit made by hand does not match it, and is caught.
- **Only commands write it.**

      anchors progress add demand "<title>" [--plan <plan>]
      anchors progress add stage <demand> "<title>"
      anchors progress add task <demand|stage> "<title>" --weight <n>
      anchors progress start <task> | done <task> | reopen <task> | cancel <task> --reason "…"
      anchors progress weigh <task> <n>
      anchors progress show [<demand>] [--json]
      anchors progress summary [--json]

  `summary` is the short view of what is in progress (the user, 2026-10-12): one line per demand
  with its progress, and one per stage when the demand has stages — never the tasks.

  A task can be added at any time — work discovered on the way is the normal case, and the progress
  falls back accordingly. A done task is never deleted: it is cancelled with a reason, or reopened.
- **A mutex for parallel agents.** Every command takes a lock (`.anchors/progress.lock`, created
  exclusively, with a timeout and stale-lock recovery — the map's lock, already measured), reads,
  changes, and writes the file atomically (a temporary file renamed over it). Two agents on the same
  clone never lose each other's change.
- **Gates.**
  - `progress-intact` — the file parses and its hash matches: no hand edit.
  - `progress-honest`, retargeted — a done task that cites a file path is confronted with the disk,
    as the per-plan files are today.
- **Where it shows.** `anchors status` shows the overall progress and each demand's, and warns while
  per-plan progress files are left to migrate; `anchors progress summary` shows the stages.
- **The monitor posts it, when asked** (the user, 2026-10-12: "por default o monitor não deve mostrar o
  summary … a não ser que receba um argumento, mas pode ser configurado como default no anchors.yaml").
  Not a monitor of its own: the project's monitor, with `--summary` (or `monitor.summary: true`),
  says the summary each time it speaks — an event, the heartbeat, its exit — and on a change of
  the file; `--summary=false` silences the block's. No monitor running, no summary.

## The migration (the user, 2026-10-12)

`anchors migrate` moves every plan's `-progress.md` into `anchors.progress.md`:

- each plan becomes a demand with `plan:` pointing at it;
- each `## <CODE>-Wnn` phase becomes a stage;
- each checkbox becomes a task, marked done when it was `[x]`, with weight 1 and flagged as an estimate
  to revise (`anchors progress weigh`);
- the TODO lines of an unfilled template are not carried over;
- the per-plan files are deleted.

Then the old mechanism goes:

- `anchors new progress` goes;
- `anchors new plan` adds the plan as a demand instead of creating a companion;
- the scan's exclusion moves to `anchors.progress.md`;
- the `merge-progress` driver merges the new file;
- the plan guide and `next`'s text point to `anchors progress`.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `PRTRP-W01` | **The file and its commands.** Format, parser and writer, the commands, the weight and its unit, the lock. | Two processes adding and finishing tasks at once lose nothing; every percentage matches the weights. |
| `PRTRP-W02` | **What reads it.** `progress-intact`, `progress-honest` retargeted, the merge driver, `status`. | A hand edit is caught; a branch merge keeps both sides' tasks. |
| `PRTRP-W03` | **The migration.** `anchors migrate` moves the per-plan files and the old mechanism goes. | On a clone of jokenpo: its progress file becomes a demand, and nothing reads a `-progress.md` any more. |

## What NOT to do

- Accept a hand edit. The format is the contract: one command per change, or the file stops being
  editable by Anchors.
- Delete done work. Cancelled and reopened, with a reason, keeps the history honest.
- Count a task without a weight. An estimate is required, even a rough one.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `PRTRP-D01` | The weight. | One number, in the unit the project names (`progress.unit`, default `points`) (the user, 2026-10-12). |
| `PRTRP-D02` | Ids across branches. | The merge driver renumbers the newer side's colliding ids (the user, 2026-10-12). |
| `PRTRP-D03` | `start`. | Records who holds the task and refuses another agent unless `--take` (the user, 2026-10-12). |
| `PRTRP-D04` | The weight of a migrated task. | 1, flagged as an estimate to revise (the user, 2026-10-12). |

## Open Decisions

None.
