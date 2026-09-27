# language: en
# @anchors
#   ref: FLWRF
#   updated_at: 2026-09-26
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

  @FLWRF-I01 @unit-level
  Scenario: What seeding writes is never outdated for the same configuration
    Given a project seeded with the default branch, and one seeded with the branch "develop"
    When the outdated pipelines are listed with the same configuration
    Then none is outdated

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
