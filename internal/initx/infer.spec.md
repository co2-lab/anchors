<!-- @anchors
  code: INPRN
  updated_at: 2026-10-01
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
records which artifacts exist, where plans and guides live, which folders hold code, the most frequent code
extensions, whether the unit sits beside the code, which attribute the project uses to mark elements for
tests, the project's language family and how its test files are named.

Inference proposes no structure. Every folder that holds code is a candidate layer, because where the
project keeps each layer is the project's decision: the candidate is the code file's folder, cut at two path
segments (`src/handlers`, `apps/mobile`), or the root for a file at the root. There is no minimum count; a
threshold dropped the small layers of a layered project, and init lets the user keep or drop each candidate.

The language is dialect, not structure: it decides how a test file is named and how a comment is written.
The family is the one the manifest at the root declares, or the family of the most frequent code extension
when there is no manifest. The test conventions are read from the project's own test files, counted by the
naming form each follows, so the patterns proposed are the ones the project already uses.

 Colocation is detected only when at least three stems pair code with a spec, a test
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
| `INPRN-B05` | Every folder holding code is a candidate code directory — the code file's folder cut at two segments, or the root for a file at the root — with no minimum count, ordered by volume, ties by name. |
| `INPRN-B06` | The code extensions are the five most frequent, most frequent first, ties by name. |
| `INPRN-B07` | Colocation is detected when at least three stems pair code with a spec, test or feature in the same place. |
| `INPRN-B08` | The test handle is the known marking attribute used most, only from five uses; on a tie, the earlier attribute of the known list wins. |
| `INPRN-B09` | The proposal carries the configuration built from it, with code layers and no artifact layer. |
| `INPRN-B10` | A test named by any convention inference knows (`_test.go`; `.test` and `.spec` with `.ts`, `.tsx`, `.js`, `.jsx`; `_test.py` and `test_*.py`; `_spec.rb` and `_test.rb`; `Test.java`, `Tests.java`, `Test.kt`, `Tests.cs`, `Test.cs`, `Test.php`) has the stem of the code it tests, so `foo_test.go` pairs with `foo.go` and `test_foo.py` with `foo.py` for colocation. |
| `INPRN-B11` | The language family is the one a manifest at the root declares — `go.mod`, `Cargo.toml`, `pyproject.toml`, `setup.py`, `requirements.txt`, `Gemfile`, `composer.json`, `package.json`, checked in that order — or, with no manifest, the family of the most frequent code extension; with neither, there is no family. |
| `INPRN-B12` | The test conventions are the naming forms the project's test files follow, most followed first, ties by name; a file follows the longest form it ends with (`.spec.tsx` before `.tsx`). |

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
