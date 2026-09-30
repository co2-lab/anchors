# language: en
# @anchors
#   ref: HDTHD
#   updated_at: 2026-09-29
#   layer: feature

@HDTHD
Feature: HeaderDateTouch — bumps the header date of the files that changed, and only of those

  @HDTHD-B01 @unit-level
  Scenario: Only a date inside the header at the top of the file is bumped
    Given files whose date sits in a line-comment header, in a closed HTML header, in a plain comment, in a string constant, and in a header below line ten
    When the touch decision is made for each
    Then only the line-comment header and the HTML header are bumped

  @HDTHD-B02 @unit-level
  Scenario: A real change and a new file are bumped to the date
    Given a file whose body changed since the last commit and a file with no previous version
    When the touch decision is made for the date 2026-09-25
    Then both are bumped and carry "updated_at: 2026-09-25"

  @HDTHD-B03 @unit-level
  Scenario: An unchanged file, a date-only change and a file already at the date are not bumped
    Given a file equal to its last commit, a file whose only change is its date, and a file already dated 2026-09-25
    When the touch decision is made for the date 2026-09-25
    Then none is bumped, each with its own reason

  @HDTHD-B04 @unit-level
  Scenario: Without the staged flag the candidates are the worktree changes and the untracked files
    Given a repository where a.ts changed and new.ts is untracked
    When touch runs for 2026-09-25
    Then it prints "bumped a.ts  (2026-09-01 → 2026-09-25)" and "bumped new.ts"
    And the unchanged d.ts keeps 2026-09-01

  @HDTHD-B05 @unit-level
  Scenario: With the staged flag the index is dated and re-staged
    Given a repository where a change to a.ts is staged
    When touch runs on the staged files for 2026-09-25
    Then the staged a.ts carries "updated_at: 2026-09-25"

  @HDTHD-B06 @unit-level
  Scenario: A staged file with changes outside the index is skipped
    Given a repository where d.ts is staged and changed again after staging
    When touch runs on the staged files
    Then d.ts is skipped for its changes outside the index
    And its later change does not reach the index

  @HDTHD-B07 @unit-level
  Scenario: In a partial commit the real index is dated too
    Given a partial commit's temporary index and real index, both holding the changed a.ts
    When touch runs on the staged files
    Then both indexes hold a.ts dated 2026-09-25

  @HDTHD-B08 @unit-level
  Scenario: The exclude globs of the flag and of the configuration add up
    Given a configuration excluding gen/** and a flag excluding d.ts, with a.ts, gen/c.ts and d.ts changed
    When touch runs
    Then a.ts is bumped and gen/c.ts and d.ts are skipped as excluded

  @HDTHD-B09 @unit-level
  Scenario: The dry run says what it would bump and writes nothing
    Given a repository where a.ts changed
    When touch runs in dry run
    Then it prints "would bump a.ts"
    And a.ts still carries 2026-09-01

  @HDTHD-B10 @unit-level
  Scenario: Each bump and each skip is listed with its reason, then the total
    Given a repository with a changed file, a date-only change and an excluded file
    When touch runs
    Then the bump shows its old and new date, and each skip names its reason

  @HDTHD-B11 @unit-level
  Scenario: Without a date the day of the run is written
    Given a repository where a.ts changed
    When touch runs in dry run without a date
    Then the date it would write is today's

  @HDTHD-B12 @unit-level
  Scenario: The pre-commit bump is on unless the project turns it off
    Given no configuration, a configuration with no touch block, one declaring pre_commit on, and one declaring it off
    When the pre-commit bump switch is read
    Then it is on in every case but the last

  @HDTHD-B13 @unit-level
  Scenario: A project root below the repository top dates its own files and nothing outside it
    Given a repository whose project root is its sub directory, with sub/x.ts staged clean, sub/y.ts staged with a change on top and a.ts staged at the top
    When touch runs with the staged flag over the sub directory, then in the worktree mode with sub/new.ts untracked
    Then x.ts is bumped and re-staged, y.ts is skipped, new.ts would be bumped, and a.ts is never named

  @HDTHD-I01 @unit-level
  Scenario: Touching twice bumps nothing the second time
    Given a repository where a.ts changed and was just touched for 2026-09-25
    When touch runs again for 2026-09-25
    Then a.ts is skipped as already dated and the total is "bumped 0 file(s)"

  @HDTHD-X01 @unit-level
  Scenario: A changed file with no dated header is neither touched nor listed
    Given a repository with a new file plain.ts that has no header
    When touch runs
    Then plain.ts is not in the output and its content is unchanged

  @HDTHD-E01 @unit-level
  Scenario: Outside a git repository touch fails
    Given a directory outside any git repository
    When touch runs
    Then it fails naming "git diff HEAD"

  @HDTHD-E02 @unit-level
  Scenario: A changed file that cannot be read is skipped and named
    Given a repository where a.ts and b.ts changed and b.ts cannot be read
    When touch runs
    Then b.ts is skipped as unreadable and a.ts is bumped

  @HDTHD-B14 @unit-level
  Scenario: Named files narrow the touch to them
    Given three changed files, one of them in a folder
    When the touch names one file and the folder, then a path outside the project, then nothing
    Then only the named file and the folder's file are dated, the outside path is an error, and with nothing named all three are
