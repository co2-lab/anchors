<!-- @anchors
  code: INCTN
  updated_at: 2026-10-01
  layer: infra
-->
# InitCatalogs — the language dialect catalog and the @TBD instruction that init seeds into judgment gates

> **Code**: `INCTN`

## Overview

`anchors init` works from two fixed catalogs that are data rather than decisions: the dialect of each
language it reads, and the instruction text it appends to the question of a judgment gate that asks about a
piece of the triad. This spec states what those catalogs must hold, and how they are read.

The dialect catalog holds what a language decides that the map needs: how a test file is named, and which
family a manifest or a code extension belongs to. It holds no structure. Where a project keeps each layer is
the project's decision, and the catalog never names a folder; a project structure per stack, written into
the configuration, pushed projects to move files to fit Anchors. Each test convention is a prefix and a
suffix: its glob matches every test file of the project, and its template puts a unit's test beside the
unit's code. A family's default convention is used only when the project has no test file to read one from.

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
| a file name | any base name | — | this unit: a name no convention matches is not a test |
| the piece named by the instruction | the piece the judgment gate asks about, as the question already names it | — | the caller (DFGTD-B09): each judgment gate passes its own piece |

## Effects

| Effect | Description |
| --- | --- |
| `INCTN-B04` | The @TBD instruction, in English like the question it closes, tells the judge not to answer `pass`, to answer `waived` naming the absence, and it names the piece the gate asks about. |
| `INCTN-B06` | A file is a test (`testConventionOf`) when its name starts with a convention's prefix and ends with its suffix and has more than both; it follows the first such convention of the catalog. |
| `INCTN-B07` | A convention's glob (`Glob`) is `**/` with its prefix, `*` and its suffix; its template (`Template`) is `{{dir}}/` with its prefix, `{{name}}` and its suffix. |
| `INCTN-B08` | A test file's unit name is its name without the convention's prefix and suffix. |
| `INCTN-B10` | A family's coverage hint (`CoverageHint`) says how its test runner emits the JUnit and lcov reports `anchors ingest` reads; a family with no hint, or no family, gets none. |
| `INCTN-B09` | With no test convention read from the project, the conventions are its family's default; with neither, there are none and no template. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INCTN-I02` | The @TBD instruction always demands checking that the declaration is still true, so a stale declaration is reported rather than obeyed. | reads the instruction and confirms it covers the stale declaration |
| `INCTN-I03` | No convention of the catalog is shadowed: a convention whose prefix and suffix also end another listed later is never reached first by a shorter one, so `.spec.tsx` is read as `.spec.tsx`, not `.tsx`. | walks every pair of conventions and checks that a longer form comes before a shorter form it contains |
| `INCTN-I04` | Every family default is one of the catalog's conventions, of that family. | walks the defaults and finds each in the catalog |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INCTN-X02` | The catalog names no folder: no convention, default or manifest carries a directory. | Where a project keeps each layer is the project's decision; a catalog that names folders is a proposed structure, and a proposed structure pushed a project to move files to fit Anchors. |

## Errors

none — the catalogs are fixed data; a name no convention matches is the normal answer (INCTN-B06), not a failure.

## Dependencies

none

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
