# language: en
# @anchors
#   code: RWFTR
#   ref: RCRWR
#   updated_at: 2026-10-03
#   layer: feature

@RCRWR
Feature: RecodeRewrite — renaming an identity code inside a text, on every surface where it appears

  @RCRWR-B01 @unit-level
  Scenario: A code is well formed only in the project's lengths and alphabet
    Given a project declaring the code lengths 4 and 5
    When "TCDTX", "MN01", "ABCDX", "ABCDE", "TCD", "ABCDEF", "tcdt" and "TC-T" are validated
    Then the first four are valid and the last four are not

  @RCRWR-B02 @unit-level
  Scenario: The header code is replaced in every header style
    Given a markdown header "code: TCDTX", a line-comment header "ref: TCDTX" and a hash-comment header "ref: TCDTX"
    When "TCDTX" is rewritten to "TCTXX"
    Then each header holds "TCTXX" and no "TCDTX"

  @RCRWR-B03 @unit-level
  Scenario: Only the old code changes in a ref list
    Given the header line "ref: TCDTX, MNDTX"
    When "TCDTX" is rewritten to "TCTXX"
    Then the line holds "TCTXX" and "MNDTX" and no "TCDTX"

  @RCRWR-B04 @unit-level
  Scenario: Scenario codes keep their suffixes
    Given the text "TCDTX-B01 TCDTX-S04 TCDTX-A04 TCDTX-FP01b TCDTX-DS-receipt-none TCDTX-VR TCDTX-RA03"
    When "TCDTX" is rewritten to "TCTXX"
    Then every code reads "TCTXX" with its suffix and 7 replacements are counted

  @RCRWR-B05 @unit-level
  Scenario: A bare mention is replaced and its neighbour kept
    Given the text "Ver a regra em TCDTX (a tela de detalhe)."
    When "TCDTX" is rewritten to "TCTXX"
    Then the text is "Ver a regra em TCTXX (a tela de detalhe)."
    And one replacement is counted

  @RCRWR-B06 @unit-level
  Scenario: The dry run classifies each occurrence
    Given a text with a "ref: TCDTX" header, the scenario code "TCDTX-S02" and a bare "TCDTX"
    When the occurrences of "TCDTX" are listed
    Then there is one scenario code, one header and one bare reference
    And a text with two of each kind lists two scenario codes, two headers and two bare references

  @RCRWR-B07 @unit-level
  Scenario: Each listed occurrence carries its line, counted from one
    Given a text with the scenario code "TCDTX-B01" on line 1 and "TCDTX-S02" on line 4
    When the occurrences of "TCDTX" are listed
    Then the first is on line 1 and the second on line 4
    And in a text with two scenario codes, then a header, then a mention, each occurrence carries its own text and line

  @RCRWR-I01 @unit-level
  Scenario: A text without the old code is left unchanged
    Given a text that holds "WXYZX" and "WXYZ-B01" but not "TCDTX", and a text already rewritten from "TCDTX"
    When "TCDTX" is rewritten to "TCTXX" in each
    Then both come back unchanged with zero replacements

  @RCRWR-X01 @unit-level
  Scenario: Longer codes and neighbours are never touched
    Given the text "TCDTXX-B01 e XTCDTX e TCDTXABCD"
    When "TCDTX" is rewritten to "TCTXX"
    Then the text is unchanged and zero replacements are counted
