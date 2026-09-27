# language: en
# @anchors
#   ref: FXIXX
#   updated_at: 2026-09-26
#   layer: feature

@FXIXX
Feature: Fix — the self-healer that applies the mechanical, safe repairs of `check --fix`

  @FXIXX-B01 @unit-level
  Scenario: Only a check with a registered fixer is fixable
    Given the check "updated-at-atual", which has a fixer, and "header-conforme", which has none
    When the fixer registry is asked which of them is fixable
    Then "updated-at-atual" is fixable
    And "header-conforme" is not

  @FXIXX-B02 @unit-level
  Scenario: A stale date on a committed file is rewritten to its last commit date
    Given a spec committed on 2024-03-05 whose header says "updated_at: 2020-01-01"
    And no pending edit to that spec
    When the fixable gates are applied to the spec
    Then the header says "updated_at: 2024-03-05"
    And one repair is reported as fixed, naming the gate and the spec

  @FXIXX-B03 @unit-level
  Scenario: A file with an uncommitted edit takes today's date
    Given a committed spec whose header says "updated_at: 2024-03-05"
    And a new line added to the spec and not committed
    When the updated_at fixer runs on the spec
    Then the header says today's date

  @FXIXX-B04 @unit-level
  Scenario: A date that already matches is left alone and reported as nothing
    Given a spec committed on 2024-03-05 whose header says "updated_at: 2024-03-05"
    When the fixable gates are applied to the spec
    Then no repair is reported

  @FXIXX-B05 @unit-level
  Scenario: Only fixable gates, the nodes they apply to and files on disk are touched
    Given a gate with no fixer, a node of a kind the fixable gate does not apply to, and a spec node whose file is not on disk
    When the fixable gates are applied
    Then no repair is reported
    And the out-of-scope file keeps its stale date

  @FXIXX-B06 @unit-level
  Scenario: Outside a git repository the file is left untouched
    Given a spec with a stale date in a directory that is not a git repository
    When the updated_at fixer runs on the spec
    Then the content is returned unchanged and nothing is marked as changed

  @FXIXX-B07 @unit-level
  Scenario: A file never committed and not edited is left untouched
    Given a spec that was never committed and that git ignores
    When the updated_at fixer runs on the spec
    Then the content is returned unchanged and nothing is marked as changed

  @FXIXX-I01 @unit-level
  Scenario: Only the date changes, the rest of the file is kept byte for byte
    Given a committed spec with a stale date and a body line "body"
    When the fixable gates are applied to the spec
    Then the file equals the original with only the date replaced

  @FXIXX-X01 @unit-level
  Scenario: A header without the date field is never given one
    Given a committed spec whose header has a code and no updated_at field
    When the updated_at fixer runs on the spec
    Then the content is returned unchanged and nothing is marked as changed

  @FXIXX-E01 @unit-level
  Scenario: A write that fails is reported as not fixed, with the cause
    Given a committed spec with a stale date whose file is read-only
    When the fixable gates are applied to the spec
    Then one repair is reported as not fixed
    And the report carries the cause of the failure
