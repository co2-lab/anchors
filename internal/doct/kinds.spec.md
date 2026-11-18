<!-- @anchors
  code: DCKND
  updated_at: 2026-10-08
  layer: apoio
-->
# DocKinds — what each kind of project documentation must answer, told to the agent that writes it

> **Code**: `DCKND`

## Overview

Naming the file is not enough. "Update the OpenAPI document" leads an agent to add the endpoint
it just wrote and stop there — no errors, no examples, no body schema. The result passes any
existence check and is useless to whoever consumes it: a contract that only describes the happy
path is not a contract.

This unit holds, for each kind of documentation the project may declare (an API contract, a C4
architecture, a data schema, a component catalogue, a decision record, a runbook), the items that
kind of document must cover and the typical trap — the way of "doing" the task while producing
something that does not serve. It renders that into the duty text an agent's work card carries.
An instruction carries its why, because an instruction without a why becomes ritual: the agent
fulfils the form and misses the content.

A kind the product does not know is still named, never refused: a project may declare a
documentation Anchors has never heard of, and refusing it would make the mechanism serve only
what was already foreseen.

## Effects

| Effect | Description |
| --- | --- |
| `DCKND-B01` | Each of the six known kinds (openapi, c4, schema, components, adr, runbook) has a title, a non-empty list of items the document must cover, and a trap. |
| `DCKND-B02` | `InstructionFor`: The kind is looked up ignoring case and surrounding spaces, so ` OpenAPI ` finds the openapi instruction. |
| `DCKND-B03` | An unknown kind gets a minimal instruction with no items and no trap, titled by the kind — or by the document's path when the kind is empty. |
| `DCKND-B04` | `Duty`: The duty text is the title followed by the path in backticks, then the declared why when there is one, then one bulleted line per item, then the trap line when there is one — in that order. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared document | a path, an optional kind and an optional why, as the project configuration declares them | — (any kind, known or not, is accepted) | this unit: an unknown kind falls back to the minimal instruction |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCKND-I01` | `KnownKinds`: The list of known kinds and the set of instructions are the same set: every listed kind has an instruction, and no instruction exists for a kind the list omits. | looks up every listed kind and verifies it is not the fallback, and compares the counts of both sets |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCKND-X01` | The unit never refuses a documentation kind. | A project may declare a documentation Anchors does not know; what it cannot instruct, it at least names. |

## Errors

none — an unknown or empty kind is normal flow answered by `DCKND-B03`, and the unit reads nothing that can fail.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
