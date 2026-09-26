<!-- @anchors
  code: TSTRT
  updated_at: 2026-09-26
  layer: gate
-->
# TestTraceable — a test linked to a feature must declare what scenario it proves

> **Code**: `TSTRT`

## Overview

Confronts a test artifact against its linked feature: **a test connected to a feature must state
what scenario it proves.**

The blind spot this gate closes is subtle and occurred twice in a real production project: the test
EXISTS, PASSES, and covers the correct behavior — yet cites no scenario identity code whatsoever.
To every relational gate, such a test is invisible:
- `feature-test-match` checks scenario-to-test alignment by scenario code: without a code in the test,
  it reports the scenarios as unimplemented, mistakenly blaming the feature rather than identifying
  the untraceable test;
- `triad-complete` verifies file presence on disk and marks the test piece satisfied;
- `tests-green` executes the test suite and verifies assertions pass.

The outcome is the worst of both worlds: the engineering work was completed and verified, yet the
pipeline claims it was not. When another developer attempts to fix the reported deficit, they write
a redundant duplicate test for the exact same behavior, having no automated way to discover that the
existing test already exercised it.

Measured in doctrine: in a real project, a shared test exercised across 6 screens used inverted digit
and letter identifiers (such as transposed tokens where specs declared proper codes). Three screens
had tests and none was tracked because code search could not reach them. In the same project, a store
test covered 5 distinct behaviors without citing a single scenario code.

What separates this gate from neighbouring gates:
- `feature-test-match`: enforces granular 1:1 matching between each scenario and a dedicated test case.
  In contrast, `test-traceable` applies the weakest possible threshold: finding at least ONE valid
  scenario code anywhere in the test file satisfies the gate. The question here is simply: "does this
  test declare itself?"
- `triad-complete`: checks only that the four triad files exist on disk, regardless of whether their
  internal contents establish traceability.
- `tests-green`: compiles and executes test suites to ensure zero failures, but remains completely
  agnostic to requirement traceability codes.

Finally, tests without a linked feature (such as unit tests for internal utilities) are not charged:
demanding scenario codes from unlinked tests would require referencing nonexistent features.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind `test` | nodes of kind `spec`, `code`, or `feature` | this unit: skips confrontation on non-test nodes |
| the dependency graph | a populated graph structure containing nodes and edges | nil graph pointer | this unit: returns Pending when the graph is nil |
| the linked feature | a feature node connected to the test via an incoming `tested-by` edge | tests with no linked feature in the graph | this unit: skips confrontation when no linked feature exists |
| the scenario codes | scenario identity codes extracted from the linked feature file | feature files declaring no scenario codes or unreadable feature paths | this unit: skips confrontation when no scenario codes are found |
| the test content | source text of the test file | non-text binary assets | the gate engine, routing text artifacts |

## Effects

| Effect | Description |
| --- | --- |
| `TSTRT-B01` | When the confronted node is not of kind test, the gate skips confrontation. |
| `TSTRT-B02` | When the dependency graph is nil, the gate returns Pending with a missing map verdict. |
| `TSTRT-B03` | When a test node has no incoming `tested-by` edge from a feature in the graph, the gate skips confrontation. |
| `TSTRT-B04` | When the linked feature file cannot be read from disk, the gate skips confrontation. |
| `TSTRT-B05` | When the linked feature file declares no scenario codes, the gate skips confrontation. |
| `TSTRT-B06` | When the test content contains at least one scenario code declared by its linked feature, the gate passes. |
| `TSTRT-B07` | A single declared scenario code in the test content is sufficient to pass, even when the linked feature defines multiple scenarios. |
| `TSTRT-B08` | When the test content contains none of the scenario codes declared by its linked feature, the gate fails. |
| `TSTRT-B09` | A test covering behavior with transposed or misspelled codes fails exact substring confrontation. |
| `TSTRT-B10` | When failing, the verdict names the linked feature path, lists up to three expected codes, and suggests the first code as a fix hint. |
| `TSTRT-B11` | Scenario codes occurring anywhere in the test file content (such as in comments or test titles) satisfy the traceability check. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TSTRT-I01` | Traceability is strictly a test requirement; non-test artifacts skip confrontation to prevent misattributing test debt to other files. | confronts non-test nodes and verifies the verdict is Skip |
| `TSTRT-I02` | Tests without a linked feature in the graph are never charged, avoiding demands for references to nothing. | confronts an unlinked test node and verifies it skips confrontation |
| `TSTRT-I03` | The acceptance threshold requires only one scenario code present in the file, preventing duplicate error reporting with scenario-matching gates. | confronts a test matching one of multiple scenarios and verifies it passes |
| `TSTRT-I04` | Confrontation without a dependency graph returns Pending rather than approving an unmeasured relationship. | confronts a test with a nil graph pointer and verifies the verdict is Pending |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TSTRT-X01` | Does not enforce a 1:1 match between each scenario and a separate test case. | Scenario-by-scenario coverage is the distinct duty of `feature-test-match`; duplicating it here would penalize the same debt twice in the pipeline. |
| `TSTRT-X02` | Does not execute test suites or inspect runtime assertion results. | Runtime execution and green status are governed exclusively by `tests-green`. |
| `TSTRT-X03` | Does not require scenario codes in standalone test files unattached to features. | Helper tests and unit tests for internal tools without business features have no specification or scenarios to cite. |
| `TSTRT-X04` | Does not inspect the semantic validity or quality of test assertions. | Detecting at least one scenario code establishes relational visibility; semantic depth is evaluated by code review and test execution. |

## Errors

Each failure the code handles is already stated as a rule of another letter; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TSTRT-E01` | REF[TSTRT-B02]: with no map the link to the feature cannot be followed, and B02 answers Pending | — | — |
| `TSTRT-E02` | REF[TSTRT-B04]: a linked feature that cannot be read is answered by B04: the confrontation is skipped | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `EdgeTestedBy`, `Graph`, `KindTest`, `Node` | core — graph structure, test nodes, and tested-by edges |
| DEP4 | `internal/scan/scan.go` | `ScenarioCodeRE` | core — scenario code regular expression grammar |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
