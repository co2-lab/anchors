# language: en
# @anchors
#   ref: GHIGT
#   updated_at: 2026-10-01
#   layer: feature

@GHIGT
Feature: GitHubIssues — the issue lifecycle on the repository's cards, when the project works on GitHub

  @GHIGT-B01 @unit-level
  Scenario: A card whose marker only starts with the key is not the issue's
    Given a search answer holding only the card of a longer key that contains the finding's key
    When the finding is opened on GitHub
    Then a new card is created instead of matching the near miss

  @GHIGT-B02 @unit-level
  Scenario: The title names the gate, the kind and the target
    Given findings of each kind for the gate "unit-complete" on "a/b.spec.md"
    When their card titles are built
    Then each names the gate, its kind's word and the target

  @GHIGT-B03 @unit-level
  Scenario: The search covers every state in the configured repository
    Given the repository "acme/x" and a finding with no card
    When the finding is opened on GitHub
    Then the search asks for all states, up to 500 cards, by the key, in "acme/x"

  @GHIGT-B04 @unit-level
  Scenario: A new card carries the marker and the labels
    Given a finding owned by the user with no card
    When it is opened on GitHub
    Then the card is created with the key marker and the labels workflow, to-do and needs-user

  @GHIGT-B05 @unit-level
  Scenario: An assumed debt's card has no flow label
    Given a finding with no card opened as a future debt
    When it is opened on GitHub
    Then the card carries neither the to-do nor the needs-user label

  @GHIGT-B06 @unit-level
  Scenario: An open card is left alone
    Given a finding whose card 12 is open
    When it is opened on GitHub
    Then nothing is created and only the search runs

  @GHIGT-B07 @unit-level
  Scenario: A closed card is reopened with the new report
    Given a finding with the detail "no feature" whose card 12 is closed
    When it is opened on GitHub
    Then card 12 is reopened and commented with "no feature"

  @GHIGT-B08 @unit-level
  Scenario: Resolving closes only an open card
    Given a finding whose card 5 is open, then closed, then missing
    When it is resolved each time
    Then card 5 is closed with a comment the first time, and nothing else runs the other times

  @GHIGT-B09 @unit-level
  Scenario: The labels applied are the ones init creates
    Given the labels anchors init creates
    When the to-do and needs-user labels applied to cards are checked
    Then both are created by init and the needs-user label is not the legacy Portuguese name

  @GHIGT-E01 @unit-level
  Scenario: A gh failure is reported with gh's output
    Given gh failing on the search, on create, on reopen, on comment, or on close, or answering unreadable JSON
    When the finding is opened or resolved
    Then the error is returned carrying "HTTP 502" when gh said it, and nothing is reported done
