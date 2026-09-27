# language: en
# @anchors
#   ref: PRFLP
#   updated_at: 2026-09-26
#   layer: feature

@PRFLP
Feature: ProgressFile — a plan's progress companion: how it is named, and how two sides of it merge

  @PRFLP-B01 @unit-level
  Scenario: Only the progress suffix marks a progress file
    Given the paths "plans/0017-mutacao-progress.md", "0001-progress.md", "plans/0017-mutacao.md", "docs/progress-notes.md" and "plans/progress.md"
    When each path is classified
    Then only the first two are progress files

  @PRFLP-B02 @unit-level
  Scenario: The companion path replaces the extension of the last segment
    Given the plan paths "plans/x.md", "a.b/plan" and "x"
    When the companion path of each is derived
    Then they are "plans/x-progress.md", "a.b/plan-progress.md" and "x-progress.md"

  @PRFLP-B03 @unit-level
  Scenario: Two sides ticking neighbouring items keep both ticks
    Given our side ticks "a.spec.md" and their side ticks "b.spec.md"
    When the sides are merged
    Then both items come out ticked and the done count is 2

  @PRFLP-B04 @unit-level
  Scenario: An item only on their side is added after our last item
    Given our side has a heading, the item "a" and a trailing line "tail"
    And their side has the items "b" ticked and "a"
    When the sides are merged
    Then "b" comes out ticked right after "a" and before "tail"
    And when our side has no item at all, "b" comes out at the end

  @PRFLP-B05 @unit-level
  Scenario: The lines that are not items survive the merge
    Given both sides have a heading, a quoted note and a section title above the item "a"
    When the sides are merged with "a" ticked on their side
    Then the heading, the note and the section title are all present and "a" is ticked

  @PRFLP-B06 @unit-level
  Scenario: A promoted tick keeps our indentation
    Given our side has the indented item "a" unticked and their side has it ticked
    When the sides are merged
    Then the item comes out ticked with the same two-space indentation

  @PRFLP-B07 @unit-level
  Scenario: The done count reads both tick letters at any indentation
    Given items ticked with "x", with "X" and nested with "x", plus one unticked item
    When the done items are counted
    Then the count is 3

  @PRFLP-I01 @unit-level
  Scenario: Merging a side with itself changes nothing
    Given a progress text with one ticked and one unticked item
    When it is merged with itself
    Then the result is the same text and the done count stays 1

  @PRFLP-I02 @unit-level
  Scenario: A merge never unticks an item
    Given our side has "a.spec.md" ticked and their side has it unticked
    When the sides are merged
    Then "a.spec.md" is still ticked

  @PRFLP-X01 @unit-level
  Scenario: An item whose text differs between the sides is kept twice
    Given our side has the item "a — first wording" and their side has "a — second wording" ticked
    When the sides are merged
    Then both wordings are in the result
