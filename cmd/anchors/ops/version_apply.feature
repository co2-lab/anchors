# language: en
# @anchors
#   code: VAPLB
#   ref: VAPLA
#   updated_at: 2026-11-19
#   layer: feature

@VAPLA
Feature: VersionApply — the version the branch asks for, brought into the repository

  @VAPLA-B01 @unit-level
  Scenario: release_pr pushes anchors/release/<product> with the bump, the release's blocks and the intents consumed, and opens its pull request into the branch
    Given review on semantic with a feat and its intent since app-v1.0.0, CHANGELOG.md from the commits and NOTES.md from the intents, archived on release
    When anchors version apply runs on review
    Then anchors/release/app on origin holds 1.1.0, both blocks and the intent archived under released/1.1.0, in chore(release): app 1.1.0, its pull request into review is opened, and the checkout is back on review

  @VAPLA-B02 @unit-level
  Scenario: A rerun updates the open pull request, and once the release is on the branch there is nothing to propose
    Given the release pull request #12 open, then a fix on review, then the release merged into review
    When apply runs again each time
    Then #12 is updated, and after the merge nothing is pushed or opened

  @VAPLA-B03 @unit-level
  Scenario: Nothing to release opens nothing
    Given review with no commit since app-v1.0.0
    When apply runs
    Then it says nothing to release since app-v1.0.0

  @VAPLA-B04 @unit-level
  Scenario: commit confirms the sources hold the version the branch asks for, and refuses otherwise
    Given apply commit, the sources at 1.0.0 while review asks for 1.1.0, then at 1.1.0
    When apply runs
    Then it fails naming both and bump, then confirms 1.1.0

  @VAPLA-B05 @unit-level
  Scenario: build writes the version into the workspace's sources, commits nothing, and prints version and tag, into $GITHUB_OUTPUT too
    Given apply build and $GITHUB_OUTPUT set
    When apply runs on review
    Then the sources hold 1.1.0, HEAD did not move, and version=1.1.0 and tag=app-v1.1.0 are printed and appended

  @VAPLA-B06 @unit-level
  Scenario: A per_version output gets the release's file, named by format.file, in place of the unreleased one, and intents consumed by delete are gone
    Given a per_version output notes with the file v{version}.md and an unreleased.md, and intents consumed by delete
    When apply runs
    Then the release branch holds notes/v1.1.0.md, no unreleased.md and no intent

  @VAPLA-E01 @unit-level
  Scenario: release_pr of a product with no source, a branch that never versions, or no gh fails naming why
    Given feat/x on none, then a product with no source, then gh failing, then no gh in the PATH
    When apply runs
    Then each fails naming why, back on review
