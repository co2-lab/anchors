<!-- @anchors
  code: PRFLO
  updated_at: 2026-09-26
  layer: gate
-->
# Profile — the verdicts of a run, gathered per gate and per node

> **Code**: `PRFLO`

## Overview

A check run produces one result per gate and per confronted artifact. Quality in Anchors is not a
single number: it is that set of verdicts. The profile gathers the raw results into the three
answers the rest of the tool needs.

Per gate, it counts how many confrontations passed, failed, skipped, stayed pending or await an AI
judgement, and how long they took: the total time and the single most expensive run. The two times
together tell a gate that is expensive by volume (cheap per target, run hundreds of times) from a gate
that is expensive per target, because the fix for each is the opposite.

For the run as a whole, it decides promotion. A failure of a blocking gate blocks it; a failure of a
non-blocking gate is still a failure (it becomes an issue) but does not block. A pending verdict does
not block by itself: in a real repository, letting every pending verdict of a blocking gate block
rejected 411 nodes at once, most of them meaning "there was nothing to confront". Only the gate knows
whether its pending means "a decision is still open", and it says so by marking the result as one that
impedes; a pending that impedes, on a blocking gate, blocks promotion.

Per node, it collapses the results into one verdict the map uses to stamp its edges: which nodes were
actually confronted, and which of them failed a blocking gate.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the raw results | any list of results, each with a gate name, a target, a verdict, the gate's blocking flag, whether a pending impedes, and a duration | — (every combination is read; an empty list is an empty profile that passes) | the runner that produced the results |

## Effects

| Effect | Description |
| --- | --- |
| `PRFLO-B01` | The aggregation (`Aggregate`) gives each gate a summary that counts its passes, failures, skips, pending verdicts and verdicts awaiting judgement, and the gate names (`GateNames`) come back sorted. |
| `PRFLO-B02` | The summary of a gate carries the total time of its confrontations and the single most expensive one. |
| `PRFLO-B03` | Every failed result is listed as a failure, whether its gate blocks or not. |
| `PRFLO-B04` | A failure of a blocking gate blocks promotion; a failure of a non-blocking gate does not. |
| `PRFLO-B05` | A pending verdict blocks promotion only when its gate is blocking and the gate marked the pending as one that impedes. |
| `PRFLO-B06` | Every result awaiting an AI judgement is listed as awaiting judgement. |
| `PRFLO-B07` | The per-node verdicts (`NodeVerdicts`) include only nodes that were confronted: a node touched only by skips or by verdicts awaiting judgement is left out, and the list is sorted by node. |
| `PRFLO-B08` | A node is marked failed when a blocking gate failed on it; a failure of a non-blocking gate leaves it confronted but not failed. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRFLO-I01` | Promotion is refused exactly when the list of blocking results is not empty. | each case of the promotion table checks the pass flag and the blocked list together |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRFLO-X01` | Does not decide what a verdict means: whether a pending impedes, and whether a gate blocks, arrive already set on the result. | Only the gate knows what it measured; the profile would otherwise have to guess, and guessing made every pending block. |

## Errors

none — the profile is a pure aggregation over results already produced: it reads no file and calls nothing that can fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/gate.go` | `Result`, `Verdict` | gate — the raw verdict of one confrontation |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
