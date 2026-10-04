# language: en
# @anchors
#   code: NTFTN
#   ref: NTFCT
#   updated_at: 2026-10-03
#   layer: feature

@NTFCT
Feature: Notifications — a message to every agent, read from one file and printed on top of `next`

  @NTFCT-B01 @unit-level
  Scenario: Local mode reads notifications.md at the project root
    Given a project root without notifications.md
    When the local message is read
    Then it is empty
    And after notifications.md is written with "hello agents", the local message is "hello agents"

  @NTFCT-B02 @unit-level
  Scenario: A file with only comments or nothing prints nothing
    Given a notifications file holding only an explanatory HTML comment, and an empty one
    When each is printed
    Then nothing is printed

  @NTFCT-B03 @unit-level
  Scenario: Github mode reads the file raw from the integration branch
    Given the repository "o/r" and the integration branch "develop"
    When the message is read from the platform
    Then the call asks for "repos/o/r/contents/notifications.md?ref=develop" with the raw media type
    And the branch "feat/x y" is escaped as "feat%2Fx+y"
    And without a branch the address carries no ref

  @NTFCT-B04 @unit-level
  Scenario: A file missing on the platform is silence
    Given a platform that answers "Not Found (HTTP 404)"
    When the message is read from the platform
    Then the message is empty and there is no error

  @NTFCT-B05 @unit-level
  Scenario: The block names the file and its source and indents the message
    Given the message "Update anchors to 0.1.179." read from "o/r@develop"
    When it is printed
    Then the output says "NOTIFICATIONS", "notifications.md" and "o/r@develop"
    And the line "  Update anchors to 0.1.179." is indented

  @NTFCT-I01 @unit-level
  Scenario: The explanatory comment is never printed
    Given a file with the comment "Write here what every agent must read" and a message
    When it is printed
    Then the message is printed and "Write here" is not

  @NTFCT-X01 @unit-level
  Scenario: Github mode does not read the agent's checkout
    Given a platform that answers the message "update the binary"
    When the message is read from the platform
    Then the text comes from the platform's answer and the call goes to the platform API

  @NTFCT-E01 @unit-level
  Scenario: A failed read is reported as one line
    Given a platform that answers "API rate limit exceeded (HTTP 403)"
    When the message is read and printed
    Then the read returns an error
    And the output is one line saying it "could not read" the file
