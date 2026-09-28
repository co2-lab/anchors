# language: en
# @anchors
#   ref: NWARN
#   updated_at: 2026-09-28
#   layer: feature

@NWARN
Feature: NewArtifact — a new artifact is born beside its unit, with a resolved identity and the sections of the project's ruler

  @NWARN-B01 @unit-level
  Scenario: An unknown kind is refused
    Given the kind widget
    When new runs
    Then it fails with unknown kind "widget"

  @NWARN-B02 @unit-level
  Scenario: The name and the output path are required
    Given a spec with no name, then one with no --out
    When new runs
    Then it asks for the unit name
    And it asks for --out

  @NWARN-B03 @unit-level
  Scenario: A spec is born with a code no unit uses
    Given a map where another unit holds the canonical code of Login
    When new spec Login runs
    Then the header code is another code
    And and a path in a layer with prefix AU gets a code starting with AU

  @NWARN-B04 @unit-level
  Scenario: A feature or test takes its identity from the sibling spec, or warns it is orphaned
    Given src/auth/Login.spec.md with code LGNSP
    When new feature Login runs with --out src/auth/Login.feature, then Orphan elsewhere
    Then the first reads LGNSP from src/auth/Login.spec.md and writes it
    And the second warns that no spec was found and it is born ORPHANED

  @NWARN-B05 @unit-level
  Scenario: Code pins the identity
    Given a plan Foundation
    When new plan runs with --code FNDTN
    Then the plan's header holds code FNDTN

  @NWARN-B06 @unit-level
  Scenario: With and without are validated against the kind's sections
    Given a spec
    When new runs with --with errors --without overview, then with --with nope
    Then errors is chosen and overview is not
    And nope is refused as an unknown section

  @NWARN-B07 @unit-level
  Scenario: A preset fixes the sections and their order, and extra sections follow
    Given the store preset, then the validation preset with --with auth
    When the sections are resolved and sorted
    Then state shape comes before invariants, actions after state shape, and open closes
    And auth comes last

  @NWARN-B08 @unit-level
  Scenario: A spec for a declarative layer is refused
    Given a layer dao with regime declarativo
    When new spec UserDao runs with --out src/dao/UserDao.spec.md
    Then it fails naming dao as a RECOGNIZED layer
    And no spec is written

  @NWARN-B09 @unit-level
  Scenario: Section titles and bodies follow the project lexicon, then its language
    Given the screen preset
    When it is rendered with lang en, es and pt-BR, with a declared overview title, with a rule type for B, and in a layer with its own title
    Then the titles follow each language, English by default
    And the declared title, the rule type's section and the layer's title win in that order
    And no Portuguese body text leaks into an English artifact

  @NWARN-B10 @unit-level
  Scenario: Features and tests follow the project's declared dialect
    Given a dialect with Gherkin language pt and family python
    When a feature and a test are rendered
    Then the feature opens with "# language: pt" and uses Funcionalidade
    And the test holds def test_calc_total( and no describe(

  @NWARN-B11 @unit-level
  Scenario: The unit regime tag comes from the project, or a visible placeholder
    Given regime mappings in Portuguese, English and the project's own words, and configs with none
    When the unit regime tag is taken
    Then it is @nivel-unit, @level-unit and @fast
    And with none it is a placeholder asking for the mapping

  @NWARN-B12 @unit-level
  Scenario: The artifact is written where out says and never overwritten
    Given a map and new spec Login with --out src/auth/Login.spec.md
    When new runs twice
    Then the first writes the file and names its code and "anchors check --changed src/auth/Login.spec.md"
    And the second fails because the file already exists

  @NWARN-B13 @unit-level
  Scenario: A plan is born with its progress companion
    Given a plan at plans/0001-foundation.md
    When new plan runs
    Then plans/0001-foundation-progress.md exists
    And no orphan warning is printed

  @NWARN-B14 @unit-level
  Scenario: List-sections prints the menu, with presets only for specs
    Given the spec and feature kinds
    When new runs with --list-sections
    Then the spec menu lists default and optional sections and the presets
    And the feature menu lists no presets

  @NWARN-B15 @unit-level
  Scenario: The target layer is the one of the unit the artifact describes
    Given a tsx screen layer and a ts catch-all layer
    When the target layer of app/P.spec.md and app/Q.spec.md is resolved with Q.ts on disk
    Then P gets screen
    And Q gets catchall

  @NWARN-I01 @unit-level
  Scenario: A refused new leaves nothing behind
    Given an empty root
    When every refused new runs
    Then the root is still empty

  @NWARN-X01 @unit-level
  Scenario: A feature references the spec's identity instead of owning one
    Given src/auth/Login.spec.md with code LGNSP
    When new feature Login runs beside it
    Then the feature's header references LGNSP

  @NWARN-B16 @unit-level
  Scenario: The help and the unknown-kind refusal name every kind, and --out is mandatory
    Given the new command
    When its help is read and it runs with the kind widget
    Then the help and the refusal name action, feature, flow, plan, product, spec and test
    And the --out help says it is mandatory and names no default

  @NWARN-B17 @unit-level
  Scenario: The new artifact enters the map at once
    Given a project with a configuration and a map
    When a spec is created with anchors new
    Then the map has its node, and the output says it was added
