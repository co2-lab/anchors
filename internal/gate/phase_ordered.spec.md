<!-- @anchors
  code: PHORP
  updated_at: 2026-09-26
  layer: gate
-->
# PhaseOrdered — plan phases and phase dependencies must be ordered and consistent

> **Code**: `PHORP`

## Overview

Confronts internal plan phase ordering and cross-artifact phase dependencies: **phases within a plan must
follow consistent non-cyclic order, and specifications citing phase prerequisites must target existing phases.**

The `needs:` declaration previously resolved execution order between separate plans, but ordering within a single
plan lived purely as informal prose (such as `### Fase 2 — a régua mecânica (depende da Fase 1)`). Natural language
prose cannot be mechanically verified. When task identification pipelines generate work cards for each specification
in a plan, all cards are born into the backlog indistinguishably, and automated claim pipelines deliver any card
without respecting prerequisite readiness.

Measured in the very first real-world usage: an automated agent was assigned the specification for a test harness
(Phase 3) while Phase 1 and Phase 2 remained completely open. Without even a `package.json` created in the repository,
there was nowhere to configure tooling. The assignment only avoided becoming lost work because a human intervened to
read the narrative plan; an agent trusting the task card would have attempted execution and failed.

Phases are therefore promoted to catalogued items with identity codes derived from the plan (such as `FNDTN-F01`).
A seeded specification declares `needs: FNDTN-F01` in its header, adopting the exact ordering keyword at phase scope.
Identity codes remain immutable even when authors revise descriptive phase titles, transforming "can this specification
be worked on now?" into an objective mechanical query.

This gate confronts three complementary structural ordering contracts:
1. **Phase order within plans (`phase-ordered`)**: Phases cannot depend on future phases, duplicate codes, or
   themselves. Plans with phase-like sections lacking catalogued codes return **Pending**.
2. **Phase references in specifications (`phase-exists`)**: A specification declaring `needs:` must target catalogued
   phases that exist in project plans, avoiding permanent task blockages.
3. **Parent hierarchy integrity (`parent-valid`)**: Artifacts declaring a `parent:` must target real entities
   (artifacts or phases) without dangling references or circular parent chains.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted plan | nodes of kind `plan` | nodes of kind `spec`, `code`, `feature`, or `test` | this unit: skips non-plan nodes for plan ordering |
| the confronted specification | nodes of kind `spec` declaring `needs:` | specifications with empty `needs:` | this unit: skips specifications without phase dependencies |
| the phase declarations | level-three markdown sections with catalogued phase codes | unnumbered or freeform prose paragraphs | this unit: extracts phases via `PlanPhases` |
| the graph model | a built graph with project nodes | nil graph pointer | this unit: returns Pending when graph is absent |
| the parent references | artifact nodes declaring a non-empty `parent:` | artifacts without parent declarations | this unit: skips artifacts with empty parent |

## Effects

| Effect | Description |
| --- | --- |
| `PHORP-B01` | The `PlanPhases` function extracts catalogued phase codes from plan headers in appearance order. |
| `PHORP-B02` | When the confronted node is not of kind plan, phase ordering skips confrontation. |
| `PHORP-B03` | When a plan contains no phase headings, phase ordering skips confrontation. |
| `PHORP-B04` | When a plan contains phase-like sections without catalogued phase codes, phase ordering returns Pending. |
| `PHORP-B05` | When plan phases declare valid backward dependencies on preceding phases, phase ordering passes. |
| `PHORP-B06` | When a plan defines duplicate phase codes, phase ordering fails. |
| `PHORP-B07` | When a phase declares a dependency on a phase code not catalogued in the plan, phase ordering fails. |
| `PHORP-B08` | When a phase declares a dependency on itself, phase ordering fails. |
| `PHORP-B09` | When a phase declares a dependency on a future phase defined later in the plan, phase ordering fails. |
| `PHORP-B10` | When a specification declares phase dependencies that exist in project plans, phase existence passes. |
| `PHORP-B11` | When a specification declares phase dependencies that do not exist in any plan, phase existence fails naming the missing phases. |
| `PHORP-B12` | When an artifact declares an existing artifact code or catalogued phase as parent, parent validation passes. |
| `PHORP-B13` | When an artifact declares a non-existent parent, self-parenting, or a circular parent chain, parent validation fails. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PHORP-I01` | Phase dependencies within a plan must be strictly acyclic and backward-directed; a phase cannot depend on itself or subsequent phases. | confronts plans with self or forward dependencies and verifies the verdict is Fail |
| `PHORP-I02` | Every declared phase or parent reference must resolve to an existing entity in the map, preventing orphaned items that disappear from dependency hierarchies. | confronts specs and artifacts with nonexistent targets and verifies verdict is Fail |
| `PHORP-I03` | Parent chains are strictly cycle-free and bounded to prevent infinite traversal loops during tree assembly. | confronts cyclical parent definitions and verifies verdict is Fail citing the cycle |
| `PHORP-I04` | Phase detection identifies structural level-three section boundaries regardless of linguistic naming variations. | confronts plans with unadorned level-three sections and verifies verdict is Pending |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PHORP-X01` | Does not mandate that small plans define catalogued phases. | Phase ordering is optional for small single-phase plans; requiring phases everywhere would impose ceremony on simple tasks. |
| `PHORP-X02` | Does not enforce timing deadlines or calendar durations for phases. | This gate validates logical ordering prerequisites; calendar scheduling belongs to external project management tools. |
| `PHORP-X03` | Does not restrict parent references to a single hierarchy kind. | Projects organize work flexibly; both artifact codes and plan phase codes are permitted as valid parent references. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern`, `Config` | core — code length pattern validation and project configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindPlan`, `KindSpec`, `Node` | core — graph model and artifact representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
