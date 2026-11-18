<!-- @anchors
  code: SYCMS
  updated_at: 2026-10-08
  layer: comando
-->
# SynthesizeCommand — two pull requests in content conflict become one card that asks for the best of each, and every end points at it

> **Code**: `SYCMS`

## Overview

A content conflict is not resolved by automation: two agents wrote different things about
the same rule, and picking a side is choosing without reading both. In the case that shaped
this command, both sides were right and about different things, and both belonged.

`synthesize` gives that wait an owner. It opens a card on the board that asks for what
neither pull request delivers alone: the best of each, in a result better than both. The
card cites both pull requests, their origin cards and titles, and the conflicting files; it
tells how to proceed and warns against silently discarding a side. When only one side is
known (the conflict is against work already merged), the card is still opened, says the
other side was not identified, and says where to look for it.

Traceability is the point. The new card is labelled for the board and under each origin
card; each pull request is commented with where the work continued and then closed; each
origin card gets a comment pointing at the new card. The card must exist before anything is
closed. Each link that fails afterwards is a warning, because the card already exists, and
the closing line reports only the pull requests that were actually closed. The card is
written in the project's language. A dry run shows the card and touches nothing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the workflow mode | github mode | any other mode | this unit: refuses before calling the host |
| the first pull request | a number, with or without a leading `#` | an empty value | this unit: refuses before calling the host |
| the second pull request | a number, or nothing when the other side is already merged | — | this unit: a one-sided card |
| the files | a comma-separated list, possibly empty | — | this unit: blank entries are dropped |

## Effects

| Effect | Description |
| --- | --- |
| `SYCMS-B01` | Outside github mode, or without `--pr-a`, the command is refused and nothing reaches the host. |
| `SYCMS-B02` | A leading `#` on a pull request number is dropped. |
| `SYCMS-B03` | The card's body cites both pull requests with their origin cards and titles, lists the conflicting files, asks for the best of each, and warns against discarding a side without saying why. |
| `SYCMS-B04` | With one side only, the body cites no empty pull request, says the other side was not identified, points at the file's history to find it, and gives the instructions in the singular. |
| `SYCMS-B05` | The card is labelled with the project's first workflow label, `anchors:to-do`, and an `under-<card>` label for each origin card. |
| `SYCMS-B06` | Each pull request is commented with where the work continued (the other pull request, or the integration branch when there is none) and closed; each origin card is commented with the new card's link. |
| `SYCMS-B07` | A comment or a close that fails is a warning on standard error naming it, and the command still succeeds. |
| `SYCMS-B08` | The closing line names only the pull requests actually closed, one or many, and lists a pull request whose close failed as still open. |
| `SYCMS-B09` | `--dry-run` prints the title and the body, and creates, comments on and closes nothing. |
| `SYCMS-B10` | The title, the body and the comments come in the project's language. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SYCMS-I01` | No pull request is closed unless the synthesis card exists. | makes the card creation fail and checks that no close was issued |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SYCMS-X01` | The card never picks a side and the command never merges or resolves the conflict. | Choosing by automation is choosing without reading; the synthesis belongs to whoever takes the card. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SYCMS-E01` | The card cannot be opened on the host. | The command fails carrying the host's message, and no pull request is commented or closed. | Closing the pull requests without the card would lose both works with no owner. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
