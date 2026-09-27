# language: en
# @anchors
#   ref: FRZEX
#   updated_at: 2026-09-26
#   layer: feature

@FRZEX
Feature: Freeze — the project is stopped in three layers with a written reason, and thawed by undoing exactly those layers

  @FRZEX-B01 @unit-level
  Scenario: A blank reason is refused
    Given a governed repository
    When freeze runs with the reason "  "
    Then it fails saying --reason is required
    And anchors.yaml has no enabled line

  @FRZEX-B02 @unit-level
  Scenario: The freeze writes two quoted lines on top and the file still loads
    Given a configuration "version: 1"
    When the reason a: "quoted" reason is written
    Then the file starts with enabled: false and the quoted reason
    And and a reason with a colon still loads as frozen with that reason

  @FRZEX-B03 @unit-level
  Scenario: The freeze is committed and pushed past refusing hooks
    Given a repository whose pre-commit, commit-msg and pre-push hooks all refuse
    When freeze runs with the reason "release blocked"
    Then the remote's last subject is "chore(anchors): freeze — release blocked"

  @FRZEX-B04 @unit-level
  Scenario: No-push leaves the frozen file as a local change
    Given a repository in local mode pushed to a remote
    When freeze runs with --no-push
    Then the remote's last subject is still "config"
    And anchors.yaml shows as changed in the working tree

  @FRZEX-B05 @unit-level
  Scenario: In github mode the freeze creates the rule and opens the issue
    Given a repository in github mode for acme/app
    When freeze runs
    Then gh is asked to POST repos/acme/app/rulesets and to create the issue with the freeze title
    And the issue link is reported
    And and in local mode gh is never called

  @FRZEX-B06 @unit-level
  Scenario: No-ruleset skips the ruleset and still opens the issue
    Given a repository in github mode
    When freeze runs with --sem-ruleset
    Then no rulesets call is made
    And the issue is created

  @FRZEX-B07 @unit-level
  Scenario: Freezing twice changes nothing and shows the reason
    Given a project frozen for "leaked credential: rotate first"
    When freeze runs again with another reason
    Then it says it is ALREADY frozen with the first reason

  @FRZEX-B08 @unit-level
  Scenario: Each failing layer is a warning and the local brake holds
    Given a repository in github mode with no remote and a gh that always fails
    When freeze and then thaw run
    Then freeze warns could not push, ruleset not created, issue not opened and succeeds with the file frozen
    And thaw warns could not push, ruleset not removed, issue not closed

  @FRZEX-B09 @unit-level
  Scenario: The thaw undoes the three layers
    Given a project frozen in github mode with rule 42 and issue 99
    When thaw runs
    Then anchors.yaml is back to the original
    And the remote's last subject starts with "chore(anchors): thaw"
    And rule 42 is deleted and issue 99 closed

  @FRZEX-B10 @unit-level
  Scenario: Thawing a project that is not frozen is a no-op
    Given a project already thawed
    When thaw runs again
    Then it says the project is not frozen

  @FRZEX-B11 @unit-level
  Scenario: With nothing on the remote the thaw removes nothing
    Given a gh that finds no rule and no issue
    When the rule is deleted and the issue closed
    Then no DELETE and no issue close is issued, and no error

  @FRZEX-B12 @unit-level
  Scenario: The deprecated Portuguese flag names still work
    Given a repository in github mode
    When freeze runs with --motivo and --sem-ruleset, and thaw with --sem-push
    Then the freeze uses the reason and skips the rule
    And the thaw does not push

  @FRZEX-B13 @unit-level
  Scenario: A config that already declares enabled is frozen with a single key and still loads
    Given a configuration with "enabled: true" and a stale freeze_reason
    When freeze runs with the reason "stop now"
    Then the file has exactly one enabled: and one freeze_reason: line
    And it loads frozen with the reason "stop now"
    And after thaw it loads not frozen

  @FRZEX-I01 @unit-level
  Scenario: Freeze and thaw give back the file byte for byte
    Given a configuration in github mode
    When it is frozen and thawed
    Then the file is exactly the original

  @FRZEX-X01 @unit-level
  Scenario: The configuration is never reserialized
    Given a configuration "version: 1"
    When the freeze is written and removed
    Then the file is "version: 1" again, with the freeze block added on top in between

  @FRZEX-E01 @unit-level
  Scenario: Freeze and thaw refuse a project without config
    Given a directory with no anchors.yaml
    When freeze and thaw run
    Then each fails with "load anchors.yaml"
