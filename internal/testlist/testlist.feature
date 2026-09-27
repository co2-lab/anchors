# language: en
# @anchors
#   ref: TSTLS
#   updated_at: 2026-09-27
#   layer: feature

@TSTLS
Feature: TestList — the project's tests, read the way the project says they are written

  @TSTLS-B01 @unit-level
  Scenario: A pattern reads each call followed by a literal as a test
    Given a test file with two calls the pattern opens, each followed by a quoted title
    When the tests are listed through the pattern
    Then each is a test with its title and the line where its call starts

  @TSTLS-B02 @unit-level
  Scenario: Titles may be quoted three ways, with escapes
    Given titles quoted with single quotes, double quotes and backticks, one with an escaped quote and one backtick title over two lines
    When the tests are listed through the pattern
    Then each title is read whole, the escape resolved, and a quoted title that breaks the line is left out

  @TSTLS-B03 @unit-level
  Scenario: A call without a literal title is left out
    Given a call the pattern opens followed by a variable
    When the tests are listed through the pattern
    Then that call is not a test

  @TSTLS-B04 @unit-level
  Scenario: The pattern's tests come in file order then position
    Given two files given out of order, each with two tests
    When the tests are listed through the pattern
    Then they come file by file, and in each file in the order they appear

  @TSTLS-B05 @unit-level
  Scenario: A script's contract output is the list of tests
    Given a script that prints two tests under the contract
    When the tests are listed through the script
    Then they are the two tests, read at the project root

  @TSTLS-B06 @unit-level
  Scenario: Output outside the contract is refused naming what is wrong
    Given outputs that break the contract in each way it can be broken
    When each is parsed
    Then each is refused with a message naming what is wrong

  @TSTLS-B07 @unit-level
  Scenario: A script's file paths are normalised to the map's form
    Given a script that names a file with a redundant segment
    When the tests are listed through the script
    Then the file is read in its clean form

  @TSTLS-B08 @unit-level
  Scenario: A source that declares nothing lists nothing
    Given a source with neither a pattern nor a script
    When the tests are listed
    Then there is no test and no error

  @TSTLS-E01 @unit-level
  Scenario: A source with both a pattern and a script is refused
    Given a source that declares both
    When the tests are listed
    Then an error names the conflict

  @TSTLS-E02 @unit-level
  Scenario: A pattern that does not compile is refused
    Given a pattern that does not compile
    When the tests are listed through it
    Then an error names the compile error

  @TSTLS-E03 @unit-level
  Scenario: A failing script is an error naming why
    Given a script that writes a reason to stderr and exits with an error
    When the tests are listed through it
    Then an error names the command and the reason
