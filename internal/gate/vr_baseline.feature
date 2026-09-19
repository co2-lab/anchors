# language: en
# @anchors
#   ref: VRBSV
#   updated_at: 2026-09-19
#   layer: feature

@VRBSV
Feature: VRBaseline — ensures visual regression scenarios have captured reference baseline images

  @VRBSV-B01 @unit-level
  Scenario: Artifacts that are not feature files leave with verdict Skip
    Given an artifact whose kind is not a feature
    When checkVRBaseline evaluates the artifact
    Then it returns Skip explaining that only features are evaluated

  @VRBSV-B02 @unit-level
  Scenario: Feature files declaring no visual regression scenarios leave with verdict Skip
    Given a feature file that does not contain any visual regression scenarios
    When checkVRBaseline evaluates the feature
    Then it returns Skip indicating absence of visual regression scenarios

  @VRBSV-B03 @unit-level
  Scenario: Visual regression scenarios having matching baseline images pass
    Given a feature declaring a visual regression scenario
    And a matching baseline image exists on disk
    When checkVRBaseline evaluates the feature
    Then it returns Pass

  @VRBSV-B04 @unit-level
  Scenario: Visual regression scenarios lacking baseline images fail naming missing codes
    Given a feature declaring a visual regression scenario
    And no matching baseline image exists on disk
    When checkVRBaseline evaluates the feature
    Then it returns Fail naming the scenario code lacking a baseline

  @VRBSV-B05 @unit-level
  Scenario: Baseline image matching supports naming variants
    Given a feature declaring a visual regression scenario
    And a baseline image with a variant suffix exists on disk
    When checkVRBaseline evaluates the feature
    Then it returns Pass recognizing the variant image

  @VRBSV-B06 @unit-level
  Scenario: Reads visual regime tag from project configuration
    Given project configuration defining a custom visual regime tag mapping
    When visualRegimeTag resolves the tag
    Then it returns the configured tag key

  @VRBSV-B07 @unit-level
  Scenario: Falls back to default visual regime tag when configuration is missing
    Given nil or empty configuration
    When visualRegimeTag resolves the tag
    Then it returns the default visual tag vr-level

  @VRBSV-B08 @unit-level
  Scenario: Multiple missing baseline scenario codes are sorted alphabetically
    Given a feature declaring multiple visual regression scenarios without images
    When checkVRBaseline evaluates the feature
    Then it returns Fail listing the missing codes in alphabetical order

  @VRBSV-B09 @unit-level
  Scenario: Only scenario codes containing visual regression suffix are matched
    Given a feature line with a visual tag co-tagging non-visual scenario codes
    When checkVRBaseline extracts visual scenarios
    Then only codes containing the visual regression suffix are extracted

  @VRBSV-I01 @unit-level
  Scenario: Every visual regression scenario declared must correspond to a baseline image
    Given visual regression scenarios declared in a feature
    When checking baseline existence against disk
    Then failure is reported if any declared visual scenario lacks an image

  @VRBSV-I02 @unit-level
  Scenario: Visual regime tag mapping treats map key as tag in feature
    Given configuration with a visual regime entry where key is tag and value is regime
    When resolving the visual regime tag
    Then the key stripped of leading at is returned as the tag to search

  @VRBSV-X01 @unit-level
  Scenario: Does not evaluate baseline staleness using disk modification timestamps
    Given existing baseline images on disk
    When checkVRBaseline executes
    Then it passes based on file presence without asserting file modification times

  @VRBSV-X02 @unit-level
  Scenario: Does not fail commits based on git commit dates of baseline images
    Given baseline images on disk with old git history
    When checkVRBaseline executes
    Then it passes without inspecting git commit dates
