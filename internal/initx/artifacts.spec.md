<!-- @anchors
  code: ARCHR
  updated_at: 2026-09-26
  layer: infra
-->
# ArtifactChoice — turns the artifacts the user chose at init into artifact layers and colocation

> **Code**: `ARCHR`

## Overview

`anchors init` always asks which artifacts the project will use, even in an empty project: the user marks
what they intend to use, not only what already exists. Inference only decides which options come
pre-checked. This unit holds the offered artifacts and applies the answer to the configuration.

The offered artifacts are spec, feature, test, guide, plan and code, always in that order. Code is a choice
with no layer of its own: the code layers come from inference and the preset, one per directory, and the
choice is what brings the gates that run on code and puts the code beside the spec in colocation. The
artifacts that inference found are pre-checked — code when inference found a code directory — and nothing
else is. Applying the choice rebuilds the artifact layers:
each chosen artifact gets its layer even when nothing was detected, each unchosen artifact loses its layer,
and a chosen artifact whose layer already exists keeps it exactly as declared. Code layers are never touched.
Guides and plans live in a folder, so their layer takes the detected folder when there is one, and the
default pattern otherwise.

Colocation is declared from the spec, which is the anchor: from the spec the code, the feature and the test
are born, beside it. The spec is never listed among the derivatives, because it is the origin. With no spec chosen
there is no anchor and no colocation, and when colocation is not wanted, or nothing derives from the spec,
the colocation declaration is removed.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the choice | a set of artifact names, possibly empty | names outside the offered artifacts | the init questions (they only offer the listed artifacts); unknown names are ignored here |
| the detected folders | the folder inference found for guides and plans, or none | — | inference (INPRN-B03, INPRN-B04) |
| the colocation answer | yes or no | — | the init questions |

## Effects

| Effect | Description |
| --- | --- |
| `ARCHR-B01` | The artifact options (`ArtifactNames`) are spec, feature, test, guide, plan and code (`CodeArtifact`), in that order. |
| `ARCHR-B02` | Each artifact inference found is pre-checked — code when inference found a code directory — and no other artifact is. |
| `ARCHR-B03` | Applying the choice (`ApplyArtifactChoice`), a chosen artifact whose layer does not exist gets one, even when nothing was detected, with the artifact's kind, its default pattern and its name as tag. |
| `ARCHR-B04` | An artifact that was not chosen loses its layer, and code layers are left as they are. |
| `ARCHR-B05` | A chosen artifact whose layer already exists keeps it as declared. |
| `ARCHR-B06` | The guide and plan layers take the detected folder when there is one, and the default pattern otherwise. |
| `ARCHR-B07` | Colocation (`ApplyColocation`) is declared with the spec as anchor, the chosen code, feature and test beside it, and the spec is never a derivative. |
| `ARCHR-B08` | No colocation is declared when it is not wanted, when the spec is not chosen, or when nothing derives from the spec. |
| `ARCHR-B09` | Choosing code creates no layer and leaving it out removes none: the code layers stay as inference, the preset and the user left them. |
| `ARCHR-B10` | Colocation writes only its own part of `derived` (the anchor and the file templates): the rest — the test handle the inference found — is kept, with colocation on and off. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `ARCHR-I01` | REF[ARCHR-B07]: the spec is the anchor of every colocation this unit declares, and never one of its derivatives | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `ARCHR-X01` | REF[ARCHR-B04]: code layers are outside this unit; the choice of which code layers to keep belongs to the init decisions (INDCN-B02) | — |

## Errors

none — applying a choice rewrites a configuration held in memory from values it receives; an unknown artifact name is simply not offered, and nothing is read or written outside the configuration.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Layer`, `Derived` | core — project configuration |
| DEP2 | `internal/initx/infer.go` | `Proposal` | infra — what inference found (INPRN) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
