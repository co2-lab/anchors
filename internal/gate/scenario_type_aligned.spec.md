<!-- @anchors
  code: STASC
  updated_at: 2026-09-19
  layer: gate
-->
# ScenarioTypeAligned — scenario classification tags must match the code nature letter

> **Code**: `STASC`

## Overview

Confronts a scenario's classification tag against the nature letter embedded in its identity code:
**the classification tag of a scenario must agree with the rule letter of its code.**

The nature of a rule is declared twice: once in the letter embedded in its identity code (such as `S`
for State or `B` for Behavior) and once in the scenario's classification tag (such as `@state` or
`@behavior`). When the two disagree, one of them lies — and no other gate catches the contradiction
because each inspects only a single facet:
- `rule-types` validates that specification sections conform to the declared rule vocabulary;
- `feature-test-match` pairs scenario identity codes and descriptions against test cases.
Neither gate cross-references the scenario's Gherkin classification tag against the code's embedded letter.

Measured in doctrine: in the project that originated this gate, 46 scenarios exhibited tag and letter
disagreements upon audit. Distinguishing between the two failure modes is essential:
- 36 were wrong tags: for instance, "Intro present is rendered in italics" had the body `When the
  component is rendered / Then I should see the intro text` — pure state without user action, under a
  correct `-S` code, but tagged `@behavior`. The remediation was correcting the tag.
- 10 were borrowed codes: for instance, "Tapping a chip selects priority" was genuine interactive
  behavior, attached to a `-S` code because the specification had not catalogued the event. The
  remediation was creating a dedicated code in the events table.

This gate detects both cases. Deciding whether the tag was misapplied or the code was borrowed requires
reading scenario steps and understanding domain context. The gate guarantees that the divergence does
not remain silent.

The gate issues a **Pending** verdict rather than a blocking Fail: type mismatch represents inherited
technical debt in codebases adopting this gate after scenarios were already written, and resolving each
finding requires case-by-case editorial judgment.

Furthermore, if the project configuration defines no `tags:` mappings in its rule types, the gate skips
silently. Guessing tag meanings across different project languages would introduce false alarms; silence
here respects projects that did not request this check.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind `feature` | nodes of kind `spec`, `code`, or `test` | this unit: skips confrontation on non-feature nodes |
| the project configuration | a configuration defining rule types and tag mappings | nil configuration or empty rule types | this unit: skips confrontation when configuration or rule types are missing |
| the tag vocabulary | tag-to-letter mappings declared under configured rule types | configuration where no rule type defines tag mappings | this unit: skips confrontation when tag mappings are absent |
| the feature content | Gherkin feature text containing scenarios with identity codes | features containing no scenarios with identity codes | this unit: skips confrontation when coded scenarios are absent |
| the scenario tags | classification tags matching configured vocabulary | non-classification tags such as test levels or priority markers | this unit: ignores tags not registered in configuration |

## Effects

| Effect | Description |
| --- | --- |
| `STASC-B01` | When the confronted node is not of kind feature, the gate skips confrontation. |
| `STASC-B02` | When the project configuration is nil or defines no rule types, the gate skips confrontation. |
| `STASC-B03` | When configured rule types define no scenario tag mappings, the gate skips confrontation. |
| `STASC-B04` | When a feature file contains no scenarios declaring identity codes, the gate skips confrontation. |
| `STASC-B05` | Scenarios with identity codes lacking a recognized rule letter are ignored during confrontation. |
| `STASC-B06` | Scenario tags not registered in configuration rule types are ignored during confrontation. |
| `STASC-B07` | When every recognized scenario tag aligns with its code's rule letter, the gate passes. |
| `STASC-B08` | A scenario tag declared under multiple rule letters passes if any of its mapped letters matches the code's letter. |
| `STASC-B09` | When a scenario is co-tagged with multiple codes, a tag matching any of those codes passes. |
| `STASC-B10` | When a recognized scenario tag disagrees with the code's rule letter, the gate returns Pending. |
| `STASC-B11` | The Pending message cites the mismatched code, rule letter, tag name, allowed letters, and shortened scenario title. |
| `STASC-B12` | Multiple mismatch findings are sorted deterministically and reported together in the Pending verdict. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `STASC-I01` | Scenario classification alignment is evaluated exclusively on feature files; non-feature nodes skip to keep confrontation scoped to where tags reside. | confronts non-feature nodes and verifies the verdict is Skip |
| `STASC-I02` | Without configured tag mappings, the gate stays silent to prevent false positives across different localization languages. | confronts a feature when configuration has no tag mappings and verifies it skips |
| `STASC-I03` | Type disagreements return Pending rather than Fail because inherited divergence requires manual editorial judgment rather than mechanical fixes. | confronts a feature with mismatched tags and verifies the verdict is Pending |
| `STASC-I04` | Co-tagged scenarios with multiple requirement codes accept tags matching any attached code to prevent false defects on valid multi-requirement scenarios. | confronts a scenario with secondary code matching the tag and verifies it passes |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `STASC-X01` | Does not enforce presence of classification tags on scenarios. | Categorization tags are optional metadata; scenarios without type tags are not penalized. |
| `STASC-X02` | Does not decide whether the tag or the code is the erroneous party when disagreement occurs. | Distinguishing whether a tag was misapplied or a code was borrowed requires reading scenario steps and understanding domain context. |
| `STASC-X03` | Does not validate or alter specifications or test code files. | Specification section validation belongs to `rule-types` and scenario-to-test code validation belongs to `feature-test-match`. |
| `STASC-X04` | Does not restrict tags from being mapped to more than one rule type letter. | Overlapping tags legitimately apply across multiple rule categories, such as state-data applying to both state and validation. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `TagLetters` | core — configuration schema and tag-to-letter mappings |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindFeature`, `Node` | core — graph model and artifact representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
