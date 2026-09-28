# language: en
# @anchors
#   ref: MPSYN
#   updated_at: 2026-09-28
#   layer: feature

@MPSYN
Feature: MapSyncForCommit — the commit carries the map a build of the commit makes

  @MPSYN-B01 @unit-level
  Scenario: The committed map is the one a build of the commit makes
    Given a committed map, a spec edited and staged, another file edited but not staged, and an untracked file
    When the hook dates the staged files, syncs the map, and the commit is made
    Then a full build of the committed files makes exactly the committed map

  @MPSYN-B02 @unit-level
  Scenario: A dated file keeps its proofs
    Given the spec had a proof at its revision before the hook dated it
    When the map is synced
    Then the spec's proof is at its new revision

  @MPSYN-B03 @unit-level
  Scenario: The index, not the tree
    Given an unstaged edit and an untracked file
    When the map is synced
    Then the edited file has its staged revision and the untracked one has no node

  @MPSYN-B04 @unit-level
  Scenario: An untracked map is left alone
    Given a project whose map git does not track
    When the map is synced
    Then nothing is written nor staged

  @MPSYN-E01 @unit-level
  Scenario: A map that cannot be written gives the error back
    Given a tracked map whose path has become a directory
    When the map is synced
    Then the error comes back and nothing is staged
