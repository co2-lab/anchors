<!-- @anchors
  code: PRBDP
  updated_at: 2026-10-08
  layer: comando
-->
# PRBody — write the lines that link a pull request to its cards, in the platform's syntax

> **Code**: `PRBDP`

## Overview

`anchors pr-body` prints the lines a pull request body needs to link it to the cards it delivers.
The link is declared in the Anchors vocabulary — the card the agent took and the findings born
under it — and the line the platform understands is generated from that. Whoever writes the pull
request does not need to know that the platform only accepts its link words in English; in a project
written in another language, that is exactly what goes wrong in silence: the pull request merges and
the card links to nothing.

The line LINKS without closing. A closing keyword makes two claims at once — the implementation is
done, and the card is done — and they are not the same: after the merge the card still crosses the
delivery columns (in test, ready to release, production) that Anchors does not track. In the
reference project, closing on merge left 71 cards closed at ready-to-test against 7 open, so the
column meant to accumulate what waits for testing showed only the residue. Whoever finishes the
pipeline closes the card.

Each requested card drags the open findings born under it, because a finding discovered doing that
work is delivered with it. The `--so-sob` switch prints only those dragged findings: it serves CI,
which already read the root cards from the body and checks the findings are there too.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the requested cards | a comma-separated list of numbers, each with or without `#` and surrounding spaces | blank entries | this unit: drops them |
| the agent's cards | the cards whose last owner is the running agent, when no list is given | an agent with no card | this unit: refuses with the two ways to name a card |
| the workflow mode | github mode | local mode, where there is no card to link | this unit: refuses local mode |

## Effects

| Effect | Description |
| --- | --- |
| `PRBDP-B01` | Outside github mode the command refuses, saying that in local mode there is no card to link. |
| `PRBDP-B02` | Every line has the form `Refs #<n>`, which links the card and does not close it; no platform syntax uses a closing word (close, fix, resolve and their forms). |
| `PRBDP-B03` | The requested cards accept the forms a person writes — `44`, `#44`, ` 44 ` — and a comma-separated list; a blank list names no card. |
| `PRBDP-B04` | Without requested cards, the roots are the cards the running agent owns on the board. |
| `PRBDP-B05` | With no card from either source the command refuses, naming `--cards` and `ANCHORS_AGENT`. |
| `PRBDP-B06` | Each root drags every open card labelled as born under it, looked up in the configured repository. |
| `PRBDP-B07` | The lines come out in numeric order, so #101 comes after #45. |
| `PRBDP-B08` | With `--so-sob`, only the dragged findings are printed and the roots are left out. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRBDP-I01` | Each card appears in at most one line, whether it came as a root, under a root, or both. | a finding that is also a root, and under two roots, is printed once |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRBDP-X01` | The command only reads the board and prints; it writes nothing to the platform. | The lines go into a body someone else writes; changing cards is the pipeline's job when the pull request moves. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRBDP-E01` | The lookup of the findings under a root fails or returns unreadable output. | That root contributes no finding, and the root itself is still linked. | A partial body still links the work that was done; refusing would leave the pull request with no link at all. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
