<!-- @anchors
  code: PRGTP
  updated_at: 2026-09-19
  layer: gate
-->
# PromotableGates — identifies clean informative gates ready for promotion to blocking

> **Code**: `PRGTP`

## Overview

Identifies informative gates whose evaluation is completely clean in a given execution profile, making them candidates for promotion to blocking gates (QUALITY section 7). When a gate is newly introduced in a project, it starts as informative so that existing debt does not block work. As the project matures and satisfies the gate's conditions, the gate remains informative unless promoted in `anchors.yaml`, measuring without actively defending the codebase. Promotion remains a human decision; this unit discovers clean gates and surfaces them as reminders directly in workflow commands (`check`, `status`, `next`).

## Signature

| Parameter | Type | Description |
| --- | --- | --- |
| `p` | `Profile` | Evaluation profile aggregating results and gate summaries across the project |

**Returns**: `[]Promotable` containing each eligible gate name and its approved node count, sorted deterministically by gate name.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `p` | any `Profile` struct, including empty gate summaries | nil profile references | Go value receiver semantics, which ensures a valid `Profile` struct is passed |

## Effects

| Effect | Description |
| --- | --- |
| `PRGTP-B01` | Returns an empty list of `Promotable` gates when the evaluation profile contains no gates. |
| `PRGTP-B02` | An informative gate with one or more passes and zero failures is included in the promotable list. |
| `PRGTP-B03` | An informative gate with one or more failures is excluded from promotion. |
| `PRGTP-B04` | An informative gate with zero passes and zero failures is excluded from promotion as having no measurement data. |
| `PRGTP-B05` | A blocking gate is excluded from promotion suggestions because it already actively defends. |
| `PRGTP-B06` | Each returned `Promotable` struct populates `Gate` with the gate identifier declared in `anchors.yaml`. |
| `PRGTP-B07` | Each returned `Promotable` struct records `Passou` equal to the count of nodes passed by that gate. |
| `PRGTP-B08` | `PromotableGates` returns promotable candidates sorted deterministically in alphabetical order of gate names. |
| `PRGTP-B09` | Multiple clean informative gates in the profile are all collected into the returned promotable list. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRGTP-I01` | An informative gate is candidate for promotion if and only if it is non-blocking, has zero failures, and has at least one pass. | evaluates gates across all combinations of blocking, pass count, and failure count, verifying exact matching |
| `PRGTP-I02` | A gate with zero passes is never classified as clean because lack of measured data must not simulate compliance. | evaluates an informative gate with skips or pendings but zero passes and verifies exclusion |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRGTP-X01` | Does not automatically promote gates or modify `anchors.yaml`. | Promotion to blocking status alters build and commit barriers, requiring deliberate human decision and team agreement. |
| `PRGTP-X02` | Does not evaluate gate execution results directly from disk or runner outputs. | Operates purely on the aggregated `Profile` summary computed by the gate runner, maintaining separation between aggregation and discovery. |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
