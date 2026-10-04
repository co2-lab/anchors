# language: en
# @anchors
#   code: CPFTC
#   ref: VRCPT
#   updated_at: 2026-10-03
#   layer: feature

@VRCPT
Feature: Captures — a visual-regression test is tied to the unit it captures

  @VRCPT-B01 @unit-level
  Scenario: A VR test captures its unit's code file and images
    Given a unit Button with code BUTTN, its Button.tsx, two baseline images, a flow BUTTN-VR-S01.yaml and a screenshot test citing BUTTN-VR-S02
    When the map is built
    Then each test has a captures edge to Button.tsx and to both images

  @VRCPT-B02 @unit-level
  Scenario: What captures nothing gets no edge
    Given a unit test citing BUTTN-B01, an image, and a flow naming ZZZZZ-VR that no spec declares
    When the map is built
    Then none of them has a captures edge

  @VRCPT-B03 @unit-level
  Scenario: The closure of a capture is one level
    Given a VR flow capturing Button.tsx, which depends on Icon.tsx
    When the flow's evidence closure is taken
    Then it holds Button.tsx and the images, and not Icon.tsx

  @VRCPT-B04 @unit-level
  Scenario: A VR code is read whole with its state
    Given a test that cites BUTTN-VR-S01
    When its codes are scanned
    Then the code is BUTTN-VR-S01

  @VRCPT-I01 @unit-level
  Scenario: The screen's change stales its capture, a component's does not
    Given a VR flow's evidence ingested over Button.tsx, which depends on Icon.tsx
    When Icon.tsx changes, and then Button.tsx changes
    Then the capture is fresh after the first change and stale after the second

  @VRCPT-B05 @unit-level
  Scenario: A contract test captures its API unit, its spec and the OpenAPI document
    Given an API unit Generate with code GENAP, a contract test citing GENAP-CT and docs/openapi.yaml
    When the map is built
    Then the test has captures edges to generate.go, generate.spec.md and docs/openapi.yaml

  @VRCPT-B06 @unit-level
  Scenario: A capture's closure reaches the uncaptured dependencies, and stops at a captured one
    Given a captured screen whose spec depends on a hook, the hook depending on a store, and the screen composing a captured component
    When the screen capture's closure is taken
    Then it holds the hook and the store, and not the component

  @VRCPT-B07 @unit-level
  Scenario: Parts Used names become composes edges
    Given a spec composing BottomSheet and Missing, and a code file BottomSheet.tsx
    When the map is built
    Then the spec composes BottomSheet.tsx, and Missing ties nothing

  @VRCPT-B08 @unit-level
  Scenario: The captures a changed file reaches
    Given the screen, hook and component above
    When the captures reaching the hook, and then the component, are asked
    Then the hook reaches the screen's capture, and the component only its own
