# language: en
# @anchors
#   code: RTFRT
#   ref: RTSRT
#   layer: feature

@RTSRT
Feature: TestsRetried — a test that passed only on a retry under load stays in sight

  @RTSRT-B01 @unit-level
  Scenario: A test file with no retry passes
    Given a record of retries of other files only
    When the gate confronts a test file
    Then it passes

  @RTSRT-B02 @unit-level
  Scenario: A test file with tests passed on retry is listed
    Given a test file with one test retried once and another three times this week, one retried 8 days ago, and a retry of another file
    When the gate confronts it
    Then it fails listing the two tests with their counts and last load, the second as unstable, and nothing of the old retry or the other file
