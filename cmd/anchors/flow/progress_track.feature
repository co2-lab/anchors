# language: en
# @anchors
#   code: PCFPR
#   ref: LDCML
#   updated_at: 2026-10-24
#   layer: feature

@LDCML
Feature: LedgerCommand — anchors progress

  @LDCML-B01 @unit-level
  Scenario: Demands, stages and tasks are added
    Given a project with a plan declaring two phases
    When a demand is added by title, the plan by --plan, a stage and a weighted task
    Then each says its id, and the plan's demand has its two stages

  @LDCML-B02 @unit-level
  Scenario: Tasks move by command
    Given a task
    When an agent starts it, another starts it without and with --take, it is done, reopened, cancelled with a reason and weighed
    Then the second agent is refused until it takes it over, and the file records who and when

  @LDCML-B03 @unit-level
  Scenario: The summary shows demands and stages
    Given a project naming its unit hours, with a demand of two stages
    When the summary is printed, and as JSON, and for an empty project
    Then it shows the whole progress in hours and one line per demand and stage, never a task, the JSON has the rows, and the empty project is told how to start

  @LDCML-B04 @unit-level
  Scenario: Show prints the tasks
    Given two demands with tasks
    When show runs for the file and for one demand
    Then every task line is printed, without the header and the sum, and only that demand's for one

  @LDCML-B05 @unit-level
  Scenario: The unit is the project's
    Given a project naming no unit
    When a task is added
    Then the file counts in points

  @LDCML-B06 @unit-level
  Scenario: A stage is renamed by command
    Given a demand with a stage holding a task
    When the stage is renamed
    Then the summary shows the new title under the same id, with its task

  @LDCML-B07 @unit-level
  Scenario: Unplanned work is marked and shown
    Given a demand with a stage holding a planned task
    When a task is added unplanned, the planned one is marked unplanned and then unmarked
    Then the summary says the unplanned weight of the demand and the stage, and show marks the task

  @LDCML-B08 @unit-level
  Scenario: Found work moves into the stage it arose in
    Given a demand with a stage of its work and a separate stage of found work holding an unplanned task
    When the task is moved to the stage of its work and the emptied stage is dropped
    Then the summary shows the unplanned weight inside that stage, and the separate stage is gone

  @LDCML-E01 @unit-level
  Scenario: A refused change writes nothing
    Given a progress file
    When a task that does not exist is done, and a demand is added with neither title nor plan
    Then each command fails with the cause, and the file is unchanged

  @LDCML-B09 @unit-level
  Scenario: The summary leaves a frozen plan out of the total
    Given two demands, one whose plan's header is frozen
    When the summary is printed
    Then the total counts the other alone, and the frozen demand's row says it is out of the total

  @LDCML-B10 @unit-level
  Scenario: A settled demand is one line
    Given a demand all done with a stage, one none started with a stage, and one with a task in progress
    When the summary and show are printed, and show of the settled one
    Then the settled demands are their own lines alone and the one in motion shows its stage and tasks, and show of the settled one details it
