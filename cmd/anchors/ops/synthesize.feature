# language: en
# @anchors
#   code: SYFTS
#   ref: SYCMS
#   updated_at: 2026-10-03
#   layer: feature

@SYCMS
Feature: SynthesizeCommand — two pull requests in content conflict become one card that asks for the best of each, and every end points at it

  @SYCMS-B01 @unit-level
  Scenario: Only github mode is accepted and the first PR is required
    Given a project in local mode, and one in github mode
    When synthesize runs in local mode with --pr-a 693, then in github mode with no --pr-a
    Then the first fails naming github mode
    And the second fails naming --pr-a
    And and gh was never called

  @SYCMS-B02 @unit-level
  Scenario: A leading # on a PR number is dropped
    Given a project in github mode
    When synthesize runs with --pr-a "#693" --pr-b 556 --dry-run
    Then the title reads "what PRs #693 and #556 deliver, in one"

  @SYCMS-B03 @unit-level
  Scenario: The card cites both PRs and cards and asks for the best of each
    Given PRs 693 and 556 with cards 666 and 198, their titles and the file x.spec.md
    When the card body is written
    Then it cites #693, #556, #666, #198, both titles and `x.spec.md`
    And it asks for the best of each
    And it warns against discarding a side without saying why

  @SYCMS-B04 @unit-level
  Scenario: A one-sided card says the other side is missing and where to look
    Given PR 693 with card 666 and no second PR
    When the card body is written
    Then it has no empty PR row
    And it says the other side was not identified and mentions git log
    And it does not say to read both closed PRs

  @SYCMS-B05 @unit-level
  Scenario: The card is labelled for the board and under each origin card
    Given PRs 693 and 556 from cards 666 and 198 in github mode
    When synthesize runs
    Then the labels anchors:under-666 and anchors:under-198 are created
    And the card is created with --label anchors --label anchors:to-do --label anchors:under-666 --label anchors:under-198

  @SYCMS-B06 @unit-level
  Scenario: Each PR is commented and closed and each origin card is pointed at the new card
    Given PRs 693 and 556 from cards 666 and 198
    When synthesize runs and the card is issue 700
    Then both PRs are commented and closed
    And cards 666 and 198 are commented with the link of issue 700
    And and with one side the PR comment names the integration branch

  @SYCMS-B07 @unit-level
  Scenario: Each failed link is a warning and the command succeeds
    Given a gh that refuses every comment and close
    When synthesize runs for PR 693 alone
    Then it succeeds
    And stderr warns that #693 did not receive the comment, was not closed, and that card #666 did not receive the link

  @SYCMS-B08 @unit-level
  Scenario: The closing line names only the PRs actually closed
    Given PRs 693 and 556 where closing 556 fails
    When synthesize runs
    Then it says PR #693 was closed and #556 is still OPEN
    And and with PR 693 alone it says PR #693 was closed, with no second PR

  @SYCMS-B09 @unit-level
  Scenario: A dry run shows the card and changes nothing
    Given PRs 693 and 556 and the file Dashboards.spec.md
    When synthesize runs with --dry-run
    Then it prints the title, both cards, both titles and the file
    And gh is never asked to create, comment, close or label

  @SYCMS-B10 @unit-level
  Scenario: The card follows the project language
    Given the project language pt-BR
    When the card body is written
    Then it reads "Dois PRs escreveram" and "## Como entregar"
    And it holds no English heading

  @SYCMS-I01 @unit-level
  Scenario: No PR is closed unless the synthesis card exists
    Given a gh that refuses to create the issue
    When synthesize runs for PRs 693 and 556
    Then no pr close is issued

  @SYCMS-X01 @unit-level
  Scenario: The card never picks a side
    Given PRs 693 and 556
    When the card body is written
    Then it asks for the best of each rather than for a choice

  @SYCMS-E01 @unit-level
  Scenario: A card that cannot be opened fails with the host's message
    Given a gh that answers HTTP 403 to issue create
    When synthesize runs
    Then it fails saying it could not open the synthesis card, with HTTP 403
