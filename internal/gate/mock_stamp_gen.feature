# language: en
# @anchors
#   ref: MKSTP
#   updated_at: 2026-09-23
#   layer: feature

@MKSTP
Feature: MockStampGenerator — writes the missing `@contract` stamps, and never rewrites one

  @MKSTP-B01 @unit-level
  Scenario: A generated stamp passes the gate, and a change to its snippet fails it
    Given a test doubling `useBalance` with no stamp
    When the generator stamps it
    Then the gate passes
    And after `useBalance` changes, the gate fails

  @MKSTP-B02 @unit-level
  Scenario: An ambiguous specifier resolves to the test's workspace, or is skipped
    Given the same module path in two workspaces
    When a test in one of them is stamped
    Then the stamp points at its own workspace's module
    And a test outside both is skipped

  @MKSTP-B03 @unit-level
  Scenario: A stamp per factory key covers only that export
    Given a double whose factory names `useBalance`
    When another export of the module changes
    Then the stamp still passes

  @MKSTP-B04 @unit-level
  Scenario: An automock gets one stamp over the whole module
    Given a double with no factory
    When it is stamped
    Then one stamp covers the whole module

  @MKSTP-B05 @unit-level
  Scenario: A third-party double is not stamped
    Given a double of a library outside the map
    When the generator runs
    Then nothing is written

  @MKSTP-B06 @unit-level
  Scenario: The header stays out of a whole-module stamp
    Given a module that opens with an @anchors header
    When it is stamped whole and only its updated_at changes
    Then the stamp still passes the gate

  @MKSTP-B07 @unit-level
  Scenario: A key added to a stamped double gets a stamp
    Given a double already stamped for one export
    When a key naming another export is added to its factory and the generator runs
    Then only the new export gets a stamp, and the existing one is untouched

  @MKSTP-I01 @unit-level
  Scenario: An existing stamp is never rewritten
    Given a test whose stamp diverges from the module
    When the generator runs
    Then the test is left untouched

  @MKSTP-I02 @unit-level
  Scenario: A repeated line never anchors a stamp
    Given a module whose first line occurs twice
    When it is stamped
    Then the anchor is a unique line and the gate passes
