# language: en
# @anchors
#   code: PRVLF
#   ref: PRVLS
#   updated_at: 2026-11-11
#   layer: feature

@PRVLS
Feature: PushReview — the reviews a push carries, as the repository keeps them

  @PRVLS-B01 @unit-level
  Scenario: A review is written as an entry
    Given an ok review of two commits with judgment answers and a note, and a failed one
    When each is written as an entry
    Then the entry has the heading, the machine line, the commits with patch ids, the judgments, the validated or failed text and the notes

  @PRVLS-B02 @unit-level
  Scenario: The file is read back into its entries
    Given a file with two reviews and a section with no machine line
    When it is read and written again
    Then it has the two reviews in order with their verdict, author and commits, and the written file reads back the same

  @PRVLS-B03 @unit-level
  Scenario: Writing a review drops the entries older than the retention
    Given entries of today, of forty days ago and one with no readable date
    When they are kept for thirty days, and with no retention
    Then the old one goes and the others stay, and with no retention all stay

  @PRVLS-B04 @unit-level
  Scenario: A commit is covered by the newest entry that names it
    Given a failed review and a newer ok review of the same commit, and a commit rebased with another hash
    When the coverage of each commit is asked
    Then the newest review covers by hash, a short hash covers, the patch id covers the rebased one, and an unnamed commit is not covered

  @PRVLS-B05 @unit-level
  Scenario: A commit that changes only the evidence owes no review
    Given commits changing only the reviews' file and the map, a code file too, and nothing
    When their exemption is asked
    Then only the first owes no review

  @PRVLS-B06 @unit-level
  Scenario: The commits of a range are read with their patch id and files
    Given a repository with three commits, the last reworded by an amend
    When the range after the first is read, and the files of the last
    Then it has the two commits oldest first with subjects, the amended one keeps its patch id, and its files are the paths it changed

  @PRVLS-E01 @unit-level
  Scenario: A range git cannot list fails naming it
    Given a repository
    When a range from a revision it does not know is read
    Then reading fails naming the range

  @PRVLS-E02 @unit-level
  Scenario: A commit git gives no files or patch id is covered by hash and owes a review
    Given a repository with an empty commit, and a hash it does not know
    When their files and patch ids are read
    Then they have no files and no patch id, and with no files they owe a review

