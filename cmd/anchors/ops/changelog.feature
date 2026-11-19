# language: en
# @anchors
#   code: CHFTC
#   ref: CLGCM
#   updated_at: 2026-11-18
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
    Then the project's template renders the release, read from the project or from a URL

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
    Given a project whose template file does not exist, or whose template URL answers nothing
    When the changelog is printed
    Then the command fails naming changelog.template

  @CLGCM-E03 @unit-level
  Scenario: A broken anchors.yaml
    Given a project whose anchors.yaml does not load
    When the changelog is printed
    Then the command fails naming the file

  @CLGCM-B09 @unit-level
  Scenario: With versioning, the releases lie between the product's tags, ordered by its scheme and headed by the version
    Given the product app tagged app-v1.9.0 and app-v1.10.0, with api-v5.0.0 between them
    When every release is printed on production, and from 1.9.0
    Then 1.10.0 comes before 1.9.0, no other product's tag shows, and from 1.9.0 lists 1.10.0 alone

  @CLGCM-B10 @unit-level
  Scenario: A branch that releases lists its releases, another the candidate under the version its strategy computes
    Given app released as app-v1.0.0 on production, review writing the changelog, and a feat since
    When the changelog is printed on review, and with no branch
    Then review lists 1.1.0 as unreleased alone, and with no branch production's 1.0.0 is listed

  @CLGCM-B11 @unit-level
  Scenario: --write writes every output of the product, and only on a branch the changelog is written on
    Given app with the outputs CHANGELOG.md and per_version notes with a template of its own
    When the changelog is written on production and on develop
    Then CHANGELOG.md holds 1.0.0 with its marker, notes/1.0.0.md is the template's, and develop is refused

  @CLGCM-E04 @unit-level
  Scenario: --from that is no tag of the product, or a template of an output that cannot be read, fails naming it
    Given app tagged app-v1.0.0, a tag v9.0.0, and an output whose template is missing
    When the changelog is printed from v9.0.0, and written
    Then the first fails saying no tag of app, the second naming the output's template field
