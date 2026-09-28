# language: en
# @anchors
#   ref: SNGTU
#   updated_at: 2026-09-28
#   layer: feature

@SNGTU
Feature: SingleTestPerUnit — a unit has one test file per test layer

  @SNGTU-B01 @unit-level
  Scenario: Two test files of one unit in one layer fail
    Given a unit tested by a .test.ts and a .test.tsx in the unit layer, and one tested in two layers
    When each unit is confronted
    Then the first fails naming the layer and both files, and the second passes

  @SNGTU-B02 @unit-level
  Scenario: A declared split passes
    Given the same unit, with @split-test and a reason in one of its files
    When it is confronted
    Then it passes

  @SNGTU-B03 @unit-level
  Scenario: Not code, no test, no map
    Given a spec node, a unit no file tests, and no map
    When each is confronted
    Then the first two are skipped and the third is pending
