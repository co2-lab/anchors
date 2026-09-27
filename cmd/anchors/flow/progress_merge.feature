# language: en
# @anchors
#   ref: PRMRP
#   updated_at: 2026-09-26
#   layer: feature

@PRMRP
Feature: ProgressMerge — the git merge driver of a plan's progress file: unite both sides, done beats pending

  @PRMRP-B01 @unit-level
  Scenario: The merge driver refuses any count other than three arguments
    Given two progress files
    When `anchors merge-progress` runs with two arguments, and again with four
    Then both runs fail

  @PRMRP-B02 @unit-level
  Scenario: The union is written into our side
    Given our side with a.spec.md done and their side with b.spec.md done, c.spec.md pending on both
    When `anchors merge-progress base ours theirs` runs
    Then our side holds a.spec.md and b.spec.md done and c.spec.md pending

  @PRMRP-B03 @unit-level
  Scenario: A missing base does not stop the merge
    Given a base path that does not exist and readable sides
    When `anchors merge-progress` runs
    Then the merge succeeds and our side holds the union

  @PRMRP-B04 @unit-level
  Scenario: The report counts the done items and what the other side added
    Given our side with one item done and their side with another item done
    When `anchors merge-progress` runs
    Then standard error says "2 done" and "(1 from the other side)"
    And when the other side adds nothing, standard error says "1 done" without "from the other side"

  @PRMRP-I01 @unit-level
  Scenario: An item done on our side stays done
    Given our side with a.spec.md done and their side with a.spec.md pending
    When `anchors merge-progress` runs
    Then our side still holds a.spec.md done

  @PRMRP-X01 @unit-level
  Scenario: The driver applies the progress reader's union
    Given our side and their side of a progress file
    When `anchors merge-progress` runs
    Then our side equals the progress reader's union of the two sides

  @PRMRP-E01 @unit-level
  Scenario: An unreadable our side is named
    Given our side that does not exist
    When `anchors merge-progress` runs
    Then it fails with "read our side"

  @PRMRP-E02 @unit-level
  Scenario: An unreadable other side is named and our side is untouched
    Given their side that does not exist
    When `anchors merge-progress` runs
    Then it fails with "read the other side"
    And our side keeps its content

  @PRMRP-E03 @unit-level
  Scenario: A result that cannot be written fails the merge
    Given our side that is read-only
    When `anchors merge-progress` runs
    Then it fails with "write the result"
