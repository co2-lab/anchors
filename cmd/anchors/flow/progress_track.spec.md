<!-- @anchors
  code: LDCML
  updated_at: 2026-10-24
  layer: comando
-->
# LedgerCommand — `anchors progress`, the only writer of the project's progress file

> **Code**: `LDCML`

## Overview

`anchors progress` keeps `anchors.progress.md` (DESIGN-progress-tracking.md): it adds demands — a plan becomes
one with a stage per phase —, stages and weighted tasks; it starts, finishes, reopens, cancels and weighs tasks;
and it shows the file, or its summary — each demand and its stages with their progress, never the tasks. Each
subcommand is one change made under the file's lock, so agents working in parallel lose nothing, and the file
keeps the shape Anchors reads back. The weight's unit is the one `anchors.yaml` names.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project | a directory with `anchors.yaml`, or below one | — | the configuration |
| who | `--by`, or this agent's identity | — | this unit |

## Effects

| Effect | Description |
| --- | --- |
| `LDCML-B01` | `add demand` adds a demand by title; with `--plan`, the plan becomes a demand titled as the plan, with one stage per phase; `add stage` and `add task --weight` add to it; each says the id it gave. (`newProgressTrackCmd`) |
| `LDCML-B02` | `start`, `done`, `reopen`, `cancel --reason` and `weigh` change one task; `start` and `done` record who (`--by`, or this agent) and today, and `start` refuses a task another holds unless `--take`. |
| `LDCML-B03` | `summary` prints the whole progress with the project's unit, then one line per demand and, indented, per stage, with the percent, the weights, the open and in-progress tasks and the weights still estimates; with nothing tracked it says how to start; `--json` prints the rows. (`printSummary`) |
| `LDCML-B04` | `show` prints every task of the file, or of one demand, as the file renders them without its header and sum line; `--json` prints the file's records. (`progressBody`) |
| `LDCML-B05` | The weight's unit is `progress.unit` of `anchors.yaml`, `points` when it names none. |
| `LDCML-B06` | `rename <id> <title>` gives a demand, a stage or a task a new title, keeping its id and history, and says it. |
| `LDCML-B07` | `add task --unplanned` adds a task marked unplanned, where it arises; `unplanned <task>` marks one added before and `--undo` takes the mark off; the summary says, on a demand or stage with unplanned work, `· <n> of <total> <unit> unplanned`. |
| `LDCML-B08` | `move <task> <demand|stage>` takes a task to another stage of its demand, and `drop <stage>` removes a stage left with no task; each says what it did. |
| `LDCML-B09` | The summary's total, the monitor's and the status's leave a frozen plan's demand out, and its row says it is frozen and out of the total. (`printSummary`, FRUNF-D10) |
| `LDCML-B10` | A settled demand — all done, or none started, nothing in progress — is one line: the summary and the monitor leave out its stages, and `show` of the whole file its stages and tasks; `show <demand>` details it, and the file and `--json` keep everything. (`printSummary`, `withoutSettledDetail`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `LDCML-E01` | A change the file refuses — an unknown id, a held task, a title the format cannot read, a hand-edited file — or `add demand` with neither title nor plan. | The command fails with the cause, and the file is not written. | A change half made, or made over a file Anchors cannot read, would break the format the file depends on. |
