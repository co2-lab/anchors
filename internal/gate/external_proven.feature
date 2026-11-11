# language: en
# @anchors
#   code: EXPRF
#   ref: EXPRS
#   updated_at: 2026-11-10
#   layer: feature

@EXPRS
Feature: ExternalProven — a file that calls a real system is proven by a test that reached it

  @EXPRS-B01 @unit-level
  Scenario: A file calling a real system with no test passed in a suite that reaches it fails
    Given a code file calling a declared system, whose test passed only in the unit suite
    When the gate confronts it, and again after the test passed in the integration suite
    Then the first fails naming the system and its layers, and the second passes

  @EXPRS-B02 @unit-level
  Scenario: Outside its subject the gate is skipped
    Given a project with no system declared, a test file, a code file that calls none, and one that waives the proof
    When the gate confronts each
    Then each is skipped

  @EXPRS-E01 @unit-level
  Scenario: Without a map the gate is pending
    Given a code file calling a declared system and no built map
    When the gate confronts it
    Then it is pending
