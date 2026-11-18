# language: en
# @anchors
#   code: PSRFA
#   ref: PSRSA
#   updated_at: 2026-11-11
#   layer: feature

@PSRSA
Feature: PushReview — the evidence each push carries: a review of its commits and the judgments they reach

  @PSRSA-B01 @unit-level
  Scenario: The section declares the review, the judgments, the file and the retention
    Given a configuration requiring a review and the judgments, with a file and a retention, one with only the review, and one with only the gates
    When each configuration is loaded
    Then the first holds them with its retention in days, the second writes to REVIEWS.md, keeps every review and asks no judgment, and the third confronts the gates alone

  @PSRSA-E01 @unit-level
  Scenario: A retention or a file that cannot hold refuses the configuration
    Given a retention that is not days, an absolute file and a file outside the project
    When each configuration is loaded
    Then each is refused naming the field
