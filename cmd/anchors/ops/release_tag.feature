# language: en
# @anchors
#   code: RLSTB
#   ref: RLSTA
#   updated_at: 2026-11-18
#   layer: feature

@RLSTA
Feature: ReleaseTag — the tag and the GitHub Release a stage that releases writes

  @RLSTA-B01 @unit-level
  Scenario: A stage that releases writes, when it passes, its product's tag on the sha and the tag's Release with the changelog since the release before
    Given production released app 1.0.0 as app-v1.0.0, then a fix with a Bug footer and the version 1.0.1
    When production-release, which releases app, is recorded on that commit
    Then the tag app-v1.0.1 is written on the commit, then its Release with the fix and where it was seen, and the stage is recorded as passed

  @RLSTA-B02 @unit-level
  Scenario: A tag already on the sha is a rerun: only a missing Release is written
    Given the tag app-v1.0.1 already on the commit, with its Release, and then without it
    When the stage is recorded again
    Then nothing is written the first time, and only the Release the second

  @RLSTA-E01 @unit-level
  Scenario: A tag on another sha, a version the sha cannot give, or GitHub not answering fails the record, recorded as failure
    Given the tag app-v1.0.1 on another commit, then GitHub answering 502, then a commit whose app/version.json holds no version
    When the stage is recorded
    Then each fails naming the tag and the commit, the tag, or the file, and the stage is recorded as failed

  @RLSTA-E02 @unit-level
  Scenario: A tag or a Release GitHub does not write, a ref it answers unread, a release it does not say, or no gh, fail the record naming the step
    Given GitHub failing to write the tag, to read the release, to write the release, and answering a ref that is not JSON, and then no gh in the PATH
    When the stage is recorded
    Then each fails naming the step and the tag, and with no gh the call says gh is not in the PATH

  @RLSTA-E03 @unit-level
  Scenario: A tag of the product the scheme cannot read stops the release body
    Given the tag app-vnext beside the releases of app
    When the stage is recorded
    Then it fails naming app-vnext
