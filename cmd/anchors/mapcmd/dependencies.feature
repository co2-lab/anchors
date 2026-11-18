# language: en
# @anchors
#   code: DPNDF
#   ref: DPNDS
#   updated_at: 2026-11-15
#   layer: feature

@DPNDS
Feature: Dependencies — everything a file depends on, the other way of impact

  @DPNDS-B01 @unit-level
  Scenario: The direct dependencies show, and --propagate goes further with the path of each item
    Given a screen that imports a util and is specified by a spec, in a navigator the app imports
    When dependencies runs on the app, and again with --propagate cascade, and with --stop-at the navigator
    Then the direct level holds the navigator, the cascade reaches the screen and the util with their paths, and the stop goes no further

  @DPNDS-B02 @unit-level
  Scenario: The files and json formats print the dependencies
    Given the same project
    When dependencies runs on the screen as files and as json
    Then the files are the spec and the util, one per line, and the json carries the origins, levels and paths

  @DPNDS-X01 @unit-level
  Scenario: The dependencies query leaves the map unchanged
    Given the same project
    When dependencies runs
    Then the map is byte for byte as it was
