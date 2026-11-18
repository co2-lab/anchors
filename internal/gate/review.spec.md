<!-- @anchors
  code: RVDUR
  updated_at: 2026-10-08
  layer: gate
-->
# ReviewsDue — the targets of a reviewed gate that no review covers at their current revision

> **Code**: `RVDUR`

## Overview

A gate that declares `review:` marks its targets to review, apart from how it measures: a judgment gate still asks its question and stamps its verdict, and the same targets are to review until a reviewer looks. A target is to review when the gate applies to it and no review is recorded at its current revision. The list informs and never blocks; the project decides when the reviewer comes.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration and the map | any, possibly absent | — | this unit: with either absent there is nothing to review |
| a gate name | a gate to list alone, or none for all | — | the caller |

## Effects

| Effect | Description |
| --- | --- |
| `RVDUR-B01` | The targets to review (`ReviewsDue`) are, for every gate that declares `review:`, each node the gate applies to with no review recorded at its current revision, with the gate's review question (`ReviewAsk`); ordered by gate and target. |
| `RVDUR-B02` | A gate with no `review:` has nothing to review, and naming a gate lists that gate's alone. |
| `RVDUR-B03` | A change to a reviewed target makes it to review again. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RVDUR-I01` | Listing what is to review changes no verdict: a judged and reviewed gate's judgment is asked with or without its review. | runs the judged gate over a target before and after a review and compares the verdict |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RVDUR-X01` | The list never blocks. | A review is a second look the project schedules; a gate that measures still blocks by its own severity. |

## Errors

none — the unit only reads the configuration and the map it is given; nothing it does can fail.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
