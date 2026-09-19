# language: en
# @anchors
#   ref: STASC
#   updated_at: 2026-09-19
#   layer: feature

@STASC
Feature: ScenarioTypeAligned — scenario classification tags must match the code nature letter

  @STASC-B01 @unit-level
  Scenario: Non-feature artifacts skip confrontation
    Given an artifact node whose kind is not feature
    When the gate confronts it
    Then it returns Skip, because scenario classification tags reside exclusively in feature files

  @STASC-B02 @unit-level
  Scenario: Missing or empty rule types in configuration skip confrontation
    Given a configuration that is nil or defines no rule types
    When the gate confronts it
    Then it returns Skip, because rule type vocabulary has not been declared

  @STASC-B03 @unit-level
  Scenario: Configuration without tag mappings skips confrontation
    Given a configuration with rule types that declare no scenario tags
    When the gate confronts it
    Then it returns Skip, avoiding heuristic guesses about tag meanings

  @STASC-B04 @unit-level
  Scenario: Features without coded scenarios skip confrontation
    Given a feature file containing no scenarios with identity codes
    When the gate confronts it
    Then it returns Skip, because there are no scenario codes to align

  @STASC-B05 @unit-level
  Scenario: Codes lacking a recognized rule letter are ignored
    Given a scenario whose identity code does not embed a rule letter
    When the gate confronts it
    Then it ignores the code without charging an alignment defect

  @STASC-B06 @unit-level
  Scenario: Unregistered scenario tags are ignored
    Given a scenario carrying tags not present in the configured rule types vocabulary
    When the gate confronts it
    Then it ignores those tags without asserting disagreement

  @STASC-B07 @unit-level
  Scenario: Scenarios whose classification tags align with code letters pass
    Given a scenario whose classification tag matches the letter of its identity code
    When the gate confronts it
    Then it returns Pass, confirming alignment between tag and rule type

  @STASC-B08 @unit-level
  Scenario: A tag mapped to multiple rule letters passes if any matches
    Given a scenario tag declared under more than one rule letter
    And a scenario code whose letter matches one of those declared letters
    When the gate confronts it
    Then it returns Pass, respecting multi-category tag definitions

  @STASC-B09 @unit-level
  Scenario: A multi-coded scenario passes when a tag matches any code
    Given a scenario co-tagged with multiple identity codes
    And a classification tag matching the letter of any attached code
    When the gate confronts it
    Then it returns Pass, recognizing valid co-tagging across requirements

  @STASC-B10 @unit-level
  Scenario: A scenario tag disagreeing with the code rule letter returns Pending
    Given a scenario whose classification tag contradicts the letter of its identity code
    When the gate confronts it
    Then it returns Pending, recording the divergence for manual review

  @STASC-B11 @unit-level
  Scenario: The Pending verdict cites details of the type divergence
    Given a scenario with mismatched classification tag and code letter
    When the gate confronts it
    Then it returns Pending with details naming the code, letter, tag, allowed letters, and title

  @STASC-B12 @unit-level
  Scenario: Multiple mismatch findings are sorted deterministically
    Given a feature file containing multiple scenario classification mismatches
    When the gate confronts it
    Then it returns Pending with findings sorted in deterministic order

  @STASC-I01 @unit-level
  Scenario: Classification alignment is evaluated only on feature artifacts
    Given a spec, code, or test artifact node confronted by the gate
    When the gate confronts it
    Then it returns Skip, ensuring alignment is verified only where tags reside

  @STASC-I02 @unit-level
  Scenario: Absence of tag mappings prevents speculative enforcement
    Given a project configuration with rule types lacking tag mappings
    When the gate confronts it
    Then it returns Skip, preventing false positives across different language conventions

  @STASC-I03 @unit-level
  Scenario: Mismatched scenario types return Pending rather than Fail
    Given a feature with mismatched scenario tags
    When the gate confronts it
    Then it returns Pending, treating discrepancies as inherited debt requiring judgment

  @STASC-I04 @unit-level
  Scenario: Secondary requirement codes prevent false mismatch reporting
    Given a scenario proving two requirements where the tag describes the secondary code
    When the gate confronts it
    Then it returns Pass, preventing valid dual-tagged scenarios from failing

  @STASC-X01 @unit-level
  Scenario: The gate does not mandate classification tags on every scenario
    Given a coded scenario containing no classification tags
    When the gate confronts it
    Then it returns Pass, treating classification tags as optional metadata

  @STASC-X02 @unit-level
  Scenario: The gate does not decide whether tag or code is erroneous
    Given a scenario with conflicting tag and code
    When the gate confronts it
    Then it reports the divergence without attempting automatic resolution

  @STASC-X03 @unit-level
  Scenario: Specifications and test files are not inspected or modified
    Given non-feature artifacts in the project
    When the gate runs
    Then it leaves non-feature files outside its confrontation scope

  @STASC-X04 @unit-level
  Scenario: Tags are permitted to map across multiple rule letters
    Given configuration mapping a single tag to several distinct rule types
    When the gate confronts scenarios using that tag
    Then it accepts the tag under any of its configured letters
