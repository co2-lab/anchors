# language: en
# @anchors
#   code: EDFSC
#   ref: ESDPS
#   updated_at: 2026-10-03
#   layer: feature

@ESDPS
Feature: EscalateDuplicate — find the open cards that already deal with the target of an escalation

  @ESDPS-B01 @unit-level
  Scenario: Without a target or a label the board is not asked
    Given a board that would answer with a card
    When the open cards about an empty target, and about "X.spec.md" with an empty label, are looked up
    Then nothing is found
    And no call reaches the board

  @ESDPS-B02 @unit-level
  Scenario: The board is searched for open cards with the label and the target
    Given the target "src/x.spec.md", the label "anchors" and the repository "acme/app"
    When the open cards about the target are looked up
    Then the board of "acme/app" is asked for open issues labelled "anchors" searching "src/x.spec.md"

  @ESDPS-B03 @unit-level
  Scenario: A hit counts when the exact target is in its title or its body
    Given a board with card 405 whose title names "apps/mobile/src/components/MetricCard.spec.md", and card 406 whose body names it
    When the open cards about that target are looked up
    Then cards 405 and 406 are found

  @ESDPS-B04 @unit-level
  Scenario: Each card found carries its number and its title
    Given a board with card 405 titled "[doc-required] Violation @ apps/mobile/src/components/MetricCard.spec.md"
    When the open cards about that target are looked up
    Then the result is "#405 [doc-required] Violation @ apps/mobile/src/components/MetricCard.spec.md"

  @ESDPS-I01 @unit-level
  Scenario: A card about a file that merely contains the target's name is not reported
    Given a board whose search returns card 999 about "apps/mobile/src/components/MetricCardList.spec.md"
    When the open cards about "apps/mobile/src/components/MetricCard.spec.md" are looked up
    Then nothing is found

  @ESDPS-X01 @unit-level
  Scenario: The lookup only reads the board
    Given a board with a card about the target
    When the open cards about the target are looked up
    Then the only call is a listing of issues

  @ESDPS-E01 @unit-level
  Scenario: A failed or unreadable lookup yields nothing
    Given a board lookup that fails, and one that answers "not json"
    When the open cards about "X.spec.md" are looked up
    Then nothing is found and no error is raised
