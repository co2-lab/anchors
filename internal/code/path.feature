# language: en
# @anchors
#   code: PTFTP
#   ref: CFPCD
#   updated_at: 2026-10-19
#   layer: feature

@CFPCD
Feature: CodeFromPath — the most meaningful unique code for a unit, given its file path

  @CFPCD-B01 @unit-level
  Scenario: Whether a name is generic is the caller's to say
    Given the path "functions/manage-metadata/handler.ts" with the name "handler"
    When its code is generated saying the name is not generic
    Then the code is made from "handler", not from the folder

  @CFPCD-B02 @unit-level
  Scenario: A generic file name takes its code from the parent folder
    Given the path "packages/backend/functions/manage-metadata/handler.ts" and the root file "handler.ts"
    When their codes are generated with nothing taken
    Then the first gets the code of "manage-metadata" and not the code of "handler"
    And the root file gets the code of "handler"

  @CFPCD-B03 @unit-level
  Scenario: The code comes from the unit's name the caller gives
    Given the paths "a/Login.spec.md", "a/Login.feature", "a/Login.test.tsx" and "a/Login.tsx", each with the unit's name Login
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
