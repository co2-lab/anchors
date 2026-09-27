# language: en
# @anchors
#   ref: SGCMS
#   updated_at: 2026-09-26
#   layer: feature

@SGCMS
Feature: SuggestCommand — the proposed fixes are listed, shown, applied or rejected, and every decision keeps its record

  @SGCMS-B01 @unit-level
  Scenario: List shows the ids of a state and points pending ones at show
    Given a repository with one pending suggestion fix-a
    When list runs, then list --state rejected after fix-a is rejected
    Then the first prints "1 suggestion(s) in pending", "fix-a" and "anchors suggest show <id>"
    And the rejected list prints fix-a without pointing at show
    And and an empty project prints "no suggestion in pending"

  @SGCMS-B02 @unit-level
  Scenario: Show prints the reason and the diff, or says it was not found
    Given a repository with one pending suggestion fix-a
    When show runs for fix-a and for nope
    Then fix-a prints its reason and "+new"
    And nope fails with "not found in pending"

  @SGCMS-B03 @unit-level
  Scenario: A dry run only checks the patch
    Given a pending suggestion whose patch turns old into new
    When apply runs with --dry-run
    Then it says the patch applies cleanly
    And a.txt still reads old and the suggestion is still pending

  @SGCMS-B04 @unit-level
  Scenario: Apply patches the file and then approves with a default reason
    Given a pending suggestion whose patch turns old into new
    When apply runs with no reason
    Then a.txt reads new
    And the suggestion is approved and its record says "patch applied without reservations"

  @SGCMS-B05 @unit-level
  Scenario: Rejecting needs a reason and keeps the record
    Given a pending suggestion fix-a
    When reject runs without a reason, then with a reason and --auto
    Then the first is refused and fix-a stays pending
    And the second moves it to rejected with the reason and the automatic marker in the record

  @SGCMS-I01 @unit-level
  Scenario: A failed apply leaves the file and the state as they were
    Given a pending suggestion whose file changed since the proposal
    When apply runs
    Then a.txt keeps the changed content and the suggestion stays pending

  @SGCMS-X01 @unit-level
  Scenario: The command decides only suggestions others proposed
    Given a repository whose only suggestion was written by the proposer
    When list runs
    Then it lists exactly that suggestion

  @SGCMS-E01 @unit-level
  Scenario: A stale patch fails before touching anything
    Given a pending suggestion whose file changed since the proposal
    When apply runs
    Then it fails saying the patch no longer matches

  @SGCMS-E02 @unit-level
  Scenario: Outside git the failure names the suggestion patch
    Given a project whose .git was removed
    When the patch is applied
    Then it fails explaining it could not apply the suggestion patch
