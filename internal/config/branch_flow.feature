# language: en
# @anchors
#   code: BFFBR
#   ref: BFSBR
#   updated_at: 2026-11-06
#   layer: feature

@BFSBR
Feature: BranchFlow — which branch merges into which, and which branches take no commit

  @BFSBR-B01 @unit-level
  Scenario: A merge is allowed when the base's entry names its head
    Given a flow where review takes develop and hotfix/*, and production takes review and hotfix/*
    When develop is merged into review, and develop into production
    Then the first is allowed, and the second is refused naming review and hotfix/*

  @BFSBR-B02 @unit-level
  Scenario: A base no entry names takes any head, unless it is protected
    Given a flow with no entry for feat/y, and production protected with no entry
    When feat/x is merged into feat/y, and review into production
    Then the first is allowed, and the second refused

  @BFSBR-B03 @unit-level
  Scenario: A branch pattern is a glob whose * does not cross /, or a regular expression after re:
    Given the patterns hotfix/* and re:release/v[0-9]+
    When they are matched against hotfix/login, hotfix/a/b, release/v12 and xrelease/v1
    Then hotfix/login and release/v12 match, and the others do not

  @BFSBR-E01 @unit-level
  Scenario: A branch pattern that is neither a glob nor a regular expression refuses the configuration
    Given a branch_flow whose pattern is re:( or [
    When the configuration is loaded
    Then it does not load, naming branch_flow
