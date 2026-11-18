# language: en
# @anchors
#   code: MMFMP
#   ref: MPMRM
#   updated_at: 2026-10-03
#   layer: feature

@MPMRM
Feature: MapMerge — the git merge driver that unites two versions of the map instead of merging text

  @MPMRM-B01 @unit-level
  Scenario: The merged map is written onto our side's file
    Given our side's map judges one edge and the other side's map judges another
    When the map merge driver runs with the base, our side and the other side
    Then our side's file holds the merged map with both judgments

  @MPMRM-B02 @unit-level
  Scenario: A node created only on the other branch reaches the merged map
    Given both sides share one node, our side has one more and the other side has two more
    When the map merge driver runs
    Then the merged map has the shared node, our node and the other side's two nodes

  @MPMRM-B03 @unit-level
  Scenario: A node both sides have keeps our side's revision
    Given the same node carries revision "r-ours" on our side and "r-theirs" on the other side
    When the map merge driver runs
    Then the merged node carries revision "r-ours"

  @MPMRM-B04 @unit-level
  Scenario: An edge created only on the other branch arrives with its judgment
    Given the other side has an edge from "c.spec.md" judged by review that our side lacks
    When the map merge driver runs
    Then the merged map has that edge with its review judgment

  @MPMRM-B05 @unit-level
  Scenario: A side with no judgment on a shared edge does not erase the other side's
    Given both sides share an edge, and only one of them holds a review judgment on it
    When the map merge driver runs, once with the judgment on our side and once on the other
    Then the merged map keeps exactly one judgment on the edge both times

  @MPMRM-B06 @unit-level
  Scenario: The driver reports on the error stream what it kept and what came from the other side
    Given our side has one judgment and the other side has one more on another edge
    When the map merge driver runs
    Then the error stream says "2 judgment(s) preserved (1 came from the other side)"

  @MPMRM-B07 @unit-level
  Scenario: Judgments of different gates on a shared edge are joined, and one gate judged on both sides keeps the latest
    Given our side judged and stamped a shared edge with "review", "atomic" (issue, older) and "tie", and the other side, later, with "rule-fulfilled", "atomic" (ok, newer) and "tie" on the same date
    When the map merge driver runs
    Then the edge carries "review", "rule-fulfilled", the other side's "atomic" and our "tie"
    And the edge keeps the other side's newer check stamp

  @MPMRM-B08 @unit-level
  Scenario: The other side's flow states and transitions and the failures of a shared node reach the merged map
    Given the other side has a flow and failures on a node both sides have, one rule observed on both sides
    When the map merge driver runs, with and without a flow on our side
    Then the result has the flow states of both sides and the other side's transition
    And the shared node has the failures of both sides, the rule observed on both keeping the latest

  @MPMRM-I01 @unit-level
  Scenario: Nothing of either side is missing from the merged map
    Given our side and the other side each have an edge and a node the other lacks
    When the map merge driver runs
    Then every node and every edge of both sides is in the merged map

  @MPMRM-X01 @unit-level
  Scenario: The base version is never read
    Given a base path that does not exist
    When the map merge driver runs with it and two readable sides
    Then the merge succeeds

  @MPMRM-E01 @unit-level
  Scenario: An unreadable side fails the merge naming the side
    Given the other side's path does not exist, or our side's file is not a map
    When the map merge driver runs
    Then it fails with "read the other side" or "read our side" respectively
    And our side's file is left as it was

  @MPMRM-E02 @unit-level
  Scenario: A call with two paths is refused
    Given only our side and the other side are passed, without the base
    When the map merge driver runs
    Then it is refused for the number of arguments
