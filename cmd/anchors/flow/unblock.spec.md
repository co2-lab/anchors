<!-- @anchors
  code: NBLCK
  updated_at: 2026-09-26
  layer: comando
-->
# Unblock — open the work card a decision demanded, linked to the card stuck waiting for a person

> **Code**: `NBLCK`

## Overview

A card marked as needing a person is not handed out by the claim while the mark is there. The
instruction the pipeline writes on it is "decide, then remove the mark", which assumes deciding is
only writing. Many decisions are not: the person decides, and the decision demands that something
change first — a spec that gains a rule, a contract that changes, a defect fixed elsewhere.

Without a link, both paths go wrong. If the person opens the change card and leaves the mark, the
original card waits forever, because its condition ("decide") was already met and nobody knows work
is missing. If the person removes the mark, the card returns to the queue and the next agent hits
the same impasse and escalates again.

`anchors unblock <card> --reason <what to change>` closes that gap. It creates the change card,
already linked to the blocked one by a per-card link label, and comments on the blocked card that
it now waits for the delivery of the new one. The blocked card keeps its mark until then, so no
agent takes it early.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the blocked card | exactly one issue number, with or without a leading `#` | no card, or more than one | this unit: refuses any other count |
| the reason | non-blank text, possibly with several paragraphs | an empty or whitespace-only reason | this unit: refuses before touching the platform |
| the file the change touches | an optional path | — | the caller: it is written as given |
| the workflow mode | github mode, with the repository and at least one workflow label | local mode, where a block is a folder, not a label | this unit: refuses local mode |

## Effects

| Effect | Description |
| --- | --- |
| `NBLCK-B01` | The command takes exactly one card; none or two are refused. |
| `NBLCK-B02` | A blank reason is refused, because whoever takes the new card does not have the context of the decision. |
| `NBLCK-B03` | Outside github mode the command refuses, saying that locally the block is not a label. |
| `NBLCK-B04` | The link label of the blocked card is created on demand before the card is created, and a failure to create it (it already exists) does not stop the command. |
| `NBLCK-B05` | The new card is titled `[unblocks #<n>]` followed by the first line of the reason, and carries the first workflow label, the to-do label and the blocked card's link label. |
| `NBLCK-B06` | The new card's body says which card its delivery unblocks, the reason, where it came from, that the original discussion need not be reopened, and to remove the needs-user mark from the blocked card on delivery; a given file adds a "Where" line, and without it there is none. |
| `NBLCK-B07` | The blocked card receives a comment naming the new card and saying its needs-user mark stays until that delivery. |
| `NBLCK-B08` | The output gives the address of the new card and says the blocked card remains stopped and now says who it waits for. |
| `NBLCK-B09` | A leading `#` on the card is stripped before it is used. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `NBLCK-I01` | The blocked card keeps its needs-user mark: the command never removes a label from it. | a full run is recorded and no call removes a label |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `NBLCK-X01` | The body is handed to the platform through a temporary file, and that file does not survive the command. | A multi-paragraph argument crosses shells differently; a file does not, and a leftover file would pile up in the temporary directory. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `NBLCK-E01` | The platform refuses to create the new card. | The command fails with "create the card" and the platform's answer, and the blocked card is not commented. | A comment pointing at a card that does not exist would send the reader nowhere. |
| `NBLCK-E02` | The comment on the blocked card fails after the new card was created. | A warning names the blocked card, and the command still succeeds. | The card exists; failing would hide what was created and what was not. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Load`, `AbsRoot`, `GitHubMode` | config — project configuration and workflow mode |
| DEP2 | `internal/initx/workflows.go` | `LabelDesbloqueia` | the link label of a blocked card |
| DEP3 | `cmd/anchors/common` | `FirstLineOfReason` | comando — the title line of a reason |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
