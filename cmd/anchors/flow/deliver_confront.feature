# language: en
# @anchors
#   code: DCFDL
#   ref: DLCND
#   updated_at: 2026-10-03
#   layer: feature

@DLCND
Feature: DeliveryConfront — confront what a delivery declares against the disk, at the moment it is declared

  @DLCND-B01 @unit-level
  Scenario: The confrontation never blocks the delivery
    Given a delivery declaring an untouched file, on a unit where an informative gate fails
    When the delivery is recorded
    Then the command succeeds and the record exists

  @DLCND-B02 @unit-level
  Scenario: Without git the output says the confrontation did not happen
    Given a project root outside any git repository, with the declared file "a.go"
    When the declared files are confronted
    Then the output says "could not confront the declared files" and "nobody verified"
    And no file is accused

  @DLCND-B03 @unit-level
  Scenario: A committed and untouched declared file is accused
    Given a git repository where "quiet.go" is committed and unchanged
    When the declared file "quiet.go" is confronted
    Then the output lists "quiet.go" under "declared files that do NOT appear in the diff"

  @DLCND-B04 @unit-level
  Scenario: Modified files and files in a new directory are not accused
    Given a git repository where "src/pricing.ts" is modified and "src/newdir/handler.ts" is in a new directory
    When those declared files are confronted
    Then neither is accused

  @DLCND-B08 @unit-level
  Scenario: A renamed file and a file under a subdirectory root are not accused
    Given a git repository where "app/src/old.ts" was renamed to "app/src/renamed.ts", "app/src/pricing.ts" is modified and "app/src/quiet.ts" is untouched
    When "src/renamed.ts", "src/pricing.ts" and "src/quiet.ts" are confronted with the root "app"
    Then only "src/quiet.ts" is accused
    And confronted from the top of the repository, "app/src/renamed.ts" is not accused

  @DLCND-B05 @unit-level
  Scenario: A tested unit without a mutation signal is warned
    Given the unit "src/pricing.ts" with the test "src/pricing.test.ts" and a map with no mutation signal for it
    When the unit is checked
    Then the warning says "NO mutation signal ingested"
    And a unit without a test, a unit with killed mutants in the map, and an empty unit get no warning

  @DLCND-B06 @unit-level
  Scenario: A failing informative gate is listed by its first sentence
    Given a project whose informative gate "always-red" fails on "src/pricing.ts" with "the loop drops the last page. Details follow here"
    When the delivery of "src/pricing.ts" is confronted
    Then the output lists "always-red @ src/pricing.ts — the loop drops the last page."
    And it does not print "Details follow here"

  @DLCND-B07 @unit-level
  Scenario: Without gates, map or delivered node no gate is listed
    Given a project with no configured gates, one with gates and no map, and one whose map lacks the delivered file
    When the delivered files are run through the gates
    Then no gate is listed in any of them

  @DLCND-I01 @unit-level
  Scenario: Not looking and finding nothing are different answers
    Given a root outside git, and a git repository where every declared file is modified
    When the declared files are compared in each
    Then the first reports the comparison did not happen
    And the second reports it happened and accuses no file

  @DLCND-X01 @unit-level
  Scenario: The confrontation writes nothing
    Given a git repository with a modified file and a committed one
    When the declared files are confronted
    Then git reports the same working tree before and after
