# language: en
# @anchors
#   code: PBFPR
#   ref: PRBDP
#   updated_at: 2026-10-03
#   layer: feature

@PRBDP
Feature: PRBody — write the lines that link a pull request to its cards, in the platform's syntax

  @PRBDP-B01 @unit-level
  Scenario: pr-body in local mode is refused
    Given a project in local mode
    When `anchors pr-body --cards 4` runs
    Then the command fails saying the command exists in github mode

  @PRBDP-B02 @unit-level
  Scenario: The link syntax links and never closes
    Given the link syntax declared for each platform
    When the first word of each syntax is read
    Then the github syntax exists and holds a place for the card number
    And no syntax starts with close, closes, closed, fix, fixes, fixed, resolve, resolves or resolved

  @PRBDP-B03 @unit-level
  Scenario: The requested cards accept the forms a person writes
    Given the requested cards "44", "#44", " 44 " and "#44 "
    When each is resolved
    Then each becomes the single card 44
    And "44, #49,50" becomes three cards
    And a blank list becomes no card

  @PRBDP-B04 @unit-level
  Scenario: Without requested cards the agent's own card is linked
    Given a project in github mode, the agent "host/dev1" and a board where that agent owns card 12
    When `anchors pr-body` runs without `--cards`
    Then the output is exactly "Refs #12"

  @PRBDP-B05 @unit-level
  Scenario: pr-body with no card from either source is refused
    Given a project in github mode and no agent set
    When `anchors pr-body` runs without `--cards`
    Then the command fails saying there is no card

  @PRBDP-B06 @unit-level
  Scenario: Each root drags the open findings born under it
    Given a board where cards 101, 45 and 50 are open and born under card 44
    When `anchors pr-body --cards "#44, 50"` runs
    Then the output links 44, 45, 50 and 101
    And each lookup asks for open cards in the configured repository

  @PRBDP-B07 @unit-level
  Scenario: The lines come out in numeric order
    Given a board where cards 101 and 45 are born under card 44
    When `anchors pr-body --cards 44` runs
    Then "Refs #45" comes before "Refs #101"

  @PRBDP-B08 @unit-level
  Scenario: The only-under switch leaves the roots out
    Given a board where card 45 is born under card 44
    When `anchors pr-body --cards 44 --so-sob` runs
    Then the output is exactly "Refs #45"

  @PRBDP-I01 @unit-level
  Scenario: A card that is both a root and a finding is linked once
    Given a board where card 50 is a requested root and also born under card 44
    When `anchors pr-body --cards "#44, 50"` runs
    Then "Refs #50" appears exactly once

  @PRBDP-X01 @unit-level
  Scenario: pr-body writes nothing to the platform
    Given a project in github mode with requested cards 44 and 50
    When `anchors pr-body` runs
    Then every call it makes is a lookup of issues

  @PRBDP-E01 @unit-level
  Scenario: A failed findings lookup contributes nothing and the root is still linked
    Given a platform whose findings lookup fails, and one whose answer is not JSON
    When the findings under card 44 are looked up, and `anchors pr-body --cards 44` runs
    Then the lookup yields no card
    And the output is exactly "Refs #44"
