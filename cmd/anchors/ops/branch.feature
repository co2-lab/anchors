# language: en
# @anchors
#   code: BCFBA
#   ref: BCSBR
#   updated_at: 2026-11-06
#   layer: feature

@BCSBR
Feature: Branch — confront merges and commits with the project's branch flow

  @BCSBR-B01 @unit-level
  Scenario: A pull request is allowed when the flow allows its head into its base, and refused naming what the base takes
    Given a flow where production takes review and hotfix/*
    When a pull request of develop into review, and one of develop into production, are checked
    Then the first passes, and the second fails naming review and hotfix/*

  @BCSBR-B02 @unit-level
  Scenario: A push to a guarded branch is allowed only as the merge of a pull request the flow allows
    Given a push to production whose sha is the merge of a pull request from review, then from develop, then of none
    When the push is checked
    Then the first passes, the second fails naming the pull request, the third fails as a direct push, and a push it cannot confront fails

  @BCSBR-B03 @unit-level
  Scenario: A commit on a protected branch is refused; another branch, or a detached HEAD, commits
    Given a flow protecting production
    When a commit is guarded on production, develop, hotfix/x and a detached HEAD
    Then production is refused, and the others commit

  @BCSBR-E01 @unit-level
  Scenario: A check with a side missing fails naming what is missing
    Given no head, base, branch or sha by flag or environment
    When branch check runs on a pull request and on a push
    Then each fails naming the flags and the variables

  @BCSBR-E02 @unit-level
  Scenario: A push that cannot be confronted fails
    Given GitHub that cannot be asked
    When a push to a guarded branch is checked
    Then the check fails

  @BCSBR-E03 @unit-level
  Scenario: A configuration that does not load fails the check and the guard
    Given an anchors.yaml whose branch_flow has a bad pattern
    When branch check and branch guard run
    Then both fail with the load error
