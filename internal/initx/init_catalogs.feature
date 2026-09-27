# language: en
# @anchors
#   ref: INCTN
#   updated_at: 2026-09-26
#   layer: feature

@INCTN
Feature: InitCatalogs — the stack preset catalog and the @TBD instruction that init seeds into judgment gates

  @INCTN-B01 @unit-level
  Scenario: A preset layer with no kind becomes a code layer
    Given a preset with a layer that declares no kind and a test layer with pattern t/** and one tag
    When the preset is turned into configuration layers
    Then the layer with no kind is a code layer
    And the test layer keeps its kind, pattern and tag

  @INCTN-B02 @unit-level
  Scenario: Looking up a preset by an unknown name finds nothing
    Given the preset catalog
    When the presets go and no-such-stack are looked up
    Then go is found
    And no-such-stack is not found and the returned preset is empty

  @INCTN-B03 @unit-level
  Scenario: The preset names are listed in catalog order
    Given the preset catalog
    When the preset names are listed
    Then there is one name per preset, in the order of the catalog

  @INCTN-B04 @unit-level
  Scenario: The @TBD instruction forbids pass, orders a waiver naming the absence, and names the piece asked about
    Given a judgment gate asking about the code
    When its @TBD instruction is produced
    Then the English text mentions @TBD, the `waived` answer and the naming of the absence, and forbids `pass`
    And the instruction produced for the test names the test

  @INCTN-B05 @unit-level
  Scenario: A modular preset declares the directory of its modules
    Given the preset catalog
    When each modular preset is inspected
    Then every modular preset has a module directory

  @INCTN-I01 @unit-level
  Scenario: Every preset has a unique name, a title, patterned layers and a test layer
    Given the preset catalog
    When every preset is inspected
    Then no name repeats, every preset has a title and layers
    And every layer has a pattern and every preset has a test layer

  @INCTN-I02 @unit-level
  Scenario: The @TBD instruction demands checking that the @TBD is still true
    Given a judgment gate asking about the code
    When its @TBD instruction is produced
    Then the text tells the judge that a marker whose piece already exists is out of date

  @INCTN-X01 @unit-level
  Scenario: No preset layer carries an identity prefix
    Given the preset catalog
    When every preset layer is inspected
    Then no layer carries a code prefix
