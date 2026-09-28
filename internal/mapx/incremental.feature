# language: en
# @anchors
#   ref: GRINC
#   updated_at: 2026-09-28
#   layer: feature

@GRINC
Feature: IncrementalMap — new files enter the map without the tree being read

  @GRINC-B01 @unit-level
  Scenario: A new file enters with the same nodes and edges the full build gives
    Given a map built without a unit's spec
    When the spec is added
    Then the map equals the full build, and no file outside the unit was read

  @GRINC-B02 @unit-level
  Scenario: The derivation links inside the unit are replaced, removals included
    Given a map of a unit with no feature
    When the feature is added, then a test, then both
    Then each map equals the full build, and the spec→test link of the unit without a feature is gone

  @GRINC-B03 @unit-level
  Scenario: A new anchor gives its declared code to its derived siblings
    Given a project whose anchor is the spec, and a map without the spec
    When the spec declaring its code is added
    Then the feature takes the spec's code, as in the full build

  @GRINC-B04 @unit-level
  Scenario: What the new file declares reaches the existing files
    Given a map without a plan that seeds an existing spec
    When the plan is added
    Then the map equals the full build, seed included

  @GRINC-B05 @unit-level
  Scenario: A known, ignored or unclassified file adds nothing
    Given a complete map
    When a known file and a file the reader does not give are added
    Then nothing is read for the known one and nothing is added for either

  @GRINC-B06 @unit-level
  Scenario: On disk the addition runs under the lock, and not without a map
    Given a project with a map and a new spec on disk, and a project with no map
    When the missing files are added in each
    Then the first map has the spec and the second project still has no map

  @GRINC-B07 @unit-level
  Scenario: The anchor that may own a new file is read, for the override its header chooses
    Given a spec anchor whose header layer chooses an override putting the test under __tests__
    When the code beside it, then the test under __tests__, then a test named after the directory, is added
    Then each map equals the full build

  @GRINC-X01 @unit-level
  Scenario: A relation an existing file declares toward the new one waits for the full build
    Given a map without a spec that an existing plan seeds
    When the spec is added
    Then the plan's seed is not in the map, and the full build has it

  @GRINC-E01 @unit-level
  Scenario: A reader that fails changes nothing
    Given a map and a reader that fails
    When a new file is added
    Then the error comes back and the map is as it was
