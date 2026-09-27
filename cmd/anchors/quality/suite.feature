# language: en
# @anchors
#   ref: STPRS
#   updated_at: 2026-09-26
#   layer: feature

@STPRS
Feature: SuiteProxy — runs the test and mutation suites the project declared, and binds their reports to the map

  @STPRS-B01 @unit-level
  Scenario: The selected suite runs at the root under a header naming it
    Given a project declaring a backend unit suite and an e2e suite
    When the test command runs for the unit layer
    Then the output shows "━━━ test [backend/unit] ━━━"
    And the e2e suite does not run

  @STPRS-B02 @unit-level
  Scenario: The report this run wrote is ingested into the map
    Given a unit suite whose command writes an lcov with one of two lines of a.go covered
    When the test command runs for the unit layer
    Then the map's a.go node carries 2 total lines and 1 covered line

  @STPRS-B03 @unit-level
  Scenario: A passing suite with no report says nothing was ingested
    Given an e2e suite that declares no report
    When the test command runs for the e2e layer
    Then it prints "[e2e] passed, but the suite declares no report — nothing was ingested."

  @STPRS-B04 @unit-level
  Scenario: A failing suite is still ingested and stops the rest
    Given a unit suite that writes an lcov and exits with failure, followed by an e2e suite
    When the test command runs for unit and e2e
    Then it fails naming the unit layer
    And the unit report reached the map
    And the e2e suite did not run

  @STPRS-B05 @unit-level
  Scenario: Without changed files the full command runs
    Given a suite declaring the full command "jest" and an incremental one
    When the command is picked with no changed file
    Then it is "jest"

  @STPRS-B06 @unit-level
  Scenario: The incremental command receives the impact path where it declares it
    Given a suite whose incremental command places the files before "--ci"
    When the command is picked for a/x.ts and b/y.ts
    Then it is "jest --findRelatedTests a/x.ts b/y.ts --ci"
    And a file of the impact path is passed as an absolute path with no backslash

  @STPRS-B07 @unit-level
  Scenario: Tests get code and tests, mutation gets only code
    Given a map where a.spec.md specifies a.go, tested by a_test.go
    When the test and the mutation commands run incrementally for a.go
    Then the test command received a.go and a_test.go as absolute paths
    And the mutation command received a.go and not a_test.go

  @STPRS-B08 @unit-level
  Scenario: The target fills the placeholder and is ignored without one
    Given a command with a target placeholder and a command without one
    When each is built with the target business-logic/dedup.ts
    Then the first reads "stryker run --mutate business-logic/dedup.ts"
    And the second is unchanged

  @STPRS-B09 @unit-level
  Scenario: A passing run chains coverage, and an empty chain does nothing
    Given a unit suite that passes
    When the test command runs chaining coverage
    Then the output shows "━━━ then: anchors coverage ━━━" and the line coverage section
    And an empty or blank chain runs nothing and returns no error

  @STPRS-B10 @unit-level
  Scenario: An impact path with no code file runs nothing
    Given a map holding only a spec, and a mutation suite with an incremental command
    When the mutation command runs incrementally for that spec
    Then it prints "the impact path reaches no code file — nothing to run."
    And the suite does not run

  @STPRS-B11 @unit-level
  Scenario: A passing run chains the check over the suite's own scope
    Given a unit suite that passes and a gate that passes
    When the test command runs chaining check, once in full and once for the changed a.go
    Then the full run's check reads "check --all" and the incremental run's check does not

  @STPRS-I01 @unit-level
  Scenario: A report older than the run is never ingested
    Given a suite whose declared report was last written before the run
    When the test command runs for that suite
    Then it prints "report OLDER than this run"
    And the map keeps the signal it had

  @STPRS-X01 @unit-level
  Scenario: The declared command runs as declared
    Given a suite declaring "yarn test:unit" with no placeholder
    When the command is built with a target
    Then it is exactly "yarn test:unit"

  @STPRS-E01 @unit-level
  Scenario: The suite commands without configuration fail
    Given a directory with no anchors.yaml
    When the test command runs
    Then it fails naming "load anchors.yaml"

  @STPRS-E02 @unit-level
  Scenario: A section with no suite shows how to declare it and fails
    Given a project with no mutation section and no tests section
    When the mutation and the test commands run
    Then the mutation command fails with "no suite declared in `mutation:`" after showing the mutation report to declare
    And the test command shows the junit and lcov reports to declare

  @STPRS-E03 @unit-level
  Scenario: A filter that names nothing declared is refused with what is declared
    Given a project declaring only the backend unit suite
    When the test command runs for the e2e layer, and for the mobile workspace
    Then the first fails naming e2e and listing "declared layers:     unit"
    And the second fails

  @STPRS-E04 @unit-level
  Scenario: Declared filters that match no suite together are refused
    Given a project declaring a backend unit suite and a mobile e2e suite
    When the test command runs for the unit layer in the mobile workspace
    Then it fails with "no suite matches layer(s) unit, workspace(s) mobile"

  @STPRS-E05 @unit-level
  Scenario: A target placeholder without a target is refused
    Given the command "npx stryker run --mutate {{target}}"
    When it is built with no target
    Then it fails asking for --target

  @STPRS-E06 @unit-level
  Scenario: The incremental mode without an incremental command is refused
    Given a suite declaring only "yarn test:unit"
    When the command is picked for a changed file
    Then it fails showing how to declare run_changed with {{files}}

  @STPRS-E07 @unit-level
  Scenario: A command line over the ceiling is refused before running
    Given the argv ceiling set to 10 characters and a longer declared command
    When the test command runs
    Then it fails with "(ceiling 10 on this platform)"

  @STPRS-E08 @unit-level
  Scenario: A chain naming another command is refused
    Given the chain "deploy"
    When the chain runs
    Then it fails saying check is what is accepted

  @STPRS-E09 @unit-level
  Scenario: The incremental mode without a map points at the map build
    Given a project with a suite and no map
    When the test command runs incrementally for a.go
    Then it fails with an error naming "anchors map build"
