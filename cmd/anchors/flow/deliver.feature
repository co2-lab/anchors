# language: en
# @anchors
#   ref: DLVRE
#   updated_at: 2026-09-26
#   layer: feature

@DLVRE
Feature: Deliver — record what a stage delivered, where the reviewer reads it, and send the author to the review

  @DLVRE-B01 @unit-level
  Scenario: A delivery missing its required parts is refused
    Given a project in local mode
    When `anchors deliver` runs without `--stage`, with the stage "deploy", with a blank intent, and without `--date`
    Then each run fails with "provide --stage", `unknown stage "deploy"`, "provide --intent" and "provide --date"

  @DLVRE-B02 @unit-level
  Scenario: The first delivery of a unit is accepted by its spec
    Given a project where only "src/pricing.spec.md" exists
    When `anchors deliver --stage spec --unit <root>/src/pricing.ts` runs, with the unit absolute
    Then the delivery succeeds
    And the output says "`src/pricing.ts` does not exist yet", names "src/pricing.spec.md" as the piece found and says the record keeps the unit "src/pricing.ts"
    And the record's unit is "src/pricing.ts" and its files list "src/pricing.ts", with no absolute path

  @DLVRE-B03 @unit-level
  Scenario: Prose with commas is one decision, and files split on commas
    Given the decision "eight rules and two invariants, in the letters the neighbours use"
    When the deliver flags are parsed with that decision, an uncovered item with a comma, and `--file a.ts,b.ts`
    Then the decision and the uncovered item are one item each, unchanged
    And the files are "a.ts" and "b.ts"

  @DLVRE-B04 @unit-level
  Scenario: Local mode writes the record under changes
    Given a project in local mode with the unit "src/pricing.ts"
    When `anchors deliver --stage code --unit src/pricing.ts --intent "implements PRICX-B01"` runs
    Then one record exists under "changes/" with the unit and the intent
    And the output says "delivery recorded: changes/"

  @DLVRE-B05 @unit-level
  Scenario: Github mode records on the card given
    Given a project in github mode and open card 77 titled "[plan] F02"
    When `anchors deliver --stage plan --unit src/pricing.ts --card 77` runs
    Then card 77 receives a comment with the intent
    And the output says "delivery recorded on issue #77 — [plan] F02"

  @DLVRE-B06 @unit-level
  Scenario: Github mode finds the card by the unit's code
    Given a map where "src/pricing.ts" has the code PRICX, and open cards 5 "[OTHER] x" and 6 "[PRICX] Implement spec — pricing"
    When `anchors deliver --stage code --unit src/pricing.ts` runs
    Then card 6 receives the record
    And the code is found for "src/pricing.feature" too, and for no unit outside the map

  @DLVRE-B07 @unit-level
  Scenario: A vendored pipeline in github mode records nothing
    Given a project in github mode and the seeded pipeline ".github/workflows/anchors-claim.yml"
    When `anchors deliver --stage code --unit .github/workflows/anchors-claim.yml` runs without `--card`
    Then the command succeeds saying there is nothing to record
    And "changes/" does not exist

  @DLVRE-B08 @unit-level
  Scenario: The next step is the review of the unit
    Given a project in local mode with the unit "src/pricing.ts"
    When the delivery is recorded
    Then the output says "anchors work review --for src/pricing.ts"

  @DLVRE-B09 @unit-level
  Scenario: The watcher hint appears only when the watcher is not running
    Given a project without the watcher's metadata, and one with it
    When a delivery is recorded in each
    Then the first output says "the watcher is not running" and the second does not

  @DLVRE-B10 @unit-level
  Scenario: The zero-decisions note appears only when nothing was declared
    Given a delivery with no decision and no gap, and one with a decision and a gap
    When each is recorded
    Then the first output says "you declared ZERO free decisions" and the second does not

  @DLVRE-I01 @unit-level
  Scenario: A refused delivery records nothing
    Given a project in local mode
    When every refused delivery has run
    Then "changes/" does not exist

  @DLVRE-X01 @unit-level
  Scenario: Github mode never falls back to changes
    Given a project in github mode and open card 77
    When `anchors deliver --card 77` records the delivery
    Then "changes/" does not exist

  @DLVRE-E01 @unit-level
  Scenario: Github mode without a code for the unit fails with the way out
    Given a project in github mode with no map
    When `anchors deliver --stage code --unit src/pricing.ts` runs
    Then the command fails with "could not find the CODE" and "--card <n>"

  @DLVRE-E02 @unit-level
  Scenario: Github mode with no open card for the code fails with the way out
    Given a map where "src/pricing.ts" has the code PRICX, and no open card
    When `anchors deliver --stage code --unit src/pricing.ts` runs
    Then the command fails naming "[PRICX]" and "reopen it"

  @DLVRE-E03 @unit-level
  Scenario: A closed card given with --card is refused
    Given card 77 which is closed
    When `anchors deliver --card 77` runs
    Then the command fails saying the card is closed
