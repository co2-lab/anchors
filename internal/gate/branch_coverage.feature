# language: en
# @anchors
#   ref: BRCOV
#   updated_at: 2026-09-28
#   layer: feature

@BRCOV
Feature: BranchCoverage — the tests take the branches the code has

  @BRCOV-B01 @unit-level
  Scenario: Below the floor fails naming the lines
    Given a file with 4 branches, 2 never taken on two lines
    When it is confronted with no floor, with a floor of 50 and with a floor of 60
    Then with no floor and at 60 it fails naming 50% and both lines, and at 50 it passes

  @BRCOV-B02 @unit-level
  Scenario: A waived branch is left out
    Given the same file with @no-branch and a reason above one line, and a bare @no-branch on the other
    When it is confronted
    Then only the unwaived line is missed

  @BRCOV-B03 @unit-level
  Scenario: A branch no test reaches is likely dead
    Given a missed branch on a line where the mutation found a mutant no test ran
    When it is confronted with a floor of 0, and again with the mutation measured at another revision
    Then the first fails as likely dead naming the line, and the second passes

  @BRCOV-B04 @unit-level
  Scenario: Nothing to measure is skipped or pending
    Given a spec, a code file with no coverage, one with stale coverage, and one whose coverage has no branch
    When each is confronted
    Then the spec and the branchless file are skipped, and the other two are pending
