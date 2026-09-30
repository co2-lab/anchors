# language: en
# @anchors
#   ref: CMCLC
#   updated_at: 2026-09-30
#   layer: feature

@CMCLC
Feature: CommonCLI — the contract every command shares: how a path becomes a node, how "not governed" is signalled, what the binary says it is

  @CMCLC-B01 @unit-level
  Scenario: The not-governed error names the path
    Given the path "docs/notes.txt" that no layer governs
    When the not-governed error is built and read
    Then its message contains "\"docs/notes.txt\"" and "not governed", and the path can be read back from the error

  @CMCLC-B02 @unit-level
  Scenario: The not-governed exit code is 3
    Given the exit code the commands use for a path that is not governed
    When a script reads it
    Then it is 3

  @CMCLC-B03 @unit-level
  Scenario: A root-relative path is kept, cleaned
    Given a project root holding "pkg/a/x.spec.md"
    When the argument "pkg/a/../a/x.spec.md" is resolved
    Then the node identifier is "pkg/a/x.spec.md"

  @CMCLC-B04 @unit-level
  Scenario: A path absent under the root is resolved from the working directory
    Given a project root that has no "ghost.ts"
    When the argument "ghost.ts" is resolved
    Then the result is the working directory's "ghost.ts" expressed relative to the root, with forward slashes

  @CMCLC-B05 @unit-level
  Scenario: A node exists only by its exact identifier
    Given a map holding "a/x.spec.md" and "a/x.ts"
    When "a/x.ts", "a/y.ts", and "a/x.ts" against no map are looked up
    Then the answers are true, false and false

  @CMCLC-B06 @unit-level
  Scenario: A task slug drops only the last extension
    Given the paths "pkg/a/x.ts", "pkg/a/x.spec.md" and "Makefile"
    When each is turned into a task slug
    Then the slugs are "pkg/a/x", "pkg/a/x.spec" and "Makefile"

  @CMCLC-B07 @unit-level
  Scenario: An unstamped build says it is a development build
    Given a binary built with no version stamp
    When its version, commit and date are read
    Then they are "dev", "none" and "unknown"

  @CMCLC-I01 @unit-level
  Scenario: Root-relative and absolute names of one file resolve to one node
    Given a project root holding "pkg/a/x.spec.md"
    When the file is resolved by its root-relative path and by its absolute path
    Then both give "pkg/a/x.spec.md"

  @CMCLC-E01 @unit-level
  Scenario: A path that cannot be related to the root is kept as given
    Given the relative root "relroot" and the absolute argument "/abs/x"
    When the argument is resolved
    Then the result is "/abs/x"

  @CMCLC-B08 @unit-level
  Scenario: A list of files is read the same way by every command
    Given arguments with a file, a comma-separated list with blanks, and a file named twice
    When they are read as files, directly and through a command that takes files
    Then each file comes once, in order, and the command is marked as taking files
