# language: en
# @anchors
#   code: EBFXP
#   ref: BREXB
#   updated_at: 2026-10-03
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

  @BREXB-B04 @unit-level
  Scenario: The active agents are counted from the snapshot
    Given a board taken at 12:00 where alice touched a card 20 minutes before and owns an older one, erin touched one exactly 30 minutes before, bob 31 minutes before, a released card was touched a minute before and dave's card was closed a minute before
    When the active agents are counted
    Then they are alice with 2 cards and erin with 1

  @BREXB-B05 @unit-level
  Scenario: A live board counts the active agents from now
    Given a live board stamped a year ago where alice touched a card 5 minutes ago, bob 45 minutes ago and idle a year ago
    When the active agents are counted
    Then only alice is active, with 1 card

  @BREXB-B06 @unit-level
  Scenario: The agent chips filter the cards by agent
    Given a board where alice and bob touched cards in the window and carol's only card is a year old
    When the agent chips are drawn and alice is selected, then no agent
    Then the chips are alice and bob, in that order, plus the chip for all, and not carol
    And with alice selected only her card shows, and with no agent every card shows
    And a board with nobody active shows "nenhum agente ativo"

  @BREXB-B07 @unit-level
  Scenario: A released card is collected with no owner
    Given a card whose last owner comment is "anchors-owner: (liberado)" and one whose owner comment is "anchors-owner: bob" followed by a report
    When the collect expression runs over them
    Then the first card has no owner and keeps its last update
    And the second card's owner is "bob"

  @BREXB-B08 @unit-level
  Scenario: With no card waiting for a person the blocked strip is hidden and empty
    Given a board with a single to-do card
    When the blocked strip is drawn
    Then the strip is hidden and has no content

  @BREXB-B09 @unit-level
  Scenario: The strip orders the waiting cards by how many cards each blocks
    Given card 99 waiting for a person and blocking two cards, and card 12 waiting for a person and blocking one card and card 99
    When the blocked strip is drawn
    Then 99 comes before 12, saying "trava 2 cards" and "trava 1 card"
    And the strip says "Esperando você · 2" and does not show the "[DEC2]" prefix

  @BREXB-B10 @unit-level
  Scenario: An escalated card in a work state appears in the strip with the state it stopped in
    Given card 311 escalated while in review and card 443 escalated while in to-do
    When the blocked strip is drawn
    Then both cards are listed under "Esperando você · 2"
    And it says "parado em in-review" and never "parado em to-do"

  @BREXB-B11 @unit-level
  Scenario: The strip tells a card waiting for a person from one waiting for decided work
    Given escalated cards 311 and 500, and card 444 that unblocks 311
    When the blocked strip is drawn
    Then 311 says "espera a entrega do #444" and 500 does not
    And card 444 is not listed in the strip

  @BREXB-B12 @unit-level
  Scenario: A hostile title is escaped in the strip
    Given a card waiting for a person titled "<img src=x onerror=alert(1)>"
    When the blocked strip is drawn
    Then the strip does not contain the raw "<img src=x" markup

  @BREXB-B13 @unit-level
  Scenario: An escalated card is marked in its column
    Given an escalated card and an ordinary card in review
    When the columns are drawn
    Then exactly one card has the class "card escalado" and the badge "esperando você"

  @BREXB-B14 @unit-level
  Scenario: The roadmap opens a card's details
    Given the board page
    When its click handling is read
    Then the roadmap is among the containers whose clicks open the details
    And the roadmap's card label carries the card number the click selects

  @BREXB-B15 @unit-level
  Scenario: A blocked card says which card blocks it
    Given the board page
    When it draws a card blocked by another
    Then the badge reads "bloqueado por #" with the blocking number and leads to that card
    And the blocking color is defined in the light palette and in both dark ones

  @BREXB-B16 @unit-level
  Scenario: Decisions and framing requests are listed apart
    Given the board page
    When it lists the cards waiting for a person
    Then one list keeps the cards without framing and another the cards with it
    And the framing list says how to send the card back to "anchors:to-do"

  @BREXB-B17 @unit-level
  Scenario: Bugs are listed apart from decisions
    Given an escalated decision and a bug card
    When the blocked strip is drawn
    Then it says "Esperando você · 1" and lists the bug under "Bugs · 1"
    And with only the bug, the strip shows "Bugs · 1" and no "Esperando você"

  @BREXB-I01 @unit-level
  Scenario: The extracted expression never carries the redirection
    Given the board pipeline carried in the binary
    When the collect expression is extracted
    Then it opens with "[" and closes with "]"
    And it does not contain "_board/board.json"

  @BREXB-I02 @unit-level
  Scenario: The page's data attributes are written and read in pairs
    Given the board page
    When the data attributes it writes and the ones it reads are collected
    Then every attribute read is written
    And every attribute written is read

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
