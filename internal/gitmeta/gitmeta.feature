# language: en
# @anchors
#   ref: GTMTG
#   updated_at: 2026-09-26
#   layer: feature

@GTMTG
Feature: GitMeta — what git knows about the files: last commit dates, pending changes, HEAD, dirty count

  @GTMTG-B01 @unit-level
  Scenario: Today is the system date
    Given the system clock
    When today is asked
    Then it is the system date formatted year-month-day

  @GTMTG-B02 @unit-level
  Scenario: The last commit date is the day of the most recent commit of the file
    Given a repository where "a.txt" was committed on 2026-01-10 and again on 2026-03-05
    When the last commit date of "a.txt" and of the never-committed "nao-existe.txt" is asked
    Then "a.txt" gives "2026-03-05" and "nao-existe.txt" gives no date

  @GTMTG-B03 @unit-level
  Scenario: The bulk reader gives each file its most recent commit day
    Given a repository where "a.txt" and "pkg/b.txt" were committed on 2026-01-10 and "a.txt" again on 2026-03-05
    When every commit date is read at once
    Then "a.txt" is "2026-03-05", "pkg/b.txt" is "2026-01-10", and a folder outside any repository gives an empty map

  @GTMTG-B04 @unit-level
  Scenario: A new file in a repository has pending changes, and the question was asked
    Given a new repository holding the uncommitted file "a.go"
    When its pending changes are asked
    Then the answer is changed and known

  @GTMTG-B05 @unit-level
  Scenario: The shortcut says yes after an edit and no after a commit
    Given a repository where "a.txt" was just committed
    When the shortcut is asked before and after editing "a.txt"
    Then it answers no and then yes

  @GTMTG-B06 @unit-level
  Scenario: HEAD is the short hash and subject, and unknown without a commit
    Given a repository with no commit, a folder outside any repository, and a repository whose last commit is "feat: the subject line"
    When HEAD is asked of each
    Then the first two are not known and the third gives its short hash and "feat: the subject line"

  @GTMTG-B07 @unit-level
  Scenario: The dirty count counts the modified files of a real repository
    Given a new empty repository
    When the dirty count is asked before and after writing "novo.txt"
    Then it is 0 and then 1

  @GTMTG-X01 @unit-level
  Scenario: A tree that could not be counted is not reported clean
    Given a folder outside any repository
    When the dirty count is asked
    Then it is negative, never 0

  @GTMTG-X02 @unit-level
  Scenario: Outside a repository the pending-change question is not known
    Given a folder outside any repository
    When the pending changes of "qualquer.go" are asked
    Then the answer is not known and not changed
