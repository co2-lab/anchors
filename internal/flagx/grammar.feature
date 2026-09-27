# language: en
# @anchors
#   ref: FLGRF
#   updated_at: 2026-09-26
#   layer: feature

@FLGRF
Feature: FlagGrammar — the fixed grammar of a feature-flag scenario's condition

  @FLGRF-B01 @unit-level
  Scenario: The three house styles say the same thing
    Given the conditions ">= 50", "gte 50", "NUM_GTE 50", "greaterThanOrEqual 50" and "Greater Than Inclusive 50"
    When each is parsed
    Then every one is the operator greater-or-equal

  @FLGRF-B02 @unit-level
  Scenario: The operand loses one layer of quotes
    Given the conditions `= "on"`, `= 'on'`, `>= 50` and `= "a \"b\" c"`
    When each is parsed
    Then the operands are "on", "on", "50" and `a \"b\" c`

  @FLGRF-B03 @unit-level
  Scenario: The absent and present cases compare against nothing, in any language
    Given the conditions "absent", "is not set", "ABSENT", "ausente", "present", "is set" and "presente"
    When each is parsed
    Then the first four are absent and the last three present, all with no operand

  @FLGRF-B04 @unit-level
  Scenario: The longest operator wins
    Given the conditions ">= 50", `starts with "beta"` and `not contains "beta"`
    When each is parsed
    Then they are greater-or-equal with "50", starts-with with "beta", and not-contains

  @FLGRF-B05 @unit-level
  Scenario: Only absent and present need no operand
    Given the operators absent, present, eq, gte, contains, matches and rollout
    When each is asked whether it needs an operand
    Then absent and present say no and the others yes

  @FLGRF-X01 @unit-level
  Scenario: Service-dependent operators and prose are refused
    Given the conditions "segmentMatch beta-users", "modulo 3" and "when the user is a beta tester"
    When each is parsed
    Then each is refused

  @FLGRF-E01 @unit-level
  Scenario: An empty condition is refused
    Given the conditions "" and "   "
    When each is parsed
    Then each is refused

  @FLGRF-E02 @unit-level
  Scenario: An operator without an operand is refused
    Given the conditions ">=" and "contains"
    When each is parsed
    Then each is refused

  @FLGRF-E03 @unit-level
  Scenario: A loose word is refused with the accepted forms
    Given the condition "maybe"
    When it is parsed
    Then it is refused with a message quoting "maybe" and showing the absent form
