# language: en
# @anchors
#   ref: CMMSC
#   updated_at: 2026-09-26
#   layer: feature

@CMMSC
Feature: CommitMsg — the commit subject is confronted with the format the changelog will read, before the commit exists

  @CMMSC-B01 @unit-level
  Scenario: The subject is the first line that is neither blank nor a comment
    Given a message file with git comment lines before and after "feat: the thing"
    When the subject is read
    Then it is "feat: the thing"

  @CMMSC-B02 @unit-level
  Scenario: Messages git generates pass
    Given subjects starting with Merge, Revert, fixup!, squash! and Reapply
    When they are confronted
    Then they pass, while "Mergeando o trabalho da semana" does not escape

  @CMMSC-B03 @unit-level
  Scenario: A conventional subject passes, with an optional scope and break mark
    Given the subjects "feat: ok", "fix(board): the refresh count" and "feat(gate)!: changes the contract"
    When they are confronted
    Then each passes

  @CMMSC-B04 @unit-level
  Scenario: A type with an uppercase letter or outside the closed list is refused
    Given the subjects "Feat: uppercase in the type" and "bugfix: a type not in the list"
    When they are confronted
    Then each is refused

  @CMMSC-B05 @unit-level
  Scenario: An empty scope is refused
    Given the subject "feat(): x"
    When it is confronted
    Then it is refused

  @CMMSC-B06 @unit-level
  Scenario: A subject over the limit is refused and the diagnosis points to the body
    Given a subject of 109 characters and one of exactly 100
    When they are confronted
    Then the first is refused with a diagnosis mentioning the BODY, and the second passes

  @CMMSC-B07 @unit-level
  Scenario: A subject ending in a period is refused
    Given the subject "feat: ends in a period."
    When it is confronted
    Then it is refused

  @CMMSC-B08 @unit-level
  Scenario: Each defect gets its own diagnosis
    Given the subjects "Feat: x", "feat(): x", "feat: x.", "bugfix: x", "no format at all" and "feat: "
    When each is confronted
    Then no two diagnoses are the same

  @CMMSC-B09 @unit-level
  Scenario: A capital letter at the start of the subject is allowed
    Given the subject "feat: SBOM leaves the ignored folder"
    When it is confronted
    Then it passes

  @CMMSC-B10 @unit-level
  Scenario: The rejection names the subject, teaches the format and lists the types
    Given a message file whose subject is "Bugfix(): x."
    When commit-msg runs on it
    Then it fails repeating "Bugfix(): x.", showing "type(scope): what changed" and listing every accepted type

  @CMMSC-B11 @unit-level
  Scenario: An empty or comment-only message passes
    Given an empty message file and one holding only comments
    When commit-msg runs on each
    Then it passes

  @CMMSC-B12 @unit-level
  Scenario: Checks run in order and only the first defect is reported
    Given the subject "Bugfix(): x." which has an uppercase type, an empty scope and a final period
    When commit-msg runs on it
    Then it reports only the uppercase type

  @CMMSC-B13 @unit-level
  Scenario: A type with nothing after the colon is refused
    Given the subject "feat: "
    When it is confronted
    Then it is refused with its own diagnosis

  @CMMSC-X01 @unit-level
  Scenario: The command only accepts or refuses
    Given an empty message file
    When commit-msg runs on it
    Then it passes and the file is not rewritten

  @CMMSC-E01 @unit-level
  Scenario: A message file that cannot be read fails naming the read
    Given a path that does not exist
    When commit-msg runs on it
    Then it fails with "read the message"
