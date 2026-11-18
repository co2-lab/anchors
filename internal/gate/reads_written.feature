# language: en
# @anchors
#   code: RDWRF
#   ref: RDWRS
#   updated_at: 2026-11-10
#   layer: feature

@RDWRS
Feature: ReadsWritten — every data a header reads is written by some file's header

  @RDWRS-B01 @unit-level
  Scenario: A data read that ties the file to no writer fails, naming it
    Given a reader whose header reads two data, one tied to its writer and one nobody writes
    When the gate confronts it, and again with only the tied one
    Then the first fails naming the unwritten data, and the second passes

  @RDWRS-B02 @unit-level
  Scenario: A file whose header reads nothing is skipped
    Given a code file whose header declares no reads
    When the gate confronts it
    Then it is skipped

  @RDWRS-E01 @unit-level
  Scenario: Without a map the gate is pending
    Given a reader and no built map
    When the gate confronts it
    Then it is pending
