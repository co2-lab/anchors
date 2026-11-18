# language: en
# @anchors
#   code: VRFRB
#   ref: VRFRA
#   updated_at: 2026-11-18
#   layer: feature

@VRFRA
Feature: VersionForward — a change to a product's version moves it forward

  @VRFRA-B01 @unit-level
  Scenario: A version kept or raised passes
    Given a product whose source holds 1.2.0 at HEAD
    When the tree keeps 1.2.0, and when it raises it to 1.3.0
    Then the gate passes both times

  @VRFRA-B02 @unit-level
  Scenario: A version that goes back, moves sideways, cannot be read or disagrees fails naming the product
    Given a product whose sources hold 1.2.0 at HEAD
    When the tree holds 1.1.0, 1.2.0+b over 1.2.0+a, 1.3, or 1.3.0 in one source and 1.2.0 in the other
    Then the gate fails each time, naming the product and the versions

  @VRFRA-B03 @unit-level
  Scenario: A version HEAD did not hold passes
    Given a product whose source file is not at HEAD
    When the tree holds 0.1.0 in it
    Then the gate passes

  @VRFRA-B04 @unit-level
  Scenario: No product with a source skips the gate
    Given a versioning block whose only product has no source
    When the gate runs
    Then it is skipped, saying there is no written version to confront
