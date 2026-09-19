# language: en
# @anchors
#   ref: SCASS
#   updated_at: 2026-09-19
#   layer: feature

@SCASS
Feature: ScenarioAsserts — scenario outcome steps must assert concrete verifiable outcomes

  @SCASS-B01 @unit-level
  Scenario: Non-feature artifacts skip confrontation
    Given an artifact whose kind is not a feature
    When the gate confronts it
    Then it returns Skip, avoiding confrontation on non-feature artifacts

  @SCASS-B02 @unit-level
  Scenario: Genuinely assertive outcome steps pass
    Given a feature file whose scenario outcome steps declare observable results
    When the gate confronts it
    Then it returns Pass, confirming that every scenario asserts an outcome

  @SCASS-B03 @unit-level
  Scenario: Tautological outcome steps citing only the code fail
    Given a feature containing an outcome step that merely repeats the requirement code
    When the gate confronts it
    Then it returns Fail, naming the offending requirement code

  @SCASS-B04 @unit-level
  Scenario: Tautological variations wrapped in linking words fail
    Given a feature with outcome steps wrapping the requirement code in standard linking words
    When the gate confronts it
    Then it returns Fail, rejecting superficial variations of the tautology

  @SCASS-B05 @unit-level
  Scenario: Outcome steps citing a code with substantive assertion pass
    Given a feature with an outcome step citing a requirement code and asserting concrete values
    When the gate confronts it
    Then it returns Pass, allowing code references when accompanied by substantive assertions

  @SCASS-B06 @unit-level
  Scenario: Multiple tautological outcome steps are reported sorted and deduplicated
    Given a feature containing multiple tautological steps across scenarios
    When the gate confronts it
    Then it returns Fail, reporting all offending codes in alphabetical order without duplicates

  @SCASS-B07 @unit-level
  Scenario: Outcome steps across recognized dialect keywords are enforced
    Given features using supported outcome keywords across different project languages
    When the gate confronts them
    Then it evaluates outcome steps consistently across configured dialects

  @SCASS-B08 @unit-level
  Scenario: Empty and comment lines are ignored
    Given a feature file containing comments and blank lines beside scenario steps
    When the gate confronts it
    Then it ignores comments and blank lines without parsing them as outcome steps

  @SCASS-B09 @unit-level
  Scenario: Non-outcome steps citing codes are ignored
    Given a feature file where setup or action steps cite requirement codes
    When the gate confronts it
    Then it evaluates only the outcome steps and skips setup steps

  @SCASS-I01 @unit-level
  Scenario: Up to two residual content words is classified as a tautology
    Given an outcome step where removing linking words leaves at most two words
    When the gate confronts it
    Then it returns Fail, preventing minimal phrasing additions from bypassing verification

  @SCASS-I02 @unit-level
  Scenario: Truly assertive outcome steps are never flagged as tautologies
    Given diverse forms of assertive outcome steps including comparisons and negative assertions
    When the gate confronts them
    Then it returns Pass, avoiding false positives on legitimate assertions

  @SCASS-I03 @unit-level
  Scenario: Language recognition covers all supported dialect alternatives
    Given outcome steps written in English, Spanish, or Portuguese
    When the gate confronts them
    Then it detects tautologies regardless of which supported language is used

  @SCASS-X01 @unit-level
  Scenario: Prose style and semantic elegance are not graded
    Given outcome steps written in diverse styles and phrasing
    When the gate confronts them
    Then it isolates tautologies mechanically without judging stylistic quality

  @SCASS-X02 @unit-level
  Scenario: Setup and trigger steps are not inspected for code citations
    Given Given and When steps that mention requirement codes
    When the gate confronts the feature
    Then it does not charge code citations in non-outcome steps

  @SCASS-X03 @unit-level
  Scenario: Scenario or test presence is not enforced by this gate
    Given a feature file without scenarios or linked tests
    When the gate confronts it
    Then it does not enforce scenario completeness or test linkages
