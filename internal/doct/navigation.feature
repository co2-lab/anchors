# language: en
# @anchors
#   code: NVFTN
#   ref: DCNAV
#   layer: feature

@DCNAV
Feature: Navigation — the app's navigation map and the code's dependency chain, compiled into pages

  @DCNAV-B01 @unit-level
  Scenario: navigation lists the screen specs and the edges the code's flags declare, each end taken to its screen
    Given an app of four screens, Home leading to GoalDetail, GoalDetail to GoalEdit, GoalEdit back, and Lonely reached by nothing
    When navigation is read
    Then it lists the four screens and the three edges with their rules

  @DCNAV-B02 @unit-level
  Scenario: An entry is marked, and a screen no entry reaches is marked unreached; with no entry, none is
    Given Home as the entry
    When the flowchart is drawn
    Then Home is a stadium, Lonely is dashed, and with no entry declared no screen is unreached

  @DCNAV-B03 @unit-level
  Scenario: ScaffoldNavigation and ScaffoldDependencies are seeded when the app has screens and dependency flags, and compile into the pages
    Given the app, and Home importing GoalDetail with a dependency flag
    When docs init and docs build run
    Then the navigation page draws the flowchart and the table, and the dependencies page lists who uses whom
