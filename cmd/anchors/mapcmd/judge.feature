# language: en
# @anchors
#   ref: JDGUE
#   updated_at: 2026-09-26
#   layer: feature

@JDGUE
Feature: Judge — records an AI's verdict on a judgment gate with the same bookkeeping as a deterministic gate

  @JDGUE-B01 @unit-level
  Scenario: The legacy spelling of the waiver is a waiver all the way to the stamp
    Given the login project with the atomic judgment gate that declares the code guide
    When the login code is judged with the verdict "dispensado" and a reason
    Then the output says WAIVED and never PASS
    And the guide-to-login edge is stamped "waived"

  @JDGUE-B02 @unit-level
  Scenario: A reason is required for fail and waived, not for pass
    Given the verdicts pass, fail and waived with and without a reason
    When each verdict is validated
    Then pass without a reason is accepted, fail and waived without one are refused, and waived with one is accepted

  @JDGUE-B03 @unit-level
  Scenario: Only the exact name of a declared judgment gate is accepted
    Given a configuration declaring the judgment gate "mock-detect-covers-dialect"
    When the gate is looked up by its name, by its legacy name and by an undeclared name
    Then only the exact name is found

  @JDGUE-B04 @unit-level
  Scenario: A fail and then a pass stamp the guide edge issue and then ok
    Given the login project with the atomic judgment gate that declares the code guide
    When the login code is judged "FAIL" with a report and later "pass"
    Then the guide-to-login edge is stamped "issue" after the fail and "ok" after the pass
    And the atomic gate counts as answered on the login code

  @JDGUE-B05 @unit-level
  Scenario: The review verdict stamps every edge of the target
    Given the login project where the login code has two edges and review is not declared as a gate
    When the login code is judged "pass" by review
    Then the output says "stamped: 2 edge(s)" and both edges carry a review judgment

  @JDGUE-B06 @unit-level
  Scenario: A target not in the map is recorded on its unit's spec
    Given the login project with no "src/login.tsx" in the map
    When "src/login.tsx" is judged "pass" by review
    Then the output says it is "recording on `src/login.spec.md`"

  @JDGUE-B07 @unit-level
  Scenario: A fail opens an issue, a repeated report changes nothing, a new report reopens it
    Given the login project with the atomic judgment gate
    When the login code is judged "fail" with a report, again with the same report, and then with a different one
    Then the first opens an issue in todo whose body is the report
    And the second says "same finding already recorded"
    And the third says "NEW finding added"

  @JDGUE-B08 @unit-level
  Scenario: A pass resolves the open issue and a waiver resolves it as a waiver
    Given the login project with an open atomic issue on the login code
    When the login code is judged "pass"
    Then the output says "previous issue resolved" and the issue moves to done
    And a waiver over an open issue says "judged WAIVED" and "previous issue resolved"

  @JDGUE-B09 @unit-level
  Scenario: Manual mode stamps the map and prints the report without writing an issue
    Given the login project in manual workflow mode
    When the login code is judged "fail" with the report "  the report itself  "
    Then the output says "no issue written (mode: manual" and prints "the report itself"
    And no issues folder exists and the guide edge is stamped "issue"
    And with the record-issues switch the issue is written

  @JDGUE-B10 @unit-level
  Scenario: A fail with a patch opens an applicable fix suggestion
    Given the login project and a patch file that replaces "x" by "y"
    When the login code is judged "fail" with that patch
    Then the output announces "anchors suggest show judge-atomic-src-login"
    And the pending suggestion "judge-atomic-src-login" carries "+y"

  @JDGUE-B11 @unit-level
  Scenario: The verdict closes its judge task and the pending list shrinks
    Given the queue holds the judge tasks "judge-review-src-login" and "judge-atomic-src-login"
    When the pending list is shown, the login code is judged by review, and the list is shown again
    Then the first list shows "2 target(s)", the review task is done, and the second list shows "1 target(s)"
    And with an empty queue the list says "no target awaiting judgment"

  @JDGUE-I01 @unit-level
  Scenario: A waiver is never stamped ok
    Given the login project with an open atomic issue on the login code
    When the login code is judged "waived" with a reason
    Then the guide-to-login edge is stamped "waived", not "ok"

  @JDGUE-X01 @unit-level
  Scenario: A waiver is never announced as a pass
    Given the login project with an open atomic issue on the login code
    When the login code is judged "waived" with a reason
    Then the output does not contain "PASS"
    And a second waiver with nothing to resolve says "○ judged WAIVED — there was nothing to confront"

  @JDGUE-E01 @unit-level
  Scenario: A judge with no target is refused
    Given the login project
    When judge runs with no target and no pending switch
    Then it fails with "provide the target"

  @JDGUE-E02 @unit-level
  Scenario: A judge with no gate is refused
    Given the login project
    When the login code is judged with no gate
    Then it fails with "--gate is mandatory"

  @JDGUE-E03 @unit-level
  Scenario: A verdict outside the three is refused
    Given the login project
    When the login code is judged "maybe" by the atomic gate
    Then it fails with "--verdict must be"

  @JDGUE-E05 @unit-level
  Scenario: A gate that is not a judgment gate is refused
    Given the login project where always-red is a run gate
    When the login code is judged "pass" by always-red
    Then it fails with "not a judgment gate"

  @JDGUE-E06 @unit-level
  Scenario: A project with no map is refused pointing at the build
    Given a directory with no map
    When the login code is judged there
    Then it fails with "anchors map build"

  @JDGUE-E07 @unit-level
  Scenario: A target whose unit has no piece in the map is refused
    Given the login project with nothing of the signup unit in the map
    When "src/signup.tsx" is judged by review
    Then it fails with "\"src/signup.tsx\" is not in the map"

  @JDGUE-E08 @unit-level
  Scenario: A map with no configuration beside it is refused
    Given a directory holding the login map and code but no anchors.yaml
    When the login code is judged by the atomic gate
    Then it fails with "load config"

  @JDGUE-E09 @unit-level
  Scenario: An unreadable patch file is refused
    Given the login project and a patch path that does not exist
    When the login code is judged "fail" with that patch
    Then it fails with "read the patch"
