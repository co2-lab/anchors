# language: en
# @anchors
#   ref: ARCHR
#   updated_at: 2026-10-01
#   layer: feature

@ARCHR
Feature: ArtifactChoice — turns the artifacts the user chose at init into artifact layers and colocation

  @ARCHR-B01 @unit-level
  Scenario: The artifact options are spec, feature, test, guide, plan and code, in that order
    Given the artifacts init offers
    When the options are listed several times
    Then every listing is spec, feature, test, guide, plan, code

  @ARCHR-B02 @unit-level
  Scenario: Each artifact inference found is pre-checked, and nothing else
    Given an inference that found specs, tests, a plans folder and a code directory
    When the detected artifacts are asked for
    Then exactly spec, test, plan and code are pre-checked
    And an inference that found features and a guides folder pre-checks exactly feature and guide
    And an inference that found nothing pre-checks nothing

  @ARCHR-B03 @unit-level
  Scenario: A chosen artifact layer is created even when nothing was detected
    Given a configuration with no layer set
    And the choice spec, feature and test
    When the choice is applied
    Then each chosen artifact has a layer of its own kind, a default pattern and its name as tag

  @ARCHR-B04 @unit-level
  Scenario: An artifact layer that was not chosen is removed, and code layers are untouched
    Given a configuration with spec, test and mobile-code layers
    And the choice spec and feature
    When the choice is applied
    Then the test layer is gone
    And the feature layer was added
    And the mobile-code layer is still there

  @ARCHR-B05 @unit-level
  Scenario: A chosen artifact layer that already exists is kept as declared
    Given a configuration whose spec layer has the pattern docs/**/*.spec.md
    And the choice spec and feature
    When the choice is applied
    Then the spec layer still has the pattern docs/**/*.spec.md

  @ARCHR-B06 @unit-level
  Scenario: The guide and plan layers take the detected directory, otherwise the default pattern
    Given guides detected in docs/guides and plans detected in docs/plans
    When guide and plan are chosen
    Then the guide pattern is docs/guides/*.md and the plan pattern is docs/plans/*.md
    And with no detected folder, the plan pattern is plans/*.md with kind plan

  @ARCHR-B07 @unit-level
  Scenario: Colocation is declared from the spec as anchor, and the spec is never a derivative
    Given colocation wanted with spec, code, feature and test chosen
    When colocation is applied
    Then the anchor is the spec, with three derivatives and no spec among them
    And the feature derivative is the spec's name with the feature extension in the same folder
    And choosing every offered artifact puts the code beside the spec, with the spec's name and the declared extension

  @ARCHR-B08 @unit-level
  Scenario: No colocation is declared when it is not wanted, when spec is not chosen, or when nothing derives from the spec
    Given a configuration with a colocation declared
    When colocation is applied without spec, then turned off, then with only spec and guide chosen
    Then each time the colocation declaration is removed

  @ARCHR-B09 @unit-level
  Scenario: Choosing code or leaving it out creates and removes no layer
    Given a configuration with a code layer named code and another code layer
    When the choice is applied with code chosen, and again with nothing chosen
    Then both code layers stay exactly as they were

  @ARCHR-B10 @unit-level
  Scenario: Colocation keeps the rest of derived
    Given an inferred test handle in the derived block
    When colocation is applied, on or off
    Then the test handle is still in the derived block

  @ARCHR-B11 @unit-level
  Scenario: A created test layer takes the project's test pattern
    Given test chosen with the project's test pattern **/*_test.go
    When the choice is applied
    Then the test layer's pattern is **/*_test.go
    And with no pattern given it is **/*.test.*

  @ARCHR-B12 @unit-level
  Scenario: The colocated test template is the project's
    Given colocation wanted with spec and test chosen and the template {{dir}}/{{name}}_test.go
    When colocation is applied
    Then the test derivative is {{dir}}/{{name}}_test.go
    And with no template given it is {{dir}}/{{name}}.test.{{ext}}
