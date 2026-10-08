<!-- @anchors
  code: RVCMR
  updated_at: 2026-10-08
  layer: comando
-->
# ReviewCommand — record a review with who looked and what they found, or list what is to review

> **Code**: `RVCMR`

## Overview

`anchors review` is where a reviewer — a person, or an agent — records a second look at a target of a gate that declares `review:`. Its product is findings: a review with findings opens an issue with the whole report, the way a failing judgment does, and each finding becomes a failing test and a fix, or is dismissed with its reason in the issue. A review without findings records who looked. There is no `waived`: a review that did not happen is still due. `--pending` lists what is to review, with the question.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the target | a node of the map | a target not in the map | this unit: refused (RVCMR-E01) |
| the gate | a gate that declares `review:` | another gate | this unit: refused naming the gate (RVCMR-E01) |
| who reviewed | any name the reviewer gives | none | this unit: refused saying what to write (RVCMR-E01) |

## Effects

| Effect | Description |
| --- | --- |
| `RVCMR-B01` | `--pending` lists what is to review, grouped by gate with its question, and how to record each; with nothing due it says so. |
| `RVCMR-B02` | A review with no findings is recorded on the target at its current revision, with who reviewed, and opens no issue. |
| `RVCMR-B03` | A review with findings is recorded the same way, marked as having found something, and opens the target's issue for the gate with the report — reopening it when it was closed. |
| `RVCMR-B04` | In manual mode a review with findings writes no issue unless `--record-issues` asks, and prints the report. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RVCMR-I01` | A review record always names who reviewed. | records with and without `--by` and checks the refusal and the record |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RVCMR-X01` | Offers no `waived`. | A review that did not happen is a review still due; recording it as done would stamp a look nobody took. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RVCMR-E01` | The target is not in the map, the gate declares no `review:`, or no reviewer is given. | The record is refused saying what to write, and nothing is written. | A review recorded on a wrong target, gate or nobody would answer for a look that did not happen there. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
