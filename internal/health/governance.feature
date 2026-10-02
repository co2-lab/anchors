# language: en
# @anchors
#   ref: GVOPG
#   updated_at: 2026-10-02
#   layer: feature

@GVOPG
Feature: GovernanceOpportunities — the doctor suggests the canonical gates and settings a project has not adopted yet

  @GVOPG-B01 @unit-level
  Scenario: Tests and code without evidence freshness get the suggestion
    Given a map with a code node and a test node, and only a tests-green gate declared
    When the opportunities are checked
    Then there is an informational sugestao-gate on evidence-fresh

  @GVOPG-B02 @unit-level
  Scenario: Code without the secret gate gets the suggestion
    Given a map with a code node and no gate declared
    When the opportunities are checked
    Then there is an informational sugestao-gate on no-secret-leaked

  @GVOPG-B03 @unit-level
  Scenario: A secret gate that does not block is suboptimal
    Given a map with a code node and no-secret-leaked declared with blocking false
    When the opportunities are checked
    Then there is a gate-subotimo on no-secret-leaked

  @GVOPG-B04 @unit-level
  Scenario: Code without dependency audit or duplication gates gets both suggestions
    Given a map with a code node and no gate declared
    When the opportunities are checked
    Then there are sugestao-gate findings on dependency-vulnerable and on no-duplication

  @GVOPG-B05 @unit-level
  Scenario: Tests without JUnit output are a suboptimal configuration
    Given a map with a test node and one test suite without JUnit output
    When the opportunities are checked
    Then there is a config-subotima on tests.junit

  @GVOPG-B06 @unit-level
  Scenario: A nil map or configuration gives nothing
    Given a nil map, or a nil configuration
    When the opportunities are checked
    Then there is no finding

  @GVOPG-B07 @unit-level
  Scenario: The quick hints are the first two opportunities
    Given a project with code and tests where five opportunities fire
    When the quick hints are asked
    Then they are the first two of the full list, and a list of one is returned whole

  @GVOPG-I01 @unit-level
  Scenario: Every opportunity is informational
    Given a project where every opportunity fires
    When the opportunities are checked
    Then every finding has severity info

  @GVOPG-X01 @unit-level
  Scenario: What the project already adopted is not suggested
    Given a map with code and tests, every suggested gate declared with the secret gate blocking, and a suite with JUnit output
    When the opportunities are checked
    Then there is no finding

  @GVOPG-B08 @unit-level
  Scenario: Every applicable undeclared catalog gate is suggested
    Given a project with a feature layer and a catalog holding feature-test-match and scenario-letter-declared presupposing rule_types
    When the governance opportunities are read
    Then both are suggested as informational, with what they measure and how to declare them
    And scenario-letter-declared's suggestion names rule_types
