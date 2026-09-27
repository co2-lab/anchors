# language: en
# @anchors
#   ref: BLCNB
#   updated_at: 2026-09-26
#   layer: feature

@BLCNB
Feature: BuildConfig — builds the configuration that inference proposes as the default for the init questions

  @BLCNB-B01 @unit-level
  Scenario: Each detected code directory becomes a code layer named after its last segment
    Given an inference that found the code directories apps/mobile and packages/backend with the extension .ts
    When the proposal is built
    Then the layers mobile-code and backend-code exist, each a code layer tagged with its own name
    And their patterns are apps/mobile/**/*.ts and packages/backend/**/*.ts

  @BLCNB-B02 @unit-level
  Scenario: With several detected extensions the pattern lists them as a set
    Given an inference that found the code directory a/b with the extensions .ts and .tsx
    When the proposal is built
    Then the pattern of b-code is a/b/**/*.{ts,tsx}

  @BLCNB-B03 @unit-level
  Scenario: A proposed code layer excludes specs, features and test files
    Given an inference that found code directories
    When the proposal is built
    Then each code layer excludes spec markdown, feature and test files

  @BLCNB-B04 @unit-level
  Scenario: Colocation is proposed only when detected, with templates only for the detected kinds
    Given an inference that detected colocation and found only features
    When the proposal is built
    Then the colocation is anchored on the spec with only the feature template
    And an inference that detected no colocation proposes none, even with features and tests found

  @BLCNB-B05 @unit-level
  Scenario: The test handle is proposed only when inference found one
    Given an inference that found the handle testID and no colocation
    When the proposal is built
    Then the proposal carries testID and no colocation template
    And an inference that found no handle proposes none
    And with colocation, the handle sits beside the colocation templates

  @BLCNB-X01 @unit-level
  Scenario: The proposal creates no artifact layer and no governs rule
    Given an inference that found specs, features, tests, guides, plans and the code directory src/app
    When the proposal is built
    Then the only layer is app-code
    And there is no governs rule
