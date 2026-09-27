# language: en
# @anchors
#   ref: ININT
#   updated_at: 2026-09-26
#   layer: feature

@ININT
Feature: InitNonInteractive — the init that an agent answers with flags: it asks in JSON, and writes only a complete, valid set of answers

  @ININT-B01 @unit-level
  Scenario: Without answers the command only asks
    Given an empty directory
    When init runs with --non-interactive and no answer
    Then it prints JSON with escrito false, precisa_descobrir true and the artifacts, colocation and workflow questions
    And no anchors.yaml is written

  @ININT-B02 @unit-level
  Scenario: The given answers reach the configuration
    Given an empty directory
    When init runs with --artifacts=spec,feature,test --colocation --workflow=manual
    Then the configuration has the spec, feature and test layers, default gates and the manual workflow
    And guides/HEADER_GUIDE.md exists

  @ININT-B03 @unit-level
  Scenario: A false flag is a deliberate no
    Given an empty directory
    When init runs with --gates=false --header=false among its answers
    Then no gate is seeded
    And no header guide is written

  @ININT-B04 @unit-level
  Scenario: The github workflow carries the repository and labels
    Given an empty directory
    When init runs with --workflow=github --repo=acme/app --labels=anchors,work
    Then the workflow is github for acme/app with the labels anchors and work

  @ININT-B05 @unit-level
  Scenario: One invalid answer refuses the whole set
    Given an empty directory
    When init runs with --artifacts=spec --workflow=github and no repository
    Then it fails
    And the JSON has escrito false and marks repo as refused
    And no anchors.yaml is written

  @ININT-B06 @unit-level
  Scenario: Defaults writes when asked to
    Given an empty directory
    When init runs with --defaults only
    Then it writes a configuration that loads

  @ININT-B07 @unit-level
  Scenario: A stack preset fills the code layers
    Given a project with a src/modules directory, with and without code in it
    When init runs with --preset=node-ts --artifacts=spec
    Then the configuration has the modules, core and common layers

  @ININT-B08 @unit-level
  Scenario: Layers prunes the other code layers
    Given a project with code in src/api and src/web
    When init runs with --layers naming one of the two
    Then only the chosen layer remains
    And and with a preset and no --layers, every preset layer remains

  @ININT-B09 @unit-level
  Scenario: The success document names the file and the next step
    Given an empty directory, then a project with code
    When the answers are accepted
    Then escrito is true and arquivo is the anchors.yaml path
    And the next step is anchors guide project for the empty one and anchors map build for the one with code

  @ININT-I01 @unit-level
  Scenario: Either the whole set is written or nothing is
    Given a set with one refused answer
    When init runs
    Then no anchors.yaml exists

  @ININT-X01 @unit-level
  Scenario: The non-interactive mode never prompts
    Given no terminal
    When init runs with --non-interactive
    Then the only output is one JSON document

  @ININT-E01 @unit-level
  Scenario: A malformed governs rule is refused with the expected form
    Given an empty directory
    When init runs with --governs guides/A.md
    Then it fails naming GUIDE=tag1,tag2

  @ININT-B10 @unit-level
  Scenario: The governs rules reach the configuration, one rule per tag
    Given an empty project
    When init runs with --governs guides/STYLE.md=backend,web and --governs guides/API.md=backend
    Then the written anchors.yaml has the governs rules API.md=backend, STYLE.md=backend and STYLE.md=web
