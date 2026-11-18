# language: en
# @anchors
#   code: PTRFF
#   ref: PTRFS
#   updated_at: 2026-11-06
#   layer: feature

@PTRFS
Feature: PathRef — the file a path of anchors.yaml names, in the project or at a URL

  @PTRFS-B01 @unit-level
  Scenario: A path reads the project's file, and a URL the external file
    Given a guide in the project and one at a pinned URL
    When each is read
    Then the local one is read from the root and the URL's from its server, a 404 being a missing file and a 500 an error

  @PTRFS-B02 @unit-level
  Scenario: {{branch}} reads the current branch's mirror, and the default when the mirror does not exist
    Given a URL with {{branch|develop}} and an external repository with develop and feature/x
    When it is read on feature/x, on feature/y and with no branch checked out
    Then feature/x reads its mirror, and the others read develop, each naming the URL that answered

  @PTRFS-B03 @unit-level
  Scenario: With the server unreachable, the last copy kept answers
    Given a URL read once on a branch with no mirror
    When the server goes down and it is read again
    Then the default's copy answers

  @PTRFS-E01 @unit-level
  Scenario: A URL that answers nothing, with or without the default, or unreachable with no copy, fails naming what was tried
    Given a URL whose mirror and default do not exist, a URL asking for a branch with none and no default, and an unreachable server with nothing kept
    When each is read
    Then the read fails naming the URLs tried
