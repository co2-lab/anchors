<!-- @anchors
  code: ESDPS
  updated_at: 2026-10-08
  layer: comando
-->
# EscalateDuplicate — find the open cards that already deal with the target of an escalation

> **Code**: `ESDPS`

## Overview

Two agents delivered the same work on the same day: one took the gate's card for a spec at 11:38;
the other, working on a different card, hit the same problem at 12:02 and opened a new card for it,
and both pull requests added the same section to the same document. The claim keeps two agents from
taking the same card; it did not keep an agent from creating a card for work already under way in
another.

Before `anchors escalate` creates a card about a file, this unit asks the board for the open cards
that already deal with that file, so the escalation can warn about them. The platform's text search
is approximate, so every hit is confirmed by the exact target appearing in the card's title or body —
otherwise a search for one spec would match another whose name merely contains it, and a warning
that points at unrelated work teaches people to ignore it.

It is an aid and never a refusal: escalating the same file twice is legitimate (two distinct problems
in one spec), and a lookup that fails yields nothing instead of stopping the escalation.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the target | the path the escalation is about, or empty when it names none | — | this unit: an empty target asks nothing |
| the label | the project's first workflow label, or empty | — | this unit: an empty label asks nothing |

## Effects

| Effect | Description |
| --- | --- |
| `ESDPS-B01` | Without a target or without a label, the board is not asked at all. |
| `ESDPS-B02` | The board of the project's repository (`workflow.repo`, passed as `--repo`) is searched for open cards carrying the label, with the target as the search text — never the repository of the current directory. |
| `ESDPS-B03` | A search hit counts when the exact target appears in its title or in its body. |
| `ESDPS-B04` | Each card found is reported as `#<number> <title>`, so the reader knows which card to open and whether it is the same work. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `ESDPS-I01` | A card whose title and body do not contain the exact target is never reported, even when the search returned it. | a search returns a card about a file whose name contains the target, and nothing is reported |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `ESDPS-X01` | The lookup only reads the board; it creates, edits and refuses nothing. | Warning is the whole job; refusing would turn a common case into stopped work. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `ESDPS-E01` | The board lookup fails, or answers with something that is not a list of cards. | Nothing is reported, and no error reaches the escalation. | The check is auxiliary; stopping the escalation for it would be worse than the duplicate it prevents. <!-- @resilient: the duplicate lookup is auxiliary, and a board that cannot be read must not stop the escalation it only advises --> |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
