<!-- @anchors
  code: SCIDS
  updated_at: 2026-09-19
  layer: gate
-->
# ScenarioIdentity — two scenarios of the same feature cannot share one code

> **Code**: `SCIDS`

## Overview

Confronts a feature against the question that decides whether its scenarios can be paired
at all: **each scenario has a code, but does each code point at ONE scenario?**

A rule legitimately has several scenarios — the happy path and its alternatives. What it
cannot have is two that are INDISTINGUISHABLE. With the same code, nothing links one
scenario to one test, and the relational gates end up comparing N titles against a single
test: at most one matches, and the others become a divergence nobody can resolve.

The way out is the scenario SUFFIX. Numbering keeps the rule legible in the prefix and
gives each case an identity of its own.

Measured on the project that originated the gate: 204 repeated codes, and 65 of them
carrying CONFLICTING kind tags on the same code — one scenario tagged as state, the other
as behaviour. That is the signature of a BORROWED code, not of a rule with two paths. One
of them proved to be a real defect: the spec defined the rule as one thing, and the second
scenario described a different behaviour that had no code of its own.

**PENDING and not FAIL**: numbering scenarios is a migration, and the gate is born over a
base that did not know the notation. Whoever already migrated stays green; whoever did not
sees what is left.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the node kind | a feature, the only artifact that carries scenarios | spec, code, test or guide | this unit: a kind outside the set leaves without a verdict |
| the artifact's content | any text, empty included | — (text with no coded scenario is a case, not an error) | this unit: nothing to confront is a skip, never a failure |
| the scenario codes | what the feature parser extracts, suffix included | a tag that is not a scenario code | `internal/gate/feature_test_match.go`, which owns the parsing of the feature |

## Effects

| Effect | Description |
| --- | --- |
| `SCIDS-B01` | Two scenarios sharing one code are reported: nothing links one of them to one test, and the relational gates compare N titles against a single test. |
| `SCIDS-B02` | The report NAMES the repeated code, so the reader does not have to scan the feature to find it. |
| `SCIDS-B03` | The report says HOW MANY scenarios share the code, which separates an accidental duplicate from a code borrowed across a whole rule. |
| `SCIDS-B04` | The report teaches the way out with the project's OWN repeated code as the example, not a generic one. |
| `SCIDS-B05` | The SUFFIX gives each scenario its own identity: numbered, two scenarios of one rule pass. |
| `SCIDS-B06` | Distinct codes pass — a rule with several scenarios is the common case and must not be accused. |
| `SCIDS-B07` | The verdict is PENDING and never a failure: numbering is a migration, and the gate is born over a base that did not know the notation. |
| `SCIDS-B08` | An artifact that is not a feature leaves without a verdict: only a feature carries scenarios. |
| `SCIDS-B09` | A feature with no coded scenario leaves without a verdict — there is nothing to confront, and that absence is another gate's charge. |
| `SCIDS-B10` | Several repeated codes in one feature are reported TOGETHER, in a stable order, so two runs over the same file produce the same message. |
| `SCIDS-B11` | Long titles are shortened in the report: the address is the code, and the title only helps recognise which scenario is which. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SCIDS-I01` | Grouping is by the COMPLETE code, suffix included. Grouping by the prefix alone would accuse exactly the projects that already did the migration the gate asks for. | confronts a feature whose scenarios share a prefix and differ in suffix, and verifies it passes |
| `SCIDS-I02` | The message is deterministic. The repeated codes are ordered before being joined, so the same feature always produces the same text and the finding does not churn between runs. | confronts a feature with several repeated codes twice and compares the messages |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SCIDS-X01` | Does not judge whether the two scenarios describe DIFFERENT behaviours. | That is the judgement the conflicting kind tags hinted at, and it is not deterministic. Here the ruler is identity: one code, one scenario — and that is decidable by reading. |
| `SCIDS-X02` | Does not look across features. | A code repeated in two different features is a different defect with a different owner; the relational gates that hold the map are the ones that can see it. Charging it here would need a graph this gate does not take. |
| `SCIDS-X03` | Does not charge the ABSENCE of a code on a scenario. | An uncoded scenario is invisible to the parser and to every relational gate. Charging it is the ruler of the gate that pairs scenarios to tests, and two gates on one defect become noise. |
| `SCIDS-X04` | Does not renumber the scenarios, even knowing the fix. | The suffix carries meaning — which case is the first, which is the alternative — and picking it is the author's call. A gate that rewrote the feature would also invalidate every test already bound to the old code. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindFeature` | core — the kind decides whether the artifact carries scenarios at all |
| DEP2 | `internal/config/config.go` | `Config` | core — the project's Structure travels with the confrontation |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
