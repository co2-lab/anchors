# language: en
# @anchors
#   code: TSFTT
#   ref: TSTLS
#   updated_at: 2026-10-16
#   layer: feature

@TSTLS
Feature: TestList — the project's tests, read the way the project says they are written

  @TSTLS-B01 @unit-level
  Scenario: A pattern reads each match as a test, titled by its title group
    Given a test file with two calls the pattern matches, each title in its title group
    When the tests are listed through the pattern
    Then each is a test with its title and the line where its match starts

  @TSTLS-B02 @unit-level
  Scenario: A pattern names a title group for each way a title is quoted
    Given a pattern with title groups for single quotes, double quotes and backticks, and a title with an escaped quote
    When the tests are listed through the pattern
    Then each title is the group that took part, read as captured, the escape kept

  @TSTLS-B03 @unit-level
  Scenario: A match whose title groups captured nothing is left out
    Given a call the pattern matches whose title is a variable, which no title group captures
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

  @TSTLS-E05 @unit-level
  Scenario: A pattern with no title group is refused
    Given a pattern that names no title group
    When the tests are listed through it
    Then an error says where to name the title

  @TSTLS-E03 @unit-level
  Scenario: A failing script is an error naming why
    Given a script that writes a reason to stderr and exits with an error
    When the tests are listed through it
    Then an error names the command and the reason

  @TSTLS-B09 @unit-level
  Scenario: The script may say where a test ends
    Given a script printing one test with an end and one without, and a pattern reading a file
    When the tests are listed
    Then the first carries its end and the others carry none

  @TSTLS-B10 @unit-level
  Scenario: A pattern scans the files through the reader given
    Given a test file whose tree and whose reader hold different tests
    When its tests are listed with the reader, and without one
    Then the reader's tests are listed with it, and the tree's without
