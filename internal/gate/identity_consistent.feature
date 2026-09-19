# language: en
# @anchors
#   ref: IDCND
#   updated_at: 2026-09-19
#   layer: feature

@IDCND
Feature: IdentityConsistent — a unit's spec identity must match its exposed testID and visual baseline

  @IDCND-B01 @unit-level
  Scenario: Confronting an artifact that is not a spec skips
    Given a node whose kind is not spec
    When the gate confronts it
    Then it returns Skip, because identity is declared in the spec

  @IDCND-B02 @unit-level
  Scenario: Confronting without a graph returns Pending
    Given a spec node and no map graph
    When the gate confronts it
    Then it returns Pending, because without a map it cannot inspect identity concordance

  @IDCND-B03 @unit-level
  Scenario: Confronting a spec without a declared code skips
    Given a spec with an empty or missing code
    When the gate confronts it
    Then it returns Skip, avoiding duplicate findings from the code presence gate

  @IDCND-B04 @unit-level
  Scenario: A testID prefix matching the spec identity code passes
    Given a code unit reached via specifies that exposes a testID with the spec's code
    When the gate confronts it
    Then it returns Pass, confirming the testID matches the spec identity

  @IDCND-B05 @unit-level
  Scenario: An orphan testID prefix with code shape fails
    Given a code unit that exposes a testID with a code-shaped prefix belonging to no unit in the map
    When the gate confronts it
    Then it returns Fail, preventing rogue acronyms from entering manual spelling dictionaries

  @IDCND-B06 @unit-level
  Scenario: A testID prefix matching another declared unit is accepted as legitimate reuse
    Given a code unit that exposes a testID prefixed with the code of another declared map unit
    When the gate confronts it
    Then it returns Pass, allowing components to identify where they appear without breaking flows

  @IDCND-B07 @unit-level
  Scenario: A visual regression baseline with a divergent code fails
    Given a visual regression baseline file alongside the unit whose code differs from the spec
    When the gate confronts it
    Then it returns Fail, because the baseline is the physical proof of this specific unit

  @IDCND-B08 @unit-level
  Scenario: Short testID prefixes of three letters or fewer pass
    Given a code unit that exposes testIDs with short prefixes of three letters or fewer
    When the gate confronts it
    Then it returns Pass, recognizing they do not have code shape

  @IDCND-B09 @unit-level
  Scenario: Inconsistency failures cite conflicting acronyms and origin files
    Given a code unit or baseline with orphan or conflicting acronyms
    When the gate confronts it
    Then it returns Fail and the verdict names each conflicting acronym and its origin file

  @IDCND-B10 @unit-level
  Scenario: When no orphan testID prefixes or baseline discrepancies exist the gate passes
    Given a spec whose code units expose only matching or reused testIDs and all baselines match
    When the gate confronts it
    Then it returns Pass, confirming full identity consistency

  @IDCND-I01 @unit-level
  Scenario: Without a map graph the relational gate never approves
    Given a spec node and a nil map graph
    When the gate confronts it
    Then it returns Pending without asserting compliance in the dark

  @IDCND-I02 @unit-level
  Scenario: Cross-unit reuse is never permitted for visual regression baselines
    Given a visual regression baseline whose code belongs to another unit in the map
    When the gate confronts it
    Then it returns Fail, because the baseline must prove this unit and not another

  @IDCND-I03 @unit-level
  Scenario: Absence of a spec code is never double-charged
    Given a spec with no code declared
    When the gate confronts it
    Then it returns Skip, leaving code absence enforcement to the code presence gate

  @IDCND-X01 @unit-level
  Scenario: Common shorthand prefixes of three letters or fewer are not scrutinized
    Given a code unit with shorthand handles like tab or btn
    When the gate confronts it
    Then it returns Pass, avoiding noisy findings across standard components

  @IDCND-X02 @unit-level
  Scenario: The gate does not enforce code presence on specs
    Given a spec node without code
    When the gate confronts it
    Then it returns Skip, because code presence is the duty of another gate

  @IDCND-X03 @unit-level
  Scenario: Components referencing parent screen codes in testIDs are not forbidden
    Given a child component exposing a testID prefixed with its parent screen code
    When the gate confronts it
    Then it returns Pass, preserving end-to-end selector navigation
