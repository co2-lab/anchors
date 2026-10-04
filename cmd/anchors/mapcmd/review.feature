# language: en
# @anchors
#   code: RVFTR
#   ref: RVCMR
#   updated_at: 2026-10-03
#   layer: feature

@RVCMR
Feature: ReviewCommand — record a review with who looked and what they found, or list what is to review

  @RVCMR-B01 @unit-level
  Scenario: Pending lists what is to review with its question
    Given a project whose rule-fulfilled gate declares a review, and a spec
    When anchors review --pending runs
    Then it lists the spec under rule-fulfilled with the review question and how to record it
    And with the spec reviewed it says there is nothing to review

  @RVCMR-B02 @unit-level
  Scenario: A review with no findings records who looked and opens no issue
    Given the spec to review
    When anchors review records it by human:ana with no findings
    Then the map holds the review by human:ana at the spec's revision, with no findings
    And no issue is written

  @RVCMR-B03 @unit-level
  Scenario: A review with findings records them and opens the issue
    Given the spec to review
    When anchors review records it by agent:openai/gpt with findings
    Then the map holds the review marked as having found something
    And the spec's issue for rule-fulfilled holds the report

  @RVCMR-B04 @unit-level
  Scenario: In manual mode the findings write no issue unless asked
    Given a project in manual mode
    When a review with findings is recorded, then again with --record-issues
    Then the first prints the report and writes no issue, and the second writes it

  @RVCMR-I01 @unit-level
  Scenario: A review record always names who reviewed
    Given the spec to review
    When a review is recorded with no --by
    Then it is refused saying what to write, and the map holds no review

  @RVCMR-X01 @unit-level
  Scenario: There is no waived
    Given the review command
    When its flags are listed
    Then none records a review as waived

  @RVCMR-E01 @unit-level
  Scenario: A wrong target, gate or no reviewer refuses the record
    Given the spec to review
    When a review is recorded for a target not in the map, and for a gate that declares no review
    Then each is refused and nothing is written
