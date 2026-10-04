# language: en
# @anchors
#   code: PCFTA
#   ref: OBPCB
#   updated_at: 2026-10-03
#   layer: feature

@OBPCB
Feature: ObligationPack — distributable sets of obligations from a norm, resolved against the project

  @OBPCB-B01 @unit-level
  Scenario: A pack file gives its metadata and its obligations
    Given the pack file "packs/privacy/lgpd.yaml" of the norm "Lei 13.709/2018"
    When it is loaded
    Then it has the name "lgpd", the domain "privacy", the jurisdiction "br", the required value "erasure_handler" and an obligation from "Art. 18, VI"

  @OBPCB-B02 @unit-level
  Scenario: A reference is a path or a name under packs
    Given the references "privacy/lgpd", "./internal/mine.yaml", "custom/other.yml" and "/abs/elsewhere/p.yaml" in the root "/r"
    When they are resolved
    Then they are "/r/packs/privacy/lgpd.yaml", "/r/internal/mine.yaml", "/r/custom/other.yml" and "/abs/elsewhere/p.yaml"

  @OBPCB-B03 @unit-level
  Scenario: Placeholders are replaced by the project's values
    Given the lgpd pack placing a duty in "{{erasure_handler}}", "{{ audit_log }}/x.go" and "static.go"
    And the project values "api/erase.go" and "api/audit"
    When the packs are loaded
    Then the places are "api/erase.go", "api/audit/x.go" and "static.go"

  @OBPCB-B04 @unit-level
  Scenario: Packs come back sorted by name
    Given the references "privacy/lgpd" and "accessibility/wcag"
    When the packs are loaded
    Then they come back as "lgpd" then "wcag"

  @OBPCB-B05 @unit-level
  Scenario: A pack of an undeclared jurisdiction is skipped with a warning
    Given a project operating only in "br" and the packs "gdpr" of "EU" and "wcag" of "global"
    When the packs are loaded
    Then only "wcag" is loaded and a warning says "pack `gdpr` belongs to jurisdiction `EU`"

  @OBPCB-B06 @unit-level
  Scenario: Global packs and projects without jurisdictions load everything
    Given a project declaring " BR " with the "lgpd" and global "wcag" packs, and a project declaring no jurisdictions with the "gdpr" pack
    When the packs are loaded
    Then every pack loads with no warning

  @OBPCB-E01 @unit-level
  Scenario: A pack without a name is refused
    Given a pack file with obligations but no name
    When it is loaded
    Then it is refused with "without `name`"

  @OBPCB-E02 @unit-level
  Scenario: A pack without obligations is refused
    Given the pack file of "empty" with no obligations
    When it is loaded
    Then it is refused with "empty: pack without obligations"

  @OBPCB-E03 @unit-level
  Scenario: An invalid or missing pack file is refused
    Given the file "bad.yaml" holding invalid YAML, and a reference to a pack that does not exist
    When each is loaded
    Then the first is refused naming "bad.yaml" and the second with the read error

  @OBPCB-E04 @unit-level
  Scenario: Unresolved placeholders refuse the whole load
    Given the lgpd pack and a project that declares no values
    When the packs are loaded
    Then the load is refused with "pack `lgpd` needs the project to declare `audit_log`, `erasure_handler`"
