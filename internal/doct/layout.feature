# language: en
# @anchors
#   ref: DCOXX
#   updated_at: 2026-09-26
#   layer: feature

@DCOXX
Feature: DocLayout — the one decision of whether a documentation page shows everything or summarizes

  @DCOXX-B01 @unit-level
  Scenario: Either measure passing its cut-off makes the selection big
    Given a cut-off of 20 units and 2000 lines
    When a selection has 21 units and 10 lines, or 1 unit and 2001 lines
    Then both selections are big
    And a selection of 1 unit and 10 lines is not

  @DCOXX-B02 @unit-level
  Scenario: The default cut-off is 20 units or 2000 lines
    Given no cut-off declared by the project
    When the default layout is taken
    Then its limits are 20 units and 2000 lines

  @DCOXX-B03 @unit-level
  Scenario: The summary sentence exists only for a big selection and names its numbers
    Given a cut-off of 20 units and 2000 lines
    When a small selection and a selection of 25 units and 90 rules are described
    Then the small one has an empty sentence
    And the big one's sentence names 25 units, 90 rules, 20 units and 2000 lines

  @DCOXX-I01 @unit-level
  Scenario: A selection exactly at the limit is not big
    Given a cut-off of 20 units and 2000 lines
    When a selection of exactly 20 units and 2000 lines is measured
    Then it is not big
    And one more unit or one more line makes it big

  @DCOXX-X01 @unit-level
  Scenario: The layout never splits a layer into pages
    Given a big selection
    When the layout decides on it
    Then the only answer is big or not big, with no page count or split
