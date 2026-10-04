# language: en
# @anchors
#   code: MSFMP
#   ref: MPSTM
#   updated_at: 2026-10-03
#   layer: feature

@MPSTM
Feature: MapStaleness — names the files of the map whose content changed after the map was built

  @MPSTM-B01 @unit-level
  Scenario: A spec edited after the map build is named as stale
    Given a project whose map was just built from a spec
    And the spec is then edited to add a rule
    When the staleness of the map is asked
    Then the spec is the one file named as stale

  @MPSTM-B02 @unit-level
  Scenario: A file removed after the map build is not named
    Given a project whose map was built from a spec
    And the spec is then deleted from disk
    When the staleness of the map is asked
    Then no file is named

  @MPSTM-B03 @unit-level
  Scenario: Staleness follows the content, not the modification time
    Given a project whose map was built from a spec
    And the spec's modification time is moved back without changing its content
    When the staleness of the map is asked
    Then no file is named
    And after the content changes with the old modification time restored, the spec is named

  @MPSTM-B04 @unit-level
  Scenario: A node with an empty recorded revision is not named
    Given a map whose spec node carries an empty revision
    And the spec on disk has content
    When the staleness of the map is asked
    Then no file is named

  @MPSTM-B05 @unit-level
  Scenario: No map names nothing
    Given a project and no map
    When the staleness of the map is asked
    Then no file is named

  @MPSTM-I01 @unit-level
  Scenario: A freshly rebuilt map is never stale
    Given a project whose map was just built from the current files
    When the staleness of the map is asked
    Then no file is named

  @MPSTM-X01 @unit-level
  Scenario: Asking about staleness leaves the map as it was
    Given a project whose map was built and whose spec was then edited
    When the staleness of the map is asked
    Then the map in memory still carries the old revision of the spec

  @MPSTM-E01 @unit-level
  Scenario: A root that cannot be walked names nothing and raises nothing
    Given a map of a project
    And a project root that does not exist
    When the staleness of the map is asked
    Then no file is named and no error is raised
