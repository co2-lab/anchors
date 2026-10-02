# language: en
# @anchors
#   ref: TSQUT
#   updated_at: 2026-10-01
#   layer: feature

@TSQUT
Feature: TaskQueue — the file-backed queue between "something changed" and "someone works on it"

  @TSQUT-B01 @unit-level
  Scenario: An enqueued task is born pending
    Given a project holding "AddItem.spec.md" and no queue
    When the task "1-spec-add" for it is enqueued
    Then it is created and the queue lists it pending

  @TSQUT-B02 @unit-level
  Scenario: The same target and step are not enqueued twice
    Given the live task "1-spec-add" for "AddItem.spec.md" and the step "implement"
    When "2-spec-add" for the same target and step is enqueued
    Then nothing is created and the pending count is 1

  @TSQUT-B03 @unit-level
  Scenario: Listing gives the live tasks sorted, and nothing without a queue
    Given a project with no queue folder, and then the tasks "2-b" and "1-a" enqueued
    When the queue is listed each time
    Then the first listing is empty and the second gives "1-a" then "2-b"

  @TSQUT-B04 @unit-level
  Scenario: A task whose target was deleted is removed
    Given the tasks for "vive.ts", which exists, and for "sonda.test.ts", which does not
    When the queue is listed
    Then only the task of "vive.ts" is listed and the other is gone from disk

  @TSQUT-B05 @unit-level
  Scenario: Claiming takes a pending task and records the worker
    Given the pending task "1-spec-add", and an empty queue in another project
    When "worker-A" claims from each
    Then the first gives "1-spec-add" claimed by "worker-A", with only its claimed file left, and the second gives nothing

  @TSQUT-B06 @unit-level
  Scenario: A pending residue of a dead claim is not served and is cleaned
    Given the task "residuo-a" whose claimed file and pending file both exist
    When "w2" claims
    Then no task is given and the pending residue is removed

  @TSQUT-B07 @unit-level
  Scenario: A done task moves to the history
    Given the claimed task "1-spec-add"
    When it is marked done
    Then the pending count is 0, it is in ".anchors/done", and the same target and step can be enqueued again

  @TSQUT-B08 @unit-level
  Scenario: Dropping deletes without history
    Given the pending task "1-doc-x"
    When it is dropped
    Then the pending count is 0 and no ".anchors/done" folder exists

  @TSQUT-B09 @unit-level
  Scenario: Old and unstamped claims are returned to pending
    Given two tasks claimed eight hours ago, and one claimed with no moment in another project
    When the queue is reclaimed
    Then all return to pending with no worker, and a returned task can be claimed again

  @TSQUT-B10 @unit-level
  Scenario: Forced reclaiming returns even recent claims
    Given a task claimed now
    When the queue is reclaimed with force
    Then 1 task is returned

  @TSQUT-B11 @unit-level
  Scenario: Recently held claims are counted
    Given two claims taken now and one taken eight hours ago
    When the recently held claims are counted before and after a reclaim, and after a forced reclaim
    Then the counts are 2, 2 and 0

  @TSQUT-B12 @unit-level
  Scenario: The next step follows the kind that changed
    Given the kinds plan, plan-draft, spec, feature, code, test, guide and an unknown kind
    When the next step is suggested for each
    Then they are spec, review-plan-draft, code, test, feature, review, review and triage, each with a reason

  @TSQUT-B13 @unit-level
  Scenario: The suggestions of the unit kinds are composable by the work command
    Given the kinds plan-draft, plan, spec, feature, code, test and guide
    When the next step is suggested for each
    Then every suggestion is a verb the work command accepts, with a reason

  @TSQUT-B15 @unit-level
  Scenario: A task whose target is an absolute path that exists is kept
    Given a task enqueued with the absolute path of a file that exists
    When the queue is listed twice
    Then both times the task is listed, and its file stays on disk

  @TSQUT-B14 @unit-level
  Scenario: The pending count counts pending and claimed tasks
    Given two pending tasks, one of which is then claimed
    When the pending count is taken
    Then it is 2

  @TSQUT-I01 @unit-level
  Scenario: Concurrent workers never claim the same task
    Given 20 pending tasks
    When 8 workers claim concurrently until the queue is empty
    Then all 20 are claimed, each exactly once

  @TSQUT-X01 @unit-level
  Scenario: A recent claim is not reclaimed
    Given a task claimed now by "worker-ativo"
    When the queue is reclaimed
    Then 0 tasks are returned

  @TSQUT-E01 @unit-level
  Scenario: Marking an unknown task done is refused
    Given a queue without the task "ghost"
    When "ghost" is marked done
    Then it is refused with "task not found"

  @TSQUT-E02 @unit-level
  Scenario: Dropping an unknown task is refused
    Given a queue without the task "nao-existe"
    When "nao-existe" is dropped
    Then an error is returned

  @TSQUT-E03 @unit-level
  Scenario: A corrupted task file is listed as triage, and can be claimed and dropped
    Given a queue holding the valid task "1-a" and the file "pending__0-bad.yaml" that is not YAML
    When the queue is listed
    Then "0-bad" and "1-a" are listed with no error, "0-bad" as a pending triage task whose reason names its file
    And "0-bad" can be claimed and then dropped, leaving only "1-a"

  @TSQUT-E04 @unit-level
  Scenario: Enqueuing an ID a live task already holds for another target is refused
    Given a live task "1-x" on "a.go"
    When a task "1-x" on "b.go" is enqueued, before and after the first is claimed
    Then both times it is refused naming "1-x", and the first task is untouched
