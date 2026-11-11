# language: en
# @anchors
#   code: PRCMF
#   ref: PRCMS
#   updated_at: 2026-11-11
#   layer: feature

@PRCMS
Feature: PushReview — the review of what a push sends, recorded by whoever pushes and confronted by the pre-push hook

  @PRCMS-B01 @PRCMS-B07 @unit-level
  Scenario: A review is written to the reviews' file and committed alone
    Given a repository with a commit to push, an old review in the file and another file changed
    When an ok review with a note is recorded over the range
    Then the file holds the new review with its commit, text and note, the old one is gone, the commit takes the file alone, and the other change stays

  @PRCMS-E01 @unit-level
  Scenario: A review with no verdict, two verdicts or no author is refused
    Given a repository with a commit to push
    When a review is recorded with no verdict, with both, with no author, and over an empty range
    Then each is refused

  @PRCMS-B02 @PRCMS-B03 @PRCMS-B04 @unit-level
  Scenario: The push asks an ok review of its commits
    Given a repository requiring reviews, with a commit to push
    When the push is confronted with no review, after a failed review, and after an ok review
    Then the first is refused naming the commit, the second shows the failed review, and the third goes through with the review's own commit

  @PRCMS-B05 @unit-level
  Scenario: The refs the hook gets say what the push sends
    Given a repository requiring reviews, with an unreviewed commit
    When the hook's refs are an update, a new branch, and a deletion
    Then the update and the new branch are refused naming the commit, and the deletion goes through

  @PRCMS-B06 @unit-level
  Scenario: Without push_review the push confronts nothing
    Given a repository with no push_review section and an unreviewed commit
    When the push is confronted
    Then it goes through

  @PRCMS-B08 @unit-level
  Scenario: The push asks the judgments its changes reach
    Given a repository requiring judgments, whose pushed change reaches a judgment gate with no verdict
    When the push is confronted, an ok review is recorded, the verdict is written to the map, and the map is committed
    Then the push and the review are refused naming the gate and target, the uncommitted verdict answers nothing, and the committed one lets the push through
