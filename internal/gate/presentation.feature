# language: en
# @anchors
#   code: PRFTB
#   ref: PRSNT
#   updated_at: 2026-10-03
#   layer: feature

@PRSNT
Feature: PresentationGates — the presentation validations, confronted

  @PRSNT-B01 @unit-level
  Scenario: The rows are read by position, and a spec with none is skipped
    Given a spec with presentation rows and a TODO row, and a spec with none
    When the rows are read and the four gates confront the second
    Then the filled rows come with their four cells, and every gate skips the spec with none

  @PRSNT-B02 @unit-level
  Scenario: Every declared value has an appearance
    Given a status prop with three declared values and rows for two, a role state covered by an otherwise row, a plan state whose table header is not a value, and a one-value prop
    When presentation-exhaustive confronts the spec
    Then it fails naming the status value no row names, and not the role

  @PRSNT-B03 @unit-level
  Scenario: One prop and one condition lead to one appearance
    Given two rows with the same prop and condition and different appearances, and two with the same appearance
    When presentation-conflict confronts the spec
    Then it fails naming both rules

  @PRSNT-B04 @unit-level
  Scenario: The text shown is a message code, not copy
    Given a row showing a quoted text with no message code, one citing a message code, and one with a single-quoted value
    When presentation-copy-single-source confronts the spec
    Then only the first is named

  @PRSNT-B05 @unit-level
  Scenario: What changes is something a test can point at
    Given rows citing a declared test identifier and a row citing none
    When presentation-observable confronts the spec
    Then only the row citing none is named
