<!-- @anchors
  code: INCTN
  updated_at: 2026-09-26
  layer: infra
-->
# InitCatalogs — the stack preset catalog and the @TBD instruction that init seeds into judgment gates

> **Code**: `INCTN`

## Overview

`anchors init` works from two fixed catalogs that are data rather than decisions: the stack presets it
offers, and the instruction text it appends to the question of a judgment gate that asks about a piece of
the triad. This spec states what those catalogs must hold, and how they are read.

The preset catalog lists the most established project structure of each stack: a unique short name, a human
title, and the code and test layers of that structure. Artifact layers (spec, feature, guide, plan) are not
part of a preset, because the artifacts are chosen apart from the stack; the test layer is the exception,
because the naming convention of tests belongs to the stack. A preset layer that declares no kind is a code
layer. A modular preset says where its modules live, and no preset layer carries an identity prefix: the
prefix is deduced per module of the real project at init time (APPRP-B03), since the modules belong to the
project and not to the stack. Presets are looked up by name, and listed in catalog order.

The @TBD instruction exists because a spec is born before its code. When the spec declares that a piece is
still to be written, a judgment gate asking whether that piece realises the rule has nothing to confront, and
the easy answer is `pass`, which leaves a stamp in the map that looks like a real verification. The
instruction tells the judge to check for the declaration first, to answer with a waiver naming the absence
instead of `pass`, and to check that the declaration is still true, because a stale declaration makes every
gate that reads it waive what it should charge. The instruction names the piece the gate asks about, so it
reads as the continuation of the question.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the preset name | any string | — | this unit: an unknown name finds nothing |
| the piece named by the instruction | the piece the judgment gate asks about, as the question already names it | — | the caller (DFGTD-B09): each judgment gate passes its own piece |

## Effects

| Effect | Description |
| --- | --- |
| `INCTN-B01` | A preset layer with no declared kind becomes a code layer; a declared kind, pattern and tags are kept. |
| `INCTN-B02` | Looking up a preset (`PresetByName`) by a name the catalog does not hold finds nothing and returns an empty preset. |
| `INCTN-B03` | The preset names (`PresetNames`) are listed in the order of the catalog. |
| `INCTN-B04` | The @TBD instruction, in English like the question it closes, tells the judge not to answer `pass`, to answer `waived` naming the absence, and it names the piece the gate asks about. |
| `INCTN-B05` | A modular preset declares the directory where its modules live. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INCTN-I01` | Every preset has a unique name, a title, at least one layer, a pattern on every layer and a test layer. | walks the whole catalog and checks each preset |
| `INCTN-I02` | The @TBD instruction always demands checking that the declaration is still true, so a stale declaration is reported rather than obeyed. | reads the instruction and confirms it covers the stale declaration |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INCTN-X01` | No preset layer carries an identity prefix. | The prefix belongs to a module of the real project and is deduced at init; a prefix fixed in the catalog would be the same for every project of the stack. |

## Errors

none — the catalogs are fixed data; the only miss, an unknown preset name, is the normal answer of the lookup (INCTN-B02) and not a failure.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Layer` | core — the layer a preset layer becomes |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
