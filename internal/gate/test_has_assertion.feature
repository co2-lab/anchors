# language: en
# @anchors
#   ref: THSAS
#   updated_at: 2026-09-28
#   layer: feature

@THSAS
Feature: TestHasAssertion — every test asserts something in its body

  @THSAS-B01 @unit-level
  Scenario: A test with no assertion fails, named by its line and title
    Given a test file with one test that asserts and one that only calls the code
    When the file is confronted
    Then it fails naming the second test's line and title, and a file whose tests all assert passes

  @THSAS-B02 @unit-level
  Scenario: The body is the block the test's line opens
    Given a bracket block closed at its opener's indentation, an indented block, and a one-line call
    When the end of each is read
    Then the bracket block ends on its closer, the indented one before the line back at its indentation, and the call on its own line

  @THSAS-B03 @unit-level
  Scenario: A multi-line literal is text, not layout
    Given a test whose body holds a backtick string and a triple-quoted string with lines at column zero, and a quoted lone backtick
    When its body is read
    Then the block runs past the literals to its real closer

  @THSAS-B04 @unit-level
  Scenario: An empty test is a label only when the gate says so
    Given a function opening with an empty labelled test and asserting below it
    When it is confronted with labels on and with labels off
    Then with labels it passes, and without them the empty test fails

  @THSAS-B05 @unit-level
  Scenario: The source's end bounds the body
    Given a script that ends a test before the line that asserts
    When the file is confronted
    Then the test fails

  @THSAS-B06 @unit-level
  Scenario: Nothing to measure is skipped
    Given a spec node, a support test, a project with no tests source, one with no assertion, and a file with no listed test
    When each is confronted
    Then each is skipped

  @THSAS-E01 @unit-level
  Scenario: Tests that cannot be listed fail the gate with the reason
    Given a tests script that fails, and an assertion that does not compile
    When a test file is confronted
    Then the gate fails naming why
