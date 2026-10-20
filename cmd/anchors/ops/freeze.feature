# language: en
# @anchors
#   code: FRFTF
#   ref: FRZEX
#   updated_at: 2026-10-20
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

  @FRZEX-B14 @unit-level
  Scenario: A freeze with a scope is recorded in the doctrine and marks every place
    Given a project with a doctrine document, two units and a file of its own
    When a unit, two rules of the other unit and the file are frozen with a reason, and then another rule
    Then the doctrine records QLDTQ-F01 with the reason, each place carries its mark, the map covers the unit's code too, and the second freeze is QLDTQ-F02

  @FRZEX-B15 @unit-level
  Scenario: Without a doctrine the marks carry the reason, and what cannot be frozen is refused
    Given a project with no doctrine, and one with two doctrine documents
    When a file is frozen in the first, and in the second with no reason, with no --in, a unit with no owner, an undefined rule, a missing file, and a part frozen twice
    Then the first marks the reason and says it is not traceable, --in records in the document it names, and each of the others is refused naming its cause

  @FRZEX-B16 @unit-level
  Scenario: A recorded freeze gains a scope by its code
    Given a unit frozen by QLDTQ-F01
    When QLDTQ-F01 is given another unit, and an unrecorded code is given one
    Then the other unit is marked with no new record, and the unrecorded code is refused

  @FRZEX-B17 @unit-level
  Scenario: Thawing a freeze removes its marks and its record, and nothing else
    Given a unit, a rule and a file frozen by QLDTQ-F01, and another rule by QLDTQ-F02
    When QLDTQ-F01 is thawed, then the other rule by its scope, then an unknown code
    Then the unit's spec is as it was, the other freeze's mark and record stay, the scope's thaw removes its mark, and the unknown code is refused

  @FRZEX-B18 @unit-level
  Scenario: The freezes are listed, each with its reason and places
    Given a project with nothing frozen, then a unit frozen and an entry that marks nothing
    When the freezes are shown, all and for one code
    Then nothing frozen is said, then each freeze has its reason and places, the orphan marks nothing, and one code lists the spec and the code it reaches

  @FRZEX-B19 @unit-level
  Scenario: Before an edit, the agent is told a file is frozen, and the edit goes on
    Given a frozen file and a spec with a frozen rule, in a project with its map
    When the hook is called for each, for a file not frozen, for a file outside the project and for broken input
    Then the frozen file is told with its freeze and reason and never denied, the spec names its frozen rule, and the others say nothing

  @FRZEX-B20 @unit-level
  Scenario: A folder is frozen by a glob written to anchors.yaml
    Given a project with a doctrine document, and a recorded freeze
    When a folder is frozen with --path, another glob is added to the same freeze, one glob is thawed with --path, and the freeze is thawed by its code
    Then the block is created with the first entry, grows by an entry, loses the thawed glob, and is gone after the thaw; show lists the globs; a glob matching no file is refused

  @FRZEX-B21 @unit-level
  Scenario: A new file in a frozen folder is told before the edit
    Given a folder frozen in anchors.yaml and no map that knows a file the agent creates there
    When the agent's hook is asked about that file
    Then it says the file is frozen, and the edit goes through

  @FRZEX-B22 @unit-level
  Scenario: The freezes section follows the project's language
    Given a project in Portuguese whose doctrine document has no freezes section, and one in English that has it already
    When a freeze is recorded in each
    Then the first gets the section as Congelamentos, the second keeps its own, and a thaw removes either
