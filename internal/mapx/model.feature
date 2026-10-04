# language: en
# @anchors
#   code: MDFTA
#   ref: GRMDG
#   updated_at: 2026-10-03
#   layer: feature

@GRMDG
Feature: GraphModel — the shape of the map file, and when a validated relation goes stale

  @GRMDG-B01 @unit-level
  Scenario: A relation never stamped is stale
    Given nodes a at r2 and b at r5, and a relation from a to b with no stamp
    When its staleness is asked
    Then it is stale

  @GRMDG-B02 @unit-level
  Scenario: A stamp is fresh while both ends keep their revisions
    Given nodes a at r2 and b at r5
    When a relation stamped with r2 and r5 and another stamped with r2 and r4 are judged
    Then the first is fresh and the second is stale

  @GRMDG-B03 @unit-level
  Scenario: A mutation scope is stale only against a recorded, different revision
    Given one scope measured at r1, one measured with no recorded revision
    When each is judged against the file at r2 and at r1
    Then only the scope measured at r1 judged against r2 is stale

  @GRMDG-I01 @unit-level
  Scenario: The map file's keys are fixed and English
    Given a graph, a node, a relation, a stamp and a judgment with every field filled
    When each is written in the map format
    Then each writes exactly the keys the spec lists

  @GRMDG-I02 @unit-level
  Scenario: Empty optional fields are left out
    Given a node with only id, kind and revision, and a relation with only its ends, type and origin
    When each is written in the map format
    Then the node writes only id, kind and rev, and the relation only from, to, type and origin

  @GRMDG-I03 @unit-level
  Scenario: The node kinds and relation types are the fixed vocabulary
    Given every node kind and every relation type of the model
    When their written values are listed
    Then they are exactly the kinds and types the spec lists

  @GRMDG-X01 @unit-level
  Scenario: The stamp's date and verdict play no part in staleness
    Given a relation stamped with both current revisions, verdict issue and an old date
    When its staleness is asked
    Then it is fresh
