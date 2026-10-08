<!-- @anchors
  code: MPRVM
  updated_at: 2026-10-08
  layer: mapa
-->
# MapReview — the reviews recorded on a node, and which of them still holds

> **Code**: `MPRVM`

## Overview

A review is a second look at what an agent decided, recorded on the node it looked at: for which gate, at which revision, by whom, when, and whether it found something. It holds while the node keeps that revision; a change to the node makes it history, and the node is to review again. The record keeps who reviewed as they said it, and does not grade it. Reviews survive a map rebuild whatever the revision, because they are the history of who looked and when.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a node id and a gate name | any | — | the caller passes the target and the gate |
| a review | a gate, a reviewer, a date and whether it found something | — | the review command builds it |

## Effects

| Effect | Description |
| --- | --- |
| `MPRVM-B01` | A node's review for a gate holds (`ReviewOf`) only when it was recorded at the node's current revision; a review of another revision, of another gate, or of a node not in the map does not. |
| `MPRVM-B02` | Recording a review (`RecordReview`) stamps it with the node's current revision and replaces the node's earlier review for the same gate, keeping the other gates' reviews; recording on a node not in the map records nothing and says so. |
| `MPRVM-B03` | A rebuild keeps each node's reviews, whatever its revision. (`PreserveStamps`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MPRVM-I01` | A node holds at most one review per gate. | records two reviews of one gate and one of another, and counts one per gate |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPRVM-X01` | Does not grade who reviewed. | Whether a reviewer is independent of the author is the project's to decide, not a rule the record imposes. |

## Errors

none — the unit reads and writes the map in memory; a node not in the map is answered by MPRVM-B02, not as a failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
