# language: en
# @anchors
#   code: STFTF
#   ref: GRPRG
#   updated_at: 2026-10-03
#   layer: feature

@GRPRG
Feature: GraphPersistence — saving and loading the map file without churn and without partial reads

  @GRPRG-B01 @unit-level
  Scenario: Saving stamps the current format and the running binary's release
    Given a graph that carries no format number, and a running binary whose release is 0.1.10
    When the graph is saved
    Then the file says version 4 and generated_by 0.1.10, and it loads back with no migration asked

  @GRPRG-B02 @unit-level
  Scenario: The saved file starts with the fixed comment header
    Given any graph
    When the graph is saved
    Then the first line of the file is the comment naming anchors.graph.yaml and anchors map build

  @GRPRG-B03 @unit-level
  Scenario: A save that changes only the writer's release leaves the file untouched
    Given a map saved by release 0.1.10
    When the same graph is saved by release dev
    Then the file on disk is byte for byte the same

  @GRPRG-B04 @unit-level
  Scenario: A real change rewrites the file with the running release
    Given a map saved by release 0.1.10 and then re-saved unchanged by release dev
    When a node is added and the graph is saved by release dev
    Then the file is rewritten and says generated_by dev

  @GRPRG-B08 @unit-level
  Scenario: A later release restamps an unchanged map, an earlier one does not
    Given a map saved by release 0.1.248
    When the same graph is saved by release 0.1.258, then by 0.1.250 and by dev
    Then the file says generated_by 0.1.258 and keeps it

  @GRPRG-B05 @unit-level
  Scenario: Loading a map in an unreadable format is refused
    Given a map file in format 1 and another in format 9
    When each is loaded
    Then each load returns the format refusal and no graph

  @GRPRG-I01 @unit-level
  Scenario: A saved graph loads back unchanged
    Given a graph with nodes, an edge with a stamp and a judgment, and a node signal
    When the graph is saved and loaded back
    Then the loaded graph equals the saved one

  @GRPRG-X01 @unit-level
  Scenario: An unreadable map yields no graph at all
    Given a map file in format 9 that also holds nodes
    When it is loaded
    Then no graph is returned, not even the nodes it recognised

  @GRPRG-E01 @unit-level
  Scenario: Loading a missing file returns the read error
    Given a path where no file exists
    When it is loaded
    Then an error is returned and no graph

  @GRPRG-E02 @unit-level
  Scenario: Loading text that is not the map returns the parse error
    Given a file whose text is not the map's structure
    When it is loaded
    Then an error is returned and no graph

  @GRPRG-E04 @unit-level
  Scenario: Saving into a missing directory returns the write error
    Given a path inside a directory that does not exist
    When a graph is saved there
    Then the write error is returned

  @GRPRG-B06 @unit-level
  Scenario: A map read from bytes is read like one from disk
    Given the bytes of a current map and of a map in an unreadable format
    When each is read from its bytes
    Then the first gives its graph and the second the format refusal naming where it came from

  @GRPRG-B07 @unit-level
  Scenario: An empty signal loads as no signal
    Given a map with a node whose signal is empty and one whose signal was measured
    When it is loaded
    Then the first has no signal and the second keeps its own
