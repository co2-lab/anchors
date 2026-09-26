<!-- @anchors
  code: PRHNP
  updated_at: 2026-09-26
  layer: gate
-->
# ProgressHonest — the progress file tells the truth about the disk

> **Code**: `PRHNP`

## Overview

The `-progress.md` is the only Anchors artifact that nothing confronted, and its exclusion
from the map is deliberate: it exists in order to CHANGE, and the gates that demand a
justification for change cannot reach it.

**But "outside the map" turned into "outside any verification", and the two are not the same
thing.** An item that cites a file PATH is trivially confrontable: the file exists, or it
does not.

Measured in the reference project: when the 17 progress files were created, the checkbox
state was carried over from the plans — and the plans were out of date. The progress of
`0002` claimed 6 open items with 7 of the 8 specs already on disk. FIVE items lied, and the
lie was transported faithfully.

**The damage runs in two directions, and the second is worse:** an open box with the file
already there produces rework — somebody redoes what is done. A ticked box with the file
absent declares a plan finished with work still to do, which is the same defect the explicit
`Closes` shut through another door.

**A gate that only looks INSIDE the file never sees what is missing from it.** That is the
third direction, and it is how this gate once passed a plan that had just gained a seeded
spec: the progress said the phase was over, and the next-card command moved on to another
plan with work declared undone.

The gate is anchored on the PLAN — which IS in the map — and confronts its companion. That
is the only way to reach a file that, by design, is not a node.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the companion progress file | a readable file beside the plan, or none | — (its absence is a case, not an error) | this unit: with no companion the verdict is a skip, because "does it exist" is a different question from "is it true" |
| a progress item | a checkbox line, with or without a quoted path | prose with no path, which is not confrontable and is not charged | this unit, by the shape of the item it can read |
| a plan seed | a checkbox line of the plan quoting a spec path | a spec named in the plan's prose, which promises nothing | this unit: the checkbox is the promise, and the prose of a revision speaks of what already exists |

> `Who guarantees` cannot be left empty nor say only "not mine": if nobody guarantees, the
> duty is orphaned — and that is exactly where the invalid input gets through.

## Effects

| Effect | Description |
| --- | --- |
| `PRHNP-B01` | An artifact that is not a plan leaves the confrontation without a verdict: the gate is anchored on the plan, which is what the map holds. |
| `PRHNP-B02` | A plan with NO companion progress file is skipped, not failed — "does the progress exist" is a different question from "is the progress true". |
| `PRHNP-B03` | The skip for a missing companion says HOW to create it, so the reader does not have to look the command up. |
| `PRHNP-B04` | A TICKED item whose file does not exist FAILS: it declares done what is not. |
| `PRHNP-B05` | An OPEN item whose file already exists fails too: it produces rework, because somebody redoes what is done. |
| `PRHNP-B06` | A spec the plan SEEDS and the progress does not list fails: a gate that only looks inside the file never sees what is missing from it. |
| `PRHNP-B07` | A checkbox item promising no file at all — the untouched template marker — fails: an eternal open box makes the plan look unfinished forever. |
| `PRHNP-B08` | The ticked-but-absent finding is reported FIRST, because it is the costliest damage: whoever reads the board decides on it. |
| `PRHNP-B09` | An item in prose, citing no path, is not charged — there is nothing to confront. |
| `PRHNP-B10` | A progress file that agrees with the disk on every item PASSES. |
| `PRHNP-B11` | The verdict NAMES each offending path, so the reader does not have to diff the file against the disk by hand. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRHNP-I01` | The companion's path has ONE definition, derived from the scanner. The scanner is what must keep the file out of the map; a second constant here would diverge from it in silence, and the gate would confront a file the scanner indexes, or hunt for one that does not exist. | derives the companion of several plan paths and verifies each against the scanner's answer |
| `PRHNP-I02` | A seed is matched by PATH, never by the item's text. The wording can be shortened, translated or moved between phases without ceasing to be the same item — what identifies it is the file. | lists a seed under a different description and verifies it is not accused |
| `PRHNP-I03` | A spec mentioned in the plan's PROSE is not a seed. The checkbox is the promise; a revision's prose speaks of what already exists, and charging it accused progress files of not listing specs they did list. | mentions a spec in prose and verifies it is not charged as missing |
| `PRHNP-I04` | A template file is never a seeded spec. The mould is the shape work is poured into, not work to be done. | seeds a template path and verifies it is not charged |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRHNP-X01` | Does not charge the EXISTENCE of the progress file. | The plan may predate the mechanism, and a dedicated command exists to create the companion. Merging "does it exist" with "is it true" would report two different debts as one finding. |
| `PRHNP-X02` | Does not judge the CONTENT of an item beyond the path it cites. | Whether the description is accurate, whether the phase makes sense, whether the order is right — all of it is judgement. The ruler here is the disk, which is deterministic and needs no opinion. |
| `PRHNP-X03` | Does not put the progress file into the map, nor demand a justification for changing it. | It exists in order to change. A gate that reached it as a node would charge every edit of the very artifact that records work, which is the opposite of what it is for. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/progress.go` | `ProgressPathFor` | core — the companion's suffix has one definition, and it belongs to the scanner |
| DEP3 | `internal/mapx/model.go` | `KindPlan` | core — the gate is anchored on the plan, and the kind is what routes the jurisdiction |
| DEP4 | `internal/config/config.go` | `Config` | core — the plan layer and its conventions are declared in the Structure |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
