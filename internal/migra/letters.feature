# language: en
# @anchors
#   ref: MGLTR
#   updated_at: 2026-09-28
#   layer: feature

@MGLTR
Feature: CodeLetters — rewriting the letter of the codes of one kind of unit

  @MGLTR-B01 @unit-level
  Scenario: The renames of the steps crossed
    Given the registered steps
    When the renames from format 4 to 5, and from 5 to 5, are gathered
    Then the first are format 5's three renames and the second none

  @MGLTR-B02 @unit-level
  Scenario: A phase, a step and a result get their new letters
    Given a text citing a plan's phase, a flow's step, an action's result and a flow's result
    When the codes are rewritten
    Then they read W, T, O and O with the same unit and number

  @MGLTR-B03 @unit-level
  Scenario: Every other code stays
    Given a text with a spec's permission, a plan's revision block, an unknown unit and a plan's other letter
    When the codes are rewritten
    Then none of them changes

  @MGLTR-B04 @unit-level
  Scenario: The rewrite says what it replaced
    Given a text citing the same phase twice and a step once
    When the codes are rewritten
    Then the report counts two and one, listed in order

  @MGLTR-B05 @unit-level
  Scenario: Rewriting twice changes nothing more
    Given a text already rewritten
    When it is rewritten again
    Then it is the same and nothing is reported
