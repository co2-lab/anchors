# language: en
# @anchors
#   code: IMFTA
#   ref: IMANM
#   updated_at: 2026-11-16
#   layer: feature

@IMANM
Feature: ImpactAnalysis — what changing one file propagates to, and what it must be confronted against

  @IMANM-B01 @unit-level
  Scenario: Changing a spec propagates down to its code, feature and test
    Given a Login spec that specifies Login code, is covered by a Login feature, which is tested by a Login test
    When the impact of changing the Login spec is analysed
    Then the propagation list is the Login feature, the Login test and the Login code, in sorted order

  @IMANM-B02 @unit-level
  Scenario: A no-propagation child is reached but the wave stops there
    Given the Login feature is marked no-propagation and tested by the Login test
    When the impact of changing the Login spec is analysed
    Then the Login feature is in the propagation list and the Login test is not

  @IMANM-B03 @unit-level
  Scenario: Changing code is validated upward against its spec and the spec's guide
    Given a spec guide that governs the Login spec, which specifies the Login code
    When the impact of changing the Login code is analysed
    Then the validation list is the Login spec and the spec guide, and the propagation list is empty

  @IMANM-B04 @unit-level
  Scenario: Climbing to a shared guide does not reach the sibling unit
    Given a spec guide that governs both the Login spec and the Home spec
    When the impact of changing the Login spec is analysed
    Then neither the Home spec nor the Home code appears in either list

  @IMANM-B05 @unit-level
  Scenario: A node with no edges has no impact
    Given a node with no edges
    When its impact is analysed
    Then both lists are empty

  @IMANM-I01 @unit-level
  Scenario: The changed node is never in its own lists, even in a cycle
    Given nodes A and B with edges from A to B and from B to A
    When the impact of changing A is analysed
    Then A appears in neither list, and B appears in both

  @IMANM-X01 @unit-level
  Scenario: The analysis leaves the graph as it was
    Given a graph with nodes, edges and a no-propagation mark
    When the impact of changing one node is analysed
    Then the graph is identical to its state before the analysis

  @IMANM-B06 @unit-level
  Scenario: The dependents are the side that depends, by level
    Given a screen imported by a navigator, tested, captured and specified, that imports a util
    When its dependents are walked to the direct and to the indirect level
    Then the direct are the navigator, the test and the capture, and the indirect add what depends on the navigator

  @IMANM-B07 @unit-level
  Scenario: A tie of types alone, a stop and @noPropagation are listed and not walked past
    Given a screen whose navigator is imported for its types by another screen
    When its dependents are walked to the end, with a stop at the navigator, and with the navigator @noPropagation
    Then the importer of types is marked and nothing past it is reached, and the stop and the marked node are not climbed from

  @IMANM-B08 @unit-level
  Scenario: The dependencies are the other end of the same ties
    Given the same screen
    When its dependencies are walked to the direct level and to the end
    Then the direct are the util and the spec, and the end reaches the spec's guide

  @IMANM-B09 @unit-level
  Scenario: The rulers are what governs the file, and the tests are those reached or holding it
    Given the same screen
    When its rulers and the tests holding it are asked
    Then the rulers are the spec and its guide, and the tests are those whose evidence holds it, or the test reached

  @IMANM-B10 @unit-level
  Scenario: A parent that imports only types is validated and not climbed past
    Given a navigator imported by another screen for its types alone, and that screen imported in turn
    When the impact of changing a module the navigator reaches is analysed
    Then the importer of types is in the validation list, and neither its own importer nor its spec is
