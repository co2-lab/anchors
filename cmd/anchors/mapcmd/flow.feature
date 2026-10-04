# language: en
# @anchors
#   code: FLFTA
#   ref: FLWOX
#   updated_at: 2026-10-03
#   layer: feature

@FLWOX
Feature: Flow — build, draw and navigate the work flows kept in the map

  @FLWOX-B01 @unit-level
  Scenario: A project with no flow declared is told so and the build is not an error
    Given a project with a map and no flow file
    When the flow build runs
    Then it prints that no flow is declared and succeeds

  @FLWOX-B02 @unit-level
  Scenario: Building the flow writes it into the map and keeps the map nodes
    Given a project with a map of five nodes, a worker flow of four steps and a pull action of three results
    When the flow build runs
    Then it reports four steps and three results
    And the map holds the seven flow states and still its five nodes

  @FLWOX-B03 @unit-level
  Scenario: A result no flow routes is named as a warning and the build still succeeds
    Given a pull action whose third result no flow routes
    When the flow build runs
    Then it warns of one result no flow handles, naming it, and succeeds

  @FLWOX-B04 @unit-level
  Scenario: The next command lists the valid exits of a step with the suggestion of each result
    Given a built worker flow whose pull step routes two results, one suggesting the work command
    When the next command runs for the pull step written in lower case
    Then it prints the step, the pull action it fits, the two valid exits with their destination titles and the suggestion
    And an exit of the work step that no result triggers shows a dash and its condition

  @FLWOX-B05 @unit-level
  Scenario: A step that fits an action no file declares is flagged
    Given a built worker flow whose work step fits an action no action file declares
    When the next command runs for the work step
    Then it warns that no action file declares it

  @FLWOX-B06 @unit-level
  Scenario: A terminal step says the work ends and a step with no exit and no terminal mark says it is stuck
    Given a built worker flow with a terminal done step and a stuck step with no exit
    When the next command runs for each of them
    Then the done step says the work ends and the stuck step says whoever arrives is stuck

  @FLWOX-B07 @unit-level
  Scenario: The text drawing keeps the file order, hides the results and sorts each step exits
    Given a built worker flow
    When the worker flow is shown as text
    Then the steps appear in the file order with the terminal marked
    And the exits of each step are sorted by result, no result is drawn as a step, and an exit with no result reads sempre with its condition

  @FLWOX-B08 @unit-level
  Scenario: The mermaid drawing emits valid identifiers, native line breaks, the stadium shape for terminals and highlights the entry and the ends
    Given a built worker flow whose work step title holds double quotes
    When the flow is shown as mermaid, left to right
    Then identifiers use underscores, the quotes become single quotes, the fitted action follows a native line break
    And the done step is a stadium, exits carry the result short names, and the pull step and the done step get the entry and end classes

  @FLWOX-B09 @unit-level
  Scenario: An unknown mermaid direction falls back to top-down
    Given a built worker flow
    When the flow is shown as mermaid with the direction sideways, and with rl and BT
    Then sideways yields a top-down diagram and rl and BT are kept in upper case

  @FLWOX-I01 @unit-level
  Scenario: What the flow build writes is what show and next read
    Given a worker flow built into the map
    When the flow is shown and a step is navigated
    Then the declared steps and exits appear as they were built

  @FLWOX-X01 @unit-level
  Scenario: Showing and navigating the flow leave the map unchanged
    Given a worker flow built into the map
    When the flow is shown as text and as mermaid and a step is navigated
    Then the map is byte for byte what it was

  @FLWOX-E01 @unit-level
  Scenario: Building a flow without a map is refused and no map is created
    Given a directory with a worker flow and no map
    When the flow build runs
    Then it fails asking to build the map first
    And no map file exists afterwards

  @FLWOX-E02 @unit-level
  Scenario: Asking the exits of a step that is not in the flow graph fails
    Given a built worker flow
    When the next command runs for a step code that is not in it
    Then it fails saying the step is not in the flow graph

  @FLWOX-E03 @unit-level
  Scenario: Showing a flow no name matches fails listing the available flows
    Given a built worker flow
    When a flow named nonexistent is shown
    Then it fails listing the available flows

  @FLWOX-E04 @unit-level
  Scenario: Showing or navigating a map with no flow asks for the flow build
    Given a project whose map holds no flow
    When the flow is shown, or a step is navigated
    Then it fails asking to run the flow build
