# language: en
# @anchors
#   ref: MPRGM
#   updated_at: 2026-09-26
#   layer: feature

@MPRGM
Feature: MapRegister — hangs the map domain's commands on the root command

  @MPRGM-B01 @unit-level
  Scenario: Registering the map domain makes each of its commands reachable from the root
    Given an empty root command
    When the map domain registers on it
    Then map, impact, ingest, judge, recode, renumber, flow and failures are each found from the root by name

  @MPRGM-X01 @unit-level
  Scenario: Registering the map domain adds no command outside it
    Given an empty root command
    When the map domain registers on it
    Then the root holds exactly those eight commands
