# language: en
# @anchors
#   code: WRFTA
#   ref: FLWRF
#   updated_at: 2026-10-03
#   layer: feature

@FLWRF
Feature: FlowWorkflows — declare the pipelines of the work flow, find what is missing or broken, and seed them without taking over what the team owns

  @FLWRF-B01 @unit-level
  Scenario: Every declared pipeline has a carried template, a role and its serialization need
    Given the flow's list of pipelines
    When each entry is checked
    Then each has a carried template and a role
    And only "anchors-gates.yml" and "anchors-board.yml" do not require serialization

  @FLWRF-B02 @unit-level
  Scenario: A pipeline is missing only when its file is absent
    Given an empty project, and then the same project with "anchors-claim.yml" holding "not even yaml"
    When the missing pipelines are listed
    Then every pipeline is missing in the empty project
    And "anchors-claim.yml" is not missing once its file exists

  @FLWRF-B03 @unit-level
  Scenario: A serial pipeline without serialization is flagged
    Given "anchors-claim.yml" with no concurrency block
    And "anchors-stale.yml" with a concurrency block that cancels the run in progress
    And "anchors-gates.yml" with no concurrency block
    When the pipelines without serialization are listed
    Then exactly "anchors-claim.yml" and "anchors-stale.yml" are flagged
    And "anchors-claim.yml" is not reported as missing

  @FLWRF-B04 @unit-level
  Scenario: Seeding writes the missing pipelines and returns them sorted
    Given an empty project, and a project whose first pipeline carries the marker with old content
    When the pipelines are seeded
    Then every pipeline is written in the empty project and the written list is sorted
    And the pipeline with the marker is rewritten with the current template and listed as written

  @FLWRF-B05 @unit-level
  Scenario: Seeding leaves a pipeline without the marker untouched
    Given a project whose first pipeline was edited by the team and has no marker
    When the pipelines are seeded
    Then the team's pipeline keeps its content
    And it is not listed as written

  @FLWRF-B06 @unit-level
  Scenario: The board page is seeded outside the pipelines folder, with the marker
    Given an empty project
    When the pipelines are seeded
    Then ".github/anchors-board.html" exists and carries the marker
    And its path is outside ".github/workflows"

  @FLWRF-B07 @unit-level
  Scenario: The board page seeding reports created, updated or unchanged
    Given a project with no board page, one with the identical page, one with an older Anchors page, and one with a page the team took over
    When the pipelines are seeded in each
    Then the outcomes are created, unchanged, updated and unchanged
    And the team's page keeps its content

  @FLWRF-B08 @unit-level
  Scenario: An intact pipeline that differs from its template is outdated
    Given a seeded project where "anchors-claim.yml" got one extra line, "anchors-stale.yml" lost its marker and "anchors-guard.yml" was removed
    When the outdated pipelines are listed
    Then only "anchors-claim.yml" is outdated

  @FLWRF-B09 @unit-level
  Scenario: The integration branch, main when none is declared, replaces every marked branch line
    Given a project whose integration branch is "develop", one that declares "main", and one that declares none
    When the pipelines are seeded for each
    Then every line the templates mark as the integration branch reads "branches: [develop] # anchors:integration-branch" in the first, with the template's indentation
    And reads "branches: [main] # anchors:integration-branch" in the other two, the resolve-queue pipeline included

  @FLWRF-B10 @unit-level
  Scenario: Anchors writes the board columns only up to READY TO TEST
    Given the board columns
    When the columns Anchors writes are asked
    Then they are "TO DO", "IN PROGRESS", "READY TO REVIEW", "IN REVIEW" and "READY TO TEST"

  @FLWRF-B11 @unit-level
  Scenario: A per-card label is its prefix followed by the card
    Given the cards "44", "556" and "900"
    When the per-card labels are built
    Then they read "anchors:under-44", "anchors:desbloqueia-44", "anchors:from-pr-556" and "anchors:blocked-by-900"

  @FLWRF-B12 @unit-level
  Scenario: The stale pipeline releases an idle owned card once
    Given in-progress cards 5, 6 and 7 idle past the window, 6 already released and the comments of 7 unreadable, and card 8 with recent progress
    When the stale pipeline runs
    Then card 5 is released naming its previous owner "a/one"
    And cards 6 and 8 are not released, and card 7 is skipped and named
    And the card list never asks for comments

  @FLWRF-B13 @unit-level
  Scenario: A rejected card is released after the shorter rework window
    Given to-do cards 5 and 6 idle for 1h30, where only 5's owner came back from a rejected review, and card 7 idle for 3h
    When the stale pipeline runs
    Then cards 5 and 7 are released
    And card 6 keeps its owner

  @FLWRF-B14 @unit-level
  Scenario: The review status follows the assigned reviewer's verdict line
    Given card 12 moved to review with "agent-b" as its owner, and a pull request that references it
    When the review job runs over each history of pull request comments
    Then "anchors-review: approved by agent-b" publishes success and "rejected" publishes failure, the last verdict winning
    And another agent's line, a line from before the move, inside a code block, mid-sentence or by someone without write access leaves it pending
    And a card back in ready-to-review is pending awaiting a reviewer, and a pull request with no card or a card never moved to review gets no status

  @FLWRF-B15 @unit-level
  Scenario: The review job is wired to verdict comments and feeds the mover
    Given the pull request checks pipeline
    When its triggers, permissions and jobs are read
    Then it listens to pull request comments and may write statuses
    And the review job runs on comments carrying "anchors-review:"
    And the mover needs the review job, runs always but never on a comment, and receives its state and reviewer

  @FLWRF-B16 @unit-level
  Scenario: A merge without the review outcome is said on the card
    Given a card whose pull request merges with the review status success, pending, failure or none
    When the mover runs on the merge
    Then the card moves to ready-to-test every time
    And only when the status was not success does the card say "merged WITHOUT the review outcome", naming the reviewer and the status or that no reviewer had been assigned

  @FLWRF-B17 @unit-level
  Scenario: A green PR publishes the review as pending
    Given an in-progress card whose pull request turned green
    When the mover runs
    Then the card moves to ready-to-review
    And the review status is published as pending "awaiting a reviewer"

  @FLWRF-B18 @unit-level
  Scenario: The claim teaches the reviewer the verdict line
    Given the claim pipeline
    When its instructions to the reviewer are read
    Then they teach "anchors-review: approved by <você>" and "anchors-review: rejected by <você>"
    And they no longer tell the reviewer to move the card on approval

  @FLWRF-B19 @unit-level
  Scenario: A verdict releases the reviewer once
    Given a card in review with the reviewer as owner
    When the reviewer posts an approved or a rejected verdict
    Then the card's owner is released
    And with no verdict yet, or with the card already released, nothing is released

  @FLWRF-B20 @unit-level
  Scenario: The closing line wins over a reference
    Given a pull request body that opens with "Refs #13" and ends with "Closes #12"
    When the review job and the mover run
    Then the review status is published for card 12
    And the merge moves card 12 and leaves card 13 alone

  @FLWRF-B21 @unit-level
  Scenario: A rejection sends the card back to its author
    Given a card implemented by "agent-a" and in review with "agent-b"
    When agent-b posts a rejected verdict
    Then the reviewer is released, the card goes to "anchors:to-do" and its owner is "agent-a" again
    And an approved verdict does none of that

  @FLWRF-I01 @unit-level
  Scenario: What seeding writes is never outdated for the same configuration
    Given a project seeded with the default branch, and one seeded with the branch "develop"
    When the outdated pipelines are listed with the same configuration
    Then none is outdated

  @FLWRF-I02 @unit-level
  Scenario: The parsed verdict line is the one the review guide teaches
    Given the review guide and the review job of the pull request checks pipeline
    When the verdict lines are compared
    Then the guide teaches "anchors-review: approved by <you>" and "anchors-review: rejected by <you>"
    And the review job parses "anchors-review: (approved|rejected) by "

  @FLWRF-X01 @unit-level
  Scenario: A file without the marker is never taken over
    Given a pipeline and a board page the team owns, with no marker
    When the pipelines are seeded and the outdated ones are listed
    Then both keep their content
    And neither is reported as outdated

  @FLWRF-E01 @unit-level
  Scenario: A pipelines folder that cannot be created fails the seeding
    Given a project where ".github" is a file
    When the pipelines are seeded
    Then the seeding fails with an error about creating the folder

  @FLWRF-E02 @unit-level
  Scenario: A pipeline that cannot be written fails the seeding
    Given a project whose pipelines folder is read-only
    When the pipelines are seeded
    Then the seeding fails with an error about writing
    And nothing is listed as written

  @FLWRF-E03 @unit-level
  Scenario: A template the binary does not carry fails the seeding
    Given a binary that carries no pipeline template
    When the pipelines are seeded
    Then the seeding fails with an error about reading the template
