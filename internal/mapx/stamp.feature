# language: en
# @anchors
#   ref: EDSTD
#   updated_at: 2026-09-28
#   layer: feature

@EDSTD
Feature: EdgeStamping — recording on each relation that it was confronted, with what result, and since when

  @EDSTD-B01 @unit-level
  Scenario: Only relations with both ends confronted are stamped
    Given a spec that specifies code, which is tested by a test
    When a round confronts only the spec and the code
    Then one relation is stamped, spec to code, and code to test has no stamp

  @EDSTD-B02 @unit-level
  Scenario: The verdict is issue when an end failed and ok when both passed
    Given a spec and its code confronted in a round
    When the spec fails in one round and both pass in another
    Then the relation's verdict is issue after the first and ok after the second

  @EDSTD-B03 @unit-level
  Scenario: A waived stamp is kept as it was
    Given a spec-to-code relation stamped waived on 2026-09-01 at old revisions
    When a round confronts both ends on 2026-09-24
    Then the stamp is still waived, dated 2026-09-01, at the old revisions, and the relation reads stale

  @EDSTD-B04 @unit-level
  Scenario: A stamp is fresh until an end moves
    Given a never-validated spec-to-code relation
    When a round stamps it and then the spec moves to a new revision
    Then the relation is fresh after the round and stale after the spec moved

  @EDSTD-B05 @unit-level
  Scenario: The date moves only when revisions or verdict change
    Given a relation stamped ok on 2026-08-30
    When it is stamped ok again on 2026-09-15, then after an end moved on 2026-09-20, then with a failed end on 2026-09-25
    Then its date is 2026-08-30, then 2026-09-20, then 2026-09-25

  @EDSTD-B06 @unit-level
  Scenario: An undated stamp takes today's date
    Given a relation stamped ok at the current revisions but with no date
    When a round stamps it ok on 2026-09-01
    Then its date is 2026-09-01

  @EDSTD-B07 @unit-level
  Scenario: One relation is stamped by its ends, a missing one answers false
    Given a spec-to-code and a code-to-test relation
    When the spec-to-code relation is stamped issue, and a relation from X to Y is stamped
    Then only spec-to-code carries the issue stamp, and stamping X to Y answers false

  @EDSTD-B08 @unit-level
  Scenario: Stamping a relation by gate records the gate and a judgment
    Given a guide-to-spec relation
    When it is stamped ok by the gate named guide-judge
    Then the stamp names guide-judge, and the spec counts as judged ok by guide-judge

  @EDSTD-B09 @unit-level
  Scenario: Stamping a node stamps every relation touching it
    Given code with one relation coming in from its spec and one going out to its test, and an unrelated relation
    When the code is stamped ok by the gate named code-judge
    Then two relations are stamped, each with a code-judge judgment, and the unrelated relation has none

  @EDSTD-B10 @unit-level
  Scenario: A judgment holds only at the revisions it was given and only for its gate
    Given a spec with one relation to its test, never judged
    When the spec is judged ok by my-gate, and later the spec moves to a new revision
    Then it was not judged before, is judged ok by my-gate after, is never judged by another gate, and is no longer judged once it moved

  @EDSTD-B11 @unit-level
  Scenario: A judgment survives the next check round
    Given a spec judged ok by my-gate
    When a check round stamps both ends of its relation
    Then the spec still counts as judged ok by my-gate

  @EDSTD-B12 @unit-level
  Scenario: One judgment per gate, its date kept when nothing changed
    Given a relation judged ok by my-gate on t0
    When it is judged ok again on t1, and then issue on t2
    Then the relation holds a single my-gate judgment, dated t0 after the second judging and issue dated t2 after the third

  @EDSTD-B13 @unit-level
  Scenario: The stale relations are listed
    Given two never-validated relations
    When one of them is stamped
    Then two relations are listed stale before and one after

  @EDSTD-B14 @unit-level
  Scenario: Judging keeps a waiver another gate recorded, and the gate that waived replaces its own
    Given the relations of the code node waived by the gate "atomic"
    When the gate "review" judges the node ok and the gate "rule-fulfilled" judges one relation as an issue
    Then every relation's stamp is still waived by "atomic" with its date, and "review" has its judgment recorded
    And when "atomic" judges the node ok the stamps read ok, and a new waiver by "review" lands

  @EDSTD-I01 @unit-level
  Scenario: The same round on the same day gives the same stamps
    Given two identical graphs
    When each is stamped with the same verdicts on 2026-08-30
    Then the two stamps are identical and dated

  @EDSTD-B15 @unit-level
  Scenario: A round carries only the stamps it changed to the map on disk
    Given a copy of the map snapshotted and stamped by a round, and the map on disk where another process changed one of the same stamps
    When the round's changes are applied to the map on disk
    Then the stamps only the round changed are applied, the other process's stamp is kept and counted, and an edge the disk does not have is left out

  @EDSTD-X01 @unit-level
  Scenario: The stamp carries the caller's date
    Given a relation to be stamped
    When a round stamps it with the date 2026-09-01
    Then the stamp's date is exactly 2026-09-01

  @EDSTD-B16 @unit-level
  Scenario: A new revision that proves nothing new keeps what was measured
    Given a file measured at one revision, with a proof, a coverage, a mutation, a test's closure, a stamp and a judgment, and a coverage of an older revision
    When the file moves to a new revision
    Then everything measured at the first moves with it, and the older coverage and the other file stay
