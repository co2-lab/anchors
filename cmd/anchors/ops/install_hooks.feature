# language: en
# @anchors
#   ref: INHKN
#   updated_at: 2026-09-26
#   layer: feature

@INHKN
Feature: InstallHooks — the git hooks that confront every commit and push with the gates and the freeze, installed without taking a hook the user wrote

  @INHKN-B01 @unit-level
  Scenario: The hooks go where git looks for them
    Given a repository without and then with core.hooksPath
    When the hooks directory is resolved
    Then it is <root>/.git/hooks by default
    And <root>/.githooks for the relative .githooks
    And the absolute path itself when absolute

  @INHKN-B02 @unit-level
  Scenario: A fresh install writes the three managed hooks, executable
    Given a governed repository with no hooks
    When install-hooks runs
    Then pre-commit, commit-msg and pre-push hold the managed scripts and are executable

  @INHKN-B03 @unit-level
  Scenario: Foreign hooks are respected unless forced
    Given a repository where the user wrote pre-commit, commit-msg and pre-push
    When install-hooks runs, then again without the pre-commit, then with --force
    Then the first run fails pointing at --force and overwrites nothing
    And the second installs the pre-commit and warns about the other two, keeping them
    And --force replaces them

  @INHKN-B04 @unit-level
  Scenario: The hooks anchors wrote, old or new, are updated on reinstall
    Given a commit-msg with the legacy anchors header, then one with the current marker but changed
    When install-hooks runs each time
    Then the commit-msg becomes the current script
    And it is not reported as foreign

  @INHKN-B05 @unit-level
  Scenario: The pre-commit refuses a commit while the remote is frozen
    Given a remote whose anchors.yaml is frozen for "rotate the key"
    When a file is committed
    Then the commit fails with "COMMIT REFUSED — the project is FROZEN" and "rotate the key"

  @INHKN-B06 @unit-level
  Scenario: A staged set with nothing governed passes the pre-commit
    Given gates that answer exit 3
    When package.json is committed
    Then the commit succeeds
    And the pre-commit says no staged file is governed and reports no failure

  @INHKN-B07 @unit-level
  Scenario: A gate failure is deferred to the commit-msg, which blocks it
    Given gates that fail
    When a file is committed
    Then the pre-commit says "gates failed. If it is deliberate"
    And the commit-msg says the commit is BLOCKED and the commit fails

  @INHKN-B08 @unit-level
  Scenario: The pre-push refuses a push while the remote is frozen
    Given a remote whose anchors.yaml is frozen and a commit made past the hooks
    When the branch is pushed
    Then the push fails with "PUSH REFUSED — the project is FROZEN"

  @INHKN-B09 @unit-level
  Scenario: The commit-msg refuses a subject the message check refuses
    Given a message check that refuses
    When a file is committed
    Then the commit fails

  @INHKN-B10 @unit-level
  Scenario: The pre-push refuses a binary older than the remote's minimum version
    Given a remote declaring min_version 2.0.0
    When the branch is pushed with a binary 1.9.0, then with 2.0.0
    Then the first push fails saying it requires 'anchors' 2.0.0 or newer
    And the second succeeds

  @INHKN-B11 @unit-level
  Scenario: Both merge drivers are registered in git config and .gitattributes
    Given a governed repository
    When install-hooks runs
    Then .gitattributes holds "anchors.graph.yaml merge=anchors-map" and "*-progress.md merge=anchors-progress" once
    And git config holds both driver commands

  @INHKN-I01 @unit-level
  Scenario: Reinstalling never duplicates an attribute line nor damages the user's lines
    Given a .gitattributes holding "*.png binary" with no final newline
    When install-hooks runs twice
    Then the file still starts with "*.png binary"
    And the map driver line appears once
    And the second run says the drivers are already configured

  @INHKN-X01 @unit-level
  Scenario: A hook the user wrote is never replaced without --force
    Given a pre-commit the user wrote
    When install-hooks runs without --force
    Then the user's pre-commit is byte for byte what it was

  @INHKN-E01 @unit-level
  Scenario: A directory without anchors.yaml gets no hook
    Given a git repository with no anchors.yaml
    When install-hooks runs
    Then it fails pointing at anchors init
    And no pre-commit exists

  @INHKN-E02 @unit-level
  Scenario: Outside a repository the error explains what needs git
    Given a directory with anchors.yaml outside git
    When install-hooks runs
    Then it fails explaining it needs git to install the pre-commit
