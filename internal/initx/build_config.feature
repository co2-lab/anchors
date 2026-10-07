# language: en
# @anchors
#   code: BCFBL
#   ref: BLCNB
#   updated_at: 2026-10-07
#   layer: feature

@BLCNB
Feature: BuildConfig — builds the configuration that inference proposes as the default for the init questions

  @BLCNB-B01 @unit-level
  Scenario: Each detected code directory becomes a code layer named after its last segment
    Given an inference that found the code directories apps/mobile and packages/backend with the extension .ts
    When the proposal is built
    Then the layers mobile-code and backend-code exist, each a code layer tagged with its own name
    And their patterns are apps/mobile/**/*.ts and packages/backend/**/*.ts
    And code directories src/handlers and lib/handlers are named src-handlers-code and lib-handlers-code
    And the root code directory is named root-code

  @BLCNB-B02 @unit-level
  Scenario: With several detected extensions the pattern lists them as a set
    Given an inference that found the code directory a/b with the extensions .ts and .tsx
    When the proposal is built
    Then the pattern of b-code is a/b/**/*.{ts,tsx}
    And with src and src/handlers detected, src covers only its own files, src/*.go, and src/handlers covers src/handlers/**/*.go
    And the root covers only its own files, *.go

  @BLCNB-B03 @unit-level
  Scenario: A proposed code layer excludes specs, features and test files
    Given an inference that found code directories
    When the proposal is built
    Then each code layer excludes spec markdown, feature and test files
    And a project whose tests follow `_test.go` excludes **/*_test.go
    And a Python project with no test yet excludes **/test_*.py
    And with no convention and no family the test files excluded are **/*.test.*

  @BLCNB-B04 @unit-level
  Scenario: Colocation is proposed only when detected, with templates only for the detected kinds
    Given an inference that detected colocation and found only features
    When the proposal is built
    Then the colocation is anchored on the spec with only the feature template
    And an inference that detected no colocation proposes none, even with features and tests found
    And an inference that detected colocation with specs, features and tests has only the feature and test templates, never a spec one

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

  @BLCNB-B06 @unit-level
  Scenario: The colocated test template follows the project's test convention
    Given an inference that detected colocation with tests following `_test.go`
    When the proposal is built
    Then the test template is {{dir}}/{{name}}_test.go
    And with tests of no known convention it is {{dir}}/{{name}}.test.{{ext}}

  @BLCNB-B07 @unit-level
  Scenario: The proposal's dialect is the family inference found
    Given an inference that found the family go
    When the proposal is built
    Then the dialect family is go
    And an inference with no family proposes no dialect

  @BLCNB-B08 @unit-level
  Scenario: The test layer's pattern is the project's test convention
    Given tests following `.spec.ts` and `.test.ts`
    When the test pattern is asked
    Then it is {**/*.spec.ts,**/*.test.ts}
    And a Go project with no test yet gets **/*_test.go
    And with no convention and no family it is **/*.test.*

  @BLCNB-B09 @unit-level
  Scenario: A new project's dialect carries what its family knows can fail
    Given a proposal for a ts project and one with no family
    When the config is built
    Then the ts project's dialect has the family's fallible patterns, and the other has no dialect
