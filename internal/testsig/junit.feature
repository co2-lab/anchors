# language: en
# @anchors
#   ref: JUIJN
#   updated_at: 2026-09-26
#   layer: feature

@JUIJN
Feature: JUnitIngest — the run's outcome per test case, and the scenario codes each case proves, read from a JUnit report

  @JUIJN-B01 @unit-level
  Scenario: A report with a suites root is read case by case
    Given a report with one suite of three cases: passed, failed and skipped
    When the report is read
    Then it has three cases

  @JUIJN-B02 @unit-level
  Scenario: A report whose root is a single suite is read
    Given a report whose root is a suite with one case
    When the report is read
    Then it has one case

  @JUIJN-B03 @unit-level
  Scenario: Nested suites are flattened
    Given a suite with one case containing a nested suite with one case
    When the report is read
    Then both cases are in the report

  @JUIJN-B04 @unit-level
  Scenario: A case without a file takes its suite's file
    Given a suite with file "a.test.ts", one case without a file and one naming "b.test.ts"
    When the report is read
    Then the first case's file is "a.test.ts" and the second's is "b.test.ts"

  @JUIJN-B05 @unit-level
  Scenario: A failure or an error marks the case failed and a skip marks it skipped
    Given cases with a failure, with an error, and marked skipped
    When the report is read
    Then the first two are failed and the third is skipped

  @JUIJN-B06 @unit-level
  Scenario: Only codes of passing cases are proven
    Given SPCRX-V01 passed, SPCRX-V02 failed and SPCRX-X01 was skipped
    When the proven codes are asked
    Then SPCRX-V01 is the only code proven, and the failed and skipped codes are not

  @JUIJN-B07 @unit-level
  Scenario: Every case's codes are seen whatever the outcome
    Given MBDT-B01 passed, MBDT-B02 failed and MBDT-B03 was skipped
    When the seen codes are asked
    Then all three are seen, and MBDT-B02 is not proven

  @JUIJN-B08 @unit-level
  Scenario: Every code in a case name is extracted in the declared vocabulary
    Given the case name "SPCRX-V01: vertical axis extends SPCRX-X01"
    When its codes are read
    Then both codes are extracted, and a letter no longer declared is not

  @JUIJN-B09 @unit-level
  Scenario: A file that is not a JUnit report yields no case
    Given a report file containing text that is not XML
    When the report is read
    Then the report has no case and no error

  @JUIJN-I01 @unit-level
  Scenario: Every proven code is a seen code
    Given a report of passed, failed and skipped cases
    When the proven and the seen codes are compared
    Then every proven code is among the seen codes

  @JUIJN-X01 @unit-level
  Scenario: A code in the suite or class name proves nothing
    Given a passing case with no code in its name, inside a suite and class named ABCDX-B01
    When the proven codes are asked
    Then no code is proven

  @JUIJN-E01 @unit-level
  Scenario: An unreadable report returns the read error
    Given a path where no report exists
    When the report is read
    Then the error says the file does not exist
