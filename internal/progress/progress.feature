# language: en
# @anchors
#   code: PFFPR
#   ref: DMLDD
#   updated_at: 2026-10-12
#   layer: feature

@DMLDD
Feature: DemandLedger — the project's progress file

  @DMLDD-B01 @unit-level
  Scenario: The file is rendered and read back
    Given a demand tied to a plan, with a stage-less task and a stage of tasks to do, in progress, done, cancelled and estimated
    When the file is rendered and parsed again
    Then it opens with its header and the whole progress, each record has its line, it closes with the sum, and it reads back the same

  @DMLDD-B02 @unit-level
  Scenario: Progress counts the done weight over the weight not cancelled
    Given tasks of weight 3 done, 2 done, 8 in progress and 5 cancelled
    When the progress is computed
    Then it is 38% of 13, and with nothing counted it is a dash

  @DMLDD-B03 @unit-level
  Scenario: Records are numbered in order
    Given a demand with two stages and tasks added to each and to the demand
    When the ids are read
    Then the stages are S1 and S2, the tasks T1 to T3 across the demand, and a weight of zero is refused

  @DMLDD-B04 @unit-level
  Scenario: A task moves through its states
    Given a task
    When one agent starts it, another tries, takes it over, it is done, reopened, cancelled and weighed
    Then the second agent is refused until it takes it over, and each state carries who, when or why

  @DMLDD-B05 @unit-level
  Scenario: Changes made at once are all kept
    Given twenty goroutines each adding a task through the lock
    When they all finish
    Then the file holds the twenty tasks

  @DMLDD-B06 @unit-level
  Scenario: The summary shows demands and stages, not tasks
    Given a demand with two stages and a demand with none
    When the summary is made
    Then it has a row per demand and per stage, with their progress and open tasks

  @DMLDD-B07 @unit-level
  Scenario: Two branches merge by id
    Given a base with one task, our side finishing it and adding T2, their side adding a different T2 and a stage
    When the sides are merged
    Then the task is done, both T2 tasks stay with theirs renumbered T3, and their stage joins

  @DMLDD-B08 @unit-level
  Scenario: A plan becomes a demand
    Given a plan titled Checkout declaring two phases
    When it is tracked
    Then the demand is titled Checkout with two stages, and tracking it again is refused

  @DMLDD-B09 @unit-level
  Scenario: Per-plan progress files are migrated
    Given a plan's -progress.md with a phase, a checked item, an open item and a TODO placeholder
    When the project is migrated
    Then the plan is a demand titled as the plan, with a stage, two tasks weighing 1 as estimates, one done, the file is deleted, and a dry run writes nothing

  @DMLDD-E01 @unit-level
  Scenario: A hand edit is refused
    Given a progress file with a line added by hand
    When it is read or changed
    Then the read fails as edited by hand and nothing is written over it

  @DMLDD-E02 @unit-level
  Scenario: What the format cannot read back is refused
    Given a title holding the separator, an empty reason and an unknown task
    When they are added or used
    Then each change is refused naming its cause
