# language: en
# @anchors
#   ref: PNDCP
#   updated_at: 2026-09-26
#   layer: feature

@PNDCP
Feature: PendingDecisions — the doctor lists the specs that still hold open decisions, the heaviest first

  @PNDCP-B01 @unit-level
  Scenario: A spec with two open decisions gives one warning carrying the count
    Given a spec whose open decisions section lists two questions
    When the pending decisions are checked
    Then there is one warning decisao-pendente on that spec whose text counts 2 decisions

  @PNDCP-B02 @unit-level
  Scenario: A section closed with none is not pending
    Given a spec whose open decisions section says none
    When the pending decisions are checked
    Then there is no finding

  @PNDCP-B03 @unit-level
  Scenario: The spec with the most open decisions comes first
    Given one spec with one open decision and another with three
    When the pending decisions are checked
    Then the spec with three comes first and the spec with one second

  @PNDCP-B04 @unit-level
  Scenario: A nil map or a map without open decisions gives nothing
    Given a nil map, and a map whose only node with open decisions is not a spec
    When the pending decisions are checked
    Then there is no finding

  @PNDCP-I01 @unit-level
  Scenario: The doctor's count is the check's count
    Given a spec with three open decisions
    When the pending decisions are checked
    Then the count in the finding equals what the check's open decisions counter reads from the file

  @PNDCP-X01 @unit-level
  Scenario: The count follows the check's rule, not a reading of its own
    Given a spec whose open decisions section is closed with none
    When the pending decisions are checked and the check's counter reads it
    Then both say zero and no finding is given

  @PNDCP-E01 @unit-level
  Scenario: A spec missing on disk is skipped and the others are still reported
    Given a map with a spec whose file is missing and a readable spec with open decisions
    When the pending decisions are checked
    Then only the readable spec is reported
