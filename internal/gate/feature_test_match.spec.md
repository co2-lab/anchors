<!-- @anchors
  code: FTMFT
  updated_at: 2026-09-19
  layer: gate
-->
# FeatureTestMatch — scenarios in feature must be implemented in test by code and description

> **Code**: `FTMFT`

## Overview

Confronts the relational edge between a feature file and the test file that executes it: **every scenario
catalogued in the feature must be implemented in the linked test.**

The gate enforces a critical **TWO-TIER VERDICT**:
1. **BY CODE (Defect / Failure)**: Every scenario code (`XXXXX-Y##`) must appear in the non-comment body
   of at least one linked test. A missing scenario code is a DEFECT indicating that an agent or developer
   skipped the scenario entirely or renamed the code without updating the test. When code is missing,
   the gate returns a blocking **Fail**.
2. **BY DESCRIPTION (Signal / Warning)**: When the scenario code exists in the test, the gate checks whether
   the test description corresponds to the scenario title. Free-form text naturally varies between human
   domain language in Gherkin ("when description and value are identical") and code identifiers in tests
   (`classifica`, `expect`). Therefore, descriptive drift is treated as an informative SIGNAL rather than
   a hard failure: the gate issues a non-blocking **Pending** warning so teams can realign descriptions
   without blocking delivery of verified code.

Neighbouring gates govern other dimensions of the triad: `triad-complete` checks the physical presence
of the triad artifacts; `spec-feature-match` ensures requirements defined in specs are covered by feature
scenarios. `feature-test-match` specifically guarantees that documented scenarios are faithfully realized
in automated test suites.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a node of kind `feature` | any node that is not a feature (returns Skip) | the gate engine, routing by the declared `on:` |
| the graph | a built dependency graph | a nil graph (returns Pending) | this unit: with no map the confrontation cannot resolve linked tests |
| the linked tests | test files reached via `tested-by` edges | features with no linked tests (returns Pending) | the map and triad discovery |
| the test surface regimes | scenarios mapped to `unit` or `integration` surfaces | scenarios scoped strictly to `e2e` or `vr` surfaces | this unit, resolving regime mappings from configuration |
| the test content | test source code stripped of comments for code presence, and with comments for description match | code in external packages or unrelated files | this unit, reading and parsing linked test files |

## Effects

| Effect | Description |
| --- | --- |
| `FTMFT-B01` | An artifact that is not a feature file skips confrontation without a verdict. |
| `FTMFT-B02` | When the dependency graph is nil, the gate returns Pending without approving. |
| `FTMFT-B03` | A feature with no scenarios declared skips confrontation. |
| `FTMFT-B04` | When no tests are linked by `tested-by` edges, the gate returns Pending. |
| `FTMFT-B05` | Scenarios scoped to non-test surfaces (such as `e2e` or `vr`) are skipped by this gate. |
| `FTMFT-B06` | A scenario whose code is completely absent from the test body fails. |
| `FTMFT-B07` | A scenario code appearing only within comments does not count as implemented and fails. |
| `FTMFT-B08` | When scenario codes match and titles are identical, the gate passes. |
| `FTMFT-B09` | When scenario codes match but test descriptions drift, the gate issues a Pending warning. |
| `FTMFT-B10` | Test titles containing quotes or nested delimiters are parsed completely without truncation. |
| `FTMFT-B11` | Multiple sibling scenario codes cited in a single test title are extracted and attributed cleanly. |
| `FTMFT-B12` | A test title shared among sibling scenarios does not require strict single-title equality. |
| `FTMFT-B13` | Test comments count towards descriptive coverage while remaining excluded from code presence. |
| `FTMFT-B14` | Scenario codes match on exact boundaries so prefix substrings of other codes do not collide. |
| `FTMFT-B15` | Go `t.Run` test bindings are recognized alongside standard test runner declarations. |
| `FTMFT-B16` | Line comments using `#` and `--` syntax are stripped when evaluating code presence. |
| `FTMFT-B17` | The exported function `RootCode` strips scenario sub-indices (`#01`) and returns the root requirement code. |
| `FTMFT-B18` | Only a regime tag the project MAPS under `regimes:` (or a canonical regime name) exempts a scenario as belonging to another surface; an unmapped tag that merely looks like a regime (`@nivel-compilacao`) leaves the scenario confronted. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FTMFT-I01` | Absence of scenario code in test implementation is always a Failure. | confronts a test missing a scenario code and verifies it returns Fail |
| `FTMFT-I02` | Descriptive divergence with valid code presence is always a Warning (Pending), never a Failure. | confronts a test with matching code but divergent description and verifies it returns Pending |
| `FTMFT-I03` | Comments are strictly separated: excluded when verifying code implementation, included when verifying semantic description. | confronts a scenario code present only in comments and verifies it fails |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FTMFT-X01` | Does not execute test suites or inspect test execution results. | This gate performs static structural and textual analysis; test execution signals are handled by test runners and ingested via `anchors ingest`. |
| `FTMFT-X02` | Does not confront E2E or visual regression scenarios. | Non-unit testing surfaces inhabit different directories and runners, governed by dedicated surfaces. |
| `FTMFT-X03` | Does not fail tests for minor natural language variations in test titles. | Human scenario language and implementation code descriptions naturally vary; classifying drift as warning prevents disruption while preserving traceability. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern`, `Config`, `DefaultRuleLetters` | core — configuration models and vocabulary patterns |
| DEP2 | `internal/config/dialect.go` | `GherkinScenarioAlternatives` | core — multilingual Gherkin scenario keywords |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — localized messages for failures and warnings |
| DEP4 | `internal/mapx/model.go` | `EdgeTestedBy`, `Graph`, `KindFeature`, `Node` | core — graph model and test relationship edges |
| DEP5 | `internal/similarity/similarity.go` | `Classify`, `Identico`, `Weights` | core — weighted term similarity analysis for title matching |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
