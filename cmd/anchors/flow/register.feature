# language: en
# @anchors
#   code: RGFTR
#   ref: FLRGF
#   updated_at: 2026-10-09
#   layer: feature

@FLRGF
Feature: FlowRegister — attach the flow domain's commands to the root command, once each

  @FLRGF-B01 @unit-level
  Scenario: The root holds exactly the eighteen flow commands after registration
    Given an empty root command
    When the flow domain registers its commands on it
    Then the root's children are exactly backfill-labels, decided, deliver, discard, done, drop, escalate, merge-progress, next, pr-body, queue, reclaim, report-bug, task-status, unblock, watch and work

  @FLRGF-B02 @unit-level
  Scenario: The watcher's controls live under watch
    Given an empty root command
    When the flow domain registers its commands on it
    Then `watch` holds start, run, status, stop, pause, resume and logs
    And none of those names is a child of the root

  @FLRGF-I01 @unit-level
  Scenario: No flow command is attached twice
    Given an empty root command
    When the flow domain registers its commands on it
    Then every child name of the root appears once

  @FLRGF-X01 @unit-level
  Scenario: The progress command is not attached to the root
    Given an empty root command
    When the flow domain registers its commands on it
    Then the root has no `progress` child
