# language: en
# @anchors
#   ref: PLPRP
#   updated_at: 2026-09-28
#   layer: feature

@PLPRP
Feature: PlanProgress — create a plan's progress file, the state that lives beside the decision and outside the map

  @PLPRP-B01 @unit-level
  Scenario: The progress file sits beside the plan
    Given the plan "plans/0017-mutation.md"
    When its progress path is computed
    Then it is "plans/0017-mutation-progress.md"

  @PLPRP-B02 @unit-level
  Scenario: One section per phase declared in the plan's headers
    Given a plan with the headers "### MTUAO-W01 — the tool and the report" and "### MTUAO-W02 — CI ingests the signal"
    And the headers "## ABCDE-W03 — two" and "#### ABCDE-W04 — four", and the out-of-range "# ABCDE-W05" and "##### ABCDE-W06"
    When its progress file is created
    Then it has the sections "## MTUAO-W01 — the tool and the report", "## MTUAO-W02", "## ABCDE-W03" and "## ABCDE-W04", each with an unchecked item
    And it has no section for ABCDE-W05 nor ABCDE-W06

  @PLPRP-B03 @unit-level
  Scenario: The phase code length follows the project's configuration
    Given a project whose code lengths are three
    And a plan with the header "### ABC-F01 — short code phase"
    When its progress file is created
    Then it has the section "## ABC-F01"

  @PLPRP-B04 @unit-level
  Scenario: A plan with no phases gets a note saying what to add
    Given a plan with no phase header
    When its progress file is created
    Then the file says the plan does not declare phases yet

  @PLPRP-B05 @unit-level
  Scenario: An existing progress file is never overwritten
    Given a plan whose progress file already holds a done item
    When its progress file is created again
    Then the creation fails
    And the progress file keeps its content

  @PLPRP-B06 @unit-level
  Scenario: new progress creates the file for an existing plan
    Given the plan "plans/0002-platform.md" with the header code PLTFR and phases PLTFR-W01 and PLTFR-W02
    When `anchors new progress --for plans/0002-platform.md` runs
    Then "plans/0002-platform-progress.md" exists, titled "# Progress — PLTFR", with the sections of both phases
    And running it again fails

  @PLPRP-B07 @unit-level
  Scenario: new progress without the plan is refused
    Given a project root
    When `anchors new progress` runs without `--for`
    Then the command fails naming `--for`

  @PLPRP-I01 @unit-level
  Scenario: The suffix is the one the scanner keeps out of the map
    Given a plan path with this unit's progress suffix
    When the scanner is asked whether it is a progress file
    Then it says yes

  @PLPRP-X01 @unit-level
  Scenario: The progress takes its code from the plan's header
    Given the plan "plans/0002-platform.md" whose header declares the code PLTFR
    When `anchors new progress --for plans/0002-platform.md` runs
    Then the progress is titled "# Progress — PLTFR"

  @PLPRP-E01 @unit-level
  Scenario: A plan without a code is refused
    Given a plan with no @anchors header
    When `anchors new progress --for p.md` runs
    Then the command fails
    And no progress file is created

  @PLPRP-E02 @unit-level
  Scenario: A plan that cannot be read is refused
    Given a project root with no plan at "plans/none.md"
    When `anchors new progress --for plans/none.md` runs
    Then the command fails with "read the plan"
