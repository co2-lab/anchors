# language: en
# @anchors
#   ref: GHRNG
#   updated_at: 2026-09-26
#   layer: feature

@GHRNG
Feature: GhRunner — running `gh` for the board so that a failure names its cause and never waits for input

  @GHRNG-B01 @unit-level
  Scenario: gh gets no input
    Given a gh that exits 7 when its input ends and echoes what it reads otherwise
    When it is run
    Then it fails at once with exit status 7

  @GHRNG-B02 @unit-level
  Scenario: Exit code 4 is explained as not authenticated
    Given a gh that exits with code 4
    When it is run for "issue list"
    Then the error says "NOT AUTHENTICATED", "gh auth login" and "interactive"

  @GHRNG-B03 @unit-level
  Scenario: Other failures get no authentication hint
    Given commands exiting with 1, 2, 3, 5 and 127, and a command that does not exist
    When their failures are explained
    Then none gets the authentication hint
