<!-- @anchors
  code: SFMSP
  updated_at: 2026-09-26
  layer: gate
-->
# SpecFeatureMatch — every requirement the spec DEFINES has at least one scenario

> **Code**: `SFMSP`

## Overview

Confronts the edge of the triad that had no watcher: **the spec declares a requirement —
is there any scenario that exercises it?**

`feature-test-match` confronts feature→test; `triad-complete` confronts that the PIECES
exist. Nobody confronted spec→feature — and that is where a silent hole lives: the spec
declares `XXXXX-X02`, the feature has no scenario carrying that tag, and the requirement
crosses the whole pipeline with nothing verifying it. **Every gate stays green**: the spec
has a code, the feature exists, the feature matches the test. The requirement simply
belongs to nobody.

Measured in a real project: **11 of 287 specs with a feature had a requirement with no
scenario.**

The ruler is the same as the other relational gates — CODE, not prose: every
`{CODE}-{letter}{NN}` the spec DEFINES must appear as a scenario tag in the feature.
**Defining is different from citing:** a spec that mentions another unit's code (in a
Dependency Table, for instance) contracts no obligation. Without that distinction the
gate would be a noise generator and would be switched off.

Honest opt-out (CONCEPT §5.1): the per-requirement waiver marker (`no-scenario`, prefixed with `@`) on the requirement's line, followed by a colon and the reason, waives
that specific requirement, with the reason written. It serves what is genuinely not
observable by scenario — and leaves the trace that it was a decision, not forgetfulness.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any map node | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` — a node that is not a spec leaves without a verdict |
| the spec's content | any text, with or without requirements | — (a spec that defines nothing is a case, not an error) | this unit: with no defined requirement the confrontation is skipped |
| the defined requirement | a code at the start of a line, a list item, a section title, or in a table's FIRST cell | a code cited in prose or in a Dependency Table, which contracts nothing | this unit, by the same "defines" grammar the `rule-types` gate uses |
| the linked features | the features reached by the `covered-by` edge, one or more | — (no feature is the ruler of `triad-complete`) | this unit: with no feature it skips, so the same defect is not reported twice |
| the waiver | the per-requirement marker (`no-scenario`) or the whole-spec one (`no-feature`), each prefixed with `@`, followed by a colon and a written reason | a bare marker with nothing after it | this unit: a waiver with no why does not waive |
| the map | a built graph, or none | — | this unit: with no graph the verdict is Pending, never Pass |

## Effects

| Effect | Description |
| --- | --- |
| `SFMSP-B01` | A requirement the spec defines and no scenario tags FAILS, and the verdict names it. |
| `SFMSP-B02` | A requirement that HAS a scenario is not accused: the verdict carries the uncovered ones only. |
| `SFMSP-B03` | With every defined requirement tagged by some scenario, the gate passes. |
| `SFMSP-B04` | A code merely CITED — in prose or in the Dependency Table — contracts no obligation, because a spec cites other units' codes all the time. |
| `SFMSP-B05` | The per-requirement marker (`no-scenario`) with a written reason waives that one requirement. |
| `SFMSP-B06` | A bare per-requirement marker, with nothing after the colon, does not waive — the waiver requires a reason. |
| `SFMSP-B07` | The whole-spec marker (`no-feature`) with a reason DRAGS the waiver to every requirement: the spec has no feature, so no requirement of it can have a scenario. |
| `SFMSP-B08` | Without that tag the same uncovered requirements keep failing — the drag cannot become a silent way of muting the gate. |
| `SFMSP-B09` | A bare whole-spec marker drags nothing, or the marker would be a switch that turns the gate off without accounting for it. |
| `SFMSP-B10` | A spec with NO feature returns Skip: that absence is the ruler of `triad-complete`. |
| `SFMSP-B11` | A spec covered by SEVERAL features has its requirements looked for across all of them — the requirement only needs to be in some. |
| `SFMSP-B12` | An artifact that is not a spec returns Skip: the gate has no jurisdiction over code, test or feature. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SFMSP-I01` | Every waiver requires a written REASON — both the per-requirement marker and the whole-spec one. A bare marker is a switch with no accounting, and silence without a why is what the gate exists to end. | confronts both bare markers and verifies each still fails |
| `SFMSP-I02` | Each gate accuses ONE thing. The missing feature belongs to `triad-complete` and the missing test to `feature-test-match`; accusing them here would print the same defect twice in the report. | confronts a spec with no feature and verifies it skips instead of failing |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SFMSP-X01` | Does not judge whether the scenario PROVES the requirement. | The ruler is the tag, which is deterministic. Whether the Given/When/Then really exercises the behaviour is judgement, and `scenario-asserts` is the ruler that looks at the shape of the assertion. Holding both here would put the same rule in two places that would diverge. |
| `SFMSP-X02` | Does not charge a code the spec merely CITES. | A spec names other units' codes in Dependency Tables, notes and cross-references, and contracts nothing by doing so. Charging them would produce an accusation per citation — and a gate that cries wolf gets switched off, which costs more than the defect it was catching. |
| `SFMSP-X03` | Does not confront feature→test. | That edge already has a watcher, `feature-test-match`. This gate exists precisely because the edge BEFORE it had none. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `SFMSP-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec` | core — jurisdiction comes from the node's KIND, and the feature is reached by the `covered-by` edge |
| DEP2 | `internal/config/config.go` | `CodeLengthPattern` | core — the identity code's shape is the project's, and the "defines" grammar is built on it |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
