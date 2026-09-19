# language: en
# @anchors
#   ref: PCJPL
#   updated_at: 2026-09-19
#   layer: feature

@PCJPL
Feature: PlanChangeJustified — a modified plan or spec must declare why it changed

  @PCJPL-B01 @unit-level
  Scenario: An artifact reached only by the impact radius is skipped
    Given an artifact that was not modified
    And an execution context where only another plan was changed
    When the gate confronts the artifact
    Then it returns Skip, because an untouched file has nothing to justify

  @PCJPL-B02 @unit-level
  Scenario: An artifact with no identity code is skipped
    Given a changed artifact that declares no identity code
    When the gate confronts the artifact
    Then it returns Skip, because identity enforcement belongs to another gate

  @PCJPL-B03 @unit-level
  Scenario: A changed plan or spec with no declared revision fails
    Given a changed plan or spec with an identity code
    And no revision declared in its content
    When the gate confronts the artifact
    Then it returns Fail, explaining how to record a revision or escalate

  @PCJPL-B04 @unit-level
  Scenario: A changed plan or spec with a valid revision passes
    Given a changed plan or spec declaring a sequential revision for its code
    When the gate confronts the artifact
    Then it returns Pass, confirming the change is justified

  @PCJPL-B05 @unit-level
  Scenario: Citing the revision of another document does not count
    Given a changed artifact citing a revision belonging to another code
    When the gate confronts the artifact
    Then it returns Fail, because another document's revision does not justify this change

  @PCJPL-B06 @unit-level
  Scenario: Revisions with non-sequential numbering fail
    Given a changed artifact with revisions skipping sequential numbers
    When the gate confronts the artifact
    Then it returns Fail, because gaps prevent determining the true revision count

  @PCJPL-B07 @unit-level
  Scenario: Multiple sequential revisions pass and cite the latest
    Given a changed artifact with multiple revisions in sequential order
    When the gate confronts the artifact
    Then it returns Pass, and the verdict names the latest revision

  @PCJPL-B08 @unit-level
  Scenario: Revisions written in various markdown formats pass
    Given a changed artifact with a revision in bold, blockquote, alert, or raw text
    When the gate confronts the artifact
    Then it returns Pass, recognizing the revision across common markdown layouts

  @PCJPL-B09 @unit-level
  Scenario: A change already explained by the plan revision mechanism passes
    Given a changed plan containing revises, revised-by, or amended-by markers
    When the gate confronts the artifact
    Then it returns Pass, avoiding redundant demands for identical information

  @PCJPL-B10 @unit-level
  Scenario: A newly created untracked file is skipped
    Given a newly created file that never existed in the repository
    When the gate confronts the artifact
    Then it returns Skip, because a new file has no previous version to justify

  @PCJPL-B11 @unit-level
  Scenario: A newly created staged file is skipped
    Given a newly created file added to the git index but not yet committed
    When the gate confronts the artifact
    Then it returns Skip, because it does not exist in the previous commit

  @PCJPL-B12 @unit-level
  Scenario: An existing committed file modified without a revision fails
    Given an existing file present in the last commit and modified on disk
    And no revision declared in its content
    When the gate confronts the artifact
    Then it returns Fail, enforcing justification for the modification

  @PCJPL-B13 @unit-level
  Scenario: Outside a git repository the gate trusts the changed list
    Given an execution occurring outside a git repository
    When the gate checks if the file changed according to git
    Then it trusts the caller's changed list rather than silencing itself

  @PCJPL-B14 @unit-level
  Scenario: RevisionsOf parses declared revisions in order
    Given document content containing multiple revisions
    When RevisionsOf parses the text
    Then it returns all matching revisions in order with code, number, and explanation

  @PCJPL-I01 @unit-level
  Scenario: Untouched nodes in the impact radius are never charged
    Given an artifact that was not modified but is present in the impact radius
    When the gate confronts the artifact
    Then it returns Skip, ensuring innocent files are never accused

  @PCJPL-I02 @unit-level
  Scenario: Revisions must be sequential from one without gaps
    Given an artifact whose largest revision number does not match the count
    When the gate confronts the artifact
    Then it returns Fail, maintaining accurate count of document modifications

  @PCJPL-I03 @unit-level
  Scenario: Files not existing in HEAD are never charged
    Given an artifact that did not exist in the last commit
    When the gate confronts the artifact
    Then it returns Skip, verifying that newly introduced files are not penalized

  @PCJPL-X01 @unit-level
  Scenario: The gate does not judge whether a correction is directional
    Given a changed plan declaring a valid revision format
    When the gate confronts the artifact
    Then it returns Pass, leaving semantic interpretation of impact to humans and agents

  @PCJPL-X02 @unit-level
  Scenario: The gate does not enforce identity code presence
    Given an artifact with missing identity code
    When the gate confronts the artifact
    Then it returns Skip, leaving code presence checks to the code presence gate

  @PCJPL-X03 @unit-level
  Scenario: Without a changed files list the gate skips confrontation
    Given a configuration without a list of changed files
    When the gate confronts the artifact
    Then it returns Skip, avoiding unjustified accusations during full runs
