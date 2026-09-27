# language: en
# @anchors
#   ref: AGCRG
#   updated_at: 2026-09-26
#   layer: feature

@AGCRG
Feature: AgentCards — the cards this agent owns, and the card a pull request declares, read from the tracker

  @AGCRG-B01 @unit-level
  Scenario: Without an agent name, a workflow or the tracker client there are no cards
    Given a tracker client that would list one card
    And either no ANCHORS_AGENT, no configuration, no workflow block, or no client on the PATH
    When the agent's cards are listed
    Then the list is empty and no error is raised

  @AGCRG-B02 @unit-level
  Scenario: The tracker is asked for the open cards of the repository owned by this agent
    Given the agent "host/worker-1" and a workflow on repository "acme/app" with the label "anchors"
    When the agent's cards are listed
    Then the tracker client is called with the repository "acme/app", the label "anchors" and the ownership comment "anchors-owner: host/worker-1"

  @AGCRG-B03 @unit-level
  Scenario: A card line becomes a card whose state loses the anchors prefix
    Given the tracker answers the line "12, Fix the parser, anchors:doing" separated by tabs
    When the agent's cards are listed
    Then the list holds card 12 titled "Fix the parser" in state "doing"

  @AGCRG-B04 @unit-level
  Scenario: Malformed card lines are skipped
    Given the tracker answers a line with no tabs and a line with an empty number among valid lines
    When the agent's cards are listed
    Then only the valid cards 12, 40 and 51 are in the list

  @AGCRG-B05 @unit-level
  Scenario: A card with no state kept as the last line
    Given the tracker answers card 12 in state "doing" and, last, card 40 with an empty state field
    When the agent's cards are listed
    Then both cards are in the list, card 40 with an empty state

  @AGCRG-B06 @unit-level
  Scenario: The issue title keeps the first line of the reason
    Given the reason "  the parser drops the header  " followed by a second line
    When the reason becomes an issue title
    Then the title is "the parser drops the header"

  @AGCRG-B07 @unit-level
  Scenario: A long title is cut to 70 with an ellipsis
    Given a reason of 80 characters and another of exactly 70
    When each reason becomes an issue title
    Then the first becomes 70 long ending in "..." and the second is kept whole

  @AGCRG-B08 @unit-level
  Scenario: The issue number is read from the last segment of the URL
    Given the URLs ".../issues/433", ".../issues/", ".../pull/12a" and the bare text "433"
    When the issue number is read from each
    Then the first gives "433" and the others give the empty string

  @AGCRG-B09 @unit-level
  Scenario: The card of a pull request is the first closing keyword at a line start
    Given a pull request body with "Closes #77" then "Refs #80" at line starts, and another body with "fixes #5" in lower case
    And a body that mentions "closes #9" only in the middle of a line
    When the card of each pull request is read
    Then the answers are "77", "5" and the empty string

  @AGCRG-B10 @unit-level
  Scenario: The pull request reference is normalised before asking the tracker
    Given the reference " #15 " on repository "acme/app"
    When the card of the pull request is read
    Then the tracker is asked to view pull request "15" of "acme/app"
    And with no repository or with a reference of only "#" the answer is empty

  @AGCRG-I01 @unit-level
  Scenario: Every returned card has a number and a state without prefix
    Given tracker lines that are well formed, malformed and without state
    When the agent's cards are listed
    Then every card in the list has a non-empty number and no state starts with "anchors:"

  @AGCRG-X01 @unit-level
  Scenario: Ownership is selected by the question sent to the tracker
    Given the agent "host/worker-1"
    When the agent's cards are listed
    Then the question sent to the tracker carries the ownership comment of "host/worker-1"

  @AGCRG-E01 @unit-level
  Scenario: A failing tracker client gives no cards
    Given a tracker client that exits with an error
    When the agent's cards are listed
    Then the list is empty and no error is raised

  @AGCRG-E02 @unit-level
  Scenario: A failing tracker client gives no pull request card
    Given a tracker client that prints "Closes #77" and exits with an error
    When the card of a pull request is read
    Then the answer is the empty string
