<!-- @anchors
  code: DCDDE
  updated_at: 2026-10-08
  layer: comando
-->
# Decided — release the card an escalation stopped for a person, once the decision became a rule

> **Code**: `DCDDE`

## Overview

`anchors escalate --for-user` stops a card: it marks it as needing a person, and the claim stops
handing it out, so the next agent does not walk back to the same doubt. What was missing was the way
back. Measured: the decision came out, the revisions were applied to four files, the decision issue
was closed — and the card stayed stopped until someone removed the mark by hand. The instruction to
remove it lived only in the body of the card the workflow opens, and the failure was the silent one:
the card is skipped by `next` and nobody notices until someone asks why the plan does not move.

`anchors decided --card <n> --resolution <revision>` closes that cycle. The resolution is mandatory
because the exit of an open decision is one: the answer becomes a rule with a code, and releasing the
card without saying which revision was born from it would leave the decision without a trace.

If the decision generated work (an unblock card is open for this card), the card waits for that
delivery and the command refuses; not being able to tell refuses too. Otherwise the command removes,
in one edit, the needs-user mark together with every blocked-by label — reading the blockers first,
because the label is the only place they are recorded — then comments the resolution with the list of
what blocked the card, and closes the open decision issues born under it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the card | a card number | no card | this unit: refuses with an example |
| the resolution | non-blank text naming the revision born from the decision | empty or whitespace-only | this unit: refuses and says why it is required |
| the workflow mode | github mode | local mode, where the decision is a file under `issues/` | this unit: refuses and says where the decision lives |

## Effects

| Effect | Description |
| --- | --- |
| `DCDDE-B01` | Without `--card` the command refuses with an example. |
| `DCDDE-B02` | Without a non-blank `--resolution` the command refuses, and the message says why: the answer becomes a RULE. |
| `DCDDE-B03` | Outside github mode the command refuses, pointing to `issues/`, where the decision lives locally. |
| `DCDDE-B04` | While a card carrying this card's unblock label is open, the command refuses, naming the pending cards and saying the needs-user label stays; nothing is edited. |
| `DCDDE-B05` | The release removes the needs-user label and every blocked-by label of the card in one edit, and announces that the claim delivers the card again. |
| `DCDDE-B06` | The card receives a comment with the resolution and, when it had blockers, the list of the cards that blocked it. |
| `DCDDE-B07` | Only the open issues carrying both the needs-user label and the card's under label are closed, each first labelled manual and then closed with the resolution. |
| `DCDDE-B08` | The output counts the decisions closed: a sentence for none, one for exactly one, and the number for more; a decision whose close fails is not counted. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCDDE-I01` | The blockers are read before the labels are removed, and the card is released before its decisions are commented or closed. | the recorded calls of a full run are in the order: read the labels, release, comment, close |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCDDE-X01` | A decision issue that is not under this card is never closed. | A decision closed by mistake is taken by omission. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCDDE-E01` | The lookup of open unblock cards fails or answers with something that is not a list. | The command refuses, naming the label it could not read and the label that stays; nothing is edited. | Not knowing is not "nothing pending": releasing could return to the queue a card whose work is not done. |
| `DCDDE-E02` | The release edit fails. | The command fails with "release card #n" and the platform's answer, announces no release, and neither comments nor closes. | A decision closed on a card still stopped is the inconsistency this command exists to undo. |
| `DCDDE-E03` | The list of decisions under the card cannot be read. | Nothing is closed, and the count is zero. | Closing without a list that was actually read could close the wrong issues. |
| `DCDDE-E04` | The card's labels cannot be read. | No blocker is reported, and the release still removes the needs-user label. | The blockers are the trace, not the release; failing to read them must not keep the card stopped. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
