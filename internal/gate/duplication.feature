# language: en
# @anchors
#   ref: DUPLC
#   updated_at: 2026-09-30
#   layer: feature

@DUPLC
Feature: Duplication — no code file holds a block copied from somewhere else

  @DUPLC-B01 @unit-level
  Scenario: A file in no clone passes
    Given a jscpd report whose clones do not involve the file
    When the check judges the file
    Then it returns Pass

  @DUPLC-B02 @unit-level
  Scenario: A file in a clone fails naming the other side, from either side
    Given a clone between two files
    When the check judges the first file and then the second
    Then each fails naming its own lines, the other file and that file's lines

  @DUPLC-B03 @unit-level
  Scenario: A clone inside one file names only the line ranges
    Given a clone whose two sides are in the same file
    When the check judges that file
    Then it fails naming the two line ranges without repeating the file

  @DUPLC-B04 @unit-level
  Scenario: Clones within the declared threshold are reported, over it they fail
    Given a .jscpd.json that declares a threshold
    When the project's duplicated percentage is at the threshold, and then over it
    Then the file's clones are Diverge, and then they fail

  @DUPLC-B05 @unit-level
  Scenario: Without a declared threshold any clone fails whatever the exit code
    Given a jscpd run that exits 0 and reports a clone, with no threshold declared
    When the check judges a file of the clone
    Then it returns Fail

  @DUPLC-B06 @unit-level
  Scenario: jscpd runs once per scan
    Given two files judged against the same root and map, and a third against a new map
    When the check judges them
    Then jscpd ran twice

  @DUPLC-B07 @unit-level
  Scenario: An absolute path in the report is read relative to the root
    Given a report that names the file by its absolute path
    When the check judges the file by its node ID
    Then it finds the clone

  @DUPLC-E01 @unit-level
  Scenario: No report leaves the check Pending naming why
    Given a jscpd run that writes no report, and one that writes a report that is not JSON
    When the check judges a file
    Then it returns Pending naming the reason, never Pass or Fail

  @DUPLC-B08 @unit-level
  Scenario: The gate runs a pinned jscpd release
    Given the command the gate builds to run jscpd
    When its arguments are read
    Then the package carries an exact version
