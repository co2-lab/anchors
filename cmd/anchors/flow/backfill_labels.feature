# language: en
# @anchors
#   code: BLFBC
#   ref: BCLBB
#   updated_at: 2026-10-03
#   layer: feature

@BCLBB
Feature: BackfillLabels — write into open cards the blocking and provenance links the board already implies

  @BCLBB-B01 @unit-level
  Scenario: backfill-labels in local mode is refused
    Given a project in local mode
    When `anchors backfill-labels` runs
    Then the command fails saying the command exists in github mode

  @BCLBB-B02 @unit-level
  Scenario: The open decisions are read in one listing
    Given a board with six open decisions
    When `anchors backfill-labels` runs
    Then the needs-user decisions are listed exactly once

  @BCLBB-B03 @unit-level
  Scenario: Every decision under an origin card holds it
    Given decisions 701 and 708 both under open card 44
    When `anchors backfill-labels` runs
    Then card 44 receives both blocked-by-701 and blocked-by-708

  @BCLBB-B04 @unit-level
  Scenario: The provenance is recovered from a cited pull request that declares a card
    Given decision 705 citing "PR #88", whose body says "Refs #50"
    And decision 706 citing "#99" without the word PR, and decision 707 citing "PR #89", whose body declares no card
    When `anchors backfill-labels` runs
    Then decision 705 receives from-pr-88 and under-50, and the output says "#705 ← PR #88, card #50"
    And decisions 706 and 707 are not touched

  @BCLBB-B05 @unit-level
  Scenario: Only an open origin card receives the block
    Given decision 701 under open card 44 and decision 702 under closed card 45
    When `anchors backfill-labels` runs
    Then card 44 receives blocked-by-701 and the output says "#44 ← blocked by #701"
    And card 45 is not touched

  @BCLBB-B06 @unit-level
  Scenario: Precedence between decisions is not inferred
    Given decision 703 under card 46, which is itself an open decision
    When `anchors backfill-labels` runs
    Then card 46 is not touched
    And the output says "#46 is also an open decision"

  @BCLBB-B07 @unit-level
  Scenario: A blocked-by label already present is not written again
    Given card 44 already carrying blocked-by-701, and decisions 701 and 708 under it
    When `anchors backfill-labels` runs
    Then card 44 is edited only to add blocked-by-708

  @BCLBB-B08 @unit-level
  Scenario: A label is created before the card is edited with it
    Given decision 701 under open card 44
    When `anchors backfill-labels` runs
    Then blocked-by-701 is created before card 44 is edited

  @BCLBB-B09 @unit-level
  Scenario: The dry run names what it would do and touches nothing
    Given the board of six open decisions
    When `anchors backfill-labels --dry-run` runs
    Then no label is created and no card is edited
    And the output says "#705 would receive", "#44 would receive" and "(nothing was touched)"

  @BCLBB-B10 @unit-level
  Scenario: The output counts what was recovered, written and skipped
    Given the board of six open decisions, and then a board whose one decision points nowhere
    When `anchors backfill-labels` runs on each
    Then the first says "1 provenance(s) recovered from the cited PR" and "1 link(s) written · 2 skipped"
    And the second says "no link to recover: 1 open decision(s)"

  @BCLBB-I01 @unit-level
  Scenario: A decision under itself is never blocked by itself
    Given decision 704 whose under label names 704
    When `anchors backfill-labels` runs
    Then card 704 is not edited

  @BCLBB-X01 @unit-level
  Scenario: Nothing is removed and nothing unsupported is written
    Given the board of six open decisions
    When `anchors backfill-labels` runs
    Then no call removes a label
    And only cards 44 and 705 are edited

  @BCLBB-E01 @unit-level
  Scenario: A failed listing fails the command
    Given a listing of decisions that fails
    When `anchors backfill-labels` runs
    Then the command fails with "list the open decisions"

  @BCLBB-E02 @unit-level
  Scenario: An unreadable listing fails the command
    Given a listing of decisions that answers "not json"
    When `anchors backfill-labels` runs
    Then the command fails with "read the list"

  @BCLBB-E03 @unit-level
  Scenario: A failed edit is reported and skipped
    Given decision 701 under open card 44, whose edit is refused with "forbidden"
    When `anchors backfill-labels` runs
    Then standard error names #44 and "forbidden"
    And the output says "0 link(s) written · 1 skipped"

  @BCLBB-E04 @unit-level
  Scenario: An origin card whose state cannot be read is skipped
    Given decision 701 under card 44, whose state cannot be read
    When `anchors backfill-labels` runs
    Then card 44 is not edited
    And the output says "0 link(s) written · 1 skipped"
