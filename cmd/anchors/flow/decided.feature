# language: en
# @anchors
#   ref: DCDDE
#   updated_at: 2026-09-26
#   layer: feature

@DCDDE
Feature: Decided — release the card an escalation stopped for a person, once the decision became a rule

  @DCDDE-B01 @unit-level
  Scenario: decided without a card is refused
    Given the resolution "ABCDE-R0001: x"
    When `anchors decided` runs without `--card`
    Then the command fails naming `--card`

  @DCDDE-B02 @unit-level
  Scenario: decided without a resolution is refused and says why
    Given card 4
    When `anchors decided --card 4` runs without a resolution, and again with a blank one
    Then both runs fail naming `--resolution`
    And the message says the answer "becomes a RULE"

  @DCDDE-B03 @unit-level
  Scenario: decided in local mode is refused
    Given a project in local mode
    When `anchors decided --card 4 --resolution "ABCDE-R0001: x"` runs
    Then the command fails pointing to `issues/`

  @DCDDE-B04 @unit-level
  Scenario: decided refuses while an unblock card is open
    Given cards 512 and 513 open with the unblock label of card 4
    When `anchors decided --card 4 --resolution "X-R0001: y"` runs
    Then the command fails naming "#512, #513" and the needs-user label
    And no issue is edited

  @DCDDE-B05 @unit-level
  Scenario: The release removes needs-user and every blocked-by label at once
    Given card 4 carrying the labels blocked-by-701 and blocked-by-702
    When `anchors decided --card 4 --resolution "PLTFR-R0004: mTLS is built in F03"` runs
    Then one edit of card 4 removes needs-user, blocked-by-701 and blocked-by-702
    And the output says "card #4 released"

  @DCDDE-B06 @unit-level
  Scenario: The comment keeps the resolution and who blocked the card
    Given card 4 carrying the labels blocked-by-701 and blocked-by-702
    When `anchors decided --card 4 --resolution "PLTFR-R0004: mTLS is built in F03"` runs
    Then card 4 receives a comment with "PLTFR-R0004" and "**Was blocked by:** #701, #702"

  @DCDDE-B07 @unit-level
  Scenario: The decisions under the card are labelled manual and closed
    Given decisions 701 and 702 open with needs-user and the under label of card 4
    When `anchors decided --card 4 --resolution "PLTFR-R0004: mTLS is built in F03"` runs
    Then the list asks for open issues with both needs-user and the under label of card 4
    And decision 701 is labelled manual and then closed with the resolution

  @DCDDE-B08 @unit-level
  Scenario: The output counts the decisions closed
    Given card 4 with no decision, then with two decisions, then with two where one fails to close
    When `anchors decided --card 4` runs in each case
    Then the output says "no open decision issue under card #4", "2 decision issues closed" and "1 decision issue closed"

  @DCDDE-I01 @unit-level
  Scenario: Blockers are read first and the card is released before its decisions close
    Given card 4 with a blocker and a decision under it
    When `anchors decided --card 4` runs
    Then the calls come in the order: read the labels of 4, edit 4, comment on 4, close the decision

  @DCDDE-X01 @unit-level
  Scenario: Only the decisions under this card are closed
    Given card 4
    When `anchors decided --card 4` runs
    Then the decisions listed are filtered by the under label of card 4

  @DCDDE-E01 @unit-level
  Scenario: An unblock lookup that fails or is unreadable refuses
    Given an unblock lookup that answers "HTTP 502", and one that answers "not json"
    When `anchors decided --card 4 --resolution "X-R0001: y"` runs
    Then the command fails naming the unblock label of card 4 and the needs-user label
    And no issue is edited

  @DCDDE-E02 @unit-level
  Scenario: A failed release fails the command
    Given a platform that refuses the edit of card 4 with "no permission"
    When `anchors decided --card 4 --resolution "X-R0001: y"` runs
    Then the command fails with "release card #4" and "no permission"
    And the output does not say released
    And nothing is commented or closed

  @DCDDE-E03 @unit-level
  Scenario: An unreadable decision list closes nothing
    Given a decision list that answers "not json"
    When the decisions under card 4 are closed
    Then zero is counted and nothing is closed

  @DCDDE-E04 @unit-level
  Scenario: Unreadable labels report no blocker
    Given a card whose labels cannot be read
    When its blockers are read
    Then no blocker is reported
    And the labels "anchors:blocked-by-1" and "anchors:blocked-by-2" are both read when the answer is readable
