# @anchors
#   code: DSFDC
#   ref: DCSYN
#   updated_at: 2026-10-03
Feature: DocsSyncForCommit — the commit carries the pages its specs produce

  @DCSYN-B01 @unit-level
  Scenario: A page left out of date is compiled and staged
    Given a committed page built from a spec, and the spec changed and staged
    When the pre-commit syncs the docs over a tree that matches the index
    Then the page is compiled from the staged spec, written and staged, and an up-to-date page is not touched

  @DCSYN-B02 @unit-level
  Scenario: With the tree ahead, the page goes straight into the index
    Given the spec staged with one text and edited again in the tree, and another governed file edited and not staged
    When the pre-commit syncs the docs
    Then the index holds the page of the staged spec, and the page on disk is as it was

  @DCSYN-B03 @unit-level
  Scenario: Nothing to compile, nothing staged
    Given a changed spec in a project that turned it off, one without templates, and one without a map
    When the pre-commit syncs the docs
    Then nothing is compiled or staged in any of them

  @DCSYN-X01 @unit-level
  Scenario: A page written by hand, or one git ignores, is never touched
    Given a page with no generator marker, and a page git ignores, both from templates the changed spec feeds
    When the pre-commit syncs the docs
    Then neither is written nor staged

  @DCSYN-X02 @unit-level
  Scenario: Nothing outside docs is written or staged
    Given a changed spec and its page
    When the pre-commit syncs the docs
    Then the only paths staged or changed on disk by the sync are under docs

  @DCSYN-I01 @unit-level
  Scenario: After the sync the gate finds the pages up to date
    Given a changed spec staged, with the tree ahead of the index
    When the pre-commit syncs the docs
    Then the pages confronted with the index are up to date

  @DCSYN-E01 @unit-level
  Scenario: A template that does not compile comes back as an error
    Given a changed spec staged and a template that does not compile
    When the pre-commit syncs the docs
    Then the error comes back naming the template, and no page is staged
