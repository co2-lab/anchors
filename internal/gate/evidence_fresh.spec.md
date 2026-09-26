<!-- @anchors
  code: EVFRV
  updated_at: 2026-09-26
  layer: gate
-->
# EvidenceFresh — the score of this test holds against TODAY's code

> **Code**: `EVFRV`

## Overview

Every other gate confronts text against text. This one confronts an EXECUTION against the
revisions it measured: **the test went green, and the stamp records the revision of every
node in its closure at that moment. If any of them moved, the score is still written and
has stopped being true.**

It is the failure that makes no noise. A test that breaks screams in the runner; a score
that ages stays green in yesterday's report, and disappears inside the conclusion of
whoever reads it. Measured in the reference app: `utils/login.yaml` is composed into 290
scripts — touching it invalidated 290 measurements, and the signal of not one of them
changed.

**The ruler for when NOT to judge matters as much as the one for when to fail.** A test with
no stamp has no score to expire: that is absence of proof, which is a different debt, with a
different name and a different fix — run it the first time, not revalidate it. Reporting
them together is the defect the edge-level `stale` has: 30 "never validated" drowned in
1419 "revision advanced", and the list stops being read.

It is blocking because the fix is known, local and cheap: run the test. There is no domain
judgement to make — unlike an undeclared letter, where both ways out are legitimate and the
choice belongs to whoever knows the product.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the map | a built graph, or none | — (no graph is a case, not an error) | this unit: with no map there is no closure to walk, and the verdict is a skip |
| the execution stamp | a signal carrying the revision the run measured, or none | a signal whose revision is empty, which counts as none | this unit: absence is a skip, never a failure — the absent proof is another gate's debt |
| the measured closure | the revisions the run recorded for every node it depends on, empty included | — (an empty closure is a case: nothing to age) | the ingestion that wrote the stamp |

> `Who guarantees` cannot be left empty nor say only "not mine": if nobody guarantees, the
> duty is orphaned — and that is exactly where the invalid input gets through.

## Effects

| Effect | Description |
| --- | --- |
| `EVFRV-B01` | An artifact that is not a test leaves the confrontation without a verdict: only a test carries an execution score. |
| `EVFRV-B02` | Without a built map the gate stays quiet: there is no closure to walk, and it will not approve what it could not look at. |
| `EVFRV-B03` | A test with NO execution stamp is skipped, never failed — a score that was never written cannot have expired. |
| `EVFRV-B04` | A test whose closure is intact PASSES. |
| `EVFRV-B05` | The passing verdict says WHAT it checked against: the size of the closure it walked is the difference between "nobody looked" and "I looked and it stands". |
| `EVFRV-B06` | A test whose DEPENDENCY advanced a revision FAILS: the score is still written and stopped holding. |
| `EVFRV-B07` | The failing verdict NAMES the culprit, so whoever fixes it knows what moved underneath. |
| `EVFRV-B08` | The failing verdict states the FIX — run the test again — because a verdict that accuses without saying what to do transfers the work to the reader. |
| `EVFRV-B09` | A test whose OWN file changed since the run is reported as such, separately from the closure. |
| `EVFRV-B10` | The culprit list is TRUNCATED at five, and the remainder is counted. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `EVFRV-I01` | Absence of proof and expired proof are NEVER reported as the same finding. They are different debts with different fixes, and merging them is what drowns the list until it stops being read. | confronts a test with no stamp and verifies the verdict is a skip, not a failure |
| `EVFRV-I02` | A test that never ran is never approved either. The skip is silence, not a stamp: approving would state a freshness nobody measured. | confronts a test with no stamp and verifies the verdict is not Pass |
| `EVFRV-I03` | Truncation never hides the size of the problem: what is not listed is COUNTED. Whoever fixes it runs the test once, regardless of how many dependencies moved — but they still learn how far it spread. | confronts a run with twenty aged dependencies and verifies the verdict reports the remainder |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `EVFRV-X01` | Does not charge the ABSENCE of a green test. | That is another gate's ruler — coverage. Charging it here would merge "never ran" with "ran and aged", which are different debts with different fixes. |
| `EVFRV-X02` | Does not RUN the test, nor judge whether the change actually broke it. | The gate measures whether the evidence still covers the current code, not whether the behaviour changed. Deciding that the diff was harmless is judgement, and the cheap fix — run it again — settles it for real instead of by opinion. |
| `EVFRV-X03` | Does not read the project's configuration. | The confronted truth — the revision moved — lives in the map, not in a convention. A ruler that depended on settings could be turned off by a default nobody chose. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `EVFRV-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph` | core — the closure walk and the culprit list belong to the map, which is what holds the revisions |
| DEP2 | `internal/config/config.go` | `Config` | core — the freshness window and the declared surfaces come from the Structure |
| DEP3 | `internal/mapx/model.go` | `KindTest` | core — only a test carries an execution score, and the kind is what routes the jurisdiction |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
