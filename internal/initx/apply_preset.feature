# language: en
# @anchors
#   ref: APPRP
#   updated_at: 2026-09-26
#   layer: feature

@APPRP
Feature: ApplyPreset — writes a stack preset's layers into the configuration and deduces one identity prefix per module

  @APPRP-B01 @unit-level
  Scenario: Applying a preset adds its layers and keeps the layers of other names
    Given a configuration holding a guide layer
    And the go preset
    When the preset is applied
    Then the configuration holds the cmd, internal, pkg and test layers
    And the guide layer is still there
    And a configuration with no layer set receives every layer of the preset

  @APPRP-B02 @unit-level
  Scenario: A layer with the same name as a preset layer is replaced by the preset's
    Given a configuration whose internal layer is a doc layer with pattern old/**
    When the go preset is applied
    Then the internal layer is the preset's code layer, without the old pattern

  @APPRP-B03 @unit-level
  Scenario: Each module receives a two-letter prefix keyed by its directory name
    Given the modules src/features/auth, src/features/audit and src/features/family/
    When the prefixes are deduced
    Then there are three prefixes of two letters, each the identity prefix of its module
    And family, despite the trailing slash, receives FM

  @APPRP-B04 @unit-level
  Scenario: A colliding prefix keeps its first letter and takes the first free second letter
    Given the modules m/auth and m/atlas, whose identity prefixes are both AT
    When the prefixes are deduced
    Then atlas, first in alphabetical order, keeps AT
    And auth receives AA

  @APPRP-B05 @unit-level
  Scenario: Applying a preset returns the prefix deduced for each detected module
    Given the node-ts preset and the detected module src/modules/family
    When the preset is applied
    Then the returned mapping is family to FM
    And with no detected module the returned mapping is empty

  @APPRP-I01 @unit-level
  Scenario: The same modules in any order give the same prefixes
    Given the modules auth, atlas and family in two different orders
    When the prefixes are deduced for each order
    Then both mappings are equal

  @APPRP-X01 @unit-level
  Scenario: Prefixes are deduced from the given paths without reading the disk
    Given the module path /does/not/exist/family, which is absent from disk
    When the prefixes are deduced
    Then family still receives FM

  @APPRP-B06 @unit-level
  Scenario: Modules sharing a folder name are each keyed by their path, with distinct prefixes
    Given the modules packages/auth, apps/auth/ and m/family
    When their prefixes are deduced
    Then apps/auth and packages/auth are each in the mapping under their path, with different prefixes
    And family stays keyed by its name with FM

  @APPRP-I02 @unit-level
  Scenario: No two modules ever share a prefix while a free one exists
    Given thirty modules whose names all start with A
    When their prefixes are deduced
    Then all thirty have a two-letter prefix and no prefix repeats
