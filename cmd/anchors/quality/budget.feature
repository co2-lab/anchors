# language: en
# @anchors
#   ref: BDGRN
#   updated_at: 2026-09-27
#   layer: feature

@BDGRN
Feature: BudgetRun — run a suite's files fastest first until a time budget is spent

  @BDGRN-B01 @unit-level
  Scenario: The plan is the timed files fastest first, then the untimed
    Given test files timed by this suite, one only another suite timed, untimed ones, a support file and a code file
    When the plan is made
    Then it lists this suite's timed files fastest first, then the untimed by path, and leaves the rest out

  @BDGRN-B02 @unit-level
  Scenario: Batches take what fits and then the untimed one by one
    Given a plan with timed files and untimed ones
    When batches are taken with a large remaining time, a small one, and one file at a time
    Then the first takes all timed files that fit, a remaining time the fastest exceeds takes nothing, and the untimed come one per batch

  @BDGRN-B03 @unit-level
  Scenario: A batch still running at the deadline is stopped with its group
    Given a command whose child process outlives the shell
    When the deadline comes while it runs
    Then the run is reported as cut and the child is stopped too

  @BDGRN-B04 @unit-level
  Scenario: A budget runs batches fastest first and reports what ran and what was left
    Given a suite whose files have recorded times and one untimed file
    When the suite runs with a budget
    Then the files reach run_changed fastest first, each batch is ingested, and the report says how many ran and were left

  @BDGRN-B05 @unit-level
  Scenario: A mutation budget runs one file per batch and records its time
    Given a mutation suite and two code files never timed
    When the suite runs with a budget
    Then each file runs alone and the map records how long its run took

  @BDGRN-B06 @unit-level
  Scenario: A failed batch does not stop the budget
    Given a suite whose first batch fails
    When it runs with a budget
    Then the next batches still run and the command fails at the end naming the suite

  @BDGRN-B07 @unit-level
  Scenario: What the budget cannot run is refused before anything runs
    Given a suite with no run_changed, one with no report, and a budget with --changed
    When each runs with a budget
    Then each is refused saying why, and no command ran

  @BDGRN-E02 @unit-level
  Scenario: A budget without a map is refused saying to build it
    Given a project whose map was never built
    When a suite runs with a budget
    Then it is refused saying to run anchors map build, and no command ran
