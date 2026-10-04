# language: en
# @anchors
#   code: SCFTA
#   ref: GSRGH
#   updated_at: 2026-10-03
#   layer: feature

@GSRGH
Feature: GherkinScenarioReader — the scenarios of a unit's feature, with their steps, for the documentation

  @GSRGH-B01 @unit-level
  Scenario: Scenarios and outlines open in any dialect, and an examples table does not
    Given a Portuguese feature with "Esquema do Cenário: por <caso>" followed by an "Exemplos:" table and a "Cenário: comum"
    And an English feature with "Scenario Outline: by <c>" and one with "Scenario: it works"
    When their scenarios are read
    Then the titles are "por <caso>", "comum", "by <c>" and "it works"
    And the examples table stays in the outline's body

  @GSRGH-B02 @unit-level
  Scenario: The first code-shaped tag is the code and the others are tags
    Given a scenario tagged "@GLCGL-B01 @nivel-unit"
    When it is read
    Then its code is "GLCGL-B01" and its tags are exactly "nivel-unit"

  @GSRGH-B03 @unit-level
  Scenario: The body is the scenario's non-blank lines until the next scenario
    Given a feature with two scenarios
    When the scenarios are read
    Then the first body contains "Então a régua não libera" and nothing of the second scenario
    And a blank line inside a scenario does not appear in its body

  @GSRGH-B04 @unit-level
  Scenario: A Background heading after a scenario does not leak into it
    Given a scenario followed by a "Contexto:" block
    When the scenarios are read
    Then the scenario's body does not contain the Background's steps

  @GSRGH-B05 @unit-level
  Scenario: The feature is found by the map edge, in either direction, not by name
    Given a spec whose feature was renamed and whose edge follows the new name
    And another project where the edge goes from the feature to the spec
    When the spec's scenarios are read
    Then both find the two scenarios

  @GSRGH-B06 @unit-level
  Scenario: The scenarios of a selection follow the spec selection and its errors
    Given a project with one spec in layer infra that has two scenarios
    When the scenarios of "layer=infra" and of "layer=nothere" are asked
    Then the first gives the two scenarios and the second returns the spec selection's error

  @GSRGH-I01 @unit-level
  Scenario: Every scenario read carries its spec's code
    Given the spec GLCGL with a linked feature
    When its scenarios are read
    Then each scenario names GLCGL as its spec

  @GSRGH-X01 @unit-level
  Scenario: The body is kept verbatim
    Given a scenario whose step is indented by four spaces
    When it is read
    Then the body line keeps its indentation and its words

  @GSRGH-E01 @unit-level
  Scenario: A linked feature missing on disk contributes no scenario and no error
    Given a spec linked to two features, one of them missing on disk
    When the spec's scenarios are read
    Then only the readable feature's scenarios come back
