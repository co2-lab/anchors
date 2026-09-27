<!-- @anchors
  code: BREXB
  updated_at: 2026-09-26
  layer: infra
-->
# BoardExposure — hand the local board the same page and the same collect contract the pipeline publishes

> **Code**: `BREXB`

## Overview

The board exists in two places: the page the board pipeline publishes, and the local board that
`board serve` shows. Two renderings of the board would drift apart over time, and nobody would know
which one is right. So the local board does not keep its own copy of anything: this unit hands it the
very page the pipeline publishes, and the very collect expression the pipeline uses to build each card
of `board.json`, read out of the pipeline carried in the binary.

The collect contract is single. Reimplementing it would create a second version that diverges at the
first new field — the owner enters the pipeline and the local board does not show it, or worse, shows it
differently. Reading it from the pipeline instead means the extraction must find the expression exactly:
the start is the collect command, and the end is where the array closes and the shell quote closes, just
before the output redirection. An earlier version cut at the first quote and swallowed the redirection,
so the expression looked complete and died only when the query ran. The collect step exists in two shapes
(a direct query option, and a slurp over the paged raw file), and both are accepted.

When the pipeline changes shape so that no expression can be found, the unit says so explicitly instead
of falling back to a built-in expression: a fallback would be exactly the second contract this unit exists
to prevent.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the board page carried in the binary | the published board page | a binary built without it | this unit: reports that the page is not carried |
| the board pipeline carried in the binary | a pipeline whose collect step writes an array expression, in either of the two shapes, before an output redirection | a pipeline of any other shape, or none | this unit: reports that the pipeline is not carried, or that the expression was not found |

## Effects

| Effect | Description |
| --- | --- |
| `BREXB-B01` | The board page (`BoardHTML`) handed to the local board is byte for byte the page the pipeline publishes. |
| `BREXB-B02` | The collect expression (`BoardCollectJQ`) handed to the local board is the array expression of the pipeline's collect step, with its card fields and the filter that drops discarded cards, and without surrounding whitespace. |
| `BREXB-B03` | The collect expression is found both in the direct query shape and in the slurp over the paged raw file, and it ends where the array and the shell quote close, before the output redirection. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `BREXB-I01` | The extracted expression never carries the shell's output redirection: it opens and closes an array and names no output file. | extracts the expression from the carried pipeline and checks its ends and the absence of the redirection target |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `BREXB-X01` | Does not keep a built-in collect expression to fall back on. | A second copy of the contract is the drift this unit exists to prevent. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BREXB-E01` | The binary does not carry the board page or the board pipeline. | An empty result and an error saying which one is not carried. | Serving an empty board would look like a board with no cards. |
| `BREXB-E02` | The carried pipeline has no collect expression of the expected shape. | An empty result and an error saying the pipeline changed shape. | The local board would otherwise read a contract different from the one that is published. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/initx/workflows.go` | `workflowsFS`, `boardFS` | infra — the pipeline and the board page carried in the binary |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
