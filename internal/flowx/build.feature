# language: en
# @anchors
#   ref: FLBLF
#   updated_at: 2026-09-26
#   layer: feature

@FLBLF
Feature: FlowBuild — assembling the work-flow graph from the project's flow and action files

  @FLBLF-B01 @unit-level
  Scenario: Actions then flows are read in a stable order, and no flows give no graph
    Given a project with no "flows/" folder, and a project with the actions "b" and "a" and the flows "z" and "m"
    When each flow graph is built
    Then the first has no graph and no error, and the second reads "a", "b", then "m", "z"

  @FLBLF-B02 @unit-level
  Scenario: Every coded heading is a state of its file
    Given the unblocking flow with five coded headings
    When the flow graph is built
    Then it has 5 states, each recorded with its title and the flow file

  @FLBLF-B03 @unit-level
  Scenario: Exits belong to the state above them
    Given the unblocking flow whose ATTACKING state lists two exits, and a flow with an exit written above every state
    When the flow graphs are built
    Then ATTACKING has two exits, the second to "DSTRV-N04" with its condition, and the stray exit belongs to nobody

  @FLBLF-B04 @unit-level
  Scenario: A line with two codes routes a result
    Given a step whose section says "- `ACTST-R01` STALE → `FLOWX-P01`"
    When the flow graph is built
    Then the step has one transition, on "ACTST-R01", to "FLOWX-P01"

  @FLBLF-B05 @unit-level
  Scenario: Terminal is declared, not inferred
    Given the unblocking flow where FIXED declares "@terminal" and MEASURED has an exit
    When the flow graph is built
    Then FIXED is terminal and MEASURED is not

  @FLBLF-B06 @unit-level
  Scenario: The piece a step fits is read in any language
    Given a step saying "Fits: `ACMAP`", a step with no piece, and a Spanish flow saying "Encaja: `ACTST`"
    When the flow graphs are built
    Then the steps fit "ACMAP", nothing, and "ACTST"

  @FLBLF-B07 @unit-level
  Scenario: A suggested reaction is recorded as a suggestion
    Given an action whose result STALE says "Sugere: `anchors map build`" and whose result FINE suggests nothing
    When the flow graph is built
    Then STALE carries the suggestion and FINE carries none
