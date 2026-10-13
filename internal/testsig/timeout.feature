# language: en
# @anchors
#   code: TFFTM
#   ref: TFSTM
#   layer: feature

@TFSTM
Feature: TimeoutFailure — telling a test that ran out of time from one that failed its assertion

  @TFSTM-B01 @unit-level
  Scenario: A failure past its timeout ran out of time
    Given the failure MIF recorded of a hook past 20 seconds
    When it is told
    Then it ran out of time

  @TFSTM-B02 @unit-level
  Scenario: A wait that ran out is told by its line
    Given the failure MIF recorded of a waitFor that rethrew its assertion, and the test file with that line
    When it is told
    Then it ran out of time

  @TFSTM-B03 @unit-level
  Scenario: A failed assertion is not a timeout
    Given an assertion failing outside a wait, a stack pointing at another file's wait, and a case that passed
    When each is told
    Then none ran out of time

  @TFSTM-B04 @unit-level
  Scenario: The tests that passed on a retry are kept
    Given a project with no record
    When two entries are added, then one more
    Then the three are read back in order, only the latest 500 are kept, and a project with no record, or an unreadable one, reads none

  @TFSTM-B05 @unit-level
  Scenario: An unreadable test file counts for nothing
    Given a failure whose stack points at a test file that does not exist
    When it is told
    Then it did not run out of time
