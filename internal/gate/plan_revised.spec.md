<!-- @anchors
  code: PLRVP
  updated_at: 2026-09-26
  layer: gate
-->
# PlanRevised — mutual revision visibility between superseded and revising plans

> **Code**: `PLRVP`

## Overview

Confronts plan nodes against declared revision relationships: **when a plan is revised, the superseded document must explicitly warn readers at the top, and the revising document is reminded until that warning exists.**

A plan is an anchor for architectural intent and execution. When planning errs — as it inevitably does once real implementation begins — editing an existing, already implemented plan destroys the historical record: it rewrites what was actually decided and executed into something that never happened.

For this reason, architectural revisions must be published as a new plan declaring which prior plan it revises. However, this discipline introduces a subtle and dangerous operational defect: out-of-order reading. A developer or agent opening an older, superseded plan has no natural indication that a newer revision exists, and will proceed to implement or follow decisions that have already been overturned. Because the old plan remains internally coherent — it was, after all, the accurate record of its own time — the reader cannot detect the obsolescence from the document alone.

This gate enforces bidirectional visibility across revised plans. When a plan is revised, it must display an explicit warning notice at its very top (within the first 40 lines) naming the revising plan, and must mark each section affected by the revision. Reciprocally, the author of a revising plan is reminded with a non-blocking pending verdict until that top notice is posted on the target file, ensuring that superseded documents are marked at the moment the revision is committed.

This gate operates in distinct territory from neighbouring gates:
- Unlike `plan-change-justified`, which evaluates whether edits inside an existing plan justify their architectural delta, this gate governs the relationship between two distinct plans across the project map.
- Unlike `phase-ordered`, which ensures linear dependency ordering of phases within a single plan, this gate addresses the lifecycle obsolescence of entire plans and phases across revision boundaries.
- Unlike `plan-seeds-valid`, which ensures seed files cited by a plan actually exist in the repository, this gate ensures that revision targets resolve to valid plan nodes and that mutual revision markers exist on disk.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind plan | nodes of kind spec, feature, code, or test | this unit: skips non-plan nodes |
| the map graph | a built graph containing project nodes | nil graph | this unit: returns pending when graph is absent |
| revision relationships | declared revision targets between plan nodes | unrelated nodes or self-referential revisions | the gate engine and map builder |
| revision notices | markdown alert blocks within the first 40 lines or metadata revision directives | unmarked prose revisions appearing after line 40 | this unit: parses document header and flags late notices |

## Effects

| Effect | Description |
| --- | --- |
| `PLRVP-B01` | When the confronted node is not a plan, the gate skips confrontation. |
| `PLRVP-B02` | When the map graph is nil, confrontation yields a pending verdict because revision relationships cannot be resolved. |
| `PLRVP-B03` | When a plan neither revises another plan nor is revised by any plan, the gate skips confrontation. |
| `PLRVP-B04` | When a revising plan declares a revision target that does not exist in the map graph, the gate fails. |
| `PLRVP-B05` | When a revising plan targets an existing plan whose file lacks a top revision notice, the revising plan receives a pending reminder. |
| `PLRVP-B06` | When a revising plan targets a plan whose file already contains the top revision notice, the pending reminder clears. |
| `PLRVP-B07` | When a revised plan lacks a top revision notice, the gate fails, reporting the revising plan identifier. |
| `PLRVP-B08` | When a revised plan places the revision notice past the first 40 lines, the gate fails because the warning arrives too late. |
| `PLRVP-B09` | When a revised plan has a top revision notice within the first 40 lines but no section amendment markers, the gate returns a pending verdict. |
| `PLRVP-B10` | When a revised plan contains both a top revision notice in the first 40 lines and section amendment markers, the gate passes. |
| `PLRVP-B11` | Both markdown alert blocks and metadata revision directives are accepted as valid notices and section amendments. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PLRVP-I01` | The top revision notice must reside within the first 40 lines of the revised document so that readers following top-down reading order encounter the warning before acting on obsolete decisions. | confronts a plan with notice placed after line 40 and verifies it fails |
| `PLRVP-I02` | A missing section amendment marker on a revised plan results in a Pending verdict rather than a hard Fail, accommodating whole-plan revisions where individual sections cannot be cleanly partitioned. | confronts a revised plan with top notice but no section markers and verifies it returns Pending |
| `PLRVP-I03` | The pending reminder on a revising plan clears as soon as the revised target file contains the required top notice, preventing permanent noise from training teams to disregard gate output. | provides the revised file with top notice and verifies the revising plan clears pending |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PLRVP-X01` | Does not mandate a specific human language for revision notices, accepting standard markdown alert callouts and language-agnostic directives. | Natural language parsing across multiple spoken languages leads to brittle false positives and false rejections. |
| `PLRVP-X02` | Does not assess the semantic accuracy or completeness of the explanatory prose written beside a revision marker. | Deterministic gates verify structural cross-references and marker presence; evaluating descriptive text quality requires human review. |
| `PLRVP-X03` | Does not enforce section amendment markers on plans that are not targeted by any revision. | Unrevised documents represent current baseline intent and require no diff annotations. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern`, `Config` | core — identity code length pattern and configuration model |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized messages for gate verdicts |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindPlan`, `Node` | core — graph structure, plan kind definition, and node representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
