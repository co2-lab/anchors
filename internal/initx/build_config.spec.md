<!-- @anchors
  code: BLCNB
  updated_at: 2026-09-26
  layer: infra
-->
# BuildConfig — builds the configuration that inference proposes as the default for the init questions

> **Code**: `BLCNB`

## Overview

What inference found is turned into a proposed configuration, which the interactive init then confirms or
adjusts part by part. The proposal only pre-fills what inference can decide from structure: the code layers,
the colocation of the triad, and the attribute the project uses to mark elements for tests.

Each detected code directory becomes one code layer, named after the last segment of the directory with a
`-code` suffix and tagged with that name. Its pattern covers the directory recursively with the detected
code extensions, and it excludes specs, features and test files, which are artifacts and not code.

Artifact layers are not proposed here: they come from the user's choice, which is always asked
(ARCHR-B03). The governs rules are left empty, because matching a guide to a tag is a semantic answer that
only the user gives.

Colocation is proposed only when inference detected it, anchored on the spec, with a template only for the
derivative kinds it found (feature, test); the spec is the anchor and never one of its own templates. The
test handle is proposed only when inference found one in use: a backend, a library, or a project that marks
nothing gets no handle, so the gates that inventory handles skip instead of accusing the whole repository of
breaking a convention the project never promised.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the inference result | code directories, code extensions, detected kinds, colocation and handle, any of them possibly empty | — | inference (INPRN) produces it |

## Effects

| Effect | Description |
| --- | --- |
| `BLCNB-B01` | Each detected code directory becomes a code layer named (`LayerNameFor`) after the directory's last segment with the `-code` suffix, tagged with that name. |
| `BLCNB-B02` | The pattern of a code layer covers its directory recursively: one detected extension is written alone, several are written as a set. |
| `BLCNB-B03` | A proposed code layer excludes specs, features and test files. |
| `BLCNB-B04` | Colocation is proposed only when inference detected it, anchored on the spec, with a template only for the derivative kinds inference found (feature, test); the spec is never among the templates. |
| `BLCNB-B05` | The test handle is proposed only when inference found one, with or without colocation; with none found, no default handle is written. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `BLCNB-I01` | REF[BLCNB-B05]: a proposal never carries a test handle the project does not use | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `BLCNB-X01` | The proposal creates no artifact layer and no governs rule. | Artifact layers come from the user's choice, always asked, and a governs rule is a semantic answer that structure cannot give. |

## Errors

none — the proposal is assembled in memory from the inference result; nothing is read, and no input is rejected.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Layer`, `Derived` | core — project configuration |
| DEP2 | `internal/initx/infer.go` | `Proposal` | infra — what inference found (INPRN) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
