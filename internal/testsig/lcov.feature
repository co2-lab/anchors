# language: en
# @anchors
#   code: LCFTA
#   ref: LCINL
#   updated_at: 2026-10-03
#   layer: feature

@LCINL
Feature: LcovIngest — line coverage per file, and the uncovered lines of a change, read from an lcov report

  @LCINL-B01 @unit-level
  Scenario: Each record becomes one file's coverage in report order
    Given a report with a record for src/a.ts and one for src/b.ts
    When the report is read
    Then it has two files, src/a.ts first

  @LCINL-B02 @unit-level
  Scenario: A line with hits is covered and a line without is not
    Given a record with lines 1 and 3 hit and line 2 with zero hits
    When the report is read
    Then lines 1 and 3 are covered, line 2 is uncovered, and the file is 2 of 3

  @LCINL-B03 @unit-level
  Scenario: Stated totals win over the line entries
    Given a record with two line entries and stated totals of 10 lines, 8 hit
    When the report is read
    Then the file is 8 of 10

  @LCINL-B04 @unit-level
  Scenario: The uncovered lines of a change are instrumented and not covered
    Given a file with line 10 covered, line 11 uncovered and line 12 covered
    When the change touches lines 10, 11 and 13
    Then only line 11 is uncovered

  @LCINL-B05 @unit-level
  Scenario: The instrumented count ignores changed lines the report did not instrument
    Given a file with lines 10, 11 and 12 instrumented
    When the change touches lines 10, 11 and 13
    Then 2 changed lines are instrumented

  @LCINL-B06 @unit-level
  Scenario: The percentage is covered over total, and zero for a file without lines
    Given a file of 8 covered out of 10 and a file of no line
    When their percentage is asked
    Then they give 80 and 0

  @LCINL-B07 @unit-level
  Scenario: A record without its end line is closed by the next record or the end of the report
    Given a record for a.ts and one for b.ts, neither closed
    When the report is read
    Then both files are read with their own lines

  @LCINL-B08 @unit-level
  Scenario: A line entry with a checksum field keeps its hit count
    Given a record with "DA:1,5,PF4Rz2r7RTliO9u6bZ7h6g" and "DA:2,0,XyZ"
    When the report is read
    Then line 1 is covered, line 2 is not, and the file is 1 of 2

  @LCINL-B09 @unit-level
  Scenario: Entries before the first source-file line belong to no file
    Given two covered line entries and totals of 9 of 9 before the first "SF:a.ts", whose record has one uncovered line
    When the report is read
    Then the report has one file, "a.ts", at 0 of 1

  @LCINL-I01 @unit-level
  Scenario: Uncovered changed lines never exceed the instrumented changed lines
    Given a file's coverage and a change with covered, uncovered and non-instrumented lines
    When both answers are asked
    Then the uncovered lines are no more than the instrumented count

  @LCINL-X01 @unit-level
  Scenario: Branch and function entries do not change the line counts
    Given a record with two line entries, one hit, and branch and function entries
    When the report is read
    Then the file is 1 of 2

  @LCINL-E01 @unit-level
  Scenario: A missing report returns the error
    Given a path where no report exists
    When the report is read
    Then the error says the file does not exist

  @LCINL-B10 @unit-level
  Scenario: Branch entries record each branch, taken or not
    Given a record with a branch taken, one with count zero, one never run, a malformed one and one before any file
    When the report is parsed
    Then the first three are recorded by line, block and branch with whether they were taken, and the others record nothing
