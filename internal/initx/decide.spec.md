<!-- @anchors
  code: INDCN
  updated_at: 2026-10-08
  layer: infra
-->
# InitDecisions — the pure decisions of init over the proposed configuration: code layers, tags and governs rules

> **Code**: `INDCN`

## Overview

The interactive init collects answers; the decisions those answers lead to live here, as pure functions
that take the proposal and the answers and return the adjusted configuration. Keeping them apart from the
prompts is what lets them be tested without a terminal.

The user is shown the proposed code layers, sorted, and chooses which to keep; pruning removes the code
layers not kept and never touches an artifact layer. The tags a guide may govern are the tags declared on
any layer, each once, sorted. Finally, each guide the user answered with a tag becomes one governs rule; a
guide left unanswered, or answered with the explicit "none" option, produces no rule. The rules come out
ordered by guide, so the written configuration does not depend on the order the answers were collected in.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | the proposed configuration, with code and artifact layers | — | inference (BLCNB) and the artifact choice (ARCHR) |
| the kept layers | a set of layer names, possibly with names that are not code layers | — | this unit: only code layers are considered |
| the governs answers | a map from guide path to a tag, the empty answer or the "none" option | — | this unit: empty and "none" are skipped |

## Effects

| Effect | Description |
| --- | --- |
| `INDCN-B01` | The code layer names (`CodeLayerNames`) are listed sorted, and only layers of kind code are listed. |
| `INDCN-B02` | Pruning (`PruneCodeLayers`) removes the code layers that were not kept and never removes an artifact layer. |
| `INDCN-B03` | The candidate tags (`Tags`) are the union of the tags of every layer, each once, sorted. |
| `INDCN-B04` | Each guide answered with a tag becomes one governs rule (`BuildGovernRules`); guides left unanswered or answered with "none" produce none; the rules are ordered by guide. |
| `INDCN-B05` | No answer produces no governs rule. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INDCN-I01` | After pruning, the listed code layers are exactly the code layers that were kept. | prunes with a keep set that also names an artifact and an unknown layer, then lists the code layers |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INDCN-X01` | REF[INDCN-B02]: this unit never removes an artifact layer; artifact layers belong to the artifact choice (ARCHR-B04) | — |

## Errors

none — every function works on values in memory and returns a value; an unknown kept name or an unanswered guide is a normal input that is skipped, not a failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
