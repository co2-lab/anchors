<!-- @anchors
  code: DMLDD
  updated_at: 2026-10-24
  layer: infra
-->
# DemandLedger — the project's progress file: demands, stages and weighted tasks, written only by Anchors

> **Code**: `DMLDD`

## Overview

`anchors.progress.md` holds the project's progress (DESIGN-progress-tracking.md): its demands — a plan in
development, or any other piece of work —, their optional stages and their tasks, each with a weight that
estimates its effort. Progress is the done weight over the weight not cancelled. The format is delimited — each
line a record with an id, the percentages computed, the content's sum on the last line —, so Anchors reads it
back deterministically and catches any hand edit. Every change is made under the file's lock and written
atomically, so agents changing it at once lose nothing. This unit also merges two branches' files by id, adds a
plan as a demand with one stage per phase, and moves the per-plan `-progress.md` files it replaces.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file | the file this unit wrote, or none | a file whose sum does not match its content, or with a line this unit does not write | this unit: refused as edited by hand |
| a title or reason | one line without the separator ` · ` | an empty, multi-line or separator-holding text | this unit: refused |
| a weight | a positive number | zero, a negative number | this unit: refused |

## Effects

| Effect | Description |
| --- | --- |
| `DMLDD-B01` | A rendered file opens with the header that says who writes it, gives the whole progress with the unit, then each demand (`## Dnn — title (plan: path) — P% (done/total)`), its stage-less tasks, its stages (`### Dnn.Sn — …`) and their tasks (`- [ ] Dnn.Tn — title · weight`, ` (estimate)` when it is one, and who and when for a task started or done, the reason for a cancelled one), and closes with the sum of the content; parsing it gives back the same file. (`Render`, `Parse`, `File`, `Demand`, `Stage`, `Task`, `State`) |
| `DMLDD-B02` | Progress is the done weight over the weight of the tasks not cancelled, rounded down, per stage, per demand and over the file; with no weight counted it is `—`. (`Totals`, `Percent`) |
| `DMLDD-B03` | Demands, stages and tasks are numbered in order — a task's number is unique across its demand, its stages included —; a task joins a demand or one of its stages, with a positive weight, at any time. (`AddDemand`, `AddStage`, `AddTask`) |
| `DMLDD-B04` | A task is started by someone — refused when another holds it, unless taken over —, done, reopened, or cancelled with its reason; a done task is never deleted; weighing a task sets its weight, which is an estimate no more. (`Start`, `Done`, `Reopen`, `Cancel`, `Weigh`, `Task`) |
| `DMLDD-B05` | A change is made under the file's lock, on the file as it is on disk, and written atomically: changes made at once by several processes are all kept; an absent file is an empty one. (`Update`, `Load`, `Path`, `FileName`) |
| `DMLDD-B06` | The summary is one row per demand and one per stage, with the done and total weight, the percent, the open and in-progress tasks and the weights still estimates — never the tasks. (`Summary`, `Row`) |
| `DMLDD-B07` | Two branches' files are merged over their base: what each side added is kept; a task both sides moved takes the further state (done, then cancelled, then in progress, then to do); a weight changed on one side only takes that side's; an id both sides created for different work is renumbered on their side. (`Merge`, `IsFile`) |
| `DMLDD-B08` | A plan becomes a demand titled as the plan's title (or the one given), with one stage per phase its headers declare; a plan already tracked is refused. (`PlanDemand`) |
| `DMLDD-B09` | A per-plan `-progress.md` beside a plan becomes its plan's demand: each heading a stage, each checkbox a task weighing 1 as an estimate, done when it was checked; the placeholder items of an unfilled template are left out; a plan already tracked keeps its demand; the migration of a project moves every such file, titles each demand as its plan, and deletes the file; and a dry run writes nothing. (`LegacyFiles`, `MigrateLegacy`, `MigrateProject`, `Migrated`, `LegacySuffix`) |
| `DMLDD-B10` | A task in progress records when it was started to the minute (`doing since 2026-10-12 21:31 by a`), and says when, in local time; a task of a file written before, with only the date, still reads, counted from the start of that day. (`StartedAt`, `Since`) |
| `DMLDD-B11` | A demand, a stage or a task takes a new title, with the same rules as a new one; its id, its tasks, their state and their history stay; an id that names nothing is refused. (`Rename`) |
| `DMLDD-B12` | A task can be marked unplanned — work the plan did not list, kept where it arose —, in its line beside the weight (`2 (unplanned)`, `2 (estimate, unplanned)`), and the mark taken off; the summary says each demand's and stage's unplanned weight, cancelled tasks left out; a merge keeps the side that changed the mark; a mark other than estimate or unplanned is refused; the flow's advice says to record such work before doing it. (`MarkUnplanned`, `UnplannedAdvice`) |
| `DMLDD-B13` | A task moves to another stage of its demand, or to the demand itself, keeping its id, state, marks and history; a stage or demand that is not its demand's is refused. A stage left with no task can be removed; one that still holds a task is refused. (`Relocate`, `DropStage`) |
| `DMLDD-B14` | Loading the file reads, for each demand with a plan, the freeze that covers the plan — a frozen folder of anchors.yaml, or the plan's own header mark —, never writing it into the file; the counted tasks leave a frozen demand's out, and its summary rows, its stages' included, carry the freeze. (`markFrozen`, `CountedTasks`, FRUNF-D10) |
| `DMLDD-B15` | A demand, and a summary row, is settled when nothing of it is in motion: all of its counted weight done, or none of it — no task in progress either way. (`Demand.Settled`, `Row.Settled`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DMLDD-E01` | The file's content does not match its sum, or holds a line this unit does not write. | The read fails (`ErrEdited`, or the line), and no change is written over it. | A hand-edited file is not one Anchors can edit deterministically; writing over it would lose the edit or keep it broken. |
| `DMLDD-E02` | A title or reason that is empty, spans lines or holds ` · `; a weight that is not positive; an id that does not exist. | The change is refused, naming the cause. | The format could not read it back. |
