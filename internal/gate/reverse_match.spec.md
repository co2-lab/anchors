<!-- @anchors
  code: RVMTR
  updated_at: 2026-09-27
  layer: gate
-->
# ReverseMatch — every scenario still has its rule, and every proven code still has its scenario

> **Code**: `RVMTR`

## Overview

The forward gates of the triad walk from the origin to the destination: every rule needs a scenario,
every scenario needs a test. Neither walks back from the destination to ask whether the origin still
exists. The cost was measured: a revert deleted a rule from the spec, the code and the test, while the
feature kept its scenario (a parallel change reintroduced it with no conflict). The scenario went on
asserting a behaviour nobody decided any more, the test stayed green proving a rule nobody declared,
and every gate was green.

This unit holds the return leg of each pair, as two gates of their own, so that whoever reads the
accusation knows it is about the destination:

- feature-spec-match confronts a feature: does every scenario code of this unit correspond to a rule
  that a spec covering this feature still defines?
- test-feature-match confronts a test: does every code the test names correspond to a scenario that a
  feature it exercises still declares?

Both were tuned against false accusations measured in a real project, because a gate that accuses
what is right teaches people to ignore it: numbered variants of one rule, data states defined by name,
the visual baseline, revision codes, and other units' codes cited to build fixtures are not orphans.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted node | a feature for feature-spec-match, a test for test-feature-match | any other kind | this unit: each gate skips the other kinds |
| the map | covering edges from specs to the feature, testing edges from features to the test | no map | this unit: both gates answer Pending |
| the linked origins | spec and feature files on disk, relative to the project root | a file that cannot be read | this unit: it contributes nothing |

## Effects

### feature-spec-match

| Effect | Description |
| --- | --- |
| `RVMTR-B01` | The gate skips a node that is not a feature, and a feature that declares no coded scenario. |
| `RVMTR-B02` | Without a map, both gates answer Pending. |
| `RVMTR-B03` | A feature no spec covers is Pending: the existence of the spec is charged by co-location, not here. |
| `RVMTR-B04` | When the covering specs define no requirement and no data state, the gate is Pending. |
| `RVMTR-B05` | A scenario code of this unit that no covering spec defines fails the gate, and the message names the code as written in the feature; when every code has its rule, the gate passes. |
| `RVMTR-B06` | A numbered variant of a scenario code is the same rule as the code without it; a variant of a rule the spec does not define is still an orphan. |
| `RVMTR-B07` | A data state cited in the feature counts as defined when a covering spec defines a state of that name, with or without the unit prefix; a state no spec defines is an orphan. |
| `RVMTR-B08` | The unit's visual baseline code is never charged here. |
| `RVMTR-B09` | The rules of every spec covering the feature count together: a scenario defined by any of them has an owner. |

### test-feature-match

| Effect | Description |
| --- | --- |
| `RVMTR-B10` | The gate skips a node that is not a test, and a test no feature exercises is Pending. |
| `RVMTR-B11` | When the features it exercises declare no coded scenario, the gate is Pending. |
| `RVMTR-B12` | A code the test names that no exercised feature declares as a scenario fails the gate, naming the code. |
| `RVMTR-B13` | A scenario declared as a numbered variant makes its rule a declared scenario for the test. |
| `RVMTR-B14` | A revision code named by the test is not read as a rule and is not charged. |
| `RVMTR-B15` | A data state a spec defines with the unit prefix is read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, `TREXXXX-DS-data-present` defines `DS-data-present`. |
| `RVMTR-B16` | A support file is not confronted as a test (Skip, saying why): it proves no scenario, so there is no feature to link it to. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RVMTR-I01` | Neither gate answers Pass when it had nothing to match against: no map, no linked origin, or an origin that declares nothing answers Pending. | runs each gate with no map, with no edge, and with an origin that declares nothing, and asserts Pending each time |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RVMTR-X01` | Codes of another unit are never charged: only this unit's codes, and in a test only the units the exercised features govern. | A scenario cites another unit's rule to say what it runs against, and a test cites one to build a fixture; charging them asks the local spec to declare someone else's rule. |
| `RVMTR-X02` | A code the test names only in a comment is not a claim of proof. | A comment is a reference, the same ruler the forward gate applies. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RVMTR-E01` | A linked spec or feature cannot be read. | It contributes no rule and no scenario; when none can be read, the gate is Pending. | An unread origin would otherwise make every code look orphaned, accusing what was never measured. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/spec_feature_match.go` | `definedRequirements` | gate — the rules a spec defines |
| DEP2 | `internal/gate/feature_test_match.go` | `parseFeatureScenarios`, `stripLineComments` | gate — the scenarios a feature declares, and a body with its comments out |
| DEP3 | `internal/gate/internal_checks.go` | `anyCodeRE` | gate — the codes a test body names |
| DEP4 | `internal/mapx/model.go` | `Graph`, `EdgeCoveredBy`, `EdgeTestedBy` | mapa — which specs cover a feature and which features a test exercises |
| DEP5 | `internal/i18n/i18n.go` | `T` | apoio — the localized messages |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
