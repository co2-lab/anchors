<!-- @anchors
  code: WRQUW
  updated_at: 2026-09-26
  layer: comando
-->
# WorkQueue — list, pull, close and discard the work, from the local queue or from the board

> **Code**: `WRQUW`

## Overview

The queue decouples "something changed" from "someone works on it": the watcher queues, and whoever
works pulls. `anchors queue` lists what is live without committing to anything; `anchors next` claims
the next item and says what to do with it; `anchors done`, `drop` and `reclaim` close, discard and
release tasks.

In local mode the queue is the task files under `.anchors/`. The claim is atomic, so two workers never
take the same task. A queue that is empty while a plan still has specs to be born is a cold start: the
next plan with missing seeds is queued and claimed, because whoever runs `next` is asking "what do I do
now", and a separate command for it would be one more step nobody remembers. Only one plan is seeded
at a time — a real repository with a hundred plans filled the queue with ten fronts nobody pulled. The
count of what a seeded plan still lacks is recomputed from the disk every time it is printed, since a
count stored in the task aged and misled.

In github mode the queue is the board, never the local files — measured: `next` answered "empty queue"
with 84 open cards, and the cycle stopped. Claiming writes who owns a card, so it requires a declared
session. The agent's own card is resumed before anything is asked; otherwise the serialized claim
pipeline is asked and awaited. The printout then says what the card asks for: a card under review asks
for a verdict and names the exact line that ends it; a spec card asks for code, feature, test and
documentation, in that order, with the target resolved through the map; and an agent that does not
decide the product is told to escalate instead of asking whoever runs it.

`anchors next` prints the project's notifications on top in both modes (see Notifications).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the worker | a name, or none (then pid@host) | — | this unit: fills the default |
| the session | `ANCHORS_SESSION`, required to claim from the board | an unset session in github mode | this unit: `next` refuses |
| the repository | `workflow.repo` in github mode | an empty repository | this unit: refuses, never infers it from the remote |
| the task to close | an id, or one filter: `--file`, `--kind` or `--all` | no id and no filter | this unit: refuses |
| the card's body | lines "Código: `CODE`" and "Unidade: `path`" written by the pipeline | — | the pipeline that opens the card |

## Effects

### The local queue

| Effect | Description |
| --- | --- |
| `WRQUW-B01` | `queue` lists every live task with its state, the file and kind that changed, the suggestion and who claimed it, then hints `reclaim` when any is claimed and `drop` when any waits for triage; an empty queue says so. |
| `WRQUW-B02` | `next` claims the next pending task for the worker and prints it with how to close it (`anchors done <id>`), and reminds how many informative gates are declared. |
| `WRQUW-B03` | On an empty queue, `next` seeds the first plan that still cites a spec not on disk, says so and claims it; when every plan is fulfilled it says the queue is empty. |
| `WRQUW-B04` | A spec a plan cites exists when its path exists, or when its file name alone matches exactly one file of the project. |
| `WRQUW-B05` | A seeded plan's task is printed with how many of its specs do not exist yet, counted at print time; a task whose reason already carries a count gets none added, and a task that is not a plan seed gets none. |
| `WRQUW-B06` | `done` closes a task by id, or in batch every task of a file, of a kind, or all of them; without an id or a filter it refuses, a filter that matches nothing says so, and an unknown id is an error. |
| `WRQUW-B07` | `drop` deletes a live task without archiving it, and dropping a missing task is an error. |
| `WRQUW-B08` | `reclaim` returns to the queue what dead workers left claimed, leaves a task claimed a moment ago with its worker, and explains a zero by counting those; `--force` takes them back too. |
| `WRQUW-B09` | The default worker is `<pid>@<host>`. |

### The board

| Effect | Description |
| --- | --- |
| `WRQUW-B10` | In github mode `next` refuses without `ANCHORS_SESSION`, saying how to set it. |
| `WRQUW-B11` | The agent's identity on the board is `<host>/<session>`; without a session it falls back to the OS user and says so on standard error. |
| `WRQUW-B12` | In github mode `next` resumes the card this agent already owns — saying to finish it before taking another — and asks nothing of the claim pipeline. |
| `WRQUW-B13` | A claim that ends without a card names the run to follow: a timeout says the claim is still pending and that running `next` again waits for it; a successful run with no card says there is no free card; neither is an error. |
| `WRQUW-B14` | The claimed card is printed with its number, title, state without its prefix, and owner. |

### What the card asks for

| Effect | Description |
| --- | --- |
| `WRQUW-B15` | A plan card asks for every spec the plan seeds; a spec card asks for code, feature, test and documentation in that order, with the target being the map's file for the body's code, or else the first folder the body names; any other card points to its body. |
| `WRQUW-B16` | A spec card lists the documentation the project requires for that unit's changes, and nothing when the project declares none. |
| `WRQUW-B17` | An agent that does not decide the product is told to escalate choices with `anchors escalate … --for-user` instead of asking whoever runs it; one that decides the product is not. |
| `WRQUW-B18` | A card under review asks for the review of its target and ends with the two verdict lines for this agent — `anchors-review: approved by <agent>` and `anchors-review: rejected by <agent>` — never with a pull request of the reviewer's own; a review card that names no unit points to the pull request that references it. |
| `WRQUW-B19` | Any other card ends by naming `anchors pr-body --cards <n>` for the pull request's body, saying it brings the `Refs` that links the card and that the card stays open — never that it closes it. |

### The role question

| Effect | Description |
| --- | --- |
| `WRQUW-B20` | A role already declared is kept without asking; without a terminal (a pipe or the null device) nothing is asked and nothing recorded, and the output says how to declare the role. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `WRQUW-I01` | A cold start seeds at most one plan, however many still have work. | two plans with missing specs are on disk, and exactly one task is seeded |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `WRQUW-X01` | `queue` claims nothing. | It is how the conversation and the human look at the work without committing to it; claiming there would take work nobody pulled. |
| `WRQUW-X02` | In github mode the repository is never inferred from the remote: without `workflow.repo`, `next` refuses. | From a fork the inferred repository is another one, and a write in the wrong place is not undone by a revert. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `WRQUW-E01` | The claim run ends in failure without a card. | `next` fails naming the run, its conclusion and how to see its log. | A failed claim is a broken pipeline, unlike a claim that found no free card. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/queue` | `List`, `Claim`, `Enqueue`, `MarkDone`, `Drop`, `Reclaim`, `ReclaimForce`, `RecentlyHeld`, `SuggestNext` | the local queue |
| DEP2 | `internal/board` | `Client.Mine`, `AskAndWait`, `ClaimOutcome` | the board and the claim pipeline |
| DEP3 | `internal/settings` | `Load`, `Save`, `HandlesUserIssues` | the agent's declared role |
| DEP4 | `internal/scan` | `Walk`, `LayerOfUnit` | scan — the plans' seeds and the unit's layer |
| DEP5 | `cmd/anchors/flow/notifications.go` | `printNotifications` | comando — the message on top of `next` |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
