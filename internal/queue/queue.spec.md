<!-- @anchors
  code: TSQUT
  updated_at: 2026-09-26
  layer: infra
-->
# TaskQueue — the file-backed queue between "something changed" and "someone works on it"

> **Code**: `TSQUT`

## Overview

The queue decouples noticing a change from working on it. The watcher enqueues a task when it classifies
a change; a worker (`anchors next`, in any terminal, for any AI client) claims the next task, works one
step, and marks it done. Tasks are plain YAML files in `.anchors/tasks/`, one per task, inspectable by
hand, and their state is in the file NAME (`pending__<id>.yaml`, `claimed__<id>.yaml`), so listing is
cheap and a claim is a single atomic file operation. Done tasks leave the live queue for the history in
`.anchors/done/`.

A claim must never be held twice: two terminals running `anchors next` at once must not work on the same
task. The claim is the exclusive creation of the claimed file, which is atomic on every platform; a
pending file left behind by a worker that died mid-claim is never served again, and it is cleaned by
whoever meets it.

A task carries a suggestion of the next step, derived from the kind of file that changed, and the
suggestion is a verb `anchors work` can compose, so whoever pulls the task can build its prompt. The
queue cleans its own noise: a task whose target file no longer exists is removed, and the same target and
step are not enqueued twice while a live task holds them.

Reclaiming returns abandoned claims to the queue. The worker is not an observable process (`anchors next`
exits once it prints the task, and an agent works afterwards), so the only evidence is time: a claim is
presumed alive for four hours. Returning it earlier is the costly mistake (two agents on the same file);
returning it later only costs waiting, and forcing is explicit.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the task ID | a stable, file-name-safe `<seq>-<kind>-<slug>` | — | the watcher, which enqueues |
| the target | a path relative to the project root, or an absolute path | — | the watcher; a target that disappears is cleaned by this unit |
| the claim moment | an RFC 3339 time, or empty | an unparsable time | this unit: an empty or unparsable moment counts as old |
| the task files | files written by this unit | hand-edited files that are not valid YAML | this unit: they are listed as triage tasks (E03) |

## Effects

| Effect | Description |
| --- | --- |
| `TSQUT-B01` | Enqueuing writes the task as `pending__<id>.yaml` in `.anchors/tasks/`, in the pending state. (`Enqueue`) |
| `TSQUT-B02` | Enqueuing a task whose target, suggested step and kind are already held by a live (pending or claimed) task creates nothing — and for a judgment the gate must match too, carried by its ID: two judgment gates on one target are two tasks, and a judgment never collapses into a review of another kind. |
| `TSQUT-B03` | Listing gives every live task sorted by ID, taking each task's state from its file name, and nothing when the queue folder does not exist. |
| `TSQUT-B04` | A task whose target file no longer exists is removed from disk when the queue is listed. |
| `TSQUT-B05` | Claiming takes the first pending task, records the worker and the claim moment, and leaves only its claimed file; an empty queue gives no task. |
| `TSQUT-B06` | A pending file whose claimed file already exists (a worker died mid-claim) is not claimed, and it is removed. |
| `TSQUT-B07` | Marking a task done moves it out of the live queue to `.anchors/done/done__<id>.yaml`, after which the same target and step may be enqueued again. (`MarkDone`) |
| `TSQUT-B08` | Dropping deletes a live task without writing any history. |
| `TSQUT-B09` | Reclaiming returns to pending, with no worker or moment, every claim older than four hours or with no claim moment; the returned tasks can be claimed again. (`ClaimIsOld`) |
| `TSQUT-B10` | Forced reclaiming returns every claimed task, however recent. (`ReclaimForce`) |
| `TSQUT-B11` | The count of recently held claims is the number of claimed tasks a reclaim would keep because they are within the window, so a reclaim of zero can explain itself. (`RecentlyHeld`) |
| `TSQUT-B12` | The next step suggested for a changed file follows its kind: a draft plan goes to plan review, a plan to spec, a spec to code, code to feature, a feature to test, a test to review, a guide to review (of each unit it governs), and any other kind to triage; every suggestion carries its reason. (`SuggestNext`) |
| `TSQUT-B13` | The suggestion for every mapped kind (a draft plan, a plan, a spec, a feature, code, a test, a guide) is a verb `anchors work` can compose; only an unmapped kind gets `triage`, the queue's marker for a step decided by hand. (`ValidWorkArtifact`) |
| `TSQUT-B14` | The pending count is the number of live tasks, pending and claimed. (`PendingCount`) |
| `TSQUT-B15` | A task whose target is an absolute path is checked at that path, so it is kept while the file exists. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TSQUT-I01` | A task is never claimed twice, however many workers claim at once. | enqueues 20 tasks, lets 8 workers claim concurrently until the queue is empty, and checks every task was claimed exactly once |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TSQUT-X01` | A plain reclaim never returns a claim taken within the last four hours. | The worker cannot be observed; a recent claim is presumed alive, and returning it puts two agents on the same file. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TSQUT-E01` | A task to mark done is neither pending nor claimed. | Refused with "task not found"; nothing moves. | Recording as done a task that is not in the queue would fabricate history. |
| `TSQUT-E02` | A task to drop is neither pending nor claimed. | Refused with "task not found". | The caller asked to discard something that is not there; saying it was dropped would hide a wrong ID. |
| `TSQUT-E03` | A task file in the queue is not valid YAML. | Listing returns it among the other tasks as a triage task with the ID and state of its file name and a reason naming the file and the parse error; it can be claimed, done or dropped by that ID. | One corrupted file must not stop every worker from reading the queue, and a skipped file would be a task nobody sees, claims or cleans. |
| `TSQUT-E04` | A task is enqueued under an ID a pending or claimed task already holds (for another target or step). | Refused naming the ID; the live task is untouched. | The ID is the file name: writing would replace the live task's file, and it would vanish from the queue with no trace. |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
