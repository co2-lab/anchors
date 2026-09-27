# language: en
# @anchors
#   ref: LCBCL
#   updated_at: 2026-09-26
#   layer: feature

@LCBCL
Feature: LocalBacklog — what is still open locally after a full check, said in two lines

  @LCBCL-B01 @unit-level
  Scenario: The issues in todo and doing are counted, with the user-owned ones apart
    Given a project with a violation and a user-owned decision in todo and one issue in doing
    When the local backlog is read
    Then it counts 2 in todo, 1 of them for the user, and 1 in doing

  @LCBCL-B02 @unit-level
  Scenario: The pending and claimed tasks are counted
    Given a task queue with one pending task and one claimed task
    When the local backlog is read
    Then it counts 1 pending and 1 claimed

  @LCBCL-B03 @unit-level
  Scenario: A project with nothing open prints no backlog
    Given a project with no issue and no task
    When the local backlog is read and printed
    Then the backlog is empty and the output is empty

  @LCBCL-B04 @unit-level
  Scenario: Only the side that has something open gets its line
    Given a backlog with pending tasks and no issue in todo or doing
    When the local backlog is printed
    Then the output holds the tasks line naming `anchors next`
    And it holds no line about the issues folders

  @LCBCL-I01 @unit-level
  Scenario: User-owned and past-window counts alone do not make a backlog
    Given a backlog whose only non-zero counts are the user-owned issues and the claims past the work window
    When it is asked whether it is empty and printed
    Then it is empty and prints nothing

  @LCBCL-X01 @unit-level
  Scenario: Reading the backlog changes no issue and no task
    Given a project with an issue in todo and a pending task
    When the local backlog is read
    Then the issue and task files are the same as before

  @LCBCL-E01 @unit-level
  Scenario: An issue folder that cannot be listed counts as zero
    Given a project whose todo state is a file instead of a folder
    When the local backlog is read
    Then todo counts zero and reading does not fail
