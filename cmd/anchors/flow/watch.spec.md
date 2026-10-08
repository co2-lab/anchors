<!-- @anchors
  code: WTCHA
  updated_at: 2026-10-08
  layer: comando
-->
# Watch — the background watcher that turns "a file changed" into "there is work in the queue"

> **Code**: `WTCHA`

## Overview

`anchors watch` observes the repository and, for each file that changes, classifies the change and
queues a task with the suggested next step. It does not run checks nor call any model: whoever works
pulls the task with `anchors next`. That is the inversion — the watcher only turns "it changed" into
"there is work".

It runs in the background and is controlled by subcommands that read and write the daemon's state in the
project: start, status, pause, resume, stop and logs; `run` is the loop itself, which the daemon runs in
the foreground.

What becomes work follows the project's structure. A governed file queues the next piece of its
unit's chain — skipping the pieces its layer waives and the pieces that already exist, because a queue
that tells a worker to rewrite a finished feature loses the trust of whoever pulls it. What is not work
never becomes a task: files outside the structure, files already gone, editor temporaries and what the
project's ignore list declares disposable. A delivery record at the root of `changes/` is the trigger
of the review: a plan's delivery asks for the review of the whole, and a unit's delivery waits until the
unit has code and test — the review attacks by execution, and there is nothing to mutate before that —
unless the unit's layer waives the test.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the changed path | a path relative to the project root | — | the loop: it relativises each event's path |
| the project's structure | the loaded configuration and map | a missing configuration or an unreadable map | `run`: refuses naming which one failed |
| the daemon's state | the pid, pause flag, metadata and log files of the project | — | the daemon runtime |

## Effects

### Controls

| Effect | Description |
| --- | --- |
| `WTCHA-B01` | `status` says "stopped" without a live pid; with one it says running or paused, with the pid, followed by the recorded metadata. |
| `WTCHA-B02` | `start` refuses while a watcher is already running. |
| `WTCHA-B03` | `pause` refuses when no watcher is running, and otherwise leaves the pause flag; `resume` removes it. |
| `WTCHA-B04` | `stop` terminates the running watcher and removes its pid file; stopping a stopped watcher is an error. |
| `WTCHA-B05` | `logs` prints the log as it is, and without a log it fails asking whether the watcher ever ran. |
| `WTCHA-B06` | `run` loads the configuration and the map before looping, and fails naming which one could not be loaded. |

### The loop

| Effect | Description |
| --- | --- |
| `WTCHA-B07` | A file created after the start becomes a task, including a file in a folder born after the start; a termination signal ends the loop cleanly, reporting it, and the pid file is removed. |
| `WTCHA-B08` | The whole tree is watched except the directories the ignore list skips, and the first-level files of a new folder are swept once, so those born before the watch took effect are not lost. |

### What becomes work

| Effect | Description |
| --- | --- |
| `WTCHA-B09` | A change in a governed file queues the next piece of its unit's chain, skipping a piece that already exists, and reports the task and the queue's size. |
| `WTCHA-B10` | A piece the unit's layer waives is skipped to the next step of the chain; the layer is the one of the unit's code file, even when the change is in a derived piece. |
| `WTCHA-B11` | A file outside the structure, a file already gone, a file the project's ignore list covers and an editor's temporary queue nothing. |
| `WTCHA-B12` | A delivery record at the root of `changes/` for a plan queues the review of the whole plan; for a unit it queues the review only once the unit has code and a test, and otherwise says the review waits for the unit to close; a record under `changes/reviewed/` queues nothing. |
| `WTCHA-B13` | A unit whose layer waives the test is reviewed as soon as its code exists, a test beside the code in any supported language closes the unit, and a record with no readable unit does not hold the review. |
| `WTCHA-B14` | A delivery record is any `.md` directly under `changes/`; one without a unit line is still a delivery. |
| `WTCHA-B15` | A task's id is the path's slug, the step and a short hash: the same path and step give the same id, another step another id, and no id holds a `/`. |
| `WTCHA-B16` | A governed file the map does not have enters it the moment the watcher sees it — read alone, with its unit —, and the watcher's copy of the map is reloaded; without a map on disk, nothing is written. |
| `WTCHA-B17` | A governed file the watcher sees changed moves to its new revision in the watcher's copy of the map, carrying its evidence when what its proofs read did not change — a flag of the chains, a spec's navigation or date —; a change to what they read leaves them stale. (`updateNodeRev`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `WTCHA-I01` | A task still in the queue is never queued twice: the same change again is reported as already in the queue. | the same change is handled twice and the queue keeps one task |
| `WTCHA-I02` | Every change is handled by the loop itself, one at a time — the debounce timer only hands the path back to it. A change still inside its debounce window when the loop ends is handled before the loop returns, and nothing is handled after it. | a change pending at SIGTERM is queued, and its node updated, by the time the loop returns; the node does not change afterwards; `go test -race` is clean |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `WTCHA-X01` | Handling a change only queues: it writes nothing outside the task queue. | Whoever works pulls the task and runs what it asks; a watcher that acted on the change itself would work behind the worker's back. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `WTCHA-E01` | REF[WTCHA-B06]: the configuration or the map cannot be loaded, and B06 fails `run` naming which | — | — |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
