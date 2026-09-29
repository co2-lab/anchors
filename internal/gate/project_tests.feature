# language: en
# @anchors
#   ref: PRJTS
#   updated_at: 2026-09-29
#   layer: feature

@PRJTS
Feature: ProjectTests — the gates read the project's tests through the source the project declares

  @PRJTS-B01 @unit-level
  Scenario: The project's own source wins over its family's, and neither declares nothing
    Given a ts project, a ts project with its own script, a project with no family, and no configuration
    When the tests source is resolved
    Then the ts pattern, the script, and nothing twice

  @PRJTS-B02 @unit-level
  Scenario: A pattern reads only the map's test files
    Given a test in a file the map lists as a test, and one in a file the map does not list
    When the project's tests are read through a pattern
    Then only the listed file's test is read

  @PRJTS-B03 @unit-level
  Scenario: The tests are read once per scan
    Given a script that counts its runs
    When the tests are read twice with the same root, map and configuration, then with a new map
    Then the script ran twice

  @PRJTS-B04 @unit-level
  Scenario: The tests of a set of files come in file order then line
    Given tests of three files out of order
    When the tests of two of the files are kept, in a given order
    Then they come in that file order, each file by line, without the third file

  @PRJTS-B05 @unit-level
  Scenario: A source error is returned with the declaration
    Given a script that fails
    When the project's tests are read
    Then the error comes back and the source counts as declared

  @PRJTS-B06 @unit-level
  Scenario: Support files are not read as tests
    Given a map with a test file and a support file, each with a test call
    When the project's tests are read and the test paths are filtered
    Then only the test file's test is read, and the support file leaves the paths

  @PRJTS-B07 @unit-level
  Scenario: Under --index the tests are listed from the commit
    Given a test file staged short and edited in the tree with a test far below
    When the project's tests are listed under --index
    Then only the staged file's tests are listed, at its lines
