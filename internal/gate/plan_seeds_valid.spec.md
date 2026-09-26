<!-- @anchors
  code: PSVPL
  updated_at: 2026-09-26
  layer: gate
-->
# PlanSeedsValid — specifications seeded in a plan must target valid governed layers

> **Code**: `PSVPL`

## Overview

Confronts artifact paths seeded by an implementation plan against the project Structure: **every specification
path that a plan promises to create must belong to a valid governed layer.**

An implementation plan seeds new artifacts to be born during execution. When a plan declares which specifications
will be created, developers and automated agents rely on those paths to execute their tasks. If a plan seeds
specifications in layers where specifications are prohibited or undefined, that structural contradiction is only
discovered downstream during execution.

Observed across three consecutive rounds in project history: an implementation plan repeatedly seeded
`packages/backend/models/metadata.spec.md — **nasce**`, even though `models/` was a recognized declarative layer
(`regime: declarativo`) which does not accept specifications by architectural definition. Each time, a full
execution round was spent rediscovering the contradiction, because downstream defenses (such as creation tools
and gate runners) only act at execution time. The root cause of the defect was the plan itself. Because the plan
is a declared layer in the repository, its structural promises can and must be verified before execution starts.

This gate confronts two structural defects in plan seeds:
1. **Seeding in a declarative layer**: Seeding a specification in a layer declared as `regime: declarativo`
   (which by definition does not carry specifications).
2. **Seeding in an undeclared layer**: Seeding a specification in a concrete repository directory that does not
   match any declared layer in the project Structure (indicating a path typo or an undeclared layer).

What separates this gate from neighbouring gates:
- It deliberately does NOT verify whether the seeded specification already exists on disk (seeds are intended
  to be created during execution).
- It does NOT check whether plan progress is synchronized or complete (which belongs to plan progress gates).
- It does NOT inspect the content or implementation steps of the plan items.
- It restricts evaluation strictly to the structural validity of the specification paths the plan promises to create.

Finally, casual mentions of template specifications (`_TEMPLATE_*.spec.md`), bare file names in prose, and
informal path fragments that do not correspond to top-level repository directories are recognized as narrative
references rather than actionable seeds.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind `plan` | nodes of kind `spec`, `code`, `feature`, or `test` | this unit: skips non-plan nodes |
| the project configuration | loaded project configuration with layer definitions | nil configuration pointer | this unit: returns Pending when configuration is nil |
| the plan content | markdown text containing backticked `.spec.md` references | plans containing no backticked specification paths | this unit: skips plans without specification seeds |
| the seed candidates | backticked specification paths with directory separators | template specifications prefixed with `_TEMPLATE` | this unit: filters out template files |
| the layer mapping | paths classified into declared layers via `scan.Classify` | paths classified as belonging to no layer | `scan.Classify`: resolves target file path against layer patterns |

## Effects

| Effect | Description |
| --- | --- |
| `PSVPL-B01` | When the confronted node is not of kind plan, the gate skips confrontation. |
| `PSVPL-B02` | When project configuration is nil, confrontation returns Pending. |
| `PSVPL-B03` | When the plan contains no seeded specification paths, the gate skips confrontation. |
| `PSVPL-B04` | Template specification references prefixed with `_TEMPLATE` are ignored as templates. |
| `PSVPL-B05` | Bare specification file names without directory paths are ignored as prose references. |
| `PSVPL-B06` | Informal path abbreviations whose top-level directory does not exist on disk are ignored as prose references. |
| `PSVPL-B07` | When all seeded specification paths target valid governed layers, the gate passes. |
| `PSVPL-B08` | Multiple target source file extensions are evaluated when resolving the governed layer. |
| `PSVPL-B09` | When a seeded specification targets a declarative layer, the gate fails citing the declarative layer. |
| `PSVPL-B10` | When a seeded specification path in a real repository directory matches no declared layer, the gate fails. |
| `PSVPL-B11` | Multiple seed defects across declarative and undeclared layers are sorted and aggregated in the failure verdict. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PSVPL-I01` | Only plan artifacts are evaluated; all other artifact kinds skip confrontation. | confronts non-plan nodes and verifies the verdict is Skip |
| `PSVPL-I02` | Declarative layers never accept specifications because declarative schema and model files are recognized rather than authored with specifications. | confronts plans seeding declarative layer paths and verifies the verdict is Fail |
| `PSVPL-I03` | Seed verification is purely structural and never inspects specification implementation or progress status. | confronts plans with non-existent future seeds in valid layers and verifies they pass |
| `PSVPL-I04` | Casual prose citations and template references are never confused with actionable artifact creation promises. | confronts plans with unpathed mentions or templates and verifies they are skipped |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PSVPL-X01` | Does not verify whether seeded specifications currently exist on disk. | A plan describes future work to be executed; seeds are expected to be created during subsequent development. |
| `PSVPL-X02` | Does not check whether the plan progress file is synchronized with execution state. | Tracking task execution status is the dedicated responsibility of plan progress gates. |
| `PSVPL-X03` | Does not enforce specification contents or scenario definitions within seeded files. | Quality and completeness of specification content belong to specification and triad gates once created. |

## Errors

Each failure the code handles is already stated as a rule of another letter; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PSVPL-E01` | REF[PSVPL-B02]: with no configuration there is nothing to confront, and B02 answers Pending | — | — |
| `PSVPL-E02` | REF[PSVPL-B06]: a path whose top directory is not on disk is answered by B06: read as prose, not charged | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project layer definitions and layer regimes |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindPlan`, `Node` | core — graph model and artifact representations |
| DEP4 | `internal/scan/scan.go` | `Classify` | core — classifies file paths into declared project layers |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
