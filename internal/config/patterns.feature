# language: en
# @anchors
#   ref: DRPTD
#   updated_at: 2026-09-26
#   layer: feature

@DRPTD
Feature: DerivedPatterns — a derived file declared as one pattern or as a list of them

  @DRPTD-B01 @unit-level
  Scenario: A single pattern written as text becomes a list of one
    Given a derived block whose code file is the text "{{dir}}/{{name}}.ts"
    When the configuration is read
    Then the code patterns are a list holding only "{{dir}}/{{name}}.ts"

  @DRPTD-B02 @unit-level
  Scenario: A list of patterns is kept whole and in order
    Given a derived block whose test file is the list "{{dir}}/{{name}}.test.ts", "__tests__/{{name}}.test.ts"
    When the configuration is read
    Then the test patterns are those two, in that order

  @DRPTD-B03 @unit-level
  Scenario: An empty list is refused, naming the cause
    Given a derived file kind declared as an empty list
    When the configuration is read
    Then reading fails with "empty pattern list"

  @DRPTD-B04 @unit-level
  Scenario: Any shape other than text or a list of text is refused
    Given a derived file kind declared as a mapping, and another as a list of mappings
    When the configuration is read
    Then reading fails for both

  @DRPTD-B05 @unit-level
  Scenario: Written back, one pattern is text and several are a list
    Given the code patterns "x.ts" and the test patterns "a.test.ts", "b.test.ts"
    When they are written back
    Then the output is "code: x.ts" followed by the two test patterns as a list

  @DRPTD-I01 @unit-level
  Scenario: What is written back reads back as the same patterns
    Given one pattern, and three patterns
    When each is written back and read again
    Then the same patterns come back, in the same order

  @DRPTD-X01 @unit-level
  Scenario: The patterns are kept as written, neither expanded nor checked as globs
    Given the list "src/[bad" and " spaced/*.ts "
    When the configuration is read
    Then reading succeeds and both come back exactly as written
