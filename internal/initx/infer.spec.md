<!-- @anchors
  code: INPRN
  updated_at: 2026-09-26
  layer: infra
-->
# InferProposal — walks the project and proposes its structure deterministically, for init to confirm

> **Code**: `INPRN`

## Overview

`anchors init` does not start from a blank form: it walks the project and proposes a structure, which the
interactive questions then confirm or adjust. The inference is structural and deterministic; no model is
involved.

The walk skips the directories that hold dependencies, build output and tool state, because their files are
not the project's. Each remaining file is classified once, in a fixed order: specs, features and tests by
their name; plans by living in a plans folder, before guides, so a plan kept under a guides folder is still a
plan; guides by living in a guides folder; code by its extension. From that classification the proposal
records which artifacts exist, where plans and guides live, which top directories hold enough code to be a
layer, the most frequent code extensions, whether the triad sits beside the code, and which attribute the
project uses to mark elements for tests.

A code directory is a top directory of up to two path segments that holds at least ten code files; fewer
is a folder, not a layer. Colocation is detected only when at least three stems pair code with a spec, a test
or a feature in the same place: a pair or two is coincidence. The test handle is detected by use, not by
dependency: among the known marking attributes, the one used most in a sample of the first code files wins,
and only from five uses; below that the presence is anecdotal and no contract is proposed. The sample is
bounded because the attribute, when a project uses it, is everywhere; reading the whole repository would
cost I/O for a signal the sample already gives.

The proposal carries the configuration built from it (BLCNB).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory that can be walked | a path that does not exist or cannot be read | this unit: the walk failure is returned (INPRN-E01) |
| the files found | any file | files under dependency, build or tool directories | this unit: those directories are not walked (INPRN-B01) |

## Effects

| Effect | Description |
| --- | --- |
| `INPRN-B01` | Dependency, build and tool directories are not walked: node_modules, .git, dist, build, vendor, .next, coverage, .expo and .anchors. |
| `INPRN-B02` | The presence of specs, features and tests is detected from their file names. |
| `INPRN-B03` | A markdown file in a plans folder is a plan, even inside a guides folder, and the first plans folder found is recorded. |
| `INPRN-B04` | The first guides folder found is recorded, with every guide file found, sorted. |
| `INPRN-B05` | A code directory is a top directory of up to two segments holding at least ten code files, and code directories are ordered by volume, ties by name. |
| `INPRN-B06` | The code extensions are the five most frequent, most frequent first, ties by name. |
| `INPRN-B07` | Colocation is detected when at least three stems pair code with a spec, test or feature in the same place. |
| `INPRN-B08` | The test handle is the known marking attribute used most, only from five uses; on a tie, the earlier attribute of the known list wins. |
| `INPRN-B09` | The proposal carries the configuration built from it, with code layers and no artifact layer. |
| `INPRN-B10` | A test named in any dialect inference knows (`.test.ts`, `.test.tsx`, `.test.js`, `.test.go`, `_test.go`, `.test.py`, `_test.py`, `.spec.ts`) has the stem of the code it tests, so `foo_test.go` pairs with `foo.go` for colocation. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INPRN-I01` | REF[INPRN-B01]: nothing under a skipped directory ever reaches the proposal, whatever it holds | — |
| `INPRN-I02` | The same tree always gives the same code extensions, code directories and layer patterns, ties included. | infers a tree whose extensions and directories all tie, twenty times, and compares each proposal with the name order |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INPRN-X01` | Only the first three hundred code files are read in search of the test handle. | The attribute, when used, appears early and everywhere; a full read would make init pay I/O for a signal the sample already gives. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `INPRN-E01` | The project root cannot be walked (it does not exist or cannot be read). | The inference fails with the walk error and no proposal. | A proposal built from a partial or missing walk would present an empty project as the truth and seed a configuration from nothing. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — the proposed configuration |
| DEP2 | `internal/initx/build_config.go` | `buildConfig` | infra — the proposal's configuration (BLCNB) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
