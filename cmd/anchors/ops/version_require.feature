# language: en
# @anchors
#   code: VRQRB
#   ref: VRQRA
#   updated_at: 2026-11-18
#   layer: feature

@VRQRA
Feature: VersionRequirement — what a stage asks of the version its sha carries, and of the fixes it contains

  @VRQRA-B01 @unit-level
  Scenario: fresh holds for a version above every release — the tags and what released: lists — and refuses one released or below
    Given app 1.0.0 tagged app-v1.0.0 and listed for production
    When fresh is asked of 1.0.0 kept, 1.1.0 listed for production with no tag, 1.1.0 listed for review, and 0.9.0
    Then 1.0.0 is released by its tag, 1.1.0 is released by the listing, 1.1.0 in review is above every release, and 0.9.0 is below 1.0.0

  @VRQRA-B02 @unit-level
  Scenario: listed holds when the file at the sha has an entry for the version matching where, and names the listing when it has not
    Given a stage asking listed where target is review, and another asking a file of its own where env is qa
    When 1.0.1 is not listed, then listed for review, and envs.yaml lists it for qa
    Then the first fails naming the listing, and the others hold naming theirs

  @VRQRA-B03 @unit-level
  Scenario: same_as_passed holds when the sha carries the version the commit that passed the stage carried
    Given the head of #9 passed review-release carrying 1.0.1
    When same_as_passed is asked of a commit at 1.0.1, of one at 1.0.2, and with no commit that passed
    Then the first holds, the second fails naming 1.0.2 and 1.0.1, and the third fails saying nothing passed

  @VRQRA-B04 @unit-level
  Scenario: expected holds when the version is the one the branch's strategy computes; manual asks it be fresh, none refuses
    Given review on semantic with a feat since app-v1.0.0, feat/* on none and maint/* on manual
    When expected is asked of 1.0.1 and 1.1.0 on review, of 1.1.0 on feat/x and maint/1, and of 1.0.0 on maint/1
    Then 1.0.1 fails as semantic asks for 1.1.0, 1.1.0 holds, feat/x never versions, maint/1 holds when fresh and fails at 1.0.0

  @VRQRA-B05 @unit-level
  Scenario: A branch's range bounds the version whatever is asked: from_base keeps the release the branch was cut from
    Given hotfix/login cut from production at app-v1.0.0, on increment patch with range from_base
    When expected is asked of 1.1.0 and of 1.0.1
    Then 1.1.0 is outside 1.0.x naming hotfix/*, and 1.0.1 is inside and as increment asks

  @VRQRA-B06 @unit-level
  Scenario: contains holds when the head of the branch on the remote is an ancestor of the sha, and names the commits missing when it is not
    Given a fix pushed to production on origin, and review going on without it
    When contains production is asked of review, then after review merges production
    Then it fails naming one commit and its subject, then holds naming the head of origin/production

  @VRQRA-B07 @unit-level
  Scenario: The version is read at the sha, not in the working tree
    Given a commit at 1.1.0 and a working tree back at 1.0.0
    When fresh is asked of the commit
    Then it holds for 1.1.0

  @VRQRA-E01 @unit-level
  Scenario: A branch the remote does not have, or a version the sha cannot give, fails the requirement naming it
    Given no branch nope on origin, and a commit whose app/version.json has no version
    When contains nope and fresh are asked
    Then the first fails naming origin/nope, and the second naming app/version.json
