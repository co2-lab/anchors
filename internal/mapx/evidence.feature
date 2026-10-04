# language: en
# @anchors
#   ref: EVFRA
#   updated_at: 2026-10-03
#   layer: feature

@EVFRA
Feature: EvidenceFreshness — a test's evidence expires when anything it exercises changes, not only its own file

  @EVFRA-B01 @unit-level
  Scenario: A test that was never ingested has no verdict
    Given a suite script node with no signal
    When the freshness of its evidence is judged
    Then there is no verdict

  @EVFRA-B02 @unit-level
  Scenario: The test's own change expires its evidence
    Given a suite script ingested at revision t1 whose file is now at revision t2
    When the freshness of its evidence is judged
    Then the verdict marks the evidence stale by its own change

  @EVFRA-B03 @unit-level
  Scenario: A composed script that changed is named as the culprit
    Given a suite script ingested with login.yaml at u1 and launchApp.yaml at l1 in its closure
    And login.yaml is now at u2 while the suite script itself did not change
    When the freshness of its evidence is judged
    Then the verdict names login.yaml as the only culprit and does not mark the test's own change

  @EVFRA-B04 @unit-level
  Scenario: Nothing moved means no verdict
    Given a suite script ingested with its closure, and no node changed revision since
    When the freshness of its evidence is judged
    Then there is no verdict

  @EVFRA-B05 @unit-level
  Scenario: A signal with no recorded closure is judged by its own file only
    Given a suite script ingested with no closure, and login.yaml changed since
    When the freshness of its evidence is judged
    Then there is no verdict, until the suite script's own file changes

  @EVFRA-B06 @unit-level
  Scenario: The closure descends transitively with current revisions
    Given a suite script that composes login.yaml at u1, which composes launchApp.yaml at l1
    When the closure of the suite script is computed
    Then it holds login.yaml at u1 and launchApp.yaml at l1

  @EVFRA-B07 @unit-level
  Scenario: A no-propagation node is in the closure but not walked through
    Given login.yaml is marked no-propagation
    When the closure of the suite script is computed
    Then login.yaml is in the closure and launchApp.yaml is not

  @EVFRA-B08 @unit-level
  Scenario: A recorded node that left the graph is not a culprit
    Given a suite script whose recorded closure names a file that is no longer in the graph
    When the freshness of its evidence is judged
    Then there is no verdict

  @EVFRA-I01 @unit-level
  Scenario: The closure never contains the test itself
    Given a suite script inside a composition chain
    When its closure is computed
    Then the suite script is not in it

  @EVFRA-X01 @unit-level
  Scenario: The closure never climbs to the spec above the test
    Given a spec with an edge to the suite script
    When the closure of the suite script is computed
    Then the spec is not in it

  @EVFRA-B09 @unit-level
  Scenario: A capture's target enters the closure and is not descended
    Given a VR flow that captures Button.tsx, which depends on Icon.tsx
    When the flow's evidence closure is taken
    Then Button.tsx is in it and Icon.tsx is not

  @EVFRA-B10 @unit-level
  Scenario: A component whose capture diverged stales the captures of who uses it
    Given a screen capture ingested at 10:00, composing a component whose capture failed at 11:00
    When the screen capture's freshness is asked
    Then it is stale, naming the component; with the component's capture failing at 09:00, or passing, it is fresh
