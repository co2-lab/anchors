# language: en
# @anchors
#   code: DHFDP
#   ref: DEPHN
#   updated_at: 2026-10-08
#   layer: feature

@DEPHN
Feature: DependencyHonored — a spec declares no dependency: the code does

  @DEPHN-B09 @unit-level
  Scenario: A spec still carrying a Dependencies table diverges, pointing at the migration
    Given a spec with a Dependencies table of two rows, a spec citing a DEPn only in a rule's uses, and a code file
    When dependency-honored runs on each
    Then the first diverges naming two rows and the migration, the second passes, and the code file leaves without a verdict
