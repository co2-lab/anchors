# language: en
# @anchors
#   ref: MPCMM
#   updated_at: 2026-09-26
#   layer: feature

@MPCMM
Feature: MapCommand — builds the dependency map from the project and answers questions about it

  @MPCMM-B01 @unit-level
  Scenario: A rebuild keeps the judgment recorded on the guide edge
    Given a project whose guide-to-login edge was judged "issue" by the atomic gate
    When the map is built again
    Then the guide-to-login edge still carries the atomic judgment "issue"

  @MPCMM-B02 @unit-level
  Scenario: A rebuild keeps the flow graph
    Given a map that carries a flow graph with one state
    When the map is built again
    Then the rebuilt map still has the flow graph with that state

  @MPCMM-B03 @unit-level
  Scenario: A lost judgment stamp is reported per gate
    Given a previous map with two review stamps and one rule-fulfilled stamp
    And a new map with one review stamp and one rule-fulfilled stamp
    When the build compares them
    Then it warns that review went "2 → 1" and does not mention rule-fulfilled
    And with no loss, or with a gain, it says nothing

  @MPCMM-B04 @unit-level
  Scenario: The edge summary shows every type, the triad's first
    Given a map with edges of types realizes, depends-on twice and specifies
    When the edge summary is printed
    Then specifies comes first, then "depends-on  2", then "realizes    1"
    And a map with no edges prints no summary

  @MPCMM-B05 @unit-level
  Scenario: The layer ambiguity warning is grouped by pair of layers
    Given two files decided in favour of shared over test and one in favour of lambdas over test
    When the ambiguity warning is printed
    Then it counts "3 file(s)", shows "shared beat test" with "(2 file(s)" and "lambdas beat test"
    And it says the wrong layer takes the file out of reach of "EVERY gate"
    And the build of a project whose test file matches both a test and a shared layer shows "shared beat test"

  @MPCMM-B06 @unit-level
  Scenario: Showing the login code lists what governs it and marks it a leaf
    Given the login project's map
    When the map shows "src/login.ts"
    Then it lists the guide and the spec above it and "(none — leaf)" below it
    And showing the login spec marks it "(none — top/root)" with its specifies and tested-by edges

  @MPCMM-B07 @unit-level
  Scenario: The orphans are the nodes with no edge
    Given the login project's map, where only the lonely guide has no edge
    When the map shows the orphans
    Then it lists "[guide] guides/LONELY.md" under a count of 1 and not the code guide

  @MPCMM-B08 @unit-level
  Scenario: The statistics count nodes by kind and edges by type
    Given the login project's map
    When the map shows the statistics
    Then it says "map: 5 nodes, 3 edges" with 2 guides, 1 code, 1 governs, 1 specifies and 1 tested-by

  @MPCMM-B09 @unit-level
  Scenario: The worklist puts rulers and specs before the code they govern
    Given the login project's map
    When the map shows the worklist
    Then the code guide and the login spec come before the login code
    And it ends with "5 node(s), in topological order."

  @MPCMM-B10 @unit-level
  Scenario: The pending worklist lists only nodes with a failing gate
    Given the login project where the login code fails the always-red gate and the atomic judgment is only pending
    When the map shows the worklist of pending nodes
    Then only the login code is listed, with "1 file(s) with pending items"

  @MPCMM-I01 @unit-level
  Scenario: A judgment survives a rebuild of the map
    Given the login code was judged by the atomic gate
    When the map is built again
    Then the atomic gate is still answered on the login code

  @MPCMM-X01 @unit-level
  Scenario: The warnings do not fail the build
    Given a project whose test file is claimed by two layers with no priority
    When the map is built
    Then the build succeeds and prints the ambiguity warning

  @MPCMM-E01 @unit-level
  Scenario: Building with no configuration points at init
    Given a directory with no anchors.yaml
    When the map is built
    Then it fails with a message pointing at "anchors init"

  @MPCMM-E02 @unit-level
  Scenario: Showing with no map points at the build
    Given a directory with no map file
    When the map shows the statistics
    Then it fails with "run `anchors map build`"

  @MPCMM-E03 @unit-level
  Scenario: Showing a file that is not in the map is refused
    Given the login project's map
    When the map shows "src/ghost.ts"
    Then it fails naming "ghost.ts" as not in the map

  @MPCMM-E04 @unit-level
  Scenario: Showing with no selector is refused
    Given the login project's map
    When the map show runs with no file and no switch
    Then it fails with "provide <file>"

  @MPCMM-E05 @unit-level
  Scenario: The pending worklist with no configuration is refused
    Given a map file with no anchors.yaml beside it
    When the map shows the worklist of pending nodes
    Then it fails with "load config"
