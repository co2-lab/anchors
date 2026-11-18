# language: en
# @anchors
#   code: MDFTM
#   ref: FLMDF
#   updated_at: 2026-10-03
#   layer: feature

@FLMDF
Feature: FlowModel — the questions a work-flow graph answers: what comes next, and what is broken

  @FLMDF-B01 @unit-level
  Scenario: From a state, only its declared exits
    Given the unblocking flow whose OPEN state has one exit and whose FIXED state is terminal
    When the exits of OPEN and of FIXED are asked
    Then OPEN has 1 exit and FIXED has none

  @FLMDF-B02 @unit-level
  Scenario: A state is found by its code
    Given a flow with the step "FLOWX-T01" titled "a later step"
    When "FLOWX-T01" and "FLOWX-T99" are looked up
    Then the first is found with its title and the second is not

  @FLMDF-B05 @unit-level
  Scenario: A result is a code whose letter is O
    Given the codes "ACTST-O01", "FLOWX-T01", "DSTRV-N02", "FLOWX-XO01" and "ACTST"
    When each is asked whether it is a result
    Then "ACTST-O01" is a result and the other four codes are not results

  @FLMDF-B03 @unit-level
  Scenario: States keep the file's order and flows are listed once
    Given a flow file declaring "FLOWX-T03" before "FLOWX-T01", and an action file "check-it"
    When the states of the flow and the list of flows are asked
    Then the states are "FLOWX-T03" then "FLOWX-T01", and the flows are the action file then the flow file

  @FLMDF-B04 @unit-level
  Scenario: The entry is the first declared step, never a result
    Given a flow file whose first step is "FLOWX-T03", and an action file declaring only results
    When the entry of each is asked
    Then the flow's entry is "FLOWX-T03" and the action file has none

  @FLMDF-B06 @unit-level
  Scenario: An action's title is its file name
    Given the action file "check-it.action.md" declaring the results of "ACTST"
    When the titles of "ACTST" and "NOPEX" are asked
    Then "ACTST" is "check-it" and "NOPEX" has none

  @FLMDF-B07 @unit-level
  Scenario: A step nobody arrives at is unreachable
    Given the unblocking flow plus the step "DSTRV-N09" no exit points at
    When the unreachable steps are asked
    Then exactly "DSTRV-N09" is unreachable, and no result is ever listed

  @FLMDF-B08 @unit-level
  Scenario: A result no flow routes is unhandled
    Given an action declaring "ACTST-O01" and "ACTST-O02", and a flow routing only "ACTST-O01"
    When the unhandled results are asked
    Then the unhandled list holds only "ACTST-O02", the result no flow routes anywhere

  @FLMDF-B09 @unit-level
  Scenario: An exit to a state that does not exist is dangling
    Given the unblocking flow plus a state whose exit points at "DSTRV-N77"
    When the dangling exits are asked
    Then exactly the exit to "DSTRV-N77" is dangling

  @FLMDF-X01 @unit-level
  Scenario: A project without flows answers nothing to every question
    Given no flow graph
    When every question is asked of it
    Then each answers nothing
