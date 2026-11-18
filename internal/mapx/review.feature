# language: en
# @anchors
#   code: RVFTB
#   ref: MPRVM
#   updated_at: 2026-10-03
#   layer: feature

@MPRVM
Feature: MapReview — the reviews recorded on a node, and which of them still holds

  @MPRVM-B01 @unit-level
  Scenario: A review holds only at the revision it looked at
    Given a node reviewed for rule-fulfilled at its current revision
    When its review for rule-fulfilled is asked, and for another gate, and after the node's revision changes
    Then the first holds, the second does not, and the third does not
    And a node not in the map has no review

  @MPRVM-B02 @unit-level
  Scenario: Recording a review replaces the same gate's and keeps the others
    Given a node with a review for g1 and one for g2
    When a new review for g1 is recorded
    Then the node has the new g1 review at its current revision and keeps the g2 one
    And recording on a node not in the map records nothing and says so

  @MPRVM-B03 @unit-level
  Scenario: A rebuild keeps the reviews whatever the revision
    Given an earlier map whose node carries a review
    When a rebuilt map is preserved from it, the node's revision changed
    Then the rebuilt node carries the review

  @MPRVM-I01 @unit-level
  Scenario: A node holds at most one review per gate
    Given two reviews of one gate and one of another, recorded on one node
    When its reviews are counted
    Then there is one per gate

  @MPRVM-X01 @unit-level
  Scenario: Who reviewed is kept as given
    Given a review by "agent:openai/gpt"
    When it is recorded
    Then it is kept as "agent:openai/gpt", with no grade
