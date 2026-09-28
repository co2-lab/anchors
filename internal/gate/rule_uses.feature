# language: en
# @anchors
#   ref: RLUSG
#   updated_at: 2026-09-28
#   layer: feature

@RLUSG
Feature: RuleUses — each rule says what it uses, and what it uses exists

  @RLUSG-B01 @unit-level
  Scenario: The sections are read in any language and by position
    Given a spec with a Portuguese Validations table, an English Rule uses table, a section the project renamed, and a TODO row
    When its rule uses are read
    Then each filled row gives its rule and uses, whatever the columns are called, and the TODO row gives nothing

  @RLUSG-B02 @unit-level
  Scenario: A uses cell lists backticked or comma-separated items
    Given a cell with backticked items and a cell with plain comma-separated items
    When their items are read
    Then the first gives the backticked items and the second the comma-separated ones

  @RLUSG-B03 @unit-level
  Scenario: A rule that does not say what it uses fails
    Given a spec whose behaviour and error have no row, whose question has none, and whose other behaviour is waived
    When rule-uses-declared confronts it, with and without letters declared
    Then it fails naming the behaviour and the error, and with letters E only it names the error

  @RLUSG-B04 @unit-level
  Scenario: Nothing to ask about is a skip
    Given a feature node, and a spec with only an open question
    When rule-uses-declared confronts them
    Then both are skipped

  @RLUSG-B05 @unit-level
  Scenario: What a rule uses must exist in the spec
    Given a spec whose rules use a contract field, a dotted field, a data-state name, a dependency, a code, and a field and a dependency it never declares
    When rule-uses-resolve confronts it
    Then it fails naming only the undeclared field and the missing dependency, each with its rule

  @RLUSG-B06 @unit-level
  Scenario: A spec whose rules say nothing yet is a skip
    Given a feature node, and a spec with no rule uses
    When rule-uses-resolve confronts them
    Then both are skipped
