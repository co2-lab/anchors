# language: en
# @anchors
#   ref: DCLDF
#   updated_at: 2026-09-26
#   layer: feature

@DCLDF
Feature: DiffChangedLines — which lines of which files a change added, read from a unified diff

  @DCLDF-B01 @unit-level
  Scenario: Added lines are recorded under the file of the new-file header
    Given a diff of src/A.tsx adding lines 11 and 12 and replacing line 22
    When the diff is read
    Then src/A.tsx has lines 11, 12 and 22 and not line 21

  @DCLDF-B02 @unit-level
  Scenario: A hunk header sets the starting line of the new side
    Given the hunk headers "+10,3", "+7" and a header without a new side
    When their new start is read
    Then they give 10, 7 and 0

  @DCLDF-B03 @unit-level
  Scenario: Path prefixes and a tab-separated timestamp are stripped
    Given a diff whose new-file header is "b/x.go" followed by a tab and a timestamp
    When the diff is read
    Then the lines are recorded under "x.go"

  @DCLDF-B04 @unit-level
  Scenario: A deleted file records nothing
    Given a diff whose new side is /dev/null and whose hunk only removes lines
    When the diff file is read
    Then the result has no file

  @DCLDF-B05 @unit-level
  Scenario: Context lines advance the numbering without being recorded
    Given a hunk starting at line 1 with a context line, an added line, a context line and an added line
    When the diff is read
    Then only lines 2 and 4 are recorded

  @DCLDF-B06 @unit-level
  Scenario: Git compares the working copy with the current commit or with a reference
    Given a repository where pkg/a.go changed line 2, gained line 4, and gone.go was deleted
    When the diff is taken without a reference, and after committing, against HEAD~1
    Then both give pkg/a.go with lines 2 and 4, and a clean tree gives nothing

  @DCLDF-I01 @unit-level
  Scenario: Removals do not shift the new-side numbering
    Given a hunk "-20,1 +22,1" that removes one line and adds one
    When the diff is read
    Then the added line is line 22

  @DCLDF-X01 @unit-level
  Scenario: A diff file is read without git
    Given a diff file outside any repository adding two lines to x.go
    When the diff file is read
    Then x.go has lines 1 and 2

  @DCLDF-E01 @unit-level
  Scenario: A directory that is not a repository fails the git diff
    Given a directory that is not a git repository
    When the diff is taken from git
    Then an error is returned

  @DCLDF-E02 @unit-level
  Scenario: A missing diff file surfaces the read error
    Given a path where no diff file exists
    When the diff file is read
    Then the error says the file does not exist
