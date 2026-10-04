# language: en
# @anchors
#   code: PLFTP
#   ref: RCPLR
#   updated_at: 2026-10-03
#   layer: feature

@RCPLR
Feature: RecodePlan — planning and applying the rename of a code across the whole project

  @RCPLR-B01 @unit-level
  Scenario: A malformed source or target code is refused
    Given a project holding the code "ABCDX"
    When the renames "abcdx" to "WXYZX", "ABCDX" to "WX" and "ABC-X" to "WXYZX" are planned
    Then each is refused as invalid

  @RCPLR-B02 @unit-level
  Scenario: Renaming a code to itself is refused
    Given a project holding the code "ABCDX"
    When "ABCDX" to "ABCDX" is planned
    Then it is refused as the same code

  @RCPLR-B03 @unit-level
  Scenario: The plan holds every file with the code, sorted, with its new content
    Given a project where "z/Foo.spec.md" and "a/Foo.tsx" hold "ABCDX" and "m/Other.tsx" does not
    When "ABCDX" to "WXYZX" is planned
    Then the plan lists "a/Foo.tsx" then "z/Foo.spec.md", both rewritten, with 3 content replacements

  @RCPLR-B04 @unit-level
  Scenario: The declared testID prefix is rewritten and counted apart
    Given a project declaring the lower testID convention, with a component holding testID "abcdx-root"
    When "ABCDX" to "WXYZX" is planned, with and without the recode block
    Then with it 1 testID is rewritten, and without it none

  @RCPLR-B05 @unit-level
  Scenario: Files named after the code are planned for renaming
    Given a project declaring the file pattern "**/{{code}}-*.yaml" and the file "flows/ABCDX-S01.yaml"
    When "ABCDX" to "WXYZX" is planned, with and without the recode block
    Then with it "flows/ABCDX-S01.yaml" is renamed to "flows/WXYZX-S01.yaml", and without it nothing is renamed

  @RCPLR-B06 @unit-level
  Scenario: A divergent testID prefix is warned about and never rewritten
    Given a project declaring the lower testID convention, whose component for "ABCDX" holds testID "zzzz-root"
    When "ABCDX" to "WXYZX" is planned
    Then no testID is rewritten, the plan carries the divergent-prefix warning, and "zzzz-root" is kept

  @RCPLR-B07 @unit-level
  Scenario: A code that appears nowhere is refused
    Given a project that does not hold the code "QQQQX"
    When "QQQQX" to "WXYZX" is planned
    Then it is refused because the code does not appear in any file

  @RCPLR-B08 @unit-level
  Scenario: Applying writes each file with its mode and performs the renames
    Given a project outside any repository with an executable "run.tsx" and the flow "flows/ABCDX-S01.yaml"
    When the plan of "ABCDX" to "WXYZX" is applied
    Then 3 changes are reported, "run.tsx" is rewritten and still executable, and the flow is renamed

  @RCPLR-B09 @unit-level
  Scenario: Moves go through git inside a repository and are plain renames outside
    Given a tracked file "a.go" in a repository, and a file "a.go" outside any repository
    When each is moved to "sub/b.go"
    Then the repository's index records a rename, and outside the file simply reaches "sub/b.go"

  @RCPLR-B10 @unit-level
  Scenario: An untracked file inside a repository is moved by a plain rename
    Given a repository holding the untracked file "loose.go"
    When it is moved to "moved.go"
    Then the move succeeds, "loose.go" is gone and "moved.go" exists

  @RCPLR-B11 @unit-level
  Scenario: A target code another unit already owns is refused
    Given a project where "a/Foo.spec.md" owns "ABCDX" and "b/Bar.spec.md" owns "WXYZX", by its header or by its scenario code "WXYZX-B01"
    When the rename of "ABCDX" to "WXYZX" is planned
    Then it is refused naming "b/Bar.spec.md"
    And the rename to "WXYZX" is planned when the other unit owns only "WXYZY"

  @RCPLR-B12 @unit-level
  Scenario: The malformed-code refusal names the lengths the project declares
    Given a project whose code lengths are 5, and then 4 and 6
    When a rename with a code of the wrong length is planned
    Then the refusal says "expected 5 characters", and then "expected 4/6 characters"

  @RCPLR-X01 @unit-level
  Scenario: A git refusal is surfaced and never bypassed
    Given a repository where "a.go" and "b.go" are both tracked
    When "a.go" is moved onto "b.go"
    Then the error says git refused, "a.go" is still there and "b.go" is unchanged

  @RCPLR-X02 @unit-level
  Scenario: Planning writes nothing
    Given a project where two files hold "ABCDX"
    When "ABCDX" to "WXYZX" is planned
    Then every file of the project still has its original content

  @RCPLR-E05 @unit-level
  Scenario: A file that cannot be written stops the apply, naming it
    Given a plan whose only file is the path of a folder
    When it is applied
    Then applying stops with an error naming the folder and 0 files written
