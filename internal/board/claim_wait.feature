# language: en
# @anchors
#   code: CWFCL
#   ref: CLWTC
#   updated_at: 2026-10-03
#   layer: feature

@CLWTC
Feature: ClaimWait — asking the claim pipeline for a card and waiting, bounded, for the answer

  @CLWTC-B01 @unit-level
  Scenario: Another agent's pending run does not count
    Given a queued claim run titled for "other/session"
    When "machine/session-1" asks and waits
    Then it dispatches its own claim once

  @CLWTC-B02 @unit-level
  Scenario: No dispatch while the agent's run is pending
    Given a claim run of "machine/session-1" in progress that later hands it the card
    When it asks and waits
    Then nothing is dispatched and the card is returned

  @CLWTC-B03 @unit-level
  Scenario: An older finished run is not the answer
    Given a finished claim run of "machine/session-1" from a previous request
    When it asks and waits and the card arrives on the third look
    Then the card is returned instead of ending on the old run

  @CLWTC-B04 @unit-level
  Scenario: The card is returned as soon as it arrives
    Given a board where the card becomes the agent's on the third look
    When the agent asks and waits
    Then card 4 is returned with exactly one dispatch

  @CLWTC-B05 @unit-level
  Scenario: A finished run ends the wait after one more look
    Given a claim run that finishes at once, with the card appearing on the second look, and another with no card at all
    When the agent asks and waits for each
    Then the first returns the card, and the second ends with the finished run and no card, having waited no time

  @CLWTC-B06 @unit-level
  Scenario: A cancelled run is asked again, bounded
    Given every claim run ends cancelled
    When the agent asks and waits
    Then the claim is dispatched 3 times and the last cancelled run is reported

  @CLWTC-B07 @unit-level
  Scenario: The wait is bounded and does not dispatch twice
    Given a claim run that never finishes, with the default bounds
    When the agent asks and waits
    Then after three minutes of five-second looks the outcome is timed out, names the run, and one claim was dispatched

  @CLWTC-E01 @unit-level
  Scenario: An unreadable run list fails before dispatching
    Given a run list that is not JSON
    When the agent asks and waits
    Then an error is returned and nothing was dispatched
