# language: en
# @anchors
#   ref: RNMBR
#   updated_at: 2026-09-26
#   layer: feature

@RNMBR
Feature: Renumber — moves the revisions a branch added when the base already took their number

  @RNMBR-B01 @unit-level
  Scenario: The default base is the remote integration branch, else the local one
    Given a branch forked from main in a repository with no remote copy of main
    When renumber runs with no base
    Then it renumbers against "main"
    And once a remote copy "origin/main" exists, it renumbers against "origin/main"

  @RNMBR-B02 @unit-level
  Scenario: Specs changed but not committed, and untracked specs, are examined
    Given a branch whose new revision PRICX-R0002 is in the working tree but not committed
    And another branch whose spec colliding with the base is not tracked yet
    When renumber runs with no file against main
    Then both colliding revisions are listed to move to the next number

  @RNMBR-B03 @unit-level
  Scenario: Only the branch's colliding revision moves to the next free number
    Given main and the branch each added PRICX-R0002 to the pricing spec after forking
    When renumber runs against main
    Then the branch's revision becomes PRICX-R0003 and main's version is not touched

  @RNMBR-B04 @unit-level
  Scenario: Citations move only on the lines the branch added
    Given a rebased branch whose pricing spec and test cite PRICX-R0002 both on a line from the base and on a line the branch added
    When renumber runs against main
    Then the branch's line cites PRICX-R0003 and the base's line still cites PRICX-R0002

  @RNMBR-B05 @unit-level
  Scenario: A changed binary file is not rewritten
    Given the branch changed a binary file that contains PRICX-R0002
    When renumber moves PRICX-R0002 to PRICX-R0003
    Then the binary file is byte for byte the same

  @RNMBR-B06 @unit-level
  Scenario: The dry run writes nothing
    Given main and the branch each added PRICX-R0002 to the pricing spec
    When renumber runs against main with the dry run switch
    Then the output shows "PRICX-R0002 → PRICX-R0003" and "(dry-run — nothing was written.)"
    And the pricing spec on disk does not contain R0003

  @RNMBR-B07 @unit-level
  Scenario: No collision says there is nothing to renumber
    Given a branch whose revisions were already renumbered
    When renumber runs against main again
    Then it says "✓ no revision this branch added collides with main — nothing to renumber."

  @RNMBR-B08 @unit-level
  Scenario: A rewritten file keeps its permission bits
    Given the branch changed an executable script that cites PRICX-R0002
    When renumber rewrites the citation
    Then the script is still executable with mode 0755

  @RNMBR-B09 @unit-level
  Scenario: With files given, only those specs are examined
    Given a branch where both the pricing spec and the tax spec collide with main
    When renumber runs with only the tax spec given
    Then only the tax spec's revision is listed to move

  @RNMBR-I01 @unit-level
  Scenario: The base's revisions keep their meaning after a renumber
    Given a rebased branch whose pricing spec holds main's PRICX-R0002 and its own PRICX-R0002
    When renumber runs against main
    Then main's PRICX-R0002 block and main's citation of it are unchanged

  @RNMBR-X01 @unit-level
  Scenario: The rewrite is left unstaged for review
    Given main and the branch each added PRICX-R0002 to the pricing spec
    When renumber runs against main
    Then the pricing spec is modified in the working tree and nothing is staged

  @RNMBR-E01 @unit-level
  Scenario: A project with no configuration is refused
    Given a directory with no anchors.yaml
    When renumber runs there
    Then it fails naming "anchors.yaml"

  @RNMBR-E02 @unit-level
  Scenario: A base with no merge base is refused naming it
    Given a repository with no branch called "no-such-branch"
    When renumber runs with that base
    Then it fails with "no merge base between no-such-branch and HEAD"

  @RNMBR-E03 @unit-level
  Scenario: An unreadable spec given by hand is refused naming it
    Given a spec path given by hand that does not exist
    When renumber runs with it and the pricing spec
    Then it fails with "read Missing.spec.md" and the pricing spec is unchanged
