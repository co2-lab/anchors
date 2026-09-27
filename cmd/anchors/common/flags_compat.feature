# language: en
# @anchors
#   ref: FLALF
#   updated_at: 2026-09-26
#   layer: feature

@FLALF
Feature: FlagAliases — a renamed command flag keeps answering to its old name

  @FLALF-B01 @unit-level
  Scenario: A value passed under the old name reaches the current flag
    Given a command with the flag "--root" and the old name "--raiz" declared as its alias
    And the command is called with "--raiz /tmp/project"
    When the aliases are resolved
    Then "--root" holds "/tmp/project"

  @FLALF-B02 @unit-level
  Scenario: The alias of a switch flag is a switch
    Given a command with the switch "--all" and the old name "--todos" declared as its alias
    And the command is called with "--todos" alone
    When the aliases are resolved
    Then the old flag is of switch type and "--all" is on

  @FLALF-B03 @unit-level
  Scenario: The old name is hidden and deprecated
    Given a command with the flag "--root" and the old name "--raiz" declared as its alias
    When the old flag is inspected
    Then it is hidden and its deprecation note mentions "--root"

  @FLALF-B04 @unit-level
  Scenario: The current name wins when both are passed
    Given a command called with "--raiz /old --root /new"
    When the aliases are resolved
    Then "--root" holds "/new"

  @FLALF-B05 @unit-level
  Scenario: A pair naming flags the command does not have is ignored
    Given a pair "missing" / "gone" that names no declared flag
    When the aliases are resolved
    Then no error is returned

  @FLALF-X01 @unit-level
  Scenario: Nothing is copied before the aliases are resolved
    Given a command called with "--raiz /tmp/project"
    When the flags have been parsed but the aliases not yet resolved
    Then "--root" still holds its default "."

  @FLALF-E01 @unit-level
  Scenario: A value the current flag cannot hold names the old flag
    Given a number flag "--limit" whose old name "--limite" was called with "many"
    When the aliases are resolved
    Then the error names "--limite"

  @FLALF-E02 @unit-level
  Scenario: An alias of a flag the command does not have stops the program
    Given a command without a flag named "nope"
    When an alias "nao" is declared for "nope"
    Then the program stops with a message naming "nope"
