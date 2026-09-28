# language: en
# @anchors
#   ref: TSRCH
#   updated_at: 2026-09-28
#   layer: feature

@TSRCH
Feature: TestReach — a test reaches the unit it says it tests

  @TSRCH-B01 @unit-level
  Scenario: A test that names what its unit defines exercises it
    Given a unit defining a function, a test calling it, and a test that calls nothing of it
    When each test is confronted
    Then the first passes and the second fails naming the unit

  @TSRCH-B02 @unit-level
  Scenario: An import of the unit's module reaches it
    Given tests importing the unit by file name and by directory, and one naming it only in a string
    When each is confronted
    Then the imports reach the unit and the string does not

  @TSRCH-B03 @unit-level
  Scenario: A test that defines the unit's names exercises a copy
    Given a test that imports the unit and also defines two of its functions
    When it is confronted
    Then it fails naming the unit and both names

  @TSRCH-B04 @unit-level
  Scenario: Short names are anyone's
    Given a unit defining a two-letter function and a test defining the same
    When the test is confronted
    Then the name neither reaches nor copies the unit

  @TSRCH-B05 @unit-level
  Scenario: The ref's unit must be reached
    Given a test whose ref names a spec governing two files, reaching one, and a test reaching none
    When each is confronted
    Then the first passes and the second fails naming the ref and the files

  @TSRCH-B06 @unit-level
  Scenario: A declared invocation reaches the unit it names
    Given an invocation pattern on the gate and tests invoking the unit by file name, by directory, and another name
    When each is confronted
    Then the first two reach the unit and the third does not

  @TSRCH-B07 @unit-level
  Scenario: Nothing to confront is skipped
    Given a spec node, a support test, a test no unit pairs with, a project with no definition, a test with no ref, and a ref governing no code
    When each is confronted
    Then each is skipped
