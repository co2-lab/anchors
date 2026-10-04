# language: en
# @anchors
#   code: RGFTC
#   ref: QLCMQ
#   updated_at: 2026-10-03
#   layer: feature

@QLCMQ
Feature: QualityCommands — the quality domain puts its twelve commands under the root command

  @QLCMQ-B01 @unit-level
  Scenario: The quality domain registers exactly its twelve commands
    Given an empty root command
    When the quality domain registers into it
    Then the root holds exactly check, verify, doctor, status, stale, coverage, report, test, mutation, stamp and touch

  @QLCMQ-B02 @unit-level
  Scenario: Each quality command is reachable by its name
    Given a root command into which the quality domain registered
    When each of the twelve names is looked up from the root
    Then each lookup finds the command of that name

  @QLCMQ-I01 @unit-level
  Scenario: No two quality commands share a name
    Given a root command into which the quality domain registered
    When the names of its commands are compared
    Then every name is distinct

  @QLCMQ-X01 @unit-level
  Scenario: Registering the quality commands prints nothing
    Given an empty root command
    When the quality domain registers into it
    Then nothing is written to the standard output
