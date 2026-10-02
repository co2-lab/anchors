<!-- @anchors
  code: DFGTD
  updated_at: 2026-10-02
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

The gates follow the artifacts they run on. A project that chose nothing is seeded with nothing, because a
gate born with nothing to measure is the impression of a defence that does not exist. Spec, feature, test and
code each bring their own gates — the spec's include the documentation, doctrine, flag and failure gates —
and plans bring only the gates that run on plans; the gates that cross two artifacts appear only when both
are chosen; guides bring their checklist gate; the gates that confront the declared parent and the justified
change run on whichever of their artifacts were chosen. Choosing every artifact `anchors init` offers seeds
the whole catalog.

What a gate carries — its question, its measure, its install hint — is written into the project's
configuration, so it is written in English.

The age of the project decides the maturation state. In an existing project a gate is born informative,
because a real project almost never meets on day one the threshold it wants to reach, and blocking at once
would stop the work. Two gates are the exception and are born blocking whatever the age, because what they
catch cannot be matured into: a leaked secret, and a contracted document that does not mention the unit that
triggers it. The flag-scenario gates follow the age like the rest: an existing project's flags were written
before the grammar was asked of them, and that is debt to mature, not a leak. In a new project the premise
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
| `DFGTD-B03` | The gates that cross spec and feature (spec-feature match, feature-spec match and scenario coverage) are seeded only when both are chosen, and the test-feature match only when test and feature are both chosen. |
| `DFGTD-B04` | Choosing only guides seeds only the guide checklist gate. |
| `DFGTD-B05` | The parent gate is seeded when spec, plan or code is chosen, and confronts exactly the chosen ones among them. |
| `DFGTD-B06` | In an existing project every gate is born informative, except two born blocking by nature: the leaked-secret and the contracted-document gates. |
| `DFGTD-B07` | In a new project every gate is born blocking, except the gates that read an ingested report. |
| `DFGTD-B08` | A gate that reads an ingested report (tests, coverage, mutation, scenario coverage, and the external-tool checks) stays informative even in a new project. |
| `DFGTD-B09` | Every judgment gate that asks about the code or a test carries the @TBD instruction; the one that asks about the proof of a permanent test waiver does not. |
| `DFGTD-B10` | The canonical declaration of a gate is found by name in the catalog of every artifact, including gates seeded only for guides or plans; an unknown name is not found. |
| `DFGTD-B11` | Loading a configuration completes a canonical gate declared by name alone with the catalog's declaration. |
| `DFGTD-B12` | The gate is of the BLOCKING class: a new project is born with it blocking, and an existing one takes it through the same maturation as every structural gate — informative until the project promotes it. It never depends on an ingested signal, so it can block from day one. |
| `DFGTD-B13` | Choosing every artifact `anchors init` offers (ARCHR-B01) seeds every gate of the catalog, the code gates included. |
| `DFGTD-B14` | The gates that run on specs, features and tests — documentation, doctrine, flags, failures, the way back of the unit, the contracted document and the justified change — are seeded without plans; the justified-change gate then runs on specs alone, and the gates that run only on plans are not seeded. |
| `DFGTD-B15` | The gate names registered for the vocabulary check (`RegisterGateNames`) are the full catalog, in catalog order. |
| `DFGTD-B16` | Choosing specs seeds `header-valid`, informative, on every kind Anchors governs — spec, feature, code, test, guide, doc, plan, product and flag —: every governed file carries the `@anchors` header. A bare `- name: header-valid` in a configuration inherits that `on:`. |
| `DFGTD-B17` | `no-duplication` is seeded as the native `duplication` check on code files, needing `npx`: one verdict per file read from jscpd's report, not a project-wide command judged by its exit code. |
| `DFGTD-B18` | The mock gates name the field each one reads: `mock-detect-covers-dialect` and `mock-stamped` presuppose `derived.mock_detect`, `mock-typed` presupposes `derived.mock_contract`; where it is not declared nothing is asked or measured (GTENG-B29), and the doctor names it (DCTRO-B27). |
| `DFGTD-B19` | `rule-fulfilled` is judged and also marked to review: the marks it judges were put by the agents who wrote the code, and a second look is a different thing from their judgment. |
| `DFGTD-B20` | The catalog carries the checkers that measure a kind of the unit and had no canonical gate — among them `evidence-fresh`, `feature-test-match`, `rule-implemented`, `scenario-identity`, `scenario-letter-declared`, `placeholder-filled` and `updated-at-atual` — each on the kinds it measures, with the field it presupposes when it needs one (`scenario-type-aligned` the `rule_types` tags), and the UI ones scoped by the `screen` and `component` tags. |
| `DFGTD-B21` | Init seeds (`SeedFor`) the default gates related to the project — a declared layer of a kind they measure, with their tags, or a kind the project chose and has no layer of yet — whose presupposed fields the configuration declares; a gate over an undeclared field is not seeded. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DFGTD-I01` | The list of seeded gates, and its order, does not change with the age of the project; only the maturation state does. | seeds the same artifacts as a new and as an existing project and compares the lists name by name |
| `DFGTD-I02` | Every gate of the full catalog has a unique name, and its id is its name. | seeds every artifact and checks names for repetition and id equality |
| `DFGTD-I03` | Every name the format migration renames a legacy gate to, followed through the later steps that rename it again, is a gate of the full catalog. | walks every rename of the migration steps, follows each target through the later steps, and looks the name it lands on up in the catalog |
| `DFGTD-I04` | Choosing plans adds only gates that run on plans. | seeds two choices with and without plans and checks every added gate runs on plans |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DFGTD-X01` | No default gate carries a legacy name. | Gate names are identifiers written into every project's configuration; a legacy name would be renamed by the migration and then stamped twice in the map. |
| `DFGTD-X02` | No default gate writes Portuguese into the project: its question, measure and install hint are English. | They are written into every new project's configuration, read by the judge and by people of any language; the catalog of translations is for what Anchors prints, not for what it writes into a project. |

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
