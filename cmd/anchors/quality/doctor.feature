# language: en
# @anchors
#   ref: HLDCH
#   updated_at: 2026-09-26
#   layer: feature

@HLDCH
Feature: DoctorCommand — the global health x-ray, and the repair of the github-mode environment

  @HLDCH-B01 @unit-level
  Scenario: The diagnosis is printed grouped by the check that found it
    Given a map with a guide nobody governs and a spec with no identity
    When the doctor runs
    Then it prints "anchors doctor — 2 nodes, 0 edges, 0 layers"
    And each finding is listed under its check with the group count
    And it closes with "(diagnosis — nothing was blocked; you decide what to reconcile)"

  @HLDCH-B02 @unit-level
  Scenario: A group with any warning is marked as a warning
    Given a report whose check "mixed" has an informational finding followed by a warning
    When the report is printed
    Then the group reads "⚠ mixed (2)" and only the warning line is marked

  @HLDCH-B03 @unit-level
  Scenario: A report with no finding says the ecosystem is sound
    Given a report with no finding
    When the report is printed
    Then it prints "✓ no systemic loose end found"

  @HLDCH-B04 @unit-level
  Scenario: The fix outside the github mode does nothing
    Given a project in local mode with a map
    When the doctor runs with the fix flag
    Then it prints "--fix: nothing to do"
    And no pipeline folder is created

  @HLDCH-B05 @unit-level
  Scenario: The fix in github mode seeds the workflow pipelines
    Given a project in github mode with an authenticated gh and no pipeline
    When the doctor runs with the fix flag
    Then it prints how many pipelines it created
    And every workflow pipeline exists on disk

  @HLDCH-B06 @unit-level
  Scenario: The fix protects the declared branches and skips one that does not exist
    Given a project in github mode declaring main and staging as protected, with staging not existing
    When the doctor runs with the fix flag
    Then gh receives the protection body for main on its standard input
    And the output reads "✓ main protected" and "staging does not exist yet"

  @HLDCH-B07 @unit-level
  Scenario: The fix disables an approval requirement the author cannot satisfy
    Given a project in github mode whose account has only write permission
    When the doctor runs with the fix flag
    Then it prints "✓ approval requirement DISABLED"

  @HLDCH-B08 @unit-level
  Scenario: The fix ensures the state labels and says the board is optional
    Given a project in github mode whose label is "anchors"
    When the doctor runs with the fix flag
    Then gh is asked to create the "anchors" label and the needs-user label
    And the output says the board is OPTIONAL

  @HLDCH-B09 @unit-level
  Scenario: The protection body carries every required field and the approvals
    Given one required approval
    When the protection body is built
    Then it holds the status checks, admin enforcement, pull-request reviews and restrictions fields
    And the pull-request review object requires 1 approval

  @HLDCH-B10 @unit-level
  Scenario: The pipelines check in local mode has nothing to check
    Given a project in local mode
    When the doctor runs with the pipelines check
    Then it prints "· local mode — no workflow pipeline to check."

  @HLDCH-B11 @unit-level
  Scenario: Missing pipelines are named and CI continues by default
    Given a project in github mode with no pipeline and no blocking declared
    When the doctor runs with the pipelines check
    Then it names the MISSING pipelines, warns that CI continues and does not fail
    And after the pipelines are seeded it says they are in place and up to date

  @HLDCH-B12 @unit-level
  Scenario: A project that declared stale pipelines as blocking fails the check
    Given a project in github mode with no pipeline and stale_pipeline_blocks declared true
    When the doctor runs with the pipelines check
    Then it fails naming "stale_pipeline_blocks: true"

  @HLDCH-I01 @unit-level
  Scenario: The diagnosis never fails the doctor
    Given a map whose diagnosis has warning findings
    When the doctor runs
    Then it returns no error

  @HLDCH-I02 @unit-level
  Scenario: Running the fix twice changes nothing the second time
    Given a project in github mode the fix has just set up
    When the doctor runs with the fix flag again
    Then it prints "--fix: the pipelines already exist and are up to date."

  @HLDCH-X01 @unit-level
  Scenario: The fix never creates a board
    Given a project in github mode with an authenticated gh
    When the doctor runs with the fix flag
    Then gh is never asked to create a project

  @HLDCH-X02 @unit-level
  Scenario: The orphan local queue and delivery records are warned about and kept
    Given a project in github mode with a pending task in the local queue and a delivery record
    When the doctor runs with the fix flag
    Then it warns about both
    And the task file and the delivery record still exist

  @HLDCH-E01 @unit-level
  Scenario: The doctor without configuration fails
    Given a directory with no anchors.yaml
    When the doctor runs
    Then it fails naming "load config"

  @HLDCH-E02 @unit-level
  Scenario: The doctor without a map fails
    Given a project with configuration and no map
    When the doctor runs
    Then it fails naming "load map"

  @HLDCH-E03 @unit-level
  Scenario: The fix without gh refuses before seeding anything
    Given a project in github mode and no gh on the PATH
    When the environment repair runs
    Then it fails naming "gh is not installed"
    And no pipeline folder is created

  @HLDCH-E04 @unit-level
  Scenario: The fix with gh not logged in refuses before seeding anything
    Given a project in github mode and a gh whose auth status fails
    When the environment repair runs
    Then it fails saying gh is not authenticated
    And no pipeline folder is created
