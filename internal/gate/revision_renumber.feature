# language: en
# @anchors
#   ref: RVRNR
#   updated_at: 2026-09-23
#   layer: feature

@RVRNR
Feature: RevisionRenumber — the revisions a branch added move to a free number when the base took theirs

  @RVRNR-B01 @unit-level
  Scenario: A revision the branch added whose number the base already uses moves to the next free number
    Given the base and the branch each added `PRICX-R0003`
    When the renumbering is planned
    Then the branch's `R0003` moves to `R0004`

  @RVRNR-B02 @unit-level
  Scenario: After a rebase, only the branch's revision moves, not the base's with the same number
    Given a rebased file holding the base's `R0003` and the branch's `R0003`
    When the renumbering is planned and applied
    Then the base's `R0003` is untouched
    And the branch's becomes `R0004`

  @RVRNR-B03 @unit-level
  Scenario: When one added revision collides, every added revision of that code moves in order
    Given the branch added `R0003` and `R0004` and the base added `R0003`
    When the renumbering is planned
    Then `R0003` moves to `R0004` and `R0004` moves to `R0005`

  @RVRNR-B04 @unit-level
  Scenario: A revision the branch added whose number the base does not use stays
    Given the branch added `R0002` and the base added nothing
    When the renumbering is planned
    Then nothing is renumbered

  @RVRNR-B05 @unit-level
  Scenario: A citation is rewritten only on a line the branch added
    Given a merge-base line and a branch line both citing `PRICX-R0003`
    When `R0003` is renamed to `R0004`
    Then the merge-base line still cites `R0003`
    And the branch line cites `R0004`

  @RVRNR-B06 @unit-level
  Scenario: The branch adding the same number twice is refused, naming the code
    Given the branch added two revisions numbered `R0002`
    When the renumbering is planned
    Then it is refused with an error naming `PRICX-R0002`

  @RVRNR-I01 @unit-level
  Scenario: A revision the base has is never renumbered, even with its explanation edited
    Given the branch edited the text of the base's `R0002`
    When the renumbering is planned
    Then nothing is renumbered

  @RVRNR-I02 @unit-level
  Scenario: Rewriting is one pass, and a chain of renames never cascades
    Given renames `R0003→R0004` and `R0004→R0005`
    When a line citing both is rewritten
    Then it cites `R0004` and `R0005`

  @RVRNR-X01 @unit-level
  Scenario: Does not rewrite a revision cited without its unit code
    Given a line citing `R0003` without its unit code
    When `PRICX-R0003` is renamed
    Then the line is unchanged
