<!-- @anchors
  code: BRCRB
  updated_at: 2026-09-26
  layer: apoio
-->
# BoardCards — reading the repository's board, where the work queue lives in github mode

> **Code**: `BRCRB`

## Overview

In `github` mode the work queue is the repository's board: each unit's work is an issue (a card) with
the Anchors label and a state label (to-do, in-progress, ready-to-review, in-review). This unit reads
that board for the agent: which card is already its own, the card of a unit, a card named by number;
and it posts the delivery record as a comment on the card, since in github mode the reviewer reads the
issue, not a file.

The agent never claims a card directly. GitHub offers no compare-and-swap, so two agents writing
ownership comments at the same time could both believe they own a card; instead the agent asks the claim
pipeline, which is serialized, and the pipeline writes the ownership comment. This unit only reads
ownership: the owner of a card is the last `anchors-owner:` comment, since ownership changes hands and
every claim stays recorded.

A card waiting for a person's decision (`needs-user`) is declined by an agent that has not declared it
decides the product: handing it to any agent would make that agent decide alone, which is what
escalation exists to prevent. When the agent resumes its own work, it finishes what is under way
(in-progress or in-review) before taking anything new. Commenting on the wrong card is worse than not
commenting, so a unit's card is found by the code in its title only, and an ambiguous code is refused.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the repository | `owner/name` | anything else | this unit: reading the board is refused |
| the labels | the project's non-empty `workflow.labels` | an empty list | this unit: reading the board is refused |
| the agent | the agent's identity as the pipeline writes it in the ownership comment | — | the caller (`anchors next`) |
| the unit code | a code, any case | an empty or blank code | this unit: the lookup is refused |

## Effects

| Effect | Description |
| --- | --- |
| `BRCRB-B01` | A card is on the board only when it carries every configured label and, when a state is asked, that state's label. |
| `BRCRB-B02` | The owner of a card is its last comment starting with `anchors-owner:`; a card with no such comment has no owner. |
| `BRCRB-B03` | A card labelled `needs-user` (or its Portuguese variant) is declined unless the agent declared it decides the product; any other card is never declined. |
| `BRCRB-B04` | The agent's own card is the first card it owns that is in progress or in review and not declined, reported with its most advanced state label; a card it owns that is to-do or ready-to-review, or a card of another owner, is not its own work. |
| `BRCRB-B05` | A unit's card is the one open card whose TITLE holds the code in brackets, in any case; the code in the body, or a longer code in brackets, does not count. |
| `BRCRB-B06` | A card named by number is returned when it is open and carries every configured label. |
| `BRCRB-B07` | Board queries through `gh api` never carry `--repo`; every other `gh` call names the repository. |
| `BRCRB-B08` | Asking for work dispatches the claim pipeline with the agent as its input. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BRCRB-E01` | The configured labels are empty. | Reading the board is refused. | Without the Anchors label the claim would pull any issue of the repository, product issues included. |
| `BRCRB-E02` | The repository is not `owner/name`. | Reading the board is refused, naming it. | The board query needs both parts; a guess would read another repository. |
| `BRCRB-E03` | The unit code is empty or blank. | The lookup is refused. | An empty code would match `[]` and return any card. |
| `BRCRB-E04` | No open card has the code in its title. | Refused, naming the bracketed code. | A delivery with nowhere to go must not be silently dropped. |
| `BRCRB-E05` | Two or more open cards have the code in their title. | Refused, listing every card and asking for `--card <n>`. | Picking the first recorded a delivery on the wrong card; the choice goes back to whoever knows it. |
| `BRCRB-E06` | The card named by number is closed, or lacks a configured label. | Refused, saying which. | A delivery there would go to a card nobody reads, or onto someone else's issue. |
| `BRCRB-E07` | The comment body is empty or blank. | Refused before calling `gh`. | A blank comment records nothing and looks like a delivery. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/board/gh.go` | `runGH` | apoio — every call to `gh` (`GHRNG`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
