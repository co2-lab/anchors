# language: en
# @anchors
#   ref: OPRGP
#   updated_at: 2026-09-26
#   layer: feature

@OPRGP
Feature: OpsRegister — the operation commands reach the CLI through one registration point

  @OPRGP-B01 @unit-level
  Scenario: The root receives the fourteen operation commands
    Given a bare root command
    When the operations are registered
    Then its commands are board, code, commit-msg, docs, freeze, generated-paths, init, install-hooks, migrate, new, settings, suggest, synthesize and thaw

  @OPRGP-I01 @unit-level
  Scenario: Each operation command is registered exactly once
    Given a bare root command
    When the operations are registered
    Then no name appears twice

  @OPRGP-X01 @unit-level
  Scenario: The registration adds commands and nothing else
    Given a bare root command
    When the operations are registered
    Then the root gains no flag, no pre-run and no run of its own
