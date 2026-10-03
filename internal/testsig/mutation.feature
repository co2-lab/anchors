# language: en
# @anchors
#   ref: MTINM
#   updated_at: 2026-10-03
#   layer: feature

@MTINM
Feature: MutationIngest — the mutation score per file, read from a Mutation Testing Elements report

  @MTINM-B01 @unit-level
  Scenario: The canonical format is read under each of its names
    Given a canonical report
    When it is read with the format "", "mutation-testing-elements", " MTE " and "Stryker"
    Then each reading finds the report's file

  @MTINM-B02 @unit-level
  Scenario: Killed and timed-out mutants count as killed
    Given a file with two killed mutants and one timed out
    When the report is read
    Then the file has 3 killed
    And the timed-out mutants are also counted apart

  @MTINM-B03 @unit-level
  Scenario: A survivor counts as survived with its line
    Given a file with one surviving mutant on line 42
    When the report is read
    Then the file has 1 survived at line 42

  @MTINM-B04 @unit-level
  Scenario: Uncovered mutants are counted apart and stay out of the score
    Given a file with two uncovered mutants, one survivor and one killed
    When the report is read
    Then it counts 2 uncovered and scores 50

  @MTINM-B05 @unit-level
  Scenario: Ignored mutants are counted apart and stay out of the score
    Given a file with one killed, one survived and two ignored mutants
    When the report is read
    Then it counts 2 ignored and scores 50

  @MTINM-B06 @unit-level
  Scenario: A mutant that failed to compile stays out of the count and the score
    Given a file with three killed, one survived, one uncovered and one compile error
    When the report is read
    Then it scores 75

  @MTINM-B07 @unit-level
  Scenario: A file where no mutant ran has no score
    Given a file whose three mutants no test covered
    When the report is read
    Then it has no score and no survivor

  @MTINM-B08 @unit-level
  Scenario: The thresholds are read from the report
    Given a report declaring the thresholds low 60 and high 80
    When the report is read
    Then the low threshold is 60 and the high is 80

  @MTINM-B09 @unit-level
  Scenario: File paths are normalized to the map's form
    Given files named "./lib/a.ts" and "/home/ci/app/src/b.ts"
    When the report is read
    Then they are keyed "lib/a.ts" and "src/b.ts"

  @MTINM-I01 @unit-level
  Scenario: A file where nothing ran has no score and keeps the ignored count
    Given a file whose two mutants were both ignored
    When the report is read
    Then it has no score, with 2 ignored and nothing killed

  @MTINM-X01 @unit-level
  Scenario: Thresholds absent from the report stay zero
    Given a report without thresholds
    When the report is read
    Then both thresholds are zero

  @MTINM-E01 @unit-level
  Scenario: An unknown format is refused naming the accepted ones
    Given the format name "stryker4s-xml"
    When a report is read under it
    Then the error lists gremlins among the accepted formats

  @MTINM-E02 @unit-level
  Scenario: A report that is not the canonical format is refused
    Given a file that is not JSON, and a gremlins report read as the canonical format
    When each is read
    Then each is refused with an error

  @MTINM-E03 @unit-level
  Scenario: A report without files is refused
    Given a JSON report with a schema version and no files
    When the report is read
    Then it is refused with an error

  @MTINM-E04 @unit-level
  Scenario: A missing report returns the read error
    Given a path where no report exists
    When the report is read
    Then the error says the file does not exist
