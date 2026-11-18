# language: en
# @anchors
#   code: TUFTS
#   ref: TSUNT
#   updated_at: 2026-10-03
#   layer: feature

@TSUNT
Feature: TestedUnits — which code a test tests, found by the project's own derivation

  @TSUNT-B01 @unit-level
  Scenario: With code as the anchor a code file tests through its own name
    Given code as the anchor and a test template beside the code
    When the tested units are found
    Then each test beside a code file is that file's

  @TSUNT-B02 @unit-level
  Scenario: With the spec as the anchor the code path gives the variables
    Given the spec as the anchor, a code template and a test template
    When the tested units are found
    Then the code file's name, recovered from the code template, finds its test

  @TSUNT-B03 @unit-level
  Scenario: An override of the code's layer places its test elsewhere
    Given a model layer whose override puts its tests under a separate tree
    When the tested units are found
    Then the model's test is found in that tree, and another layer's is found beside its code

  @TSUNT-B04 @unit-level
  Scenario: A glob test template matches the tests it covers, a code template is read literally
    Given a test template with a glob, a code template with a glob, and a code template under a bracketed directory
    When the tested units are found
    Then the glob test template finds the tests it matches, the glob code template matches no ordinary path, and the bracketed directory matches itself

  @TSUNT-B05 @unit-level
  Scenario: Only the map's tests are answered, each unit once, in order
    Given a derivation that points at a test the map lacks, two templates that find the same unit, and a project without a derivation
    When the tested units are found
    Then the missing test is not answered, the unit is listed once, and the project without a derivation answers nothing
