# language: en
# @anchors
#   ref: VRSTC
#   updated_at: 2026-10-03
#   layer: feature

@VRSTC
Feature: VRStatesCovered — every state of a visual unit is proven by visual regression

  @VRSTC-B01 @unit-level
  Scenario: Specs that have nothing to cover leave without a verdict
    Given a spec node, a code file with no spec beside it, a spec with no code, and a spec with no state
    When vr-states-covered confronts each
    Then each leaves without a verdict

  @VRSTC-B02 @unit-level
  Scenario: The unit's feature must declare the VR scenario
    Given a screen spec with states and a feature without the BUTTN-VR scenario
    When vr-states-covered confronts the unit's code
    Then it fails naming BUTTN-VR and the visual-regime tag

  @VRSTC-B03 @unit-level
  Scenario: Every state needs its baseline image, in any image format
    Given a screen spec with states S01 and S02 and only Button.BUTTN-VR-S01.svg beside it
    When vr-states-covered confronts the unit's code
    Then it fails naming BUTTN-S02 and the name to save it under

  @VRSTC-B04 @unit-level
  Scenario: A capture test is a flow named by the VR code or a screenshot test beside the unit
    Given a spec whose only test naming BUTTN-VR is a baseline image
    When vr-states-covered confronts the unit's code, then again with a flow BUTTN-VR.yaml, then with Button.vr.test.ts naming BUTTN-VR
    Then the first fails saying no test captures BUTTN-VR and the other two find the capture

  @VRSTC-B05 @unit-level
  Scenario: Scenario, images and capture pass
    Given a screen spec with states, its VR scenario, one image per state and a capture flow
    When vr-states-covered confronts the unit's code
    Then it passes

  @VRSTC-B06 @unit-level
  Scenario: The State letter comes from the project's rule types
    Given rule_types where the letter E is the term "Estado"
    When the State letter is read
    Then it is E, and S without rule types

  @VRSTC-B07 @unit-level
  Scenario: vr-baseline accepts other image formats
    Given a VR scenario whose baseline is a jpg
    When vr-baseline confronts the feature
    Then it passes

  @VRSTC-I01 @unit-level
  Scenario: One state's image never covers another
    Given a spec with states S01 and S02 and the image of S01 only
    When vr-states-covered confronts the unit's code
    Then it fails naming S02

  @VRSTC-E01 @unit-level
  Scenario: A code file with no spec beside it leaves without a verdict
    Given Button.styles.ts with no Button.styles.spec.md
    When vr-states-covered confronts it
    Then it leaves without a verdict, saying there is no spec beside it

  @VRSTC-E02 @unit-level
  Scenario: A unit with no feature has no VR scenario
    Given a visual unit with states and no feature file
    When vr-states-covered confronts the unit's code
    Then it fails naming the missing VR scenario

  @VRSTC-E03 @unit-level
  Scenario: A state whose image cannot be found counts as without one
    Given a visual unit whose S02 image is not on disk
    When vr-states-covered confronts the unit's code
    Then it fails naming S02

  @VRSTC-E04 @unit-level
  Scenario: An unreadable test beside the unit is not the capture
    Given a test beside the unit that the map lists and the disk does not have
    When vr-states-covered confronts the unit's code
    Then it fails saying no test captures the VR

