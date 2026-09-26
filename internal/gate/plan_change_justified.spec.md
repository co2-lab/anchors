<!-- @anchors
  code: PCJPL
  updated_at: 2026-09-26
  layer: gate
-->
# PlanChangeJustified — a modified plan or spec must declare why it changed

> **Code**: `PCJPL`

## Overview

Confronts a modified plan or spec against the declaration of its change: **planning makes mistakes,
and whoever implements discovers them — but fixing the plan silently makes the project walk towards
a destination nobody chose.**

Drift is the greatest risk, and it is SILENT by construction: no STATE gate can catch it, because
the corrected plan or spec is perfectly valid — the inconsistency was removed. What exposes the defect
is not the state of the file, but the UNJUSTIFIED CHANGE.

For this reason, this gate inspects the diff rather than just the content: any plan or spec that
appears among the changed files must carry a declared revision (`-R0001`). Without it, the gate blocks.
The ruler relies on the explicit judgment of whoever made the change:
- An INNOCUOUS correction (wording, example, typo, an ambiguity with only one viable interpretation) —
  correct the text and register the sequential revision. The gate verifies that the revision exists.
- A correction that CHANGES DIRECTION, or any doubt about whether it does — do not fix it in-place.
  Escalate via `anchors escalate` (`anchors:needs-user`), halting delivery of the card until a decision is made.

The gate separates silent drift from transparent adaptation: it does not judge whether a correction is
innocuous or directional, but ensures that a human or agent made the conscious choice and recorded it in
the document itself instead of letting drift happen by omission.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any map node | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the execution configuration | a `Config` carrying the list of changed paths, or nil | — (nil config skips confrontation) | this unit: without the changed list the gate goes quiet, never accuses |
| the file status | a file confirmed as modified by git relative to the last commit | untracked or newly created files, files untouched by git, or files only in the impact radius | this unit, by querying git status and git cat-file |
| the revision format | revision lines matching `[CODE]-R[0001]` with sequential numbering | unnumbered or freeform markers | this unit, parsing via `RevisionsOf` |

## Effects

| Effect | Description |
| --- | --- |
| `PCJPL-B01` | An artifact not listed among the changed files is skipped, even if reached by the impact radius. |
| `PCJPL-B02` | An artifact that has no identity code is skipped without failing: identity enforcement belongs to another gate. |
| `PCJPL-B03` | A changed plan or spec with no declared revision fails, and the verdict explains how to record the revision or escalate. |
| `PCJPL-B04` | A changed plan or spec with a valid sequential revision for its own identity code passes. |
| `PCJPL-B05` | Citing the revision of another document does not count as justifying this document's change. |
| `PCJPL-B06` | Revisions whose numbering has gaps or is not strictly sequential from 1 fail. |
| `PCJPL-B07` | Multiple sequential revisions pass, and the verdict cites the most recent revision explaining the change. |
| `PCJPL-B08` | Revision markers formatted in bold, blockquotes, raw lines, or GitHub alerts pass. |
| `PCJPL-B09` | A change already explained by the plan revision mechanism (`revises:`, `@revised-by`, `@amended-by`) passes without requiring duplicate notation. |
| `PCJPL-B10` | A newly created untracked file has no previous state to justify and is skipped. |
| `PCJPL-B11` | A newly created staged file has no commit history and is skipped. |
| `PCJPL-B12` | An existing committed file that is modified without a revision fails. |
| `PCJPL-B13` | When run outside a git repository, the gate trusts the caller's changed list rather than silencing itself. |
| `PCJPL-B14` | The exported function `RevisionsOf` extracts all declared revisions in order with code, number, and explanation. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PCJPL-I01` | Only files genuinely modified are charged. A file present in the impact radius but not in the changed list is never accused. | confronts an untouched node inside an impact radius and verifies it returns Skip |
| `PCJPL-I02` | Numbering must be strictly sequential starting at 1. Holes or repetitions prevent determining how many times the document changed. | confronts non-sequential revisions and verifies it returns Fail |
| `PCJPL-I03` | Newly created files are never charged. A file that did not exist in the previous commit has no past state to justify. | confronts untracked and staged new files and verifies both return Skip |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PCJPL-X01` | Does not distinguish between innocuous corrections and directional changes. | That distinction requires semantic understanding and contextual judgment belonging to a human or agent; the gate only verifies that a conscious decision was recorded. |
| `PCJPL-X02` | Does not enforce identity codes on files without one. | Code presence is enforced by `spec-has-code`; duplicating that check here would produce two findings for a single defect. |
| `PCJPL-X03` | Does not run or enforce revisions during full-project checks (`--all`). | In a full run there is no diff context, and charging files that never needed revision would penalize units designed correctly from the start. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern`, `Config` | core — code length patterns and configuration for changed files |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and diagnostic messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `Node` | core — graph structure and artifact node representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
