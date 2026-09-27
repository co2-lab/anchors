# language: en
# @anchors
#   ref: BREXB
#   updated_at: 2026-09-26
#   layer: feature

@BREXB
Feature: BoardExposure — hand the local board the same page and the same collect contract the pipeline publishes

  @BREXB-B01 @unit-level
  Scenario: The local board page is the published page
    Given the board page carried in the binary
    When the local board asks for its page
    Then the page is byte for byte the carried board page
    And it reads "board.json" and has the columns, detail and roadmap areas

  @BREXB-B02 @unit-level
  Scenario: The collect expression comes out of the pipeline
    Given the board pipeline carried in the binary
    When the local board asks for the collect expression
    Then the expression names the fields "number:", "title:", "state:", "owner:" and "ownership:"
    And it filters out "anchors:discarded"
    And it has no surrounding whitespace

  @BREXB-B03 @unit-level
  Scenario: Both collect shapes are accepted and cut before the redirection
    Given a collect step written as a direct query "--jq '[.[] | {n: .number}]' > _board/board.json"
    And a collect step written as "jq -s '[.[] | {n: .number}]' _board/raw.jsonl > _board/board.json"
    When the collect expression is extracted from each
    Then both give "[.[] | {n: .number}]"

  @BREXB-I01 @unit-level
  Scenario: The extracted expression never carries the redirection
    Given the board pipeline carried in the binary
    When the collect expression is extracted
    Then it opens with "[" and closes with "]"
    And it does not contain "_board/board.json"

  @BREXB-X01 @unit-level
  Scenario: A collect step of another shape yields no expression
    Given a collect step written as "gh issue list --jq '.[] | .x' > out"
    When the collect expression is extracted
    Then no expression is found, rather than a built-in one

  @BREXB-E01 @unit-level
  Scenario: A binary without the board files fails loudly
    Given a binary that carries neither the board page nor the board pipeline
    When the local board asks for its page and its collect expression
    Then the page request fails with "the board HTML is not embedded"
    And the expression request fails with "the board pipeline is not embedded"

  @BREXB-E02 @unit-level
  Scenario: A pipeline that changed shape fails naming the change
    Given a carried board pipeline whose collect step is "gh issue list --jq '.[] | .x' > out"
    When the collect expression is looked for in it
    Then no expression is found, which is the condition the request fails on instead of returning one
