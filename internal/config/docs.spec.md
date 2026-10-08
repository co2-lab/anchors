<!-- @anchors
  code: DCRQA
  updated_at: 2026-10-08
  layer: config
-->
# DocsRequired — which aggregate documentation a change of a unit obliges to touch

> **Code**: `DCRQA`

## Overview

Anchors always charged what lives INSIDE a unit: the spec, the feature, the test, the code.
Some documentation lives outside every unit and is fed by many of them: the API contract,
the component catalogue, the data schema, the architecture. A new endpoint that never
reaches the API contract is invisible to whoever consumes it, and no unit gate sees that,
because the contract is not a file of the unit.

The project declares these documents, and declaring is the act: a project that declares
none is charged none. Each declared document names its file, why it exists, and its
triggers — the layers or unit codes whose change obliges touching it. This unit answers the
question the agent's card and the duty listing ask: given the layer of a change, and
optionally the unit's code, which declared documents are owed?

A trigger accepts both granularities because the layer alone gets it wrong. Measured in the
reference project: its infrastructure layer has nine units and only one touches the data
schema; charging the schema to all nine teaches the agent to ignore the warning. A layer
trigger charges the whole layer; a code trigger charges that unit and not its neighbours. A
document with no trigger (the architecture, which changes when the structure changes, not
when a unit does) is never owed by a unit change.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project configuration | a loaded configuration, with or without a docs block, or none at all | — | this unit: no configuration or no block means nothing is owed |
| the layer of the change | any layer name, in any case | — | this unit: triggers are compared ignoring case and spaces |
| the unit code | a unit code, or nothing when the question is about the layer in the abstract | — | this unit: no code, or an empty one, answers for the layer alone |

## Effects

| Effect | Description |
| --- | --- |
| `DCRQA-B01` | A trigger naming a layer charges every change in that layer, whatever the unit (`TriggeredBy`). |
| `DCRQA-B02` | A trigger naming a unit code charges that unit, and not its neighbours in the same layer. |
| `DCRQA-B03` | Asked without a unit code, the answer is the layer's alone: no code trigger, not even a blank one, can match an absent code (`RequiredFor`). |
| `DCRQA-B04` | A documentation with no trigger is never owed by a unit change. |
| `DCRQA-B05` | A trigger matches ignoring case and the spaces around it. |
| `DCRQA-B06` | Every declared documentation can be listed; a project with no docs block, or no configuration, owes none and lists none (`AllRequiredDocs`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCRQA-I01` | Naming the unit never removes a documentation the layer alone owes: the answer with a code contains the answer without one. | asks for several layers with and without codes and checks every document owed by the layer is still owed with the code |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCRQA-X01` | A trigger is matched whole, never as part of a longer or shorter name. | A partial match would charge a layer or a unit the project never named, which is the noise that teaches the agent to ignore the duty. |

## Errors

none — an undeclared block, a document without triggers and a change that triggers nothing are answers, not failures.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
