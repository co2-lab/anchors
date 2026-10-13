# language: en
# @anchors
#   code: PCFPR
#   ref: LDCML
#   updated_at: 2026-10-12
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

  @LDCML-E01 @unit-level
  Scenario: A refused change writes nothing
    Given a progress file
    When a task that does not exist is done, and a demand is added with neither title nor plan
    Then each command fails with the cause, and the file is unchanged
