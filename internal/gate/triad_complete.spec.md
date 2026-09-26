<!-- @anchors
  code: TRCMT
  updated_at: 2026-09-26
  layer: gate
-->
# TriadComplete — the pieces that realize a spec EXIST

> **Code**: `TRCMT`

## Overview

Confronts a spec of a GOVERNED layer against the simplest question of the triad: **do the pieces
that realize it exist?** The code it specifies, the feature that covers it, and the test
that proves it.

It exists because the relational gates FAIL OPEN by construction. Without a test linked, the
feature↔test confrontation returns "nothing to confront yet" instead of failing; without code,
the dependency one likewise. The side effect is grave: a lone spec, with no
implementation at all, crosses ALL the gates and the pipeline concludes "can promote" — the green
certifying work that does not exist.

This gate closes the hole from the positive side. Instead of asking "do the pieces match?" — which
requires that they exist —, it asks "do the pieces exist?".

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any map node | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the map | a built graph, or none | — (an absent map is a case, not an error) | this unit: without a map the verdict is undetermined, never approval |
| the target's layer | a governed layer, a recognized one, or none declared | — | this unit, by the regime declared in the Structure |
| the per-unit waiver | the mark with a written reason beside it | a bare mark, without a reason | this unit: a waiver without a why does not waive |

## Effects

| Effect | Description |
| --- | --- |
| `TRCMT-B01` | An artifact that is not a spec leaves without a verdict: only the spec has a triad to demand. |
| `TRCMT-B02` | A RECOGNIZED layer (declarative regime) leaves without a verdict: it has neither spec nor triad by definition. |
| `TRCMT-B03` | Without a map the verdict is UNDETERMINED. Approving without being able to look would be asserting what was not measured. |
| `TRCMT-B04` | A spec with the three pieces linked passes. |
| `TRCMT-B05` | A spec missing some piece fails, and the verdict NAMES which ones are missing and where each one is born. |
| `TRCMT-B06` | The layer may waive a piece as a block, declared in the Structure. |
| `TRCMT-B07` | The unit may waive a piece in the spec itself, with a written reason — it is the granularity the per-layer waiver does not reach. @realizes SAIDA-R01 |
| `TRCMT-B08` | A piece declared TO BE DEVELOPED (`@TBD`) leaves the verdict UNDETERMINED, never approved: it is DEBT, and it stays visible until someone writes it. Measured before the fix: `@TBD` shared a bucket with the permanent waiver, and a spec declaring three unwritten pieces came out green — the honest declaration erased the pending work from the radar. @realizes SAIDA-R02 @realizes SAIDA-R03 |
| `TRCMT-B09` | A waiver written on a line that DEFINES a rule — its table row, its heading (`### CODE-S01: …`) or its bullet — is that rule's, not the unit's: only a marker on a line of its own waives a piece of the triad. Measured before the fix: MIF had 308 triad failures, 293 of them per-rule `@no-code:` in table rows read as a unit waiver of code, feature and test. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TRCMT-I01` | The test is reached in TWO hops — spec → feature → test —, because the one that points at the test is the feature. Checking the test directly on the spec would accuse the whole project of missing tests. | links the triad in two hops and verifies that the gate considers it complete |
| `TRCMT-I02` | Waiving the test and writing a scenario in the feature is a CONTRADICTION, and fails. The two assertions do not coexist: either the scenario is real and someone must prove it, or it should not exist. | declares the waiver, links a feature with a scenario, and verifies the failure |
| `TRCMT-I03` | Waiving the test requires saying WHERE the proof is, and the place has to exist. An orphaned reference fails. | declares the waiver pointing at a nonexistent target and verifies the failure |
| `TRCMT-I04` | The waiver holds only for the declared piece. Waiving one never waives the others. | declares the waiver of one piece and verifies that the rest remain demanded |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TRCMT-X01` | Does not confront whether the pieces MATCH one another — only whether they exist. | Matching is the work of the relational gates. This one exists precisely because they fail open when the piece does not exist; doing both here would duplicate the ruler. |
| `TRCMT-X02` | Does not judge the QUALITY of any piece. | An empty test satisfies this gate, and that is correct: the ruler here is EXISTENCE. The one that confronts the content is another gate, and confusing the two would make this one fail for a reason it does not know how to measure. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph` | core — the pieces are edges, and without the graph there is nothing to look at |
| DEP2 | `internal/config/config.go` | `Config` | core — the layer's regime and the block waiver are declared in the Structure |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
