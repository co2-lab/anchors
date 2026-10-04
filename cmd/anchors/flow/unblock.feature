# language: en
# @anchors
#   code: UNFTN
#   ref: NBLCK
#   updated_at: 2026-10-03
#   layer: feature

@NBLCK
Feature: Unblock — open the work card a decision demanded, linked to the card stuck waiting for a person

  @NBLCK-B01 @unit-level
  Scenario: Unblock with no card or with two cards is refused
    Given a project in github mode
    When `anchors unblock --reason x` runs with no card, and again with cards 311 and 312
    Then both runs fail and no card is created

  @NBLCK-B02 @unit-level
  Scenario: Unblock with a blank reason is refused
    Given a project in github mode
    When `anchors unblock 311 --reason " "` runs
    Then the command fails naming `--reason`

  @NBLCK-B03 @unit-level
  Scenario: Unblock in local mode is refused
    Given a project in local mode
    When `anchors unblock 311 --reason x` runs
    Then the command fails saying the command exists in github mode

  @NBLCK-B04 @unit-level
  Scenario: The link label is created on demand, and an existing one does not stop the command
    Given a project in github mode whose platform refuses to create the link label of card 311 because it exists
    When `anchors unblock 311 --reason x` runs
    Then the link label of 311 is created before the card
    And the new card is still created without error

  @NBLCK-B05 @unit-level
  Scenario: The new card carries the unblock title and the three labels
    Given a project in github mode whose first workflow label is "anchors"
    When `anchors unblock 311` runs with the reason "the dispatch loop needs try/catch per token" followed by a second paragraph
    Then the card is titled "[unblocks #311] the dispatch loop needs try/catch per token"
    And it carries the labels "anchors", "anchors:to-do" and the link label of 311

  @NBLCK-B06 @unit-level
  Scenario: The body says what to do, where it came from and how the blocked card returns
    Given blocked card 311, the reason "add the retry rule" and the file "src/a.ts"
    When the body of the unblock card is written
    Then it says "unblocks #311", carries the reason, "**Where:** `src/a.ts`", "Where it came from" and "What to do"
    And it says to remove the label `anchors:needs-user` from #311
    And without a file the body has no "Where:" line

  @NBLCK-B07 @unit-level
  Scenario: The blocked card is told which card it waits for
    Given a platform that answers the card creation with the address of issue 512
    When `anchors unblock 311 --reason x` runs
    Then card 311 receives a comment naming the address of issue 512
    And the comment says the needs-user label STAYS

  @NBLCK-B08 @unit-level
  Scenario: The output gives the new card and the state of the blocked one
    Given a platform that answers the card creation with the address of issue 512
    When `anchors unblock 311 --reason x` runs
    Then the output says "unblock card created:" with that address
    And the output says "#311 remains stopped"

  @NBLCK-B09 @unit-level
  Scenario: The leading hash of the card is stripped
    Given a project in github mode
    When `anchors unblock "#311" --reason x` runs
    Then the comment is made on card 311 by its bare number

  @NBLCK-I01 @unit-level
  Scenario: The blocked card keeps its needs-user label
    Given a project in github mode
    When `anchors unblock 311 --reason x` runs to the end
    Then no call removes a label

  @NBLCK-X01 @unit-level
  Scenario: The temporary body file does not survive the command
    Given a project in github mode
    When `anchors unblock 311 --reason x` runs
    Then the body file handed to the card creation no longer exists

  @NBLCK-E01 @unit-level
  Scenario: A card the platform refuses to create fails the command
    Given a platform that answers the card creation with "create failed"
    When `anchors unblock 311 --reason x` runs
    Then the command fails with "create the card" and "create failed"
    And card 311 is not commented

  @NBLCK-E02 @unit-level
  Scenario: A failed comment on the blocked card is a warning
    Given a platform that creates the card and refuses the comment
    When `anchors unblock 311 --reason x` runs
    Then the command succeeds
    And the output says it could not comment on #311
