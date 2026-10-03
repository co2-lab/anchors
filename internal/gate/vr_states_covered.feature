# language: en
# @anchors
#   ref: VRSTC
#   updated_at: 2026-10-03
#   layer: feature

@VRSTC
Feature: VRStatesCovered — each state of a visual unit tied to its visual regression, both ways

  @VRSTC-B01 @unit-level
  Scenario: Nothing to confront leaves every gate without a verdict
    Given a spec node, a code file with no spec beside it, and a spec with no code
    When the four VR gates confront each
    Then each leaves without a verdict

  @VRSTC-B02 @unit-level
  Scenario: Every state needs a VR scenario
    Given a spec with states S01 and S02 and a feature with a VR scenario of S01 only
    When vr-states-covered confronts the unit's code
    Then it fails naming BUTTN-S02 and the regime tag

  @VRSTC-B03 @unit-level
  Scenario: Every VR scenario needs a VR test
    Given VR scenarios of S01 and S02 and a test naming BUTTN-VR-S01 only
    When vr-scenarios-tested confronts the unit's code
    Then it fails naming BUTTN-VR-S02 as untested

  @VRSTC-B04 @unit-level
  Scenario: Every VR scenario needs its baseline image, in any image format
    Given VR scenarios of S01 and S02, tests of both, and only Button.BUTTN-VR-S01.svg
    When vr-scenarios-tested confronts the unit's code
    Then it fails naming BUTTN-VR-S02 and the name to save it under

  @VRSTC-B05 @unit-level
  Scenario: A VR scenario of a state the spec does not register
    Given a VR scenario of S03 in a unit whose spec registers S01 and S02
    When vr-scenarios-of-states confronts the unit's code
    Then it fails naming BUTTN-S03

  @VRSTC-B06 @unit-level
  Scenario: A VR scenario of no state
    Given a VR scenario tagged @BUTTN-VR alone
    When vr-scenarios-of-states confronts the unit's code
    Then it fails naming BUTTN-VR as a scenario of no state

  @VRSTC-B07 @unit-level
  Scenario: A VR test of a scenario the feature does not declare
    Given a test naming BUTTN-VR-S01 and BUTTN-VR-S03, and a VR scenario of S01 only
    When vr-tests-of-scenarios confronts the unit's code
    Then it fails naming BUTTN-VR-S03 and the test

  @VRSTC-B08 @unit-level
  Scenario: A VR scenario carries the regime tag and the state's code, in either form
    Given scenarios tagged "@BUTTN-S01 @vr-level" and "@BUTTN-VR-S02 @vr-level", and "@BUTTN-S01 @unit-level"
    When the unit is read
    Then the VR scenarios are of S01 and S02, and the unit-level one is not a VR scenario

  @VRSTC-B09 @unit-level
  Scenario: Which tests are of the unit
    Given a flow in a folder named Button, a test beside the unit, a flow named by the code, an unrelated test, and an image
    When the unit is read
    Then the first three are of the unit, and the unrelated test and the image are not

  @VRSTC-B10 @unit-level
  Scenario: The State letter comes from the project's rule types
    Given rule_types where the letter E is the term "Estado"
    When the State letter is read
    Then it is E, and S without rule types

  @VRSTC-B11 @unit-level
  Scenario: vr-baseline accepts other image formats
    Given a VR scenario whose baseline is a jpg
    When vr-baseline confronts the feature
    Then it passes

  @VRSTC-I01 @unit-level
  Scenario: One state never answers for another
    Given S01 with scenario, test and image, and S02 with none
    When the four VR gates confront the unit's code
    Then the failures name S02 alone

  @VRSTC-E01 @unit-level
  Scenario: A code file with no spec beside it leaves without a verdict
    Given Button.styles.ts with no spec of its own
    When vr-states-covered confronts it
    Then it leaves without a verdict, saying there is no spec beside it

  @VRSTC-E02 @unit-level
  Scenario: A unit with no feature has no VR scenario
    Given a visual unit with states and no feature file
    When vr-states-covered confronts the unit's code
    Then it fails naming every state

  @VRSTC-E03 @unit-level
  Scenario: A baseline that cannot be found counts as none
    Given a VR scenario whose image is not on disk
    When vr-scenarios-tested confronts the unit's code
    Then it fails naming that scenario's image

  @VRSTC-E04 @unit-level
  Scenario: A test that cannot be read is read by its path
    Given a flow BUTTN-VR-S01.yaml the map lists and the disk does not have
    When vr-scenarios-tested confronts the unit's code
    Then the flow still counts as the VR test of S01
