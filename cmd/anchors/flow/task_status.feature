# language: en
# @anchors
#   code: TSFTS
#   ref: TSSTT
#   updated_at: 2026-10-03
#   layer: feature

@TSSTT
Feature: TaskStatus — discover what the machine knows about the task at hand, so the agent's report does not have to

  @TSSTT-B01 @unit-level
  Scenario: The working tree is read from the root
    Given a git repository on branch "main" with an untracked file
    When the task state is collected
    Then the branch is "main" and the tree is not clean

  @TSSTT-B02 @unit-level
  Scenario: A card given by number carries its state, and a closed card reads closed
    Given card 303 open with the label "anchors:in-progress", and card 9 closed with the label "anchors:ready-to-review"
    When each card is read by number
    Then card 303's state is "anchors:in-progress"
    And card 9's state is "closed" and its labels still include "anchors:ready-to-review"

  @TSSTT-B03 @unit-level
  Scenario: Without a number the agent's own card is found on the board
    Given the session "dev7" and a board where card 42 is owned by this host's dev7 and card 41 by someone else
    When `anchors task-status` runs without `--card`
    Then the report is about card 42

  @TSSTT-B04 @unit-level
  Scenario: The decisions waiting on a person are listed
    Given open cards 701 "which currency" and 702 "which region" waiting on a person
    When the task state is collected in github mode
    Then both are listed with their numbers and titles

  @TSSTT-B05 @unit-level
  Scenario: Running checks are running, not failed
    Given a pull request with a succeeded, an in-progress, a queued, a failed, a pending commit status and a successful commit status
    When the pull request of the branch is read
    Then three are running, one failed and two passed

  @TSSTT-B06 @unit-level
  Scenario: Only the lock's own reversal counts, reduced to a clean first line
    Given card 303 with a reversal written by "github-actions" and a comment by "bob" starting with the same marker
    When the task state is collected
    Then exactly one undone move is reported: "Reverted: closed by hand (by bob)"
    And a comment by the automation that does not start with the marker is not a reversal

  @TSSTT-B07 @unit-level
  Scenario: The command prints the report
    Given the agent's card 42 titled "[ABCDE] mine"
    When `anchors task-status` runs
    Then the output has "Task  #42 · mine"

  @TSSTT-B08 @unit-level
  Scenario: The turn-ended event carries numbers and vocabulary only
    Given a card in review titled "secret title" and a pull request with three failed checks of four
    When the turn-ended event is sent
    Then it carries "in-review", "open" and the failed-checks count
    And its attribute names are the English identifiers "card_state", "pr_state", "has_card", "clean_tree", "checks_failed" and "checks_running"
    And it does not carry "secret title" nor "anchors:in-review"

  @TSSTT-I01 @unit-level
  Scenario: The check classes add up to the total
    Given a pull request with six checks of every kind
    When they are classified
    Then running, failed and passed add up to six

  @TSSTT-X01 @unit-level
  Scenario: Every lookup names the configured repository and the root's branch
    Given a project configured for "acme/app" whose root is on branch "feat-303"
    When the task state is collected for card 303
    Then the card, waiting-list, pull request and comments lookups all pass "--repo acme/app"
    And the pull request looked up is "feat-303"

  @TSSTT-X02 @unit-level
  Scenario: Local mode looks up no card
    Given a project in local mode and a board that would answer with card 303
    When the task state is collected for card 303
    Then there is no card and no card lookup is made

  @TSSTT-E01 @unit-level
  Scenario: Failed and unreadable lookups leave the part absent
    Given a platform that answers "not json", and then one that fails every call
    When each source is read, and the task state is collected
    Then no card, no waiting list, no pull request and no undone move are reported
