# language: en
# @anchors
#   ref: MNVRM
#   updated_at: 2026-09-26
#   layer: feature

@MNVRM
Feature: MinVersion — whether the running binary meets the minimum version the project declares

  @MNVRM-B01 @unit-level
  Scenario: Versions are ordered part by part as numbers
    Given the minimum 0.1.84
    When the running versions 0.1.85, 0.2.0, 0.1.83 and 0.1.9 are compared with it
    Then 0.1.85 and 0.2.0 meet it, and 0.1.83 and 0.1.9 do not

  @MNVRM-B02 @unit-level
  Scenario: The tag prefix v is accepted on either side
    Given the minimum v0.1.84 and the running version 0.1.84, and the reverse
    When they are compared
    Then the running version meets the minimum, with no error

  @MNVRM-B03 @unit-level
  Scenario: With no minimum declared, any running binary satisfies it
    Given no declared minimum
    When the running versions 0.0.1, dev, empty and "anything" are compared with it
    Then each meets it, with no error

  @MNVRM-B04 @unit-level
  Scenario: A running version that cannot be ordered does not satisfy, and says why
    Given the minimum 0.1.84
    When the running versions dev, empty, 0.1 and 1.2.3.4 are compared with it
    Then none meets it, and each comparison returns an error

  @MNVRM-B05 @unit-level
  Scenario: A declared minimum that is not MAJOR.MINOR.PATCH is refused, dev included
    Given a project whose configuration declares min_version dev, latest, 0.1, 0.1.84-rc1 or 0.1.x
    When the configuration is loaded
    Then the load fails naming min_version

  @MNVRM-B06 @unit-level
  Scenario: An absent or well-formed minimum is accepted
    Given the declared minimums empty, 0.1.84, v0.1.84 and 10.20.30
    When each is checked at load
    Then none is refused

  @MNVRM-B07 @unit-level
  Scenario: The refusal names the expected format, an example and what a bad value silences
    Given the declared minimum latest
    When it is checked at load
    Then the message contains "MAJOR.MINOR.PATCH", "0.1.84", the quoted value and "silences"

  @MNVRM-I01 @unit-level
  Scenario: Swapping the two versions always flips the order
    Given the versions 0.1.9, 0.1.84, v0.2.0, 1.0.0 and 0.99.99
    When every pair is compared in both directions
    Then the two answers are always opposite, and zero only for the same version

  @MNVRM-X01 @unit-level
  Scenario: A pre-release is never compared as if it were its final release
    Given the minimum 0.1.84
    When the running versions 0.1.84-rc1 and 0.1.85-rc1 are compared with it
    Then neither meets it, and each comparison returns an error
