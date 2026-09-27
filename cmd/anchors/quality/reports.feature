# language: en
# @anchors
#   ref: RPRTS
#   updated_at: 2026-09-26
#   layer: feature

@RPRTS
Feature: Reports — markdown perspectives on what Anchors already measures, written into docs

  @RPRTS-B01 @unit-level
  Scenario: A single perspective is written to docs or to the chosen file
    Given a configured project with a map
    When the tests report is generated to custom/tests.md, and the structure report with no destination
    Then it prints "report generated: custom/tests.md" and that file holds the test report
    And docs/anchors-structure-report.md holds the structure report

  @RPRTS-B02 @unit-level
  Scenario: The all command writes every perspective and an index into docs/anchors
    Given a configured project with a map
    When every report is generated
    Then it prints "6 report(s) + index generated in docs/anchors/"
    And docs/anchors/index.md links each of the six perspective files, which all exist

  @RPRTS-B03 @unit-level
  Scenario: Every perspective of a configured project opens with the same header and closes with a footer
    Given a configured project with a map, rendered at 2026-01-02 03:04
    When each perspective is rendered
    Then each opens with its title and the shared generated-on and do-not-edit lines
    And each closes with an Anchors footer line

  @RPRTS-B04 @unit-level
  Scenario: The tests perspective merges execution by layer and warns about failures and stale signals
    Given a test file with 5 passed and 1 failed in unit and 2 passed and 3 skipped in e2e, whose signal is stale
    When the tests perspective is rendered
    Then it shows the e2e and unit rows and the total "7 | 1 | 3"
    And it warns "1 failing test(s)" and "1 test file(s) with a STALE signal"

  @RPRTS-B05 @unit-level
  Scenario: The tests perspective counts only measured specs and lists the unproven scenarios
    Given a measured spec proving BBBBB-B01 of BBBBB-B01 and BBBBB-B02, and an unmeasured spec
    When the tests perspective is rendered
    Then it reads "1/2 requirements proven" across 1 measured spec
    And it lists b.spec.md with BBBBB-B02 unproven and reports 1 spec still NOT measured

  @RPRTS-B06 @unit-level
  Scenario: The tests perspective lists the files below 70% of lines and the coverage regressions
    Given x.go at 40% of lines after 60% before, and y.go at 90%
    When the tests perspective is rendered
    Then x.go is listed below 70% and as "60% → 40% (-20)"
    And y.go is not listed

  @RPRTS-B07 @unit-level
  Scenario: The tests perspective with nothing ingested says so in each section
    Given a map with a spec and a code file and no signal
    When the tests perspective is rendered
    Then it says no execution result was ingested, no spec was measured yet, and no coverage dropped
    And it claims no failing test and no stale signal

  @RPRTS-B08 @unit-level
  Scenario: The quality perspective gives the verdict per gate and the divergences
    Given a blocking gate that fails on both code files and a judgment gate
    When every report is generated
    Then the quality report shows "| always-fails | **blocking** | 0 | 2 | 0 |" and "✗ **Barred**"
    And it lists "[BLOCKS] `always-fails` @ `x.go`" and "Awaiting AI judgment (2)"

  @RPRTS-B09 @unit-level
  Scenario: The structure perspective counts nodes by kind, the governance and the identity findings
    Given a map of 8 nodes and 2 edges where GUIDE.md governs two files and a spec has no code
    When every report is generated
    Then the structure report counts "| code | 2 |" and "| spec | 3 |", "8 nodes, 2 edges.", "| `GUIDE.md` | 2 |"
    And it lists nocode.spec.md under "Missing identity (1)"

  @RPRTS-B10 @unit-level
  Scenario: The configuration perspective lists what anchors.yaml declares and what it misses
    Given a configuration with two layers, a governance rule and two gates, and a guide nobody governs
    When every report is generated
    Then the configuration report lists the layers, the rule, each gate with its kind and force, and "not configured" co-location
    And it lists LONE.md as a guide with no governance and spec as a kind with no gate

  @RPRTS-B11 @unit-level
  Scenario: The issues perspective splits the open issues by who must act, and lists the tasks
    Given a user-owned issue in todo, an issue in doing, a deferred issue and one live task
    When every report is generated
    Then the issues report counts one per state and lists decide-pricing under "Waiting on YOU", fix-x under "Open" and later under "Deferred"
    And it lists "1 live task(s)" with x.go

  @RPRTS-B12 @unit-level
  Scenario: The inconsistencies perspective lists every health finding by check and the failing gates
    Given a guide nobody governs and a blocking gate that fails on both code files
    When every report is generated
    Then the inconsistencies report has "### guide-sem-governo (1)" with LONE.md
    And it has "## From the quality gates (2 failures)"

  @RPRTS-B13 @unit-level
  Scenario: A finding section takes the findings of its check, caps the list at 25, and is absent when empty
    Given 27 findings of a check starting with the section's name and one of another check
    When the finding section is built
    Then it is titled with 27, lists 25 and "… and 2 more", and leaves out the other check
    And a section for a check with no finding is empty

  @RPRTS-B14 @unit-level
  Scenario: Without anchors.yaml the perspectives say what is missing instead of inventing
    Given an empty map and no configuration
    When the quality, configuration and issues perspectives are rendered
    Then they read "_No gate declared in anchors.yaml._", "_anchors.yaml not found._" and "_Empty queue._"
    And the issues perspective prints no empty list

  @RPRTS-I01 @unit-level
  Scenario: An issue waiting on the user is listed only as the user's
    Given a user-owned issue decide-pricing in todo
    When every report is generated
    Then the issues report names decide-pricing exactly once

  @RPRTS-X01 @unit-level
  Scenario: Generating the reports leaves the map as it was
    Given a configured project with a map
    When every report is generated
    Then the map file has the same content as before

  @RPRTS-E01 @unit-level
  Scenario: The reports without a map point at the map build
    Given a configured project with no map
    When the tests report and every report are generated
    Then both fail naming the map load and "anchors map build"

  @RPRTS-E02 @unit-level
  Scenario: A destination that cannot be written fails the report
    Given a destination path under an existing file
    When the tests report is generated there
    Then it fails
