# language: en
# @anchors
#   ref: PCSDP
#   updated_at: 2026-09-26
#   layer: feature

@PCSDP
Feature: PackSeeding — copy the compliance packs carried in the binary into the project, never over an adapted one

  @PCSDP-B01 @unit-level
  Scenario: Seeding copies every carried pack byte for byte
    Given an empty project folder
    When the packs are seeded
    Then every pack the binary carries is written under "packs/<domain>/<file>"
    And each written file is identical to the carried pack

  @PCSDP-B02 @unit-level
  Scenario: An adapted pack is preserved, not overwritten
    Given a project whose "packs/privacy/lgpd.yaml" was adapted by the team
    When the packs are seeded
    Then "packs/privacy/lgpd.yaml" keeps the team's content
    And it is reported as preserved and not as created

  @PCSDP-B03 @unit-level
  Scenario: The created and preserved lists come back sorted
    Given a project with one adapted pack
    When the packs are seeded
    Then the created list is "packs/accessibility/wcag.yaml", "packs/health/hipaa.yaml", "packs/payment/pci-dss.yaml", "packs/privacy/ccpa.yaml", "packs/privacy/gdpr.yaml" in that order

  @PCSDP-B04 @unit-level
  Scenario: The available packs are grouped by domain
    Given the packs the binary carries
    When the available packs are listed
    Then "privacy" lists "privacy/ccpa", "privacy/gdpr", "privacy/lgpd" in that order
    And "accessibility", "health" and "payment" list "accessibility/wcag", "health/hipaa" and "payment/pci-dss"

  @PCSDP-I01 @unit-level
  Scenario: Seeding twice is the same as seeding once
    Given a project folder already seeded
    When the packs are seeded again
    Then nothing is created
    And all six packs are reported as preserved

  @PCSDP-X01 @unit-level
  Scenario: Seeding does not filter by the adopted jurisdiction
    Given a project folder with no pack adopted
    When the packs are seeded
    Then the packs of every domain are on disk

  @PCSDP-E01 @unit-level
  Scenario: A folder that cannot be written fails the seeding
    Given a project where "packs" is a file, so its folders cannot be created
    When the packs are seeded
    Then the seeding returns an error
