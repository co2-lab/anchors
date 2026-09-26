<!-- @anchors
  code: RVRNR
  updated_at: 2026-09-26
  layer: gate
-->
# RevisionRenumber — the revisions a branch added move to a free number when the base took theirs

> **Code**: `RVRNR`

## Overview

A revision is numbered by the file it revises: `PRICX-R0003` is the third recorded change of
`PRICX`. Two open pull requests that revise the same spec each take the next free number — and
both take the SAME one. The second to merge lands a `R0003` that already means something else
on the base: two revisions answering to one code, every citation of it ambiguous, and
`plan-change-justified` seeing the count and the highest number disagree.

The number stays, because it is what says "changed three times". The collision is resolved
where it happens: on the branch that arrives second, at rebase. This unit is the engine behind
`anchors renumber` — text in, text out; the command reads the three versions from git and writes
the result.

**Only what the branch added moves.** A revision the base already has is history other people
have read and cited; renumbering it would point every one of those citations at the wrong
change. The branch's revision is the one nobody outside the branch has seen yet. The same
reasoning decides which citations move: only those on a line the branch added.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the revision-bearing file | a spec or plan, in three versions: the branch's, the merge base's, the base's | a file the branch did not change | the command: it passes the Markdown files the branch changed, or the ones named |
| a revision | a `{CODE}-R000N` line in the format `plan-change-justified` reads | a revision cited in prose | `revisionRE`, shared with `plan-change-justified` |
| a citation | the full code, `{CODE}-R000N`, anywhere in a line | the short form (`R0003`) | this unit: only the full code is unambiguous across files |
| the lines the branch added | lines of the branch's version absent from the merge base's | — | this unit: compares line sets |

## Effects

| Effect | Description |
| --- | --- |
| `RVRNR-B01` | `PlanRenumber`: a revision the branch added whose number the base already uses moves to the next free number. |
| `RVRNR-B02` | After a rebase, with the base's revision and the branch's sharing a number in the same file, only the branch's moves. |
| `RVRNR-B03` | When one added revision collides, every revision the branch added for that code moves, in the order of its old number, after the highest number in use — so the sequence stays in order. |
| `RVRNR-B04` | A revision the branch added whose number the base does not use stays as it is. |
| `RVRNR-B05` | `RewriteRevisionCitations`: a citation is rewritten only on a line the branch added; a line that existed at the merge base keeps its code. |
| `RVRNR-B06` | A branch that adds the same number twice is refused, naming the code: a citation of it could not be told apart. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RVRNR-I01` | A revision the base has is never renumbered — including one whose explanation the branch edited. | edits the text of a base revision on the branch and verifies nothing is planned |
| `RVRNR-I02` | The rewrite is one pass: a chain of renames (`R0003→R0004`, `R0004→R0005`) never cascades. | rewrites a line citing both codes and verifies each moved once |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RVRNR-X01` | Does not rewrite a revision cited without its unit code. | `R0003` alone names a revision only inside its own spec; rewriting it in another file would guess which unit it meant. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/plan_change_justified.go` | `revisionRE` | gate — a revision is recognised exactly as the gate that counts them recognises it |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
