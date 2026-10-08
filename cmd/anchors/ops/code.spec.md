<!-- @anchors
  code: CDCMC
  updated_at: 2026-10-08
  layer: comando
-->
# CodeCommand — a new unit gets an identity code that no other unit in the map already owns, and the codes in use are listed from the map

> **Code**: `CDCMC`

## Overview

Every specifiable unit carries a short, stable identity code that prefixes its rules and
crosses spec, feature and test. Two units with the same code would propagate each other's
changes, so the code must be chosen against the codes already in use, before the spec is
written. `code` answers that question from the map.

Given a unit name, it suggests a code free in the map. When the argument is a path inside a
layer that declares a module prefix, the code starts with that prefix; when it is a path
whose basename is generic, the identity comes from the parent directory; otherwise it comes
from the name. When the canonical code for the name is taken, it suggests another free code
and says who owns the canonical one. `--check` answers whether a proposed code is free,
naming the owners when it is not.

`code list` enumerates the codes the map knows, from the map's identity field rather than a
text search: a grep for a code pattern matches prose, comments and file names, and depends on
the code length each project declares. The output is one line per code with the folders that
use it, sorted, so it pipes cleanly; a JSON form gives consumers each code's file, kind,
title and work order fields. `code list --check` confronts each DECLARED code's length with
the project's accepted lengths and proposes the canonical code for each one outside them;
codes that are only cited (fixtures, examples) are counted and left alone. It changes
nothing: each proposal ends as the `anchors recode` command that applies it, because a
rename rewrites the whole project and deserves its own dry-run.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the unit name | a bare name, or a path to a unit or one of its artifacts | nothing (with no `--check`) | this unit: fails asking for the name |
| the proposed code | any text; compared in uppercase | — | this unit |
| the map | a loadable map | a missing or broken map | this unit: fails asking for `anchors map build` |
| the config | a loadable `anchors.yaml` | a missing one (for `list`) | this unit: `list` fails; `code` works without a prefix |
| the path filter | a path prefix | — | this unit: filters by prefix of the folder |

## Effects

| Effect | Description |
| --- | --- |
| `CDCMC-B01` | With nothing taken, the suggestion for a name is its canonical code. |
| `CDCMC-B02` | When the canonical code is taken, the suggestion is another free code, and the output names who owns the canonical one. |
| `CDCMC-B03` | A path inside a layer that declares a module prefix gets a code starting with that prefix, and the output says so. |
| `CDCMC-B04` | A path whose basename is generic takes its identity from the parent directory, and the output says so. |
| `CDCMC-B05` | `--check` compares the code in uppercase; a taken code names the folders that own it and fails with the collision error, a free one is reported free. |
| `CDCMC-B06` | `code list` prints one line per code, `code<TAB>folders`, sorted by code, with a code used by several files of the same folder on one line; the summary goes to standard error. |
| `CDCMC-B07` | `--in` keeps only the folders under the given prefix; when nothing matches, the output names the filter instead of claiming the project has no code. |
| `CDCMC-B08` | `--json` prints each code with its folders, file, kind, title, needs, parent and revises; an empty list is `[]`. |
| `CDCMC-B09` | A unit's title is its first `# ` heading in a markdown file, with the text up to the dash dropped; a file that is not markdown has no title. |
| `CDCMC-B10` | `list --check` accuses only DECLARED codes whose length is outside the accepted lengths, proposing the canonical code of the unit's name, resolved against the conforming codes; a cited code is counted and not accused; the check fails when anything diverges and passes with a summary otherwise. |
| `CDCMC-B11` | The unit name of a path drops the `.spec.md`, `.feature` and `.test.<ext>` suffixes, or the plain extension. |
| `CDCMC-B12` | A map with no identity node makes `list` say so and point at `anchors map build`. |
| `CDCMC-B13` | Every run reads its codes, owning files and declared identities from its own map: a second run in the same process never inherits what a previous map declared. |
| `CDCMC-B14` | When the length check accuses codes, it ends with the exact `anchors recode <old> <new>` command for each one, to be reviewed in its dry-run and applied with `--apply`; the check itself changes nothing and has no `--fix`. |
| `CDCMC-B15` | A name shaped like an identity code — capitals and digits at a length `code_lengths` allows — still gets its suggestion, and the command also says whether that code is free or who uses it, pointing at `--check`; a name with lower case gets no such note. (`codeShapeRE`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CDCMC-I01` | REF[CDCMC-B02]: a suggested code is never one the map already holds — the taken canonical is always adjusted | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CDCMC-X01` | The codes come only from the map's identity field; no file is searched for a code pattern. | A text search matches mentions, comments and file names, and depends on the project's code length. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CDCMC-E01` | `code` runs with no name and no `--check`, or the map cannot be read. | It fails asking for the unit name, or asking to run `anchors map build`. | A code suggested without the map could collide with any unit. |
| `CDCMC-E02` | `code list` runs where `anchors.yaml` or the map cannot be loaded. | It fails naming which one. | The length check and the listing need both; an empty answer would read as "no code in use". |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
