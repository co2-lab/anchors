# language: en
# @anchors
#   code: SYNFT
#   ref: MGSYN
#   updated_at: 2026-10-16
#   layer: feature

@MGSYN
Feature: Declared syntax — what the core knew of a language, written into the project at format 8

  @MGSYN-B01 @unit-level
  Scenario: The comments of the examples covering the project's files are added
    Given a project with Go, TS, YAML and Zig files, a spec and a feature, declaring no comments
    When the syntax to declare is read
    Then the Go, TS and YAML entries are added, and .zig is named to declare by hand
    And a project declaring its comments gets none

  @MGSYN-B02 @unit-level
  Scenario: The family's lexicon fills what the project left empty
    Given a Go project declaring its own loop
    When the syntax to declare is read
    Then its loop is kept and every other field comes from the Go example
    And with no family only the naming conventions are added, and env-declared is named

  @MGSYN-B03 @unit-level
  Scenario: Inside the project's own tests only what the old binary filled is added
    Given a ts project declaring only its assertion, one with its own script, and a Go one with an untitled pattern
    When the syntax to declare is read
    Then the first gets the family's pattern and the block shape, the second only the block shape
    And the untitled pattern is named

  @MGSYN-B04 @unit-level
  Scenario: The gates' shapes are added for the gates that read them
    Given a Go project declaring route-exists and mock-typed with its own mock detection, a layer of YAML flows and one already headless
    When the syntax to declare is read
    Then the ts route pattern and the ts mock shapes are added, the declared detection kept
    And only the flows layer becomes headless

  @MGSYN-B05 @unit-level
  Scenario: The additions are written at the end of each block, the rest kept
    Given an anchors.yaml with comments, a tests block and a flow-style layer
    When the additions are written into it
    Then its comments stay, the tests block takes the block shape, the layer takes headless inline
    And the file reads back with nothing left to add

  @MGSYN-B06 @unit-level
  Scenario: A missing section is written whole, a flow-style one refused
    Given an anchors.yaml with no dialect nor derived, and one with a flow-style dialect
    When the additions are written into each
    Then the first gets both sections appended, and the second is refused naming the dialect

  @MGSYN-B07 @unit-level
  Scenario: The example runners the project's suites run with are added
    Given suites run with jest, maestro and stryker, and a maestro runner declared
    When the additions are written
    Then jest and stryker are appended to the project's runners, and maestro is kept as declared
    And a project with no monitor block gets one
