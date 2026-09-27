# language: en
# @anchors
#   ref: DSCRD
#   updated_at: 2026-09-26
#   layer: feature

@DSCRD
Feature: Discard — take off the board a card that no longer makes sense, without deleting it

  @DSCRD-B01 @unit-level
  Scenario: Discard without any card argument is refused
    Given a project in github mode
    When `anchors discard --reason obsolete` runs with no card
    Then the command fails and no call reaches the platform

  @DSCRD-B02 @unit-level
  Scenario: Discard with a blank reason is refused
    Given a project in github mode
    When `anchors discard --reason "  " 42` runs
    Then the command fails naming `--reason`

  @DSCRD-B03 @unit-level
  Scenario: Discard in local mode is refused
    Given a project in local mode
    When `anchors discard --reason obsolete 42` runs
    Then the command fails saying the command exists in github mode

  @DSCRD-B04 @unit-level
  Scenario: The discard label is created before any card is touched, and an existing label does not stop it
    Given a project in github mode whose platform refuses to create the discard label because it exists
    When `anchors discard --reason obsolete 42` runs
    Then the first call creates the discard label
    And card 42 is still discarded without error

  @DSCRD-B05 @unit-level
  Scenario: A card is labelled, then commented, then closed
    Given card 42 on a platform that accepts every call
    When card 42 is discarded
    Then the first call adds the discard label, the second comments, and the third closes

  @DSCRD-B06 @unit-level
  Scenario: The comment carries the reason and the way back
    Given card 42 and the reason "spec deleted in the cleanup of plan 0004"
    When card 42 is discarded
    Then the comment contains that reason
    And the comment names the discard label to remove to bring the card back

  @DSCRD-B07 @unit-level
  Scenario: The leading hash is stripped and a blank argument is skipped
    Given a project in github mode
    When `anchors discard --reason obsolete "#42" " "` runs
    Then card 42 is edited by its bare number
    And no call is made for the blank argument

  @DSCRD-B08 @unit-level
  Scenario: A card that is already closed is still discarded
    Given card 42 whose close is refused by the platform
    When card 42 is discarded
    Then no error is returned

  @DSCRD-B09 @unit-level
  Scenario: Each discarded card is reported on standard output
    Given cards 42 and 43, where 43 cannot receive the label
    When `anchors discard --reason obsolete 42 43` runs
    Then the output says "#42 discarded"
    And the output does not say "#43 discarded"

  @DSCRD-I01 @unit-level
  Scenario: A card is never closed before its label and reason are recorded
    Given card 43 whose label fails and card 42 whose comment fails
    When both are discarded
    Then neither card 42 nor card 43 is closed

  @DSCRD-X01 @unit-level
  Scenario: Discarding never deletes the issue
    Given card 42 on a platform that accepts every call
    When card 42 is discarded
    Then no call deletes the issue

  @DSCRD-E01 @unit-level
  Scenario: A card whose label fails is named in the error while the others are discarded
    Given cards 42 and 43, where the platform answers "issue not found" for the label of 43
    When `anchors discard --reason obsolete 42 43` runs
    Then the command fails with "1 card(s) were not discarded", naming #43 and "issue not found"
    And card 43 is neither commented nor closed
    And card 42 is closed

  @DSCRD-E02 @unit-level
  Scenario: A card whose reason cannot be recorded is not closed
    Given card 42 whose comment is refused with "rate limited"
    When card 42 is discarded
    Then the error says "record the reason" and "rate limited"
    And card 42 is not closed
