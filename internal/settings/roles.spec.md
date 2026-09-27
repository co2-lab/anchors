<!-- @anchors
  code: AGRLG
  updated_at: 2026-09-27
  layer: apoio
-->
# AgentRoles — who is who in a project, and the capabilities each role carries

> **Code**: `AGRLG`

## Overview

A project with several contributors has more than one kind of work, and not everyone may do all of it.
The person running an agent declares a ROLE, and the agent's capabilities derive from it: the role is
what a person can say about themselves ("I am a dev"), and it survives Anchors gaining new capabilities,
which are born mapped to the existing roles.

The capability that motivated roles is deciding the product: acting on escalated cards that wait for a
product decision. Only the product owner and the architect decide the product, and only the architect
decides the structure (layers, boundaries, dependencies), since the two questions need different
context. QA and the reviewers execute and attack what exists; they decide nothing. Each role that
reviews carries a distinct lens, because "review" is not one thing: whoever hunts data leaks and whoever
hunts queries in loops read the same code with different questions.

A role is not hierarchy and not repository permission. An unknown role, such as one invented in a
hand-written settings file, can do nothing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the typed role name | a role name or a common abbreviation, in English or Portuguese, any case, with spaces | anything else | this unit: it is not recognised, and the caller asks again |
| the declared role | one of the known roles | an invented role | this unit: it can do nothing |

## Effects

| Effect | Description |
| --- | --- |
| `AGRLG-B01` | The known roles are dev, qa, reviewer-security, reviewer-performance, reviewer, product-owner and architect, listed in a stable order, each with a title, a sentence of what it does and at least one capability. (`KnownRoles`) |
| `AGRLG-B02` | Only the product owner and the architect decide the product. |
| `AGRLG-B03` | The architect decides the structure; the product owner does not. |
| `AGRLG-B04` | Dev, QA and the reviewers decide neither the product nor the structure, and all of them can review. |
| `AGRLG-B05` | QA and each reviewer role carry a review lens, all distinct; the dev carries none. |
| `AGRLG-B06` | A typed role is recognised by its name or a common abbreviation, in English or Portuguese, ignoring case and surrounding spaces; anything else is not recognised. (`ParseRole`) |
| `AGRLG-B07` | A role lists its capabilities in alphabetical order, whatever order they were declared in, so the same role always shows them the same way. (`Caps`) |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `AGRLG-X01` | An unknown role, including a known name in another case, has no capability at all. | A hand-written settings file with an invented role must not unlock anything: the cost of erring open is someone deciding without authority, invisible after the fact. |

## Errors

none — an unknown role is a value that can do nothing (X01), and an unrecognised name is an empty answer the caller asks again about.

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
