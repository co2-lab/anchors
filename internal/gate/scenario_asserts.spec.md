<!-- @anchors
  code: SCASS
  updated_at: 2026-09-19
  layer: gate
-->
# ScenarioAsserts — scenario outcome steps must assert concrete verifiable outcomes

> **Code**: `SCASS`

## Overview

Confronts feature files against tautological step definitions: **the outcome step of a scenario must assert an observable result, rather than merely repeating the requirement code.**

A scenario in Gherkin exists to govern the test: it commits the system to an observable outcome before anyone writes the assertion code. When an outcome step degenerates into a tautology — such as stating that the requirement code is verified — it asserts nothing. Two scenarios for entirely different requirements become indistinguishable, and the definition of what constitutes verification migrates completely into the test implementation, which is the very artifact the scenario was meant to govern.

The operational risk of tautological scenarios is structural: when nobody between the specification and the test is required to write down the expected outcome, the resulting test is written without a discriminating case. A mutation that completely removes the production logic can leave the test suite green, because the test itself was never anchored to a concrete outcome. In a measured reference project, 359 scenarios suffered from this exact defect, having been mechanically copied from templates.

This gate prevents this structural defect by isolating outcome steps and mechanically checking whether the step consists solely of a requirement identity code wrapped in linking words.

This gate operates in distinct territory from neighbouring gates:
- Unlike `feature-test-match`, which verifies scenario-to-test alignment by comparing scenario codes and titles against test titles, this gate inspects the inner substance of outcome steps to ensure they assert real behavior.
- Unlike `spec-complete`, which verifies that requirements are structurally catalogued in specifications, this gate ensures that scenarios realizing those requirements commit to verifiable outcomes.
- Unlike natural language review gates, this gate does not judge prose style; it applies a deterministic filter that flags steps where removing linking words leaves only the requirement identity.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind feature | nodes of kind spec, code, test, plan, or config | this unit: skips non-feature nodes |
| the feature content | feature file text containing scenario definitions | binary files or non-feature text | the gate engine, routing text artifacts |
| outcome steps | lines starting with localized outcome keywords | setup lines, action lines, and comments | this unit: filters lines by outcome keyword set |
| step assertions | statements declaring observable state or behavior | bare requirement codes surrounded only by linking words | this unit: strips codes and linking words to detect tautologies |

## Effects

| Effect | Description |
| --- | --- |
| `SCASS-B01` | When the confronted node kind is not a feature, the gate skips confrontation. |
| `SCASS-B02` | When all outcome steps in a feature assert observable results, the gate passes. |
| `SCASS-B03` | When an outcome step merely asserts that a requirement code is verified, the gate fails, citing the code. |
| `SCASS-B04` | Tautological outcome steps wrapped with linking words and up to two residual words fail. |
| `SCASS-B05` | An outcome step citing a requirement code while asserting three or more descriptive content words passes. |
| `SCASS-B06` | Multiple tautological outcome steps within a feature are all reported, deduplicated and sorted. |
| `SCASS-B07` | Outcome steps written in any recognized dialect keyword are evaluated according to configured project settings. |
| `SCASS-B08` | Empty lines and comment lines in feature files are ignored and never evaluated as outcome steps. |
| `SCASS-B09` | Non-outcome steps such as setup and action steps are ignored even when citing requirement codes. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SCASS-I01` | Up to two residual content words beyond linking words around a code is classified as a tautology, preventing trivial phrasing variations from bypassing the rule. | confronts variations with linking words and verifies failure |
| `SCASS-I02` | Truly assertive outcome steps are never flagged as tautologies, protecting against false positives that would encourage teams to disable the gate. | confronts diverse assertive outcome formats and verifies pass |
| `SCASS-I03` | Language recognition covers all supported dialect alternatives, preventing inherited multi-language features from being silently bypassed. | confronts outcome steps in English, Spanish, and Portuguese and verifies consistent enforcement |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SCASS-X01` | Does not evaluate the semantic accuracy or elegance of prose beyond the mechanical removal of linking words. | Subjective stylistic judgment belongs to human review and language models, while structural gates remain deterministic. |
| `SCASS-X02` | Does not inspect setup and action steps for requirement code references. | Setup and trigger steps legitimately reference requirement context without needing to assert system outcomes. |
| `SCASS-X03` | Does not enforce the presence of scenarios or tests. | Scenario presence is governed by feature-to-spec alignment gates and triad completeness rules. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern`, `Config` | core — identity code length pattern and dialect configuration |
| DEP2 | `internal/config/dialect.go` | `GherkinThenAlternatives` | core — set of supported outcome keywords across languages |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — localized defect messages |
| DEP4 | `internal/mapx/model.go` | `Graph`, `KindFeature`, `Node` | core — graph model and node kind representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
