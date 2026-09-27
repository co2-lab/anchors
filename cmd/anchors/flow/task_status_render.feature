# language: en
# @anchors
#   ref: TSRTS
#   updated_at: 2026-09-26
#   layer: feature

@TSRTS
Feature: TaskStatusReport — the format of the round's report: where the task is, the verdict, what is missing, and what comes next

  @TSRTS-B01 @unit-level
  Scenario: The sections come in the order that decides
    Given a state with a card, a reversion, a pull request, a branch and a decision waiting on a person
    When the report is rendered
    Then "Task", "UNDONE", "PR", "Git", "Waiting on a person's decision", "What I proved", "What was left out" and "Next" appear in that order

  @TSRTS-B02 @unit-level
  Scenario: The card line drops the code and translates the state
    Given card 301 "[X] screen" under review owned by "host/dev1", and a state with no card
    When each report is rendered
    Then the first has "Task  #301 · screen", "under review" and "owner: host/dev1"
    And the second says "no card found — provide `--card N`"

  @TSRTS-B03 @unit-level
  Scenario: A reversion appears before the verdict, with the way to authorise it
    Given card 483 whose close by hand was reverted
    When the report is rendered
    Then "UNDONE" appears before the pull request line, with what was undone
    And the report names the label "anchors:manual"
    And it says the merge moves the card to ready-to-test, never that the merge closes it
    And without a reversion there is no "UNDONE" section

  @TSRTS-B04 @unit-level
  Scenario: The pull request line never hides missing or mixed checks
    Given a pull request with no check, one with three passed and one failed, and one with four passed
    When each report is rendered
    Then the first pull request line says "no check ran"
    And the second lists the failed before the passed
    And the third reads "4/4 passou"
    And a state without a pull request says "PR    none for this branch"

  @TSRTS-B05 @unit-level
  Scenario: The working-tree line says what is not yet shared
    Given a dirty tree, a clean tree with two unpushed commits, and a clean pushed tree on branch "impl-x"
    When each report is rendered
    Then the lines say "uncommitted change", "2 unpushed commit(s)" and "up to date with the remote"

  @TSRTS-B06 @unit-level
  Scenario: The decisions waiting on a person have their own section
    Given decisions 99 "[DEC1] which vocabulary?" and 12 "[DEC2] where is the limit?"
    When the report is rendered
    Then "Waiting on a person's decision" lists #99 and #12 without their codes
    And it says "no agent resolves these"
    And without such decisions the section does not exist

  @TSRTS-B07 @unit-level
  Scenario: The two gaps are always explicit
    Given a state with no card and a clean tree
    When the report is rendered
    Then it has "What I proved" and "What was left out"

  @TSRTS-B08 @unit-level
  Scenario: Uncommitted and unpushed work come first
    Given a card in progress with an uncommitted change and an unpushed commit
    When the next step is derived
    Then the first step names the uncommitted change and the second "git push"
    And no step says to open the pull request

  @TSRTS-B09 @unit-level
  Scenario: Without a pull request the step depends on where the card is
    Given a card in progress on a clean pushed tree, a card under review, a card ready for review and a card to do, none with a pull request
    When the next step is derived for each
    Then the card in progress is told to open the PR, whose `anchors pr-body` lines link the card without closing it
    And the cards under or ready for review are told `gh pr list --search` finds theirs
    And the card to do is not told to open a PR

  @TSRTS-B12 @unit-level
  Scenario: The check classes are printed in the user's language
    Given a pull request with one failed, one running and one passed check
    When the report is rendered in English, in Portuguese and in Spanish
    Then the pull request line reads "1 failed, 1 running, 1 passed", "1 reprovou, 1 em curso, 1 passou" and "1 falló, 1 en curso, 1 pasó"

  @TSRTS-B10 @unit-level
  Scenario: With an open pull request the step follows the checks
    Given an open pull request with no check, one with a running check, one with a failed check, and one with all passed
    When the next step is derived for each
    Then they say "do not trust the green that does not exist", "--watch" with "do not end the turn here", "work of THIS card", and "the review is missing"

  @TSRTS-B11 @unit-level
  Scenario: A closed card and an idle state have their steps
    Given a closed card, and a state with nothing to do
    When the next step is derived for each
    Then the first says "`anchors next` asks the claim pipeline for the next one"
    And the second says "`anchors status` says where the project is"

  @TSRTS-I01 @unit-level
  Scenario: A running check never reads as passed
    Given a pull request with three passed checks and one running
    When the report is rendered
    Then "4/4" does not appear and the running check is shown

  @TSRTS-X01 @unit-level
  Scenario: Rendering looks nothing up
    Given a platform that records every call
    When a full report is rendered
    Then no call was made
