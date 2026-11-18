# language: en
# @anchors
#   code: EMFXM
#   ref: EXMCH
#   updated_at: 2026-10-03
#   layer: feature

@EXMCH
Feature: ExamplesMatch — every Examples row is a case its tests run

  @EXMCH-B01 @unit-level
  Scenario: A row the tests do not run fails
    Given an outline with three rows, two carried by two tests citing its code and one not
    When the feature is confronted
    Then it fails naming the code, the third row's line and its missing value

  @EXMCH-B02 @unit-level
  Scenario: A value is a whole token with its case
    Given a value inside a longer word, one with another case, one quoted in the cell, and an empty cell
    When each is looked for
    Then only the whole token with its case is found, the quotes are dropped and the empty cell asks nothing

  @EXMCH-B03 @unit-level
  Scenario: A label column is display text
    Given an outline with a label column in English and one in Portuguese beside the key column
    When the feature is confronted
    Then only the key column is looked for

  @EXMCH-B04 @unit-level
  Scenario: The rows are the coded outline's tables
    Given a Portuguese outline with two examples tables, a step after a table, and an uncoded scenario with a table
    When the tables are read
    Then the rows under both headers belong to the coded outline and the uncoded one has none

  @EXMCH-B05 @unit-level
  Scenario: The test's body is read as the assertion gate reads it
    Given a labelled test whose values are in the function around it
    When it is confronted with labels on and off
    Then with labels it passes and without them it fails

  @EXMCH-B06 @unit-level
  Scenario: Nothing to confront is skipped
    Given a spec, a feature without examples, a project with no tests source, and a feature no test cites
    When each is confronted
    Then each is skipped

  @EXMCH-E01 @unit-level
  Scenario: Tests that cannot be listed fail the gate with the reason
    Given a tests script that fails
    When a feature with examples is confronted
    Then the gate fails naming why
