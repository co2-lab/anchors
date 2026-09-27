# language: en
# @anchors
#   ref: CFPCD
#   updated_at: 2026-09-26
#   layer: feature

@CFPCD
Feature: CodeFromPath — the most meaningful unique code for a unit, given its file path

  @CFPCD-B01 @unit-level
  Scenario: Generic file names are recognised in any case
    Given the names "handler", "index", "resource", "Handler" and "INDEX", and "metadata", "Login" and "useTransactions"
    When each is checked for being generic
    Then the first five are generic and the last three are not

  @CFPCD-B02 @unit-level
  Scenario: A generic file name takes its code from the parent folder
    Given the path "packages/backend/functions/manage-metadata/handler.ts" and the root file "handler.ts"
    When their codes are generated with nothing taken
    Then the first gets the code of "manage-metadata" and not the code of "handler"
    And the root file gets the code of "handler"

  @CFPCD-B03 @unit-level
  Scenario: Artifact suffixes are dropped from the base name
    Given the paths "a/Login.spec.md", "a/Login.feature", "a/Login.test.tsx" and "a/Login.tsx"
    When their codes are generated with nothing taken
    Then all four get the same code

  @CFPCD-B04 @unit-level
  Scenario: A normal file name gets the generated code when free
    Given the path "packages/backend/models/metadata.spec.md" with nothing taken
    When its code is generated
    Then it is the generated code of "metadata"

  @CFPCD-I01 @unit-level
  Scenario: Same-named units get distinct, deterministic codes
    Given the code of "packages/backend/models/metadata.spec.md" already taken
    When the code of "packages/backend/repositories/metadata.spec.md" is generated twice
    Then both answers are equal, differ from the first unit's code, and are not taken
