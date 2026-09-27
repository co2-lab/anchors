# language: en
# @anchors
#   ref: TLCNT
#   updated_at: 2026-09-26
#   layer: feature

@TLCNT
Feature: TelemetryConfig — the opt-out is easy to find, easy to get right, and the environment overrides the file

  @TLCNT-B01 @unit-level
  Scenario: Every form someone would try turns telemetry off
    Given the environment variable set to "off", "OFF", "Off", "0", "false", "FALSE", "no" or " off "
    When the opt-out is evaluated
    Then telemetry is off for each of them

  @TLCNT-B02 @unit-level
  Scenario: With nothing declared telemetry is on
    Given the environment variable blank and no value in the file
    When the opt-out is evaluated
    Then telemetry is on

  @TLCNT-B03 @unit-level
  Scenario: The environment overrides the file in both directions
    Given the environment says "on" and the file says "off"
    When the opt-out is evaluated
    Then telemetry is on
    And with the environment saying "off" and the file saying "on", telemetry is off

  @TLCNT-B04 @unit-level
  Scenario: Without the environment the file decides
    Given the environment variable blank
    When the file says "off"
    Then telemetry is off
    And when the file says "on", telemetry is on

  @TLCNT-I01 @unit-level
  Scenario: A value means the same in the environment and in the file
    Given each of the off-words and the word "on"
    When each is evaluated once through the environment and once through the file
    Then both answers are the same for every word

  @TLCNT-X01 @unit-level
  Scenario: A word outside the closed list keeps telemetry on
    Given the environment or the file set to "yes", "disabled" or "of"
    When the opt-out is evaluated
    Then telemetry stays on
