# language: en
# @anchors
#   code: ENFTN
#   ref: DCENV
#   updated_at: 2026-10-03
#   layer: feature

@DCENV
Feature: Environment — the project's environment variables page, compiled from its specs

  @DCENV-B01 @unit-level
  Scenario: Every declared variable is listed by name, in any language
    Given an English spec declaring PAYMENTS_URL and a Portuguese one declaring LOG_LEVEL, and a TODO row
    When envVars is read
    Then it lists LOG_LEVEL and PAYMENTS_URL, in that order, and no placeholder

  @DCENV-B02 @unit-level
  Scenario: A variable two units read is listed once, with both
    Given both specs declaring PAYMENTS_URL
    When envVars is read
    Then PAYMENTS_URL is listed once with both units

  @DCENV-B03 @unit-level
  Scenario: The environment page is seeded and compiled
    Given a project whose specs declare variables
    When the scaffolds are written and the docs built
    Then docs/environment.md has a row per variable naming the units that read it

  @DCENV-E01 @unit-level
  Scenario: A placeholder row is no variable
    Given a spec whose Environment Variables table keeps the template's TODO row
    When envVars is read
    Then the TODO row is not listed

