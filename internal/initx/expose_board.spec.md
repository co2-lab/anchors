<!-- @anchors
  code: BREXB
  updated_at: 2026-10-08
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

Because the page it hands is the one both boards show, this unit also states what that page does with the
cards, and its tests run the page's script. The page answers two questions the columns alone do not. The
first is what is stuck: a card waiting for a person has no column of its own, so a strip lists the cards
waiting for a decision, the ones waiting only for framing and the bugs, each apart, ordered by how many
cards each one blocks, and a blocked card names the card that holds it. The second is who is working right
now: one chip per agent owning an open card touched in the last 30 minutes, counted from the moment the
board was taken, or from now on a live board, and a released card has no agent.

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
| `BREXB-B04` | The page counts as active each agent that owns an open card touched within the 30 minutes before the moment the board was taken, the edge included, with how many of its open cards it owns; a released card has no owner and a closed card is not counted. |
| `BREXB-B05` | On a live board the page counts the 30-minute window from now, not from the moment stamped on the board. |
| `BREXB-B06` | The page shows one chip per active agent, sorted by name, plus a chip for all agents, or a quiet line when nobody is active; selecting an agent shows only that agent's cards, and selecting none shows every card. |
| `BREXB-B07` | The collect expression gives a card whose last owner comment releases it an empty owner, takes the owner from the first line of the last owner comment otherwise, and carries the card's last update. |
| `BREXB-B08` | With no card waiting for a person the blocked strip is hidden and empty. |
| `BREXB-B09` | The strip lists the cards waiting for a person, the one that blocks more cards first, saying how many cards each blocks and how many wait, with the title stripped of its code prefix; a card that is itself waiting for a person is not counted as blocked. |
| `BREXB-B10` | A card escalated to a person enters the strip whatever its work state, and the strip says the state it stopped in unless that state is to-do. |
| `BREXB-B11` | An escalated card that another card unblocks is shown waiting for that card's delivery, and the unblocking card does not enter the strip. |
| `BREXB-B12` | The strip escapes the card titles, so markup in an issue title is shown as text and never run. |
| `BREXB-B13` | In the columns, an escalated card, and only it, carries a class of its own and a badge saying it waits for you. |
| `BREXB-B14` | Clicking a card's label in the roadmap opens its details, as a click on a card does in the other tabs. |
| `BREXB-B15` | A blocked card shows a badge with the number of the card that blocks it, which leads to that card, in a color defined for the light palette and for both dark ones. |
| `BREXB-B16` | The cards waiting for a person are split into those that ask a decision and those that ask only framing, and the framing list says how to send a card back to to-do. |
| `BREXB-B17` | A bug card is not counted as waiting for you: it is listed in a strip of its own. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `BREXB-I01` | The extracted expression never carries the shell's output redirection: it opens and closes an array and names no output file. | extracts the expression from the carried pipeline and checks its ends and the absence of the redirection target |
| `BREXB-I02` | Every data attribute the page's handlers look for is written by one of its renders, and every data attribute a render writes is looked for by a handler. | collects the data attributes the page writes and the ones it reads, with the hyphenated and camel-case spellings made one, and compares the two sets |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `BREXB-X01` | Does not keep a built-in collect expression to fall back on. | A second copy of the contract is the drift this unit exists to prevent. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BREXB-E01` | The binary does not carry the board page or the board pipeline. | An empty result and an error saying which one is not carried. | Serving an empty board would look like a board with no cards. |
| `BREXB-E02` | The carried pipeline has no collect expression of the expected shape. | An empty result and an error saying the pipeline changed shape. | The local board would otherwise read a contract different from the one that is published. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
