# language: en
# @anchors
#   code: VRFTV
#   ref: VPFVR
#   updated_at: 2026-10-03
#   layer: feature

@VPFVR
Feature: VerifyPhaseFacade — one invocation per phase, delegated to the check pipeline

  @VPFVR-B01 @unit-level
  Scenario: The staged scope is the added, copied, modified and renamed files of the index
    Given a repository where a.ts is modified and staged, new.ts is added, b.ts is deleted and untracked.ts is not added
    When the staged files are listed
    Then the list is exactly a.ts and new.ts

  @VPFVR-B02 @unit-level
  Scenario: Nothing staged has nothing to verify
    Given a repository with an empty index
    When verify runs over the staged files
    Then it prints "nothing staged — nothing to verify." and succeeds

  @VPFVR-B03 @unit-level
  Scenario: Verify hands the files to check in a child process
    Given a repository with a.ts staged
    When verify runs over the staged files in the pre-commit phase
    Then the child process receives "check --root <root> --changed a.ts --phase pre-commit --deterministic --only-issues"

  @VPFVR-B04 @unit-level
  Scenario: An automatic phase asks check for computable gates and issues only
    Given the phase pre-commit, pre-push or ci
    When the check invocation is built
    Then it carries the deterministic and only-issues flags

  @VPFVR-B05 @unit-level
  Scenario: The manual phase, or no phase, asks check for the full report
    Given the manual phase, or no phase
    When the check invocation is built
    Then it carries neither the deterministic nor the only-issues flag

  @VPFVR-B06 @unit-level
  Scenario: The pre-commit over the index dates the staged files first
    Given a repository with a.ts staged and dated 2026-09-01
    When verify runs over the staged files in the pre-commit phase
    Then it prints "· touch: dated 1 staged file(s) today (updated_at): a.ts"

  @VPFVR-B07 @unit-level
  Scenario: A project that turned pre-commit dating off gets no dating
    Given a project whose configuration declares touch pre_commit false, with b.ts staged
    When verify runs over the staged files in the pre-commit phase
    Then nothing is dated and b.ts keeps "updated_at: 2026-09-01"

  @VPFVR-B08 @unit-level
  Scenario: A staged file with changes outside the index is named, not dated
    Given a repository where a.ts is staged and changed again after staging
    When verify runs over the staged files in the pre-commit phase
    Then it prints "· touch: a.ts not dated — it has changes outside the index"

  @VPFVR-B09 @unit-level
  Scenario: The facade's flags reach check unchanged
    Given the pre-commit phase, a commit message file, the category types, skip-slow and no-record
    When the check invocation is built for a.md and b.md
    Then it carries each of those flags and one changed-file argument per file

  @VPFVR-B10 @unit-level
  Scenario: A project root below the repository top hands check its own staged files, by its own paths
    Given a repository whose project root is its sub directory, with sub/x.ts and a.ts staged
    When verify runs with the staged flag over the sub directory
    Then the child check gets "--changed x.ts" and nothing else

  @VPFVR-I01 @unit-level
  Scenario: Verify never asks check for both the full sweep and a file list
    Given a verify over a list of files, and another over the full sweep
    When the check invocations are built
    Then the first has no full-sweep flag and the second has no changed-file argument

  @VPFVR-X01 @unit-level
  Scenario: Verify holds no verdict of its own
    Given a repository with a.ts staged
    When verify runs
    Then the work is done by a child "check" invocation and verify returns what it returned

  @VPFVR-E01 @unit-level
  Scenario: The child's not-governed exit stays not-governed
    Given a child check that exits with code 3
    When verify runs over package.json
    Then verify returns the not-governed error

  @VPFVR-E02 @unit-level
  Scenario: Any other failing exit of the child stays a failure
    Given a child that exits with code 1
    When its result is translated
    Then it is an error, and not the not-governed one

  @VPFVR-E03 @unit-level
  Scenario: Verify with no scope is refused
    Given a repository
    When verify runs with neither staged, changed nor all
    Then it fails with "specify --staged, --changed <file> or --all"

  @VPFVR-E04 @unit-level
  Scenario: The staged scope outside a git repository is refused
    Given a directory outside any git repository
    When the staged files are listed
    Then it fails

  @VPFVR-E05 @unit-level
  Scenario: A dating failure warns and does not stop the verify
    Given dating the staged files failed with "boom"
    When the dating result is printed
    Then it prints "· touch: could not date the staged files (boom)"

  @VPFVR-B11 @unit-level
  Scenario: Over the index the check reads the index
    Given a verify over the staged files and one over named files
    When the check's arguments are built
    Then the first asks for --index and the second does not
