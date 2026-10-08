# language: en
# @anchors
#   code: STFTE
#   ref: EDSTD
#   updated_at: 2026-10-08
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

  @EDSTD-B17 @unit-level
  Scenario: A declared change keeps what was proven, and the lines only when asked
    Given a spec proven at two earlier revisions, a code file with coverage and mutation, and a test whose closure holds the spec
    When their evidence is kept, with and without lines
    Then the proofs, closures and stamps move to the current revision, the coverage moves only with lines, and each declaration is recorded

  @EDSTD-B18 @unit-level
  Scenario: A repair carries only the evidence that held at the file's revision before it
    Given a file whose coverage and proof are at its content and whose mutation ran on an older one
    When check --fix repairs it with no line moved
    Then the coverage, the proof and the closures move to the new content, and the mutation stays stale

  @EDSTD-B19 @unit-level
  Scenario: A rebuild carries the evidence of a file whose evidence revision held
    Given a proven spec, a covered code file, a test whose closure reaches it and an edge stamped on both
    When the spec's revision moves with its evidence held, and the code gains a flag line, a flag at a line's end, or a change
    Then the spec's proofs stand, the test stays fresh and the stamp follows, coverage goes along only when the lines held, and a change carries nothing

  @EDSTD-B20 @unit-level
  Scenario: A spec whose rules alone changed carries its evidence with those rules' scenarios stale
    Given a spec proven for two rules, the second with two variants
    When the second rule's definition changes, and then something outside the rules
    Then the second rule's variants are stale and the first's proof stands, and the change outside drops the proofs

  @EDSTD-B21 @unit-level
  Scenario: A stale scenario is fresh again when a run proves it, or the author keeps the evidence
    Given a spec whose two variants of a rule are stale
    When a suite proves one again, another stops proving the other, and the author keeps the evidence
    Then the first is fresh, the second leaves the stale ones, and the author's declaration clears what was left

  @EDSTD-B22 @unit-level
  Scenario: Advancing a node to its content carries the evidence its revisions show unchanged
    Given a code file and a test whose closure reaches it
    When the file is advanced to its content with a flag, then with a change, then an unknown file
    Then the flag carries the evidence, the change does not, and the unknown file is left alone

  @EDSTD-B23 @unit-level
  Scenario: A map written before evidence revisions gets them from the content at its revision
    Given a proven spec in a map with no evidence revisions, edited since in its date and history alone
    When the build reads a content that is not the node's revision, then the one that is
    Then the first is not taken, the second carries the proof, and a map with evidence revisions is not read
