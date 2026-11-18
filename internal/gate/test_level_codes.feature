# language: en
# @anchors
#   code: TLCFT
#   ref: TLVCD
#   updated_at: 2026-10-03
#   layer: feature

@TLVCD
Feature: TestLevelCodes — each scenario references only codes its test level accepts

  @TLVCD-B01 @unit-level
  Scenario: A node that is not a feature is skipped
    Given a node whose kind is not feature
    When the gate confronts it
    Then it returns Skip

  @TLVCD-B02 @unit-level
  Scenario: A project with no filter per level is skipped
    Given a configuration that declares no levels on a test-level-codes gate
    When the gate confronts a feature
    Then it returns Skip

  @TLVCD-B03 @unit-level
  Scenario: A feature with no filtered level is skipped
    Given a filter declared for a level no scenario of the feature uses
    When the gate confronts the feature
    Then it returns Skip

  @TLVCD-B04 @unit-level
  Scenario: A level with allow accepts only matching codes
    Given a visual level that allows only VR codes
    When a scenario of that level references a rule code
    Then the gate fails naming the rule code

  @TLVCD-B05 @unit-level
  Scenario: A level with exclude refuses matching codes even when allowed
    Given a level that excludes VR codes
    When a scenario of that level references a VR code
    Then the gate fails naming the VR code

  @TLVCD-B06 @unit-level
  Scenario: Every code of the scenario is confronted without its suffix
    Given a scenario whose second code is refused and whose first carries a scenario suffix
    When the gate confronts the feature
    Then it names the second code, and the first without its suffix

  @TLVCD-B07 @unit-level
  Scenario: The failure names each refused code with its level, sorted
    Given two refused codes under two levels
    When the gate confronts the feature
    Then the message lists both with their levels in order

  @TLVCD-B08 @unit-level
  Scenario: Accepted codes pass
    Given scenarios whose codes their levels accept
    When the gate confronts the feature
    Then it returns Pass

  @TLVCD-B09 @unit-level
  Scenario: Levels come from the gate's own entries, and two entries add their lists
    Given two test-level-codes entries excluding different codes for the same level
    When the gate confronts scenarios carrying each code
    Then both exclusions hold, and levels declared on another check's entry are not read
