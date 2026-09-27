<!-- @anchors
  code: TSRTS
  updated_at: 2026-09-26
  layer: comando
-->
# TaskStatusReport — the format of the round's report: where the task is, the verdict, what is missing, and what comes next

> **Code**: `TSRTS`

## Overview

The format is the reason `anchors task-status` exists, and its order is the order in which each piece
of information decides something. First where the task is: the card and what its state means, before
any detail. Then what was undone: a move the state lock reverted changes what the agent thinks it did,
so it comes before the verdict. Then the verdict: the pull request and its checks — the piece missing
from the report that motivated the command, since an open pull request whose checks nobody read looks
like delivered work — and the working tree. Then what is missing, with the decisions waiting on a
person set apart, because it is the only kind of pending work that does not move by continuing to work.
Then the two lines the machine cannot know, left as explicit gaps: a named gap gets filled, an absent
section goes unnoticed. Last, the next step, derived from the state and not from the agent's intention.

A check verdict is never summarised in one word when the checks disagree ("passed" with three of four
is a report that lies by omission), and the worst class comes first so whoever reads fast hits the
problem, not what went right.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the task state | a card or none, reversions, a pull request or none, the branch and tree, the decisions waiting on a person | — | the discovery (TaskStatus) |

## Effects

### Sections

| Effect | Description |
| --- | --- |
| `TSRTS-B01` | The sections come in this order: the card, what was undone, the pull request, the working tree, the decisions waiting on a person, the two gaps, and the next step. |
| `TSRTS-B02` | The card line shows its number and its title without the leading `[CODE]`, then what its state means ("under review", "in progress", "closed", …) and its owner when known; without a card it says none was found and to provide `--card N`. |
| `TSRTS-B03` | The undone section appears only when there are reversions: it lists each and says how to authorise a deliberate move with the `anchors:manual` label. |
| `TSRTS-B04` | The pull request line says "no check ran" when there are no checks; one class alone reads as count over total; mixed classes are listed one by one, the failed first, then the running, then the passed. Without a pull request it says there is none for the branch. |
| `TSRTS-B05` | The working-tree line shows the branch and says "uncommitted change", the number of unpushed commits, or "up to date with the remote". |
| `TSRTS-B06` | The waiting section appears only when decisions wait on a person: it lists each number and title without its code, and says no agent resolves these. |
| `TSRTS-B07` | "What I proved" and "What was left out" always appear, as gaps for the agent to fill. |

### The next step

| Effect | Description |
| --- | --- |
| `TSRTS-B08` | An uncommitted change comes first, then an unpushed commit, before any other step. |
| `TSRTS-B09` | Without a pull request: a card in progress with a clean, pushed tree is told to open the pull request; a card under review or ready for review is told its pull request is on another branch and how to find it; no other card is told to open one. |
| `TSRTS-B10` | With an open pull request: no check → trigger the checks and do not trust a green that does not exist; any running → wait with `--watch` and do not end the turn; any failed → it is work of THIS card, fix it; all passed → the review is missing. |
| `TSRTS-B11` | A closed card is told that `anchors next` asks for the next one; when no step applies, the fallback points to `anchors status` and `anchors next`. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TSRTS-I01` | A running check is never reported as passed: a mix of passed and running checks never reads as all of them passing. | three passed and one running are rendered, and "4/4" does not appear |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TSRTS-X01` | Rendering only formats the state it receives; it looks nothing up. | The discovery and the format fail differently; a format that looked things up would abort the report it exists to guarantee. |

## Errors

none — the renderer receives a state and formats it; an absent part of the state is itself a line of the report (no card, no pull request), not a failure.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/flow/task_status.go` | `taskState` | comando — the state it formats |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
