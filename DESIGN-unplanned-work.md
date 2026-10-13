<!-- @anchors
  code: UNPLW
  layer: desenho
-->

# Unplanned work — what the agent does on the way, outside the plan, is tracked too

> APPROVED on 2026-10-13. Decisions D01–D03 below (the user: "durante o desenvolvimento do plano
> apareceram algumas atividades que o agente fez mas não estavam no plano, então ele não registra
> nem no progress… queria que isso fosse uma sugestão, fizesse parte do fluxo no anchors pra não
> ter que deixar cada agente/projeto tirando suas próprias conclusões").
>
> IMPLEMENTED on 2026-10-13. Measured with headless Claude Code sessions, each holding a task
> whose way crossed a bug outside it:
> - In a project `init` set up, asked only to implement the task and fix the bug, the agent
>   followed AGENTS.md. It recorded the fix as `D01.T2 … · 1 (unplanned)` and finished both.
> - With the command's help alone, it did the same.

## The problem

Measured at MIF on 2026-10-13, while a plan was developed:

- The agent did work the plan did not list: fixes found on the way, a missing helper, a test
  the plan forgot.
- It recorded none of it in the progress file. The progress stayed as if that work had not
  happened, and its weight never entered the count.
- The user had to tell the agent to add a stage for it ("Não planejadas") and record the tasks,
  so the progress was computed again.
- Anchors says nothing about such work. Each agent, and each project, decides alone, or does
  nothing.

## The design

### 1. The mark: `--unplanned`

    anchors progress add task D01.S2 "fix the date parser" --weight 2 --unplanned
    anchors progress unplanned D01.T31          # a task already added
    anchors progress unplanned D01.T31 --undo

- **The task stays where it arose,** in the demand or the stage it belongs to, marked in its
  line like an estimate: `- [x] D01.T31 — fix the date parser · 2 (unplanned) · done …`. New
  tasks that appear on the way are not set apart, so they are not hidden (the user, D01).
- `unplanned` marks a task added before, such as MIF's 20 tasks; `--undo` takes the mark off.
- **`move <task> <stage>`** takes a task to the stage it arose in, keeping its id and history.
  **`drop <stage>`** removes a stage left with no task. At MIF the found work stood in a stage of
  its own, set apart from the stages it arose in. The user wants it inside them, so each stage
  shows how much of it was not planned. Without `move`, the only way was to cancel the tasks and
  add them again.

### 2. Where it shows

- **`summary`** says, on each demand and stage with unplanned work, its weight against the
  whole: `D01 Checkout — 41% (75/182) · 12 open · 20 of 182 points unplanned`.
- **`show`** marks each task.
- A plan that keeps growing on the way is visible, not hidden in its percentage.

### 3. The flow says it

The agent learns it where it learns the rest:

- **The plan guide** (`anchors guide plan`) and **the work prompt** (`anchors work`):
  > Work the plan does not list, done for it, is recorded before it is done:
  > `anchors progress add task <demand> "<what>" --weight <n> --unplanned`. The plan is not
  > rewritten for it; the progress counts it.
- **The CONTRIBUTING.md and AGENTS.md** that `init` generates say the same, in one line, beside
  `anchors progress`.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `UNPLW-W01` | **The mark.** `add task --unplanned`, `unplanned <task> [--undo]`, the mark in the task's line, the share in `summary`. | A task added unplanned, and one marked after, read back marked; undone, unmarked; the summary says the unplanned weight of the demand and of the stage; a file written before still reads. |
| `UNPLW-W02` | **The flow.** The plan guide, the work prompt, the generated CONTRIBUTING.md and AGENTS.md. | Each says it; a headless session told to fix something outside a plan records it with `--unplanned` before doing it. |

## What NOT to do

- **Rewrite the plan for every finding.** The plan says what was decided; the progress says what
  was done. Unplanned work belongs in the second.
- **Guess unplanned work from commits.** What the agent did is the agent's to say; a gate that
  infers it from the files changed would guess wrong.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `UNPLW-D01` | One unplanned stage per demand, or a mark per task? | **A mark per task** (the user, 2026-10-13: "são tasks novas que aparecem sem serem planejadas, não quero esconder elas, quero dar visibilidade"). |
| `UNPLW-D02` | The mark's word. | `unplanned` (the user, 2026-10-13). |
| `UNPLW-D03` | Does the summary say the unplanned share? | Yes (the user, 2026-10-13), on each demand and stage with unplanned work. |
