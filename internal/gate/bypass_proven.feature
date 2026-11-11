# language: en
# @anchors
#   code: BYPRF
#   ref: BYPRS
#   updated_at: 2026-11-10
#   layer: feature

@BYPRS
Feature: BypassProven — a path the code skips while a test runs is run by a suite that reaches it

  @BYPRS-B01 @unit-level
  Scenario: A bypass with no test passed in a suite that runs its path fails
    Given a hook that returns early under the end-to-end flag, whose test passed only in the unit suite
    When the gate confronts it, and again after the test passed in the integration suite
    Then the first fails naming the bypass's line and the layer, and the second passes

  @BYPRS-B02 @unit-level
  Scenario: Outside its subject the gate is skipped
    Given a project with no bypass declared, a test file, and a hook whose bypass is waived on the line above
    When the gate confronts each
    Then each is skipped

  @BYPRS-E01 @unit-level
  Scenario: Without a map the gate is pending
    Given a hook with a bypass and no built map
    When the gate confronts it
    Then it is pending
