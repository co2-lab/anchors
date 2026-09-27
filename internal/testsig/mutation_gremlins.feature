# language: en
# @anchors
#   ref: GRING
#   updated_at: 2026-09-27
#   layer: feature

@GRING
Feature: GremlinsIngest — the mutation score per file, read from a gremlins report

  @GRING-B01 @unit-level
  Scenario: The listed files are read into one result per file
    Given a gremlins report listing src/domain/bankaccount.go
    When it is read with the format gremlins
    Then the result has that file

  @GRING-B02 @unit-level
  Scenario: Killed and timed-out mutants count as killed
    Given a file with one KILLED and one TIMED OUT mutant
    When the report is read
    Then the file has 2 killed
    And the timed-out mutants are also counted apart

  @GRING-B03 @unit-level
  Scenario: A mutant that lived counts as survived with its line
    Given a file with one LIVED mutant on line 42
    When the report is read
    Then line 42 is among the survived lines

  @GRING-B04 @unit-level
  Scenario: Not viable, runnable, skipped and unknown statuses stay out of the score
    Given a file with one killed, one lived, and one each of NOT VIABLE, RUNNABLE, SKIPPED and an unknown status
    When the report is read
    Then it scores 50 with 1 killed and 1 survived

  @GRING-B05 @unit-level
  Scenario: The status is matched ignoring case and spaces
    Given a file with the statuses "killed", " TimedOut ", "Timed out" and "Lived"
    When the report is read
    Then it has 3 killed and 1 survived

  @GRING-B06 @unit-level
  Scenario: The thresholds stay zero
    Given a gremlins report
    When it is read
    Then both thresholds are zero

  @GRING-B07 @unit-level
  Scenario: File paths are normalized as in the canonical reading
    Given files named "./cmd/a.go" and "/home/ci/svc/src/b.go"
    When the report is read
    Then they are keyed "cmd/a.go" and "src/b.go"

  @GRING-I01 @unit-level
  Scenario: The score is killed over killed plus survived
    Given files of three killed and one lived, and of one killed and three lived
    When the report is read
    Then they score 75 and 25

  @GRING-B08 @unit-level
  Scenario: A mutant no test covered is counted apart and does not enter the score
    Given a file with one KILLED and one NOT COVERED mutant
    When the report is read
    Then it has 1 no-coverage mutant, no survivor, and scores 100

  @GRING-B09 @unit-level
  Scenario: A file where no mutant ran scores 100
    Given a file whose mutants are NOT VIABLE and SKIPPED, and a file whose only mutant is NOT COVERED
    When the report is read
    Then both score 100, and the second keeps its no-coverage count

  @GRING-X01 @unit-level
  Scenario: The report's own efficacy figure is not used
    Given a report stating an efficacy of 40.0 whose statuses give two killed of three that ran
    When it is read
    Then the score is two thirds of one hundred, not 40.0

  @GRING-E01 @unit-level
  Scenario: A canonical-format report under the gremlins format is refused naming format
    Given a canonical-format report
    When it is read with the format gremlins
    Then the error mentions format

  @GRING-E02 @unit-level
  Scenario: A report without files is refused
    Given a gremlins report whose file list is empty
    When it is read
    Then it is refused with an error
