<!-- @anchors
  code: DFGTD
  updated_at: 2026-09-26
  layer: infra
-->
# DefaultGates — the gates a project is born with, by artifact and by project age, and the canonical gate catalog

> **Code**: `DFGTD`

## Overview

A project should not have to write its gates by hand. This unit decides which gates `anchors init` seeds,
and in which maturation state, from two answers: which artifacts the project uses, and whether the project
is new or already exists. The same list, with every artifact turned on, is the canonical catalog of gates:
the one place where each gate's wording lives, which configuration loading uses to complete a gate the
project declared by name alone.

The gates follow the artifacts. A project that chose nothing is seeded with nothing, because a gate born
with nothing to measure is the impression of a defence that does not exist. Spec, feature and test each
bring their own gates; the gates that cross two artifacts appear only when both are chosen; guides bring
their checklist gate; the gate that confronts the declared parent runs on whichever of spec, plan and code
were chosen.

The age of the project decides the maturation state. In an existing project a gate is born informative,
because a real project almost never meets on day one the threshold it wants to reach, and blocking at once
would stop the work. Five gates are the exception and are born blocking whatever the age, because what they
catch cannot be matured into: a leaked secret, a contracted document that does not mention the unit that
triggers it, and the grammar, completeness and existence of flag scenarios. In a new project the premise
inverts: there is no debt to accommodate, so every gate is born blocking and stops the first deviation,
when fixing costs least. The exception there is the gates that read an ingested report (tests, coverage,
mutation, scenario coverage, and the external-tool checks of dependencies, duplication, licences, cycles,
dead code and spelling): with no report yet they would block for lack of data, not for a defect, so they
stay informative until the project promotes them.

A judgment gate that asks whether the code or a test realises a rule carries the @TBD instruction
(INCTN-B04), so a piece the spec declared as still to be written is waived by name instead of approved.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the chosen artifacts | a set of artifact names, possibly empty | names no gate depends on | this unit: they seed nothing |
| the project age | new or existing | — | the init flow, from whether the project already has content |
| the gate name to resolve | any string | — | this unit: an unknown name is not found |

## Effects

| Effect | Description |
| --- | --- |
| `DFGTD-B01` | A project with no artifact chosen is seeded with no gate, new or existing. |
| `DFGTD-B02` | A project with spec, feature and test is born with the gates of each, among them the spec completeness, feature non-emptiness, green tests, line coverage and scenario coverage gates. |
| `DFGTD-B03` | The gates that cross spec and feature (spec-feature match and scenario coverage) are seeded only when both are chosen. |
| `DFGTD-B04` | Choosing only guides seeds only the guide checklist gate. |
| `DFGTD-B05` | The parent gate is seeded when spec, plan or code is chosen, and confronts exactly the chosen ones among them. |
| `DFGTD-B06` | In an existing project every gate is born informative, except five born blocking by nature: the leaked-secret, contracted-document and the three flag-scenario grammar, completeness and existence gates. |
| `DFGTD-B07` | In a new project every gate is born blocking, except the gates that read an ingested report. |
| `DFGTD-B08` | A gate that reads an ingested report (tests, coverage, mutation, scenario coverage, and the external-tool checks) stays informative even in a new project. |
| `DFGTD-B09` | Every judgment gate that asks about the code or a test carries the @TBD instruction; the one that asks about the proof of a permanent test waiver does not. |
| `DFGTD-B10` | The canonical declaration of a gate is found by name in the catalog of every artifact, including gates seeded only for guides or plans; an unknown name is not found. |
| `DFGTD-B11` | Loading a configuration completes a canonical gate declared by name alone with the catalog's declaration. |
| `DFGTD-B12` | The gate is of the BLOCKING class: a new project is born with it blocking, and an existing one takes it through the same maturation as every structural gate — informative until the project promotes it. It never depends on an ingested signal, so it can block from day one. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DFGTD-I01` | The list of seeded gates, and its order, does not change with the age of the project; only the maturation state does. | seeds the same artifacts as a new and as an existing project and compares the lists name by name |
| `DFGTD-I02` | Every gate of the full catalog has a unique name, and its id is its name. | seeds every artifact and checks names for repetition and id equality |
| `DFGTD-I03` | Every canonical name the format migration renames a legacy gate to is a gate of the full catalog. | walks every rename of the migration steps and looks each target up in the catalog |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DFGTD-X01` | No default gate carries a legacy name. | Gate names are identifiers written into every project's configuration; a legacy name would be renamed by the migration and then stamped twice in the map. |

## Errors

none — seeding builds a list in memory from the answers it receives, and resolving a name that is not in the catalog is the normal "not found" answer (DFGTD-B10), not a failure.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Gate`, `Bool`, `SetCanonicalGateResolver` | core — gate declarations and the canonical merge |
| DEP2 | `internal/config/vocabulary.go` | `RegisterGateNames` | core — the default gate names for the vocabulary check |
| DEP3 | `internal/initx/tbd_instruction.go` | `tbdInstruction` | infra — the @TBD instruction (INCTN) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
