# language: en
# @anchors
#   ref: STCMS
#   updated_at: 2026-09-26
#   layer: feature

@STCMS
Feature: SettingsCommand — one agent's local decisions, declared with a date and kept out of the project's configuration

  @STCMS-B01 @unit-level
  Scenario: A decision without a date is refused
    Given a project with no local settings
    When role and user-issues run without --date
    Then each fails naming --date
    And no settings file is written

  @STCMS-B02 @unit-level
  Scenario: An unknown role or answer is refused naming it
    Given a project with no local settings
    When role runs with "wizard" and user-issues with "maybe"
    Then each fails naming the value given
    And no settings file is written

  @STCMS-B03 @unit-level
  Scenario: Declaring a role records it with the agent and the date
    Given a project whose settings hold a legacy escalated-cards answer
    When role runs with "architect" and --date 2026-09-26
    Then the settings hold role architect, date 2026-09-26 and the agent
    And the legacy answer is gone
    And the output names .anchors/settings.yaml

  @STCMS-B04 @unit-level
  Scenario: The escalated-cards answer is recorded with the agent and the date
    Given a project with no local settings
    When user-issues is answered "nao" with --date 2026-09-26
    Then the settings hold the answer false, the date 2026-09-26 and the agent

  @STCMS-B05 @unit-level
  Scenario: The role is asked on the terminal when not given
    Given the terminal answers "qa"
    When role runs with no argument
    Then the recorded role is qa

  @STCMS-B06 @unit-level
  Scenario: An unclear reply is asked again and never assumed
    Given the terminal answers "talvez" then "nao"
    When user-issues runs with no argument
    Then it points out it did not understand "talvez"
    And it records false
    And and three unclear replies make it fail with no recognized answer

  @STCMS-B07 @unit-level
  Scenario: Show lists the role's capabilities
    Given a project whose agent declared the architect role
    When show runs
    Then it prints every capability of the architect role

  @STCMS-B08 @unit-level
  Scenario: Show without a role teaches how to declare one
    Given a project with no role declared
    When show runs
    Then it prints "anchors settings role --date" and no capabilities

  @STCMS-I01 @unit-level
  Scenario: A role declaration leaves one source for the escalated-cards question
    Given settings holding a legacy escalated-cards answer
    When a role is declared
    Then reading the settings back finds no legacy answer

  @STCMS-X01 @unit-level
  Scenario: The decisions go to the local settings file
    Given a project
    When a role is declared
    Then the output names .anchors/settings.yaml as where it was recorded

  @STCMS-E01 @unit-level
  Scenario: A closed input fails instead of assuming no
    Given a terminal input that is already closed
    When user-issues asks for the answer
    Then it fails
