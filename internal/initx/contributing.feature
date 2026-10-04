# language: en
# @anchors
#   code: CNFTD
#   ref: CNGDC
#   updated_at: 2026-10-03
#   layer: feature

@CNGDC
Feature: ContributingGuide — render the project's CONTRIBUTING.md from the configuration init writes

  @CNGDC-B01 @unit-level
  Scenario: The whole guide is a title, a seeding note and the section
    Given a configuration
    When the whole guide is rendered
    Then it starts with "# Contributing", says it was seeded from anchors.yaml, and ends with the section

  @CNGDC-B02 @unit-level
  Scenario: The section states the order of the work
    Given a configuration
    When the section is rendered
    Then it opens with "## Working with Anchors"
    And it says the spec is the anchor and a bug fix writes the rule, then the scenario, then the test

  @CNGDC-B03 @unit-level
  Scenario: Where each piece lives
    Given spec, feature and test layers and a colocation of feature and test
    When the section is rendered
    Then it lists the spec, feature and test patterns in that order
    And the colocation templates beside the spec
    And a configuration with no artifact layer lists no piece

  @CNGDC-B04 @unit-level
  Scenario: The declared code layers, or that there are none
    Given the code layers handlers-code and api-code
    When the section is rendered
    Then it lists api-code then handlers-code, each with its pattern
    And with no code layer it says no code layer is declared yet

  @CNGDC-B05 @unit-level
  Scenario: The queue the daily commands name
    Given a github workflow on acme/app
    When the section is rendered
    Then anchors next names the issues of acme/app
    And without a workflow it names the local queue

  @CNGDC-B06 @unit-level
  Scenario: What blocks and what informs
    Given a blocking gate tests-green and an informative gate spec-complete
    When the section is rendered
    Then tests-green bars the commit and spec-complete only informs
    And with no blocking gate it says every gate only informs

  @CNGDC-B07 @unit-level
  Scenario: The project's guide folder
    Given the guide folder docs/guides
    When the section is rendered
    Then it names docs/guides/
    And with no guide folder no folder is named

  @CNGDC-I01 @unit-level
  Scenario: The guide never names what the configuration does not hold
    Given a configuration with the code layer core-code and the gate tests-green
    When the section is rendered
    Then the only layer listed is core-code and the only gate named is tests-green

  @CNGDC-X01 @unit-level
  Scenario: Rendering writes nothing to disk
    Given an empty working folder
    When the guide is rendered
    Then the folder is still empty and the guide comes back as text
