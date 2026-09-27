# language: en
# @anchors
#   ref: TLEVT
#   updated_at: 2026-09-26
#   layer: feature

@TLEVT
Feature: TelemetryEvent — a decision event has a name from a closed vocabulary, a caller-stamped instant and attributes

  @TLEVT-B01 @unit-level
  Scenario: A new event keeps its name and attributes
    Given the name "claim.served" and the attributes candidates 14 and skipped 2
    When the event is created
    Then its name is "claim.served" and its attributes are candidates 14 and skipped 2

  @TLEVT-B02 @unit-level
  Scenario: An event without attributes carries an empty set
    Given no attributes
    When the event is created
    Then its attribute set exists and is empty

  @TLEVT-B03 @unit-level
  Scenario: The instant comes from the caller's clock
    Given a clock fixed at 2026-09-14 18:30 UTC
    When the event is created
    Then its instant is 2026-09-14 18:30 UTC

  @TLEVT-I01 @unit-level
  Scenario: The vocabulary is the five declared names
    Given the declared event names
    When their wire values are read
    Then they are exactly "claim.served", "claim.empty", "check.finished", "escalate.raised" and "turn.ended"

  @TLEVT-X01 @unit-level
  Scenario: The system clock is never read
    Given a clock fixed at a date years in the past
    When the event is created
    Then its instant is that past date, not the current time
