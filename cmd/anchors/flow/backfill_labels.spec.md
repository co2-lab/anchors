<!-- @anchors
  code: BCLBB
  updated_at: 2026-10-08
  layer: comando
-->
# BackfillLabels — write into open cards the blocking and provenance links the board already implies

> **Code**: `BCLBB`

## Overview

A new label does not reach the past. `escalate --for-user` began to mark the origin card as blocked by
the decision it opened; the cards that stopped before only carry the needs-user mark, so the board says
they wait without saying for whom, and the claim cannot check whether the decision came out. Measured:
29 open decisions, and unblocking the dependent cards meant rereading them one by one.

The link is recoverable because the escalation always wrote the origin card's number in the new
card's under label: the relation exists, in the opposite direction. `anchors backfill-labels` reads
every open decision in one listing, builds the reverse map, and writes the blocked-by label that was
missing on each origin card. For a decision born without an under label while reviewing someone
else's pull request, it recovers the provenance from the body when the body cites the pull request
and that pull request declares its card.

It never guesses. A decision that points at no origin card stays as it is. A closed origin card
receives nothing. An origin card that is itself an open decision receives nothing either: the under
label says where a card was born, not which decision must be answered first, and measured, that
inference was wrong in four of five cases — precedence between decisions is declared by hand. The
dry run shows what would be written without touching anything.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the open decisions | the open cards carrying the needs-user label, with their labels | — | the platform listing |
| the pull request mention | the words `PR #<n>` in a decision's body | a bare `#<n>` in prose | this unit: only the form with the word counts |
| the pull request's card | a line starting with Refs, Closes, Fixes or Resolves followed by `#<n>` in its body | a pull request that declares no card | this unit: such a decision is left untouched |
| the workflow mode | github mode | local mode, where the link is a folder | this unit: refuses local mode |

## Effects

| Effect | Description |
| --- | --- |
| `BCLBB-B01` | Outside github mode the command refuses, saying that locally the link lives in the `issues/` folder. |
| `BCLBB-B02` | Every open decision is read in a single listing, not one call per card. |
| `BCLBB-B03` | A decision carrying `under-<n>` holds origin card `n`; every decision pointing at the same origin holds it, and a decision whose under label names itself holds nothing. |
| `BCLBB-B04` | A decision with no under label whose body cites `PR #<n>`, where that pull request declares a card, receives `from-pr-<n>` and `under-<card>`; a bare `#<n>` in the body, or a cited pull request that declares no card, recovers nothing. |
| `BCLBB-B05` | An open origin card receives `blocked-by-<decision>` for each decision holding it; a closed origin card receives nothing and is counted as skipped. |
| `BCLBB-B06` | An origin card that is itself an open decision receives nothing: the command says precedence between decisions is not inferred and to put the label by hand, and counts it as skipped. |
| `BCLBB-B07` | A blocked-by label the origin card already carries is not written again. |
| `BCLBB-B08` | Each label is created on demand before the card is edited with it. |
| `BCLBB-B09` | The dry run prints what each card would receive, says nothing was touched, and makes no write. |
| `BCLBB-B10` | The output ends with the count of provenances recovered (when any) and of links written and skipped; with no origin card at all it says there is no link to recover and how many decisions are open. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `BCLBB-I01` | No card is ever made blocked by itself. | a decision whose under label names itself is listed, and no call edits it |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `BCLBB-X01` | The command never removes a label and never writes a link that no label or cited pull request supports. | An invented blocker would make the claim hold a card for a decision nobody linked to it, with nothing to release it. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BCLBB-E01` | The listing of open decisions fails. | The command fails with "list the open decisions". | Without the list there is nothing to recover. |
| `BCLBB-E02` | The listing answers with something that is not a list of cards. | The command fails with "read the list". | An unread list cannot be taken for an empty board. |
| `BCLBB-E03` | The edit that writes a label fails. | The card and the platform's answer are reported on standard error, and the link counts as skipped, not written. | One refused card must not stop the others, and the count must not claim a write that did not happen. |
| `BCLBB-E04` | The state of an origin card cannot be read. | The origin card receives nothing and counts as skipped. | Writing a block on a card whose state is unknown could label a closed card. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
