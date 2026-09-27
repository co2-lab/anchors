# language: en
# @anchors
#   ref: BRCRB
#   updated_at: 2026-09-26
#   layer: feature

@BRCRB
Feature: BoardCards — reading the repository's board, where the work queue lives in github mode

  @BRCRB-B01 @unit-level
  Scenario: A card needs every configured label and the asked state
    Given the labels "anchors" and "team-x" and a board with a card missing "team-x", a to-do card and an in-progress card
    When the in-progress cards are listed
    Then only the in-progress card carrying both labels is listed

  @BRCRB-B02 @unit-level
  Scenario: The owner is the last ownership comment
    Given a card commented "anchors-owner: maq-a/sessao-1", then another comment, then "anchors-owner: maq-b/sessao-2", and a card with no comment
    When their owners are read
    Then the first is owned by "maq-b/sessao-2" and the second has no owner

  @BRCRB-B03 @unit-level
  Scenario: An escalated card goes only to whoever decides the product
    Given a to-do card labelled "anchors:needs-user", one labelled "anchors:precisa-do-usuario", and an ordinary to-do card
    When an agent that declared nothing, one that declared no and one that declared yes consider each
    Then the escalated cards are declined by the first two and taken by the third, and the ordinary card is declined by none

  @BRCRB-B04 @unit-level
  Scenario: The agent resumes its own work under way
    Given a board where the agent owns an escalated in-progress card, a to-do card, and a card in review and to-do, and another agent owns an in-progress card
    When the agent's own card is asked
    Then the card in review is returned with the state in-review

  @BRCRB-B05 @unit-level
  Scenario: A unit's card is found by the code in its title
    Given open cards titled "[GLCGL] Implementar spec", "[ALTPL] depends on GLCGL" with "[GLCGL]" in its body, and "[GLCGLX] other"
    When the card of "glcgl" is looked up
    Then the card "[GLCGL] Implementar spec" is returned

  @BRCRB-B06 @unit-level
  Scenario: An open Anchors card is returned by number
    Given issue 801 open with the label "anchors"
    When card 801 is asked
    Then card 801 is returned

  @BRCRB-B07 @unit-level
  Scenario: Board queries do not name the repository, other calls do
    Given the repository "o/r"
    When the agent's card is asked and card 801 is asked
    Then the board query carries no "--repo" and the issue view carries "--repo o/r"

  @BRCRB-B08 @unit-level
  Scenario: Asking for work dispatches the claim pipeline
    Given the agent "maq/1"
    When it asks for work
    Then the claim workflow is dispatched with "agent=maq/1"

  @BRCRB-E01 @unit-level
  Scenario: Empty labels are refused
    Given a client of "org/repo" with no labels
    When the to-do cards are listed
    Then the listing is refused

  @BRCRB-E02 @unit-level
  Scenario: A repository that is not owner/name is refused
    Given a client of "just-a-name" with the label "anchors"
    When the cards are listed
    Then the listing is refused naming "just-a-name"

  @BRCRB-E03 @unit-level
  Scenario: An empty code is refused
    Given the codes "", "   " and a tab
    When their cards are looked up
    Then each lookup is refused

  @BRCRB-E04 @unit-level
  Scenario: A code no title holds is refused
    Given open cards titled "[ALTPL] other"
    When the card of "RIMRD" is looked up
    Then it is refused naming "[RIMRD]"

  @BRCRB-E05 @unit-level
  Scenario: An ambiguous code is refused naming every card
    Given open cards 966 and 887 both titled with "[RIMRD]"
    When the card of "RIMRD" is looked up
    Then it is refused naming "#966", "#887" and "--card"

  @BRCRB-E06 @unit-level
  Scenario: A closed card or a non-Anchors issue is refused by number
    Given issue 801 closed, and then open without the "anchors" label
    When card 801 is asked each time
    Then the first is refused as closed and the second as not an Anchors card

  @BRCRB-E07 @unit-level
  Scenario: An empty comment is refused
    Given the comment body "   "
    When it is posted on issue 1
    Then it is refused
