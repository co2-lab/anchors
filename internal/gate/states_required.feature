# language: en
# @anchors
#   code: STRQF
#   ref: STRQS
#   updated_at: 2026-11-13
#   layer: feature

@STRQS
Feature: StatesRequired — every spec of a layer registers the states the project requires of it

  @STRQS-B01 @unit-level
  Scenario: A required state the spec neither registers nor waives fails
    Given screens must register loading, error, empty and offline, and a screen spec registering loading and error, waiving offline with a reason and naming empty only in prose
    When the gate confronts it, and again with an empty state registered
    Then the first fails naming empty, and the second passes

  @STRQS-B02 @unit-level
  Scenario: Outside its subject the gate is skipped
    Given a project requiring no state, a spec of a layer with none required, and a code file
    When the gate confronts each
    Then each is skipped

  @STRQS-B03 @unit-level
  Scenario: A state is registered where it is defined, by its name
    Given a screen spec with no empty state, whose message and data rows cite another state and say "empty"
    When states-required confronts it, and again with a row, a list item or a heading that defines the empty state
    Then the first fails naming empty, the definitions pass, and a definition whose name is another word fails
