# language: en
# @anchors
#   ref: SCLTE
#   updated_at: 2026-09-26
#   layer: feature

@SCLTE
Feature: Escalate — open the right card for a change the plan, the spec or the tool needs, and stop the work only when it must

  @SCLTE-B01 @unit-level
  Scenario: escalate without a reason is refused
    Given a project in github mode
    When `anchors escalate --card 44` runs with no reason
    Then the command fails and no card is created

  @SCLTE-B02 @unit-level
  Scenario: Contradictory exits are refused before any call
    Given a project in github mode
    When `anchors escalate` runs with `--bug --for-user`, with `--bug --unsure`, and with `--blocking` alone
    Then the first two fail with "`--bug` is not a decision" and the third with "`--blocking` goes with `--bug`"
    And no call reaches the platform

  @SCLTE-B03 @unit-level
  Scenario: escalate in local mode is refused
    Given a project in local mode
    When `anchors escalate x` runs
    Then the command fails saying the command exists in github mode

  @SCLTE-B04 @unit-level
  Scenario: The card of the reviewed pull request is the origin
    Given pull request 556 whose body has the lines "Closes #198" and "Refs #200"
    When `anchors escalate --reviewing-pr "#556" "the loop truncates"` runs
    Then the new card carries under-198 and from-pr-556
    And standard error says "PR #556 declares card #198"
    And a body with "Refs #735", "closes #42" or "Fixes #7" at the start of a line names 735, 42 and 7, while "the card refs #99 is another's" names none

  @SCLTE-B05 @unit-level
  Scenario: The one card in hand is the origin, and two are ambiguous
    Given the agent "host/dev1" holding only card 44
    When `anchors escalate something` runs
    Then the new card carries under-44 and standard error says "born under card #44 (from `anchors-owner`)"
    And when the agent holds cards 44 and 45, the command fails with "2 cards in hand (#44, #45)" and creates nothing

  @SCLTE-B06 @unit-level
  Scenario: A finding with no origin is created and warned
    Given pull request 556 whose body names no card, and no agent card
    When `anchors escalate --reviewing-pr 556 "the loop truncates"` runs
    Then standard error says "PR #556 declares no card" and "WITHOUT provenance"
    And the new card carries no under label

  @SCLTE-B07 @unit-level
  Scenario: The title is the exit's prefix and the reason's first line
    Given the reason "the plan misses migrations" followed by a second line
    When it is escalated as an ordinary card, as a decision, as a framing question and as a bug
    Then the titles start with "[plan]", "[decision]", "[framing]" and "[bug]" followed by "the plan misses migrations"
    And a reason longer than seventy characters becomes a title of at most seventy ending in "..."

  @SCLTE-B08 @unit-level
  Scenario: The labels follow the exit
    Given a project whose first workflow label is "anchors"
    When a finding is escalated as an ordinary card, as a decision, as a framing question and as a bug
    Then every card carries "anchors" and "anchors:to-do"
    And the decision adds needs-user, the framing question adds needs-user and needs-framing, and the bug adds the bug label

  @SCLTE-B09 @unit-level
  Scenario: The origin card and the reviewed pull request are both labels
    Given origin card 44
    When `anchors escalate --card 44 "the plan misses migrations"` runs
    Then the label under-44 is created and carried by the new card
    And an escalation reviewing pull request 556 whose card is 198 carries both from-pr-556 and under-198

  @SCLTE-B10 @unit-level
  Scenario: The body says why and how to go on for each exit
    Given the reason "The spec asks for a cache; the plan said there would be none", the target "plans/0001-foundation.md" and card 12
    When the decision body, the framing body, the ordinary body and the bug body are written
    Then the decision body carries the reason, the target, "R0001", the needs-user label and "#12"
    And the framing body says "THE FIRST QUESTION IS THE FRAMING"
    And the ordinary body names "anchors:under-12", "same PR" and "--for-user" and does not mention needs-user
    And the bug body says "nothing to decide" and, when blocking, "Card #966 waits for it"

  @SCLTE-B11 @unit-level
  Scenario: A target with open cards is warned about, three at most
    Given five open cards returned for "src/x.ts", four of which name it
    When `anchors escalate --card 44 --about src/x.ts dup` runs
    Then standard error lists "#1 a src/x.ts", "#2 b" and "#3 c src/x.ts"
    And it does not list "#4 d"

  @SCLTE-B12 @unit-level
  Scenario: The output names the kind and the address
    Given a platform that answers the creation with the address of issue 900
    When a finding, a decision and a bug are escalated
    Then the output says "finding recorded:", "decision opened:" and "bug reported:" followed by that address

  @SCLTE-B13 @unit-level
  Scenario: A decision stops the origin card and prints the way back
    Given origin card 44 and a platform that creates issue 900
    When `anchors escalate --card 44 --for-user "which currency?"` runs
    Then card 44 receives needs-user and blocked-by-900, and a comment "Stopped: there is an open decision"
    And the output says "card #44 stopped until the decision" and "anchors decided --card 44"
    And a framing question with `--unsure` stops card 44 the same way

  @SCLTE-B14 @unit-level
  Scenario: A blocking bug holds the card by blocked-by alone
    Given origin card 44 and a platform that creates issue 900
    When `anchors escalate --card 44 --bug --blocking "claim picks the wrong card"` runs
    Then card 44 receives blocked-by-900 and a comment that a bug blocks it
    And the output says "stopped until the bug is fixed"

  @SCLTE-B15 @unit-level
  Scenario: An ordinary card and a non-blocking bug let the origin card go on
    Given origin card 44
    When a finding and a non-blocking bug are escalated from it
    Then card 44 is not edited
    And it receives the comments "Plan change recorded from this work" and "Bug reported from this work (the card goes on)"

  @SCLTE-B16 @unit-level
  Scenario: The new card's number is read only from an all-digit address tail
    Given the addresses ".../issues/900", ".../issues/900" with a line break, ".../issues/abc", ".../issues/" and "no url"
    When the number is read from each
    Then they give "900", "900", and nothing for the other three

  @SCLTE-I01 @unit-level
  Scenario: Without an origin card no other card is touched
    Given pull request 556 whose body names no card, and no agent card
    When `anchors escalate --reviewing-pr 556 "the loop truncates"` runs
    Then no card is edited and no card is commented

  @SCLTE-X01 @unit-level
  Scenario: A bug never carries needs-user
    Given origin card 44
    When `anchors escalate --card 44 --bug --blocking "claim picks the wrong card"` runs
    Then neither the new card nor card 44 receives needs-user

  @SCLTE-E01 @unit-level
  Scenario: A card the platform refuses to create fails the command
    Given a platform that answers the creation with "label not found"
    When `anchors escalate --card 44 x` runs
    Then the command fails with "open the issue" and "label not found"

  @SCLTE-E02 @unit-level
  Scenario: A card that cannot be stopped is warned about
    Given a platform that refuses to edit card 44
    When `anchors escalate --card 44 --for-user "which currency?"` runs
    Then the command succeeds
    And the output says "could not label card #44" and not "stopped until the decision"
