# language: en
# @anchors
#   code: FXFTF
#   ref: FXIXX
#   updated_at: 2026-10-03
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

  @FXIXX-B08 @unit-level
  Scenario: The repair detail is written in the project's language
    Given a committed spec whose updated_at is stale
    When the fixer repairs it in an English and in a Portuguese project
    Then the detail of the repair is the English text in the first and the Portuguese text in the second

  @FXIXX-B09 @unit-level
  Scenario: The fix writes the missing header from the map
    Given a code file its spec specifies, a test its feature tests, a guide, a script with a shebang, a header with no identity, and a code file of no unit
    When the header fix runs
    Then the code and the test get the ref of their unit, the guide its layer, the script's header goes below the shebang, the identity-less header gets its line
    And the code file of no unit is left as it is

  @FXIXX-B10 @unit-level
  Scenario: The fix gives a file with an identity and no code of its own a code unique in the map
    Given a code file whose header refs its unit and has no code of its own
    When the header fix runs, then again
    Then the header gets a five-character code that is not the unit's, beside the ref, and the second run changes nothing
