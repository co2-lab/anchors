# language: en
# @anchors
#   code: STFTC
#   ref: PRSTP
#   updated_at: 2026-10-03
#   layer: feature

@PRSTP
Feature: ProjectStatus — where the project stands in the cycle, and the one next step

  @PRSTP-B01 @unit-level
  Scenario: A directory with no git repository stops at git init
    Given a directory outside any git repository
    When the status runs
    Then it names `git init` as the next step
    And it does not mention the DISCOVER phase

  @PRSTP-B02 @unit-level
  Scenario: Without the git binary status warns and goes on
    Given a git repository with nothing in it and no git binary on the PATH
    When the status runs
    Then it prints "⚠ git not installed"
    And it goes on to name the DISCOVER phase

  @PRSTP-B03 @unit-level
  Scenario: A project with neither PROJECT.md nor configuration is not started
    Given an empty git repository
    When the status runs
    Then it names the DISCOVER phase as the NEXT STEP

  @PRSTP-B04 @unit-level
  Scenario: A project with PROJECT.md and no configuration is sent to init
    Given a git repository with PROJECT.md and no anchors.yaml
    When the status runs
    Then it prints "✓ PROJECT.md exists" and names `anchors init`

  @PRSTP-B05 @unit-level
  Scenario: A configured project with no map is sent to the map build
    Given a project with PROJECT.md and anchors.yaml and no map
    When the status runs
    Then it prints "○ no map" and does not show the queue

  @PRSTP-B06 @unit-level
  Scenario: The configuration and the map are counted, and clean informative gates are named
    Given a project with one layer, one informative gate that passes everywhere and a map of one node
    When the status runs
    Then it prints "✓ anchors.yaml (1 layers, 1 gates)" and "✓ map (1 nodes, 0 edges)"
    And it prints "○ 1 informative gate(s) CLEAN: always-green"

  @PRSTP-B07 @unit-level
  Scenario: The local queue names the first pending step in the cycle's order
    Given local projects with work in doing, in todo, only tasks, and nothing
    When the status runs on each
    Then each names only its first pending step: doing, then todo, then `anchors next`, then nothing pending

  @PRSTP-B08 @unit-level
  Scenario: An assembled local project with no work is sent to the first plan
    Given a local project with PROJECT.md, anchors.yaml and an empty map
    When the status runs
    Then it prints "queue: local" and "project is assembled and has no work yet"

  @PRSTP-B09 @unit-level
  Scenario: The github queue with missing pipelines stops at the doctor fix
    Given a github-mode project for acme/app with no workflow pipeline
    When the status runs
    Then it prints "queue: GitHub (acme/app, label [anchors])" and that pipelines are missing
    And it does not state the pull-request flow

  @PRSTP-B10 @unit-level
  Scenario: The github queue states the pull-request flow and the protected branches
    Given a github-mode project with its pipelines, integration branch develop and protected branches develop and main
    When the status runs
    Then it prints "work enters via PR to `develop` · protected: develop, main"

  @PRSTP-B11 @unit-level
  Scenario: The agent's own open cards come before claiming new work
    Given a github-mode project where the agent host/session owns card #12 "Fix login" in progress
    When the status runs as that agent
    Then it lists "#12 Fix login [in-progress]" under "YOU ALREADY HAVE WORK"
    And it does not tell the agent to claim new work

  @PRSTP-B12 @unit-level
  Scenario: Without an agent identity the next step is to claim work
    Given a github-mode project with its pipelines and no agent name in the environment
    When the status runs
    Then it tells the reader to ask the claim pipeline for work

  @PRSTP-B13 @unit-level
  Scenario: A github project with no real work is sent to the first plan
    Given a github-mode project with its pipelines whose map holds only a guide
    When the status runs
    Then it prints "project is assembled and has no work yet"

  @PRSTP-I01 @unit-level
  Scenario: Status never names a step beyond the first one missing
    Given a git repository with PROJECT.md and no anchors.yaml
    When the status runs
    Then the output does not mention the map build

  @PRSTP-X01 @unit-level
  Scenario: Status leaves the project as it found it
    Given a local project with an issue in todo and a pending task
    When the status runs
    Then every file of the project is as it was, and none was added

  @PRSTP-E01 @unit-level
  Scenario: A configuration that does not load fails the status
    Given a project whose anchors.yaml holds an unknown key
    When the status runs
    Then it fails naming "load anchors.yaml"
