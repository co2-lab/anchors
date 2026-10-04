# language: en
# @anchors
#   code: VCFTV
#   ref: GTVCG
#   updated_at: 2026-10-03
#   layer: feature

@GTVCG
Feature: GateVocabulary — the list of default gate names, injected into the configuration layer

  @GTVCG-B01 @unit-level
  Scenario: With no source registered, the default gate names are absent, not a failure
    Given no source of default gate names is registered
    When the default gate names are asked for
    Then the answer is empty, and nothing fails

  @GTVCG-B02 @unit-level
  Scenario: With a source registered, the default gate names are the ones it gives
    Given a registered source giving "spec-complete" and "rule-fulfilled"
    When the default gate names are asked for
    Then they are "spec-complete" and "rule-fulfilled"

  @GTVCG-B03 @unit-level
  Scenario: The letters of plans, flows and actions
    Given the canonical rule letters
    When the phase, step and outcome letters are read
    Then they are W, T and O, only the phase letter is canonical, and step and outcome are no spec letter

  @GTVCG-I01 @unit-level
  Scenario: The answer always comes from the source registered last
    Given a source giving "first" registered, and then a source giving "second"
    When the default gate names are asked for
    Then they are "second"

  @GTVCG-X01 @unit-level
  Scenario: The names are not cached: each question asks the registered source again
    Given a registered source that counts how often it is asked
    When the default gate names are asked for twice
    Then the source was asked twice
