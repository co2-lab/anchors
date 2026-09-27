# language: en
# @anchors
#   ref: CLGCM
#   updated_at: 2026-09-27
#   layer: feature

@CLGCM
Feature: anchors changelog — the technical changelog, printed or written

  @CLGCM-B01 @unit-level
  Scenario: The latest release by default
    Given a repository with two releases
    When the changelog is printed with no range flag
    Then only the latest release is printed

  @CLGCM-B02 @unit-level
  Scenario: From a tag, or every release
    Given a repository with three releases
    When the changelog is printed from the first tag, and with all
    Then the releases after the tag, and every release, come newest first

  @CLGCM-B03 @unit-level
  Scenario: What is not released yet
    Given a repository with a commit after its last tag, and one with no tag
    When the changelog is printed with unreleased, and the untagged one without it
    Then the unreleased heading is on top in both

  @CLGCM-B04 @unit-level
  Scenario: Writing the incremental file
    Given a project in Portuguese with two releases
    When the changelog is written twice
    Then CHANGELOG.md holds the Portuguese title and the release, and the second run leaves it as it was

  @CLGCM-B05 @unit-level
  Scenario: One file per release
    Given a project in per_version mode with a release already written by hand
    When the changelog is written with every release and the unreleased one
    Then the missing release and the unreleased one get files and the hand-written one is kept

  @CLGCM-B06 @unit-level
  Scenario: The project's template and language
    Given a project in Spanish with its own template
    When the changelog is printed
    Then the project's template renders the release

  @CLGCM-B07 @unit-level
  Scenario: No anchors.yaml
    Given a repository with no anchors.yaml
    When the changelog is written
    Then CHANGELOG.md is written with the defaults

  @CLGCM-B08 @unit-level
  Scenario: Nothing to list
    Given a release with only chores
    When the changelog is printed
    Then nothing is printed and the error stream says why

  @CLGCM-E01 @unit-level
  Scenario: An unknown start tag
    Given a repository with one tag
    When the changelog is printed from a tag it does not have
    Then the command fails naming it

  @CLGCM-E02 @unit-level
  Scenario: A template that cannot be read
    Given a project whose template file does not exist
    When the changelog is printed
    Then the command fails naming changelog.template

  @CLGCM-E03 @unit-level
  Scenario: A broken anchors.yaml
    Given a project whose anchors.yaml does not load
    When the changelog is printed
    Then the command fails naming the file
