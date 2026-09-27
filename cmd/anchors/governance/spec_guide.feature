# language: en
# @anchors
#   ref: SPGDS
#   updated_at: 2026-09-26
#   layer: feature

@SPGDS
Feature: SpecGuide — the project's own spec guide, instantiated with its dialect and a complete example

  @SPGDS-B01 @unit-level
  Scenario: An empty example code defaults to LOGI
    Given a configuration and an empty example code
    When the spec guide is rendered
    Then the example carries "code: LOGI" and "### LOGI-B01"

  @SPGDS-B02 @unit-level
  Scenario: The guide shows the three rule forms and a complete example
    Given the example code LOGIX
    When the spec guide is rendered
    Then it contains "### LOGIX-B01", "| `LOGIX-B02` |" and "- **LOGIX-B03**"
    And a complete example opening with "<!-- @anchors" and "code: LOGIX"

  @SPGDS-B03 @unit-level
  Scenario: The section titles come from the catalogue the generator uses
    Given a project with no language declared
    When the spec guide is rendered
    Then the example's overview, rules and open decisions headings are the catalogue's titles, such as "## Open Decisions"

  @SPGDS-B04 @unit-level
  Scenario: The rule letters are the project's when declared and the canonical ones otherwise
    Given a project that declares the rule type I as Invariant, and a project that declares none
    When the spec guide is rendered for each
    Then the first lists "| `I` | Invariant |"
    And the second lists the canonical letters, such as "`B` behaviour"

  @SPGDS-B05 @unit-level
  Scenario: The code length is stated only when the project declares it
    Given a project that declares code lengths 4 and 5, and a project that declares none
    When the spec guide is rendered for each
    Then the first says "The identity code has 4 or 5 character(s) in this project"
    And the second does not mention code_lengths

  @SPGDS-B06 @unit-level
  Scenario: The guide starts from the command that generates the skeleton
    Given any configuration
    When the spec guide is rendered
    Then "anchors new spec" and "--list-sections" appear before the format section

  @SPGDS-X01 @unit-level
  Scenario: A project that declares its rule types is not offered the canonical letters
    Given a project that declares its own rule types
    When the spec guide is rendered
    Then it does not mention the canonical letters
