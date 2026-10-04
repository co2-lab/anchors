# language: en
# @anchors
#   code: STFTG
#   ref: MSCMG
#   updated_at: 2026-10-03
#   layer: feature

@MSCMG
Feature: MigrationStepChain — one step per format version, chained in order, and a hole in the chain is an error

  @MSCMG-B01 @unit-level
  Scenario: Steps registered out of order are kept sorted
    Given an empty registry
    When the steps producing formats 4, 2 and 3 are registered in that order
    Then the registry lists them as 2, 3 and 4

  @MSCMG-B02 @unit-level
  Scenario: The steps between two formats come in ascending order
    Given the steps producing formats 3 and 2 registered in that order
    When the steps from format 1 to format 3 are asked for
    Then the answer is the step producing 2 followed by the step producing 3

  @MSCMG-B03 @unit-level
  Scenario: A file already at the target needs no step
    Given the registered steps
    When the steps from format 2 to format 2 are asked for
    Then the answer is an empty list and no error

  @MSCMG-B05 @unit-level
  Scenario: A step may rename code letters by kind
    Given a step declaring a letter renamed for plans, and a file on the format before it
    When the file is migrated through that step
    Then the step carries the rename and the file only gets the new version

  @MSCMG-B04 @unit-level
  Scenario: A key some step renames is reported as renamed
    Given the registered steps
    When the keys "julgamentos", "rule_marking", "no_such_key" and "judgments" are asked about
    Then "julgamentos" and "rule_marking" are renamed keys
    And "no_such_key" and "judgments" are not

  @MSCMG-I01 @unit-level
  Scenario: The listed steps are a copy of the registry
    Given the registered steps
    When a caller changes the reason of the first step in the list it received
    Then listing the steps again shows the original reason

  @MSCMG-X01 @unit-level
  Scenario: Only the steps inside the interval are returned
    Given the steps producing formats 2, 3 and 4
    When the steps from format 2 to format 3 are asked for
    Then the answer is only the step producing 3

  @MSCMG-E01 @unit-level
  Scenario: A hole in the chain is an error naming the missing format
    Given the steps producing formats 2 and 4, with no step producing 3
    When the steps from format 1 to format 4 are asked for
    Then the answer is an error that names format 3

  @MSCMG-E02 @unit-level
  Scenario: A target beyond the last step is an error naming the first missing format
    Given the steps producing formats 2 and 3
    When the steps from format 1 to format 5 are asked for
    Then the answer is an error that names format 4
