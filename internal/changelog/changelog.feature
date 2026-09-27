# language: en
# @anchors
#   ref: CHNGL
#   updated_at: 2026-09-27
#   layer: feature

@CHNGL
Feature: Changelog — the releases of a project, read from its commits

  @CHNGL-B01 @unit-level
  Scenario: A breaking change is never silent
    Given a refactor with a bang and a feat with a BREAKING CHANGE footer
    When the commits are classified
    Then both are breaking changes and the footer's text follows the subject

  @CHNGL-B02 @unit-level
  Scenario: A feat is a feature
    Given a feat with a scope
    When the commits are classified
    Then it is a feature with its scope and hash

  @CHNGL-B03 @unit-level
  Scenario: A fix with a Bug footer is a bug fixed, without it a fix
    Given a fix with a Bug footer and a fix without one
    When the commits are classified
    Then the first is a bug fixed carrying where it was seen and the second a fix

  @CHNGL-B04 @unit-level
  Scenario: The internal types and free-form subjects are left out
    Given a chore, a refactor and a subject outside Conventional Commits
    When the commits are classified
    Then the release is empty

  @CHNGL-B05 @unit-level
  Scenario: Only the last paragraph holds footers
    Given a fix whose Bug line is in the middle of the body
    When the commits are classified
    Then it is a fix, not a bug fixed

  @CHNGL-B06 @unit-level
  Scenario: Each release holds the commits after the previous tag
    Given a repository with two tags and commits before each
    When the releases of both tags are read
    Then the newest comes first and each holds only its own commits

  @CHNGL-B07 @unit-level
  Scenario: What came after the last tag is unreleased
    Given a repository with a commit after its last tag
    When the releases are read with unreleased
    Then a release with no version holding that commit is on top

  @CHNGL-B08 @unit-level
  Scenario: The built-in template translates its headings
    Given a release with a feature and a bug fixed, and a translating function
    When it is rendered with the built-in template
    Then the version, the translated sections and where the bug was seen are written, and the empty sections are not

  @CHNGL-B09 @unit-level
  Scenario: A written release is marked by its version
    Given releases written into an empty file
    When the versions it holds are listed
    Then each written version is listed

  @CHNGL-B10 @unit-level
  Scenario: New releases go on top and the file keeps what it holds
    Given a file with a title, a hand-edited release and old hand-written text
    When a newer release and the one it holds are written
    Then only the newer one is added, below the title, and the rest is as it was

  @CHNGL-B11 @unit-level
  Scenario: The unreleased block is replaced
    Given a file with an unreleased block above a release
    When a new unreleased block is written
    Then the old one is gone and the release below is kept

  @CHNGL-B12 @unit-level
  Scenario: An empty file gets a title
    Given no file
    When a release is written
    Then the file starts with the title

  @CHNGL-E01 @unit-level
  Scenario: A template that does not parse is an error
    Given a template with an unclosed action
    When a release is rendered
    Then the error says it is the changelog template

  @CHNGL-E02 @unit-level
  Scenario: A git failure names the command
    Given a directory that is not a repository
    When its tags are read
    Then the error names the git command
