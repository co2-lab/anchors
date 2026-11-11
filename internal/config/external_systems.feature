# language: en
# @anchors
#   code: EXSFA
#   ref: EXSSA
#   updated_at: 2026-11-10
#   layer: feature

@EXSSA
Feature: ExternalSystems — the real systems the code talks to, and the suites that reach them

  @EXSSA-B01 @unit-level
  Scenario: A system declares its calls, the layers that reach it, its version and a recording's age
    Given a configuration declaring a system with a call pattern, two layers, a version and a max age
    When the configuration is loaded
    Then the system holds them, its pattern compiled and its age in days

  @EXSSA-E01 @unit-level
  Scenario: A system that cannot hold refuses the configuration
    Given systems with no name, a repeated name, an invalid pattern, no layer, and an age that is not days
    When each configuration is loaded
    Then each is refused naming the field
