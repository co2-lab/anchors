<!-- @anchors
  code: AGCRG
  updated_at: 2026-10-08
  layer: comando
-->
# AgentCards — the cards this agent owns, and the card a pull request declares, read from the tracker

> **Code**: `AGCRG`

## Overview

The workflow commands need to know which cards belong to the agent that is running, and which card a
pull request is about. Both answers live in the issue tracker, not in the repository: a card is owned
by whoever wrote the LAST ownership comment on it, and a pull request names its card in its body with a
closing keyword. This unit asks the tracker's command-line client for those facts and turns its
tab-separated answer into cards (number, title, state).

It also shapes what goes back to the tracker: an issue title is one line of bounded length, and an
issue URL yields the issue number.

Everything here is best effort. The tracker is outside the process, and a command that only wanted a
hint (the cards of this agent, the card of a pull request) must not stop because the client is not
installed, not authenticated or offline: every missing ingredient answers "nothing" instead of an error.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the agent name | the `ANCHORS_AGENT` environment variable, trimmed | an empty or blank value | this unit: no agent name means no cards |
| the configuration | a project config with a workflow block that names the repository and at least one label | no config, no workflow block, a workflow with no label | this unit for a missing config or workflow (no cards); the caller for the label list, which today checks the GitHub mode first |
| the tracker client | the `gh` executable on the PATH | a PATH without it, or a client that exits with an error | this unit: both answer no cards / no card |
| the client's card lines | one card per line: number, title and state label separated by tabs | lines without three fields, or with an empty number | this unit: skips them |
| the reason text | any text, possibly multi-line | — | this unit: keeps the first line only |
| the issue URL | a URL whose last path segment is the issue number | a URL ending in a slash, a segment with non-digits, text with no slash | this unit: answers the empty string |
| the pull request reference | a number, with or without `#` and surrounding spaces, and a repository | an empty repository or an empty reference | this unit: answers the empty string without asking the tracker |

## Effects

| Effect | Description |
| --- | --- |
| `AGCRG-B01` | With no agent name, no configuration, no workflow block or a workflow with no label, or with no tracker client on the PATH, the agent's card list is empty and no error is raised. |
| `AGCRG-B02` | The tracker is asked for the open cards of the workflow repository carrying the workflow's first label, keeping only those whose last ownership comment names this agent. |
| `AGCRG-B03` | Each card line becomes a card with its number, its title and its state, the state losing the `anchors:` prefix of its label. |
| `AGCRG-B04` | A line that does not have exactly three fields, or whose number is empty, is skipped; the other cards are kept. |
| `AGCRG-B05` | A card with no state label is kept wherever it appears, including as the last line of the answer. |
| `AGCRG-B06` | `FirstLineOfReason`: an issue title made from a reason keeps only the reason's first line, trimmed of surrounding spaces. |
| `AGCRG-B07` | `FirstLineOfReason`: a title longer than 70 characters is cut so that, with a trailing ellipsis, it is exactly 70 characters long; a title of 70 characters or less is kept whole. Characters, not bytes: a multi-byte character is never split, and the title stays valid UTF-8. |
| `AGCRG-B08` | `NumeroDaIssue`: the issue number is the last path segment of an issue URL, trimmed, and only when it is made of digits alone; otherwise the answer is empty. |
| `AGCRG-B09` | `CardDoPR`: the card a pull request declares is the number after the first `Refs`, `Closes`, `Fixes` or `Resolves #N` that starts a line of its body, in any letter case; a keyword in the middle of a line does not count. |
| `AGCRG-B10` | `CardDoPR`: the pull request reference is trimmed and loses its leading `#` before the tracker is asked; with no repository or no reference the answer is empty. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `AGCRG-I01` | Every card that reaches the caller comes from a well-formed line: it has a non-empty number, and its state never carries the `anchors:` prefix. | feeds well-formed, malformed and state-less lines and compares the exact card list |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `AGCRG-X01` | Does not decide ownership itself: the choice of the last ownership comment is part of the question sent to the tracker. | The comments live in the tracker; fetching them all to filter locally would multiply the traffic for no gain. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `AGCRG-E01` | The tracker client exits with an error while listing the agent's cards. | No cards, no error. | The list is a hint for the next step; an offline or unauthenticated client must not stop the command that asked. <!-- @resilient: the card list is only a hint for the next step, and an offline or unauthenticated tracker is an ordinary state of a workstation, not a defect to report --> |
| `AGCRG-E02` | The tracker client exits with an error while reading a pull request body. | No card (empty answer). | Same reason: the caller treats "no card" as "no link declared" and moves on. <!-- @resilient: a PR body that cannot be read is read as no card linked, which is the answer the caller already handles; the tracker being unreachable is not this lookup's to report --> |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
