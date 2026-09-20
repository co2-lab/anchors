# language: en
# @anchors
#   ref: PLCFL
#   updated_at: 2026-09-20
#   layer: feature

@PLCFL
Feature: PlaceholderFilled — the skeleton the generator emits must be FILLED IN

  @PLCFL-B01 @unit-level
  Scenario: A raw skeleton fails confrontation
    Given a freshly generated artifact where placeholder markers remain in place
    When the gate confronts the artifact
    Then it returns Fail, because the artifact was generated and never written

  @PLCFL-B02 @unit-level
  Scenario: A header field whose value is the marker fails
    Given an artifact whose header field holds a placeholder marker
    When the gate confronts the artifact
    Then it returns Fail, because a placeholder is not a valid field value

  @PLCFL-B03 @unit-level
  Scenario: A table cell holding only the marker fails
    Given an artifact containing a table cell that holds only the marker
    When the gate confronts the artifact
    Then it returns Fail, because a rule code without content says nothing

  @PLCFL-B04 @unit-level
  Scenario: A title or body line opening with the marker fails
    Given an artifact where a heading or body line opens with the marker
    When the gate confronts the artifact
    Then it returns Fail, because the title or statement was never filled in

  @PLCFL-B05 @unit-level
  Scenario: The verdict names what was left behind
    Given an artifact containing unfilled placeholders
    When the gate confronts the artifact
    Then the failure verdict explicitly names each placeholder left behind

  @PLCFL-B06 @unit-level
  Scenario: An artifact with every marker replaced passes
    Given an artifact where every generator marker has been replaced with written content
    When the gate confronts the artifact
    Then it returns Pass, charging only generator leftovers and nothing else

  @PLCFL-I01 @unit-level
  Scenario: A section written on purpose to list pending work is not accused
    Given an artifact containing an intentional pending work section
    When the gate confronts the artifact
    Then it returns Pass, distinguishing honest task tracking from unwritten artifacts

  @PLCFL-X01 @unit-level
  Scenario: The gate does not judge the quality of replacement text
    Given an artifact where all markers are replaced with low-quality or trivial sentences
    When the gate confronts the artifact
    Then it returns Pass, leaving semantic quality judgment to other gates

  @PLCFL-X02 @unit-level
  Scenario: A marker in running prose is not accused
    Given an artifact whose prose mentions a marker outside any value position
    And every header field, table cell and rule title is written
    When the gate confronts the artifact
    Then it returns Pass, because the position is the evidence and not the word
