# language: en
# @anchors
#   ref: DCRQD
#   updated_at: 2026-09-19
#   layer: feature

@DCRQD
Feature: DocRequired — the aggregated document the unit must feed

  @DCRQD-B01 @unit-level
  Scenario: A mandatory document that does not exist fails
    Given a project declaring a document as mandatory for this unit's layer
    And that document does not exist on disk
    When the gate confronts the unit
    Then it returns Fail

  @DCRQD-B02 @unit-level
  Scenario: A document that exists and does not mention the unit fails
    Given the mandatory document exists and never names this unit
    When the gate confronts the unit
    Then it returns Fail, because existence alone would approve an empty file created
      to silence the gate

  @DCRQD-B03 @unit-level
  Scenario: A mention by the identity code counts as documented
    Given the mandatory document cites the unit's identity code
    When the gate confronts the unit
    Then it returns Pass

  @DCRQD-B04 @unit-level
  Scenario: A mention by the file name also counts
    Given the mandatory document cites the unit's file name and not its code
    When the gate confronts the unit
    Then it returns Pass, because the document speaks of the unit either way

  @DCRQD-B05 @unit-level
  Scenario: Satisfying one of two duties is not enough
    Given two documents declared mandatory for this unit's layer
    And only one of them mentions the unit
    When the gate confronts the unit
    Then it returns Fail, because each document is charged on its own

  @DCRQD-B06 @unit-level
  Scenario: Without a declaration nothing is charged
    Given a project that declares no mandatory document
    When the gate confronts the unit
    Then it returns Pass, because the ruler is what the project committed to, not what
      one supposes it owes

  @DCRQD-B07 @unit-level
  Scenario: A layer with no trigger is not charged
    Given a mandatory document whose trigger names another layer
    And the document exists and does not mention this unit
    When the gate confronts the unit
    Then it returns Pass, because this layer triggers no duty

  @DCRQD-B08 @unit-level
  Scenario: Aggregated, the verdict is one per document
    Given three units of a triggering layer and one mandatory document
    When the gate runs over the whole project
    Then it reports ONE verdict for that document, not one per unit

  @DCRQD-I01 @unit-level
  Scenario: The duty starts from the spec, not from the code
    Given a unit whose spec and code both exist
    When the gate confronts the project
    Then the charge lands on the spec, because the spec is what declares the unit

  @DCRQD-I02 @unit-level
  Scenario: The layer used is the UNIT's, not the node's
    Given a spec whose node layer is the spec layer
    And the unit it describes belongs to a triggering layer
    When the gate confronts it
    Then the duty is charged, because reading the node's layer would charge every spec
      of the project the same duty, or none

  @DCRQD-I03 @unit-level
  Scenario: Without a map the aggregated verdict is skipped
    Given no graph built
    When the gate runs aggregated
    Then it does not approve, because approving without being able to look would stamp
      what was never measured

  @DCRQD-X01 @unit-level
  Scenario: The gate does not understand the document's content
    Given the mandatory document cites the unit and describes it wrongly
    When the gate confronts the unit
    Then it returns Pass, because the gate separates "not documented" from "documented" —
      judging the quality of the documentation is another ruler

  @DCRQD-X02 @unit-level
  Scenario: The gate does not decide which documents are mandatory
    Given a project whose Structure declares no duty for this layer
    And a document that a reviewer would consider obviously required
    When the gate confronts the unit
    Then it returns Pass, because inventing duties would charge what nobody committed to
