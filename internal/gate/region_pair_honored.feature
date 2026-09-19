# language: en
# @anchors
#   ref: RPHRG
#   updated_at: 2026-09-19
#   layer: feature

@RPHRG
Feature: RegionPairHonored — every opened source region must close with its own identity code

  @RPHRG-B01 @unit-level
  Scenario: Confronting an artifact that is neither code nor test skips
    Given an artifact whose kind is not code or test
    When the gate confronts it
    Then it returns Skip, avoiding confrontation on non-code artifacts

  @RPHRG-B02 @unit-level
  Scenario: Code or test artifacts without region markers skip
    Given a code or test artifact that contains no region markers
    When the gate confronts it
    Then it returns Skip, because region delimitation is optional

  @RPHRG-B03 @unit-level
  Scenario: Balanced regions closing with matching identity codes pass
    Given source code where every opened region closes with its own identity code
    When the gate confronts it
    Then it returns Pass, confirming all source regions are well-formed

  @RPHRG-B04 @unit-level
  Scenario: An opened region that is never closed fails
    Given source code containing a region opening marker without a corresponding close
    When the gate confronts it
    Then it returns Fail, reporting the line and missing close marker

  @RPHRG-B05 @unit-level
  Scenario: An end region marker without an opening marker fails as an orphan close
    Given source code containing an end region marker without a preceding opening marker
    When the gate confronts it
    Then it returns Fail, reporting the line of the orphan close

  @RPHRG-B06 @unit-level
  Scenario: An end region marker closing with a different code fails
    Given source code where a region closes with a code different from the open region
    When the gate confronts it
    Then it returns Fail, citing both the expected code and the found code

  @RPHRG-B07 @unit-level
  Scenario: Multiple pairing errors are ordered sequentially by line number
    Given source code with multiple pairing defects across different lines
    When the gate confronts it
    Then it returns Fail with defects sorted in ascending order of line number

  @RPHRG-I01 @unit-level
  Scenario: Region absence is never charged as a failure
    Given a project code unit lacking region markers
    When the gate confronts it
    Then it returns Skip, allowing incremental adoption across the codebase

  @RPHRG-I02 @unit-level
  Scenario: Pairing defects produce a blocking Fail verdict
    Given code containing invalid region nesting or unclosed markers
    When the gate confronts it
    Then it returns Fail rather than Pending, because syntax defects require direct repair

  @RPHRG-I03 @unit-level
  Scenario: End markers must explicitly match opening codes
    Given code where an end marker carries another requirement code
    When the gate confronts it
    Then it returns Fail, preventing silent interval corruption across neighboring requirements

  @RPHRG-X01 @unit-level
  Scenario: The gate does not mandate region markers across all source files
    Given standard code files without region markers
    When the gate confronts them
    Then it skips confrontation without requiring markers

  @RPHRG-X02 @unit-level
  Scenario: Region markers inside specifications or documentation are ignored
    Given a specification file containing region syntax examples
    When the gate confronts it
    Then it returns Skip, avoiding false alarms on documentation text

  @RPHRG-X03 @unit-level
  Scenario: Internal code semantics inside regions are not evaluated
    Given well-formed region boundaries containing arbitrary code statements
    When the gate confronts it
    Then it validates pairing structure without inspecting internal language semantics
