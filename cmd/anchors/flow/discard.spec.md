<!-- @anchors
  code: DSCRD
  updated_at: 2026-09-26
  layer: comando
-->
# Discard — take off the board a card that no longer makes sense, without deleting it

> **Code**: `DSCRD`

## Overview

`anchors discard` marks one or more cards as discarded. Closing a card is not enough to take it off
the board: the board keeps closed cards because the roadmap draws the delivered past from them, so a
test card or a finding about a file that no longer exists stays there forever, closed and without a
possible parent. In the reference project, twelve of the twenty-seven roadmap roots were that kind of
permanent noise.

Discarding is a soft delete. The card receives the discard label, a comment with the reason and the
way back, and is closed. It stays on the platform, so whoever looks the number up still finds that the
question existed and why it was dropped — the difference between "resolved" and "discarded".

The command exists only where a card is an issue (github mode). It works card by card: a card that
fails is named in the error while the others are still discarded.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the cards | one or more issue numbers, with or without a leading `#` | no argument at all; blank arguments | this unit: refuses an empty list, skips a blank argument |
| the reason | any non-blank text | an empty or whitespace-only reason | this unit: refuses before touching the platform |
| the workflow mode | github mode, with the repository configured | local mode, where a card is a file | this unit: refuses and says to delete the card file |

## Effects

| Effect | Description |
| --- | --- |
| `DSCRD-B01` | Without any card argument the command refuses to run. |
| `DSCRD-B02` | A blank reason is refused before anything is written, because a card that leaves the board without saying why is indistinguishable from a lost one. |
| `DSCRD-B03` | Outside github mode the command refuses and says that, locally, the card file is deleted instead. |
| `DSCRD-B04` | The discard label is created on demand before any card is touched, and a failure to create it (it already exists) does not stop the command. |
| `DSCRD-B05` | Each card is handled in a fixed order: the discard label, then the reason comment, then the close. |
| `DSCRD-B06` | The comment left on the card carries the reason and names the label to remove to bring the card back. |
| `DSCRD-B07` | A leading `#` is stripped from each card, and a blank argument is skipped. |
| `DSCRD-B08` | A card whose close fails (a card already closed) still counts as discarded. |
| `DSCRD-B09` | Each card discarded is reported on standard output as off the board and still on the platform. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DSCRD-I01` | A card is closed only after its label and its reason were both recorded. | a failed label and a failed comment are scripted, and neither card is closed |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DSCRD-X01` | The command never deletes the issue. | Deleting would lose the trace that the question existed, which is what distinguishes discarded from resolved. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DSCRD-E01` | The discard label cannot be applied to a card. | That card is neither commented nor closed, the other cards are still discarded, and the command fails naming the card and the platform's answer. | A reason recorded on a card that stays on the board would say two contrary things at once. |
| `DSCRD-E02` | The reason comment cannot be recorded on a card. | That card is not closed, and the failure is reported as "record the reason" with the platform's answer. | A closed card with no reason is indistinguishable from a lost one. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Load`, `AbsRoot`, `GitHubMode` | config — project configuration and workflow mode |
| DEP2 | `internal/initx/workflows.go` | `LabelDiscarded` | the name of the discard label |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
