# language: en
# @anchors
#   code: SPFTS
#   ref: UNCDN
#   updated_at: 2026-10-07
#   layer: feature

@UNCDN
Feature: UnitCodes — the identity code of a unit, read from a header, from the map, or from the codes a file names

  @UNCDN-B01 @unit-level
  Scenario: The header code is read in any comment style
    Given headers whose code line is bare inside a markup block, or behind "//", "#", " *" or "<!--"
    When the header code is read from each
    Then the codes "RLSGR", "ALFDL", "HSHCM", "BLKCM" and "HTMLC" are read

  @UNCDN-B02 @unit-level
  Scenario: A code of the wrong length or case is not read
    Given the project code length of five and the lines "code: ABCD", "code: ABCDEF" and "some code: lower"
    When the header code is read from each
    Then every answer is the empty string

  @UNCDN-B03 @unit-level
  Scenario: The exact map node answers with its code
    Given a map where "pkg/beta.ts" carries the code "BETAX"
    When the code of the unit "pkg/beta.ts" is asked
    Then the answer is "BETAX"

  @UNCDN-B04 @unit-level
  Scenario: A node of the same stem answers when the exact node has no code
    Given a map where "pkg/alpha.spec.md" carries "ALPHA" and "pkg/alpha.ts" carries no code
    When the codes of "pkg/alpha.ts" and of the absent "pkg/gamma.ts" are asked
    Then the answers are "ALPHA" and the empty string

  @UNCDN-B05 @unit-level
  Scenario: Without a map a unit has no code
    Given a project directory with no map file
    When the code of the unit "pkg/beta.ts" is asked
    Then the answer is the empty string

  @UNCDN-B06 @unit-level
  Scenario: The codes of a file are the unit's own, once each, in order
    Given a file naming "ALPHA-B01", "ALPHA-B02", "BETAX-B07" and "ALPHA-B01" again
    When its codes are asked for the unit "ALPHA", and then for no unit
    Then the first answer is "ALPHA-B01", "ALPHA-B02" and the second holds three codes

  @UNCDN-I01 @unit-level
  Scenario: No listed code belongs to another unit or repeats
    Given a file naming its own codes twice and another unit's code once
    When its codes are asked for its unit
    Then the list holds each of its own codes once and none of the other unit

  @UNCDN-E01 @unit-level
  Scenario: A file that cannot be read is an error
    Given a path to a file that does not exist
    When its codes are asked for the unit "ALPHA"
    Then an error is returned

  @UNCDN-B07 @unit-level
  Scenario: A data state written bare is the unit's own code
    Given a spec with a rule, a bare data state, a prefixed data state of its own and one of another unit
    When the unit's codes are listed
    Then the rule, the prefixed data state and the bare one as the unit's code are listed, and the other unit's is not
