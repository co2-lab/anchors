# language: en
# @anchors
#   code: ATFFF
#   ref: AFTSF
#   updated_at: 2026-10-11
#   layer: feature

@AFTSF
Feature: AffectedTests — which tests to run for what changed

  @AFTSF-B01 @unit-level
  Scenario: A changed rule reaches its tests
    Given a spec whose rule B01 changed and whose title changed, and tests titled with B01 and B02
    When the affected tests are read
    Then only the B01 test is affected

  @AFTSF-B02 @unit-level
  Scenario: A changed region reaches its tests, and a change outside reaches the closure
    Given a code file with two regions, one changed, and then a change outside them
    When the affected tests are read
    Then the first reaches the tests of the changed region's code, and the second every test whose closure holds the file

  @AFTSF-B03 @unit-level
  Scenario: A changed test runs whole, and a frozen file reaches nothing
    Given a changed test file and a changed frozen code file
    When the affected tests are read
    Then every test of the test file is affected, and none for the frozen file

  @AFTSF-B04 @unit-level
  Scenario: The tests are printed for the runner
    Given affected tests with a regex metacharacter in a title
    When they are printed as lines, regex, files and json
    Then each format lists them in file and line order, the regex escaped and anchored

  @AFTSF-B05 @unit-level
  Scenario: With no tests source, the test files are printed
    Given a project with no dialect.tests and a test file run whole
    When the affected tests are printed
    Then the file is printed, and stderr says the tests by code are left out

  @AFTSF-E01 @unit-level
  Scenario: An unknown format is refused
    Given affected tests
    When they are printed as tap
    Then the command fails naming the four formats
