<!-- @anchors
  code: PRFLP
  updated_at: 2026-09-26
  layer: scan
-->
# ProgressFile — a plan's progress companion: how it is named, and how two sides of it merge

> **Code**: `PRFLP`

## Overview

A plan is a DECISION; its progress is STATE. They live in two files so that marking a phase done
never counts as changing the plan — the difference the plan-change gates exist to preserve. The
progress companion of a plan carries a fixed suffix, and this unit is the one definition of it: the
scan keeps such files out of the map, and every consumer that looks for a plan's companion derives
its path here, so no second copy of the suffix can drift.

The unit also merges two sides of a progress file, which is what the git merge driver runs. Two
branches that tick neighbouring checkboxes of the same plan conflicted every time, and the manual
resolution was always the same. The single rule is that a ticked box wins over an unticked one:
marking done records a fact (the spec was delivered, the check stamped it, the PR merged), and a
merge that unticked it would erase the fact. The merge starts from our side, so every line that is
not a checkbox item — headings, prose, sections — stays exactly as it was, and adds what only the
other side has.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a path to classify | any slash-separated path | — | this unit: any path is answered |
| a plan path | a slash-separated path, with or without an extension | — | the caller passes the plan's path |
| the two sides of a merge | any text; lines of the form `- [ ] text`, `- [x] text` or `- [X] text`, optionally indented, are items | — | this unit: a line that is not an item is kept as text |

## Effects

| Effect | Description |
| --- | --- |
| `PRFLP-B01` | `IsProgressFile`: A path is a progress file only when it ends with the suffix `-progress.md`; `progress` elsewhere in the name, or a file named just `progress.md`, is not. |
| `PRFLP-B02` | `ProgressPathFor`: The companion of a plan is its path with the extension of the last segment replaced by `-progress.md`; a dot in a directory name is not an extension, and a path with no extension gets the suffix appended. |
| `PRFLP-B03` | `MergeProgress`: In a merge, an item ticked on either side comes out ticked. |
| `PRFLP-B04` | An item present only on their side is added after the last item of our side, or at the end when our side has no item. |
| `PRFLP-B05` | Every line of our side that is not an item comes out unchanged and in its place. |
| `PRFLP-B06` | A ticked mark promoted from their side keeps our side's indentation. |
| `PRFLP-B07` | `ProgressDone`: The done count is the number of items ticked with `x` or `X`, at any indentation. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRFLP-I01` | Merging a side with itself changes nothing, and the done count is the same before and after. | merges a file with itself and compares the text and the count |
| `PRFLP-I02` | A merge never unticks: an item ticked on our side is ticked in the result whatever their side says. | merges a ticked item against the same item unticked |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRFLP-X01` | An item's identity is its text after the checkbox; an item whose text differs between the sides is two items, and both are kept. | The text comes from the plan and nobody rewrites it when ticking; guessing that two different texts are the same item could merge distinct deliveries. |

## Errors

none — the functions work on text alone: every path gets an answer, and every line of a side is either an item or text kept as is.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
