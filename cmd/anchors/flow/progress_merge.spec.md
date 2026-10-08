<!-- @anchors
  code: PRMRP
  updated_at: 2026-10-08
  layer: comando
-->
# ProgressMerge — the git merge driver of a plan's progress file: unite both sides, done beats pending

> **Code**: `PRMRP`

## Overview

Two branches delivering specs of the same plan mark neighbouring checkboxes of that plan's progress
file, and a textual merge asks for manual resolution every time. Measured: three consecutive pull
requests conflicted on the same progress file, with the identical mechanical resolution in all three,
while the map file — which has its own merge driver — never conflicted in the same period.

`anchors merge-progress` is that driver for progress files. Git calls it with the base, our side and
their side; it unites the two sides and writes the result into our side, which is where git expects
it. The union rule belongs to the progress reader (a done item beats a pending one, because marking
done records a fact that a merge must not erase); this unit applies it and reports what it did.

The base is not read: there is no legitimate "unmark" coming from a merge. Whoever wants to unmark
edits the file, and that edit arrives as one side.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the arguments | exactly three paths: base, ours, theirs (git's `%O %A %B`) | any other count | this unit: refuses |
| our side and their side | readable progress files | a path that cannot be read | this unit: fails naming the side |
| the base | any path, even a missing one | — | this unit: it is never read |

## Effects

| Effect | Description |
| --- | --- |
| `PRMRP-B01` | The driver takes exactly three arguments; any other count is refused. |
| `PRMRP-B02` | The united result is written into our side, where an item done on either side is done. |
| `PRMRP-B03` | The base is not read: a missing base does not stop the merge. |
| `PRMRP-B04` | Standard error reports how many items are done after the merge, and how many came from the other side when it added any. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRMRP-I01` | An item done on our side stays done after the merge, whatever the other side says. | our side marks an item done that their side leaves pending, and the result keeps it done |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRMRP-X01` | The rule that unites two progress files is not decided here; the driver applies the progress reader's union. | One rule, one owner: the gate that confronts progress and the driver that merges it must agree on what "done" is. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRMRP-E01` | Our side cannot be read. | The driver fails with "read our side" and the path. | Git needs to know which side broke to fall back to a manual merge. |
| `PRMRP-E02` | Their side cannot be read. | The driver fails with "read the other side" and the path, and our side is left untouched. | A half-merged file written on a failed read would lose our side. |
| `PRMRP-E03` | The result cannot be written into our side. | The driver fails with "write the result". | Git must not take an unwritten merge for a resolved one. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
