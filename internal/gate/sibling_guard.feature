# language: en
# @anchors
#   ref: SBGRD
#   updated_at: 2026-09-19
#   layer: feature

@SBGRD
Feature: SiblingGuard — sibling functions treat the same parameter consistently

  @SBGRD-B01 @unit-level
  Scenario: An artifact that is not code leaves without a verdict
    Given a node whose kind is spec, feature or test
    When the gate confronts it
    Then it returns Skip, because the gate reads function bodies

  @SBGRD-B02 @unit-level
  Scenario: Without a declared dialect the verdict is undetermined
    Given a project that declares no dialect for its stack
    When the gate confronts a code file
    Then it returns Pending, because reading code it cannot recognise would stamp
      what it never checked

  @SBGRD-B03 @unit-level
  Scenario: Fewer than three siblings on the same parameter is left alone
    Given a module exporting two functions that receive the same parameter
    And one of them guards it while the other does not
    When the gate confronts the module
    Then it accuses nothing, because below the bar asymmetry is not evidence

  @SBGRD-B04 @unit-level
  Scenario: The sibling that does not guard is accused
    Given a module exporting three functions that receive the same parameter
    And two of them guard it and the third does not
    When the gate confronts the module
    Then it returns Fail against the third

  @SBGRD-B05 @unit-level
  Scenario: When every sibling guards, nothing is accused
    Given a module exporting three functions that all guard the same parameter
    When the gate confronts the module
    Then it accuses nothing, because there is no asymmetry to report

  @SBGRD-B06 @unit-level
  Scenario: When no sibling guards, nothing is accused either
    Given a module exporting three functions and none of them guards the parameter
    When the gate confronts the module
    Then it accuses nothing, because a module that never guards is a decision,
      not an oversight

  @SBGRD-B07 @unit-level
  Scenario: The verdict names the function and the parameter
    Given a module whose majority guards and one function does not
    When the gate confronts the module
    Then the verdict names that function and the parameter at stake, so the reader
      does not diff the module by hand

  @SBGRD-B08 @unit-level
  Scenario: A waiver with a written reason silences the accusation
    Given a module whose majority guards and one function does not
    And that function carries a waiver marker followed by the reason
    When the gate confronts the module
    Then it returns Pass, because the sibling that legitimately delegates the check
      says so where whoever reads the function will see it

  @SBGRD-I01 @unit-level
  Scenario: The gate never judges what the guard does
    Given a module whose sibling guards differ in what they actually check
    And one function applies no guard at all
    When the gate confronts the module
    Then only the absence is reported, because understanding the guard would require
      understanding the domain

  @SBGRD-I02 @unit-level
  Scenario: The three conservatism conditions hold together
    Given modules that each satisfy only some of the three conditions
    When the gate confronts each of them
    Then none is accused, because any single condition alone would produce the false
      positive that costs the gate its credibility

  @SBGRD-X01 @unit-level
  Scenario: The gate does not invent what an exported function looks like
    Given a project declaring a dialect whose export shape differs from another stack
    When the gate confronts a file written in that other shape
    Then it accuses nothing, because a pattern hardcoded here would recognise one
      ecosystem and report green over every other

  @SBGRD-X02 @unit-level
  Scenario: A single function in isolation is not accused
    Given a module exporting one function that does not guard its parameter
    When the gate confronts the module
    Then it accuses nothing, because with no siblings there is no asymmetry — and
      asymmetry is the whole evidence
