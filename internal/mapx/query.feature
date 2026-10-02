# language: en
# @anchors
#   ref: GRQRG
#   updated_at: 2026-10-02
#   layer: feature

@GRQRG
Feature: GraphQueries — read-only questions over the loaded map

  @GRQRG-B01 @unit-level
  Scenario: A guide governs only the targets of its governs edges, sorted
    Given guide A governs y.tsx and x.tsx, and a spec specifies x.tsx
    When what guide A governs is asked
    Then the answer is x.tsx then y.tsx

  @GRQRG-B02 @unit-level
  Scenario: The governance summary counts governs edges only
    Given guide A governs two files and a spec specifies one of them
    When the governance summary is computed
    Then guide A counts 2 and the spec does not appear

  @GRQRG-B03 @unit-level
  Scenario: The neighbourhood is one step each way, sorted
    Given a spec governed by a guide, specifying code and covered by a feature, with the edges stored out of order
    When the neighbourhood of the spec is asked
    Then the incoming list is the guide's edge and the outgoing list is the code edge then the feature edge

  @GRQRG-B04 @unit-level
  Scenario: Orphans are the nodes with no edge, sorted
    Given nodes zeta and alpha with no edge, among connected nodes
    When the orphans are asked
    Then the answer is alpha then zeta

  @GRQRG-B05 @unit-level
  Scenario: Statistics count nodes and edges by kind and type
    Given a graph of 5 nodes and 3 edges
    When the statistics are computed
    Then they report 5 nodes, 3 edges, one spec, one code, one governs and one specifies

  @GRQRG-B06 @unit-level
  Scenario: Parents come before their children
    Given a guide that governs a spec, which specifies code and is covered by a feature, which is tested by a test
    When the parents-first order is computed
    Then the guide precedes the spec, the spec precedes the code and the feature, and the feature precedes the test

  @GRQRG-B07 @unit-level
  Scenario: A depends-on edge imposes no order
    Given node b depends on node a, and no vertical edge links them
    When the parents-first order is computed
    Then a comes before b, as the alphabet says

  @GRQRG-B08 @unit-level
  Scenario: Cycle members are appended at the end
    Given A specifies B, B specifies A, and C stands alone
    When the parents-first order is computed
    Then the order is C, then A, then B

  @GRQRG-I01 @unit-level
  Scenario: The order does not depend on how nodes are stored
    Given the same graph stored with its nodes in two different orders
    When the parents-first order of each is computed
    Then both orders are identical

  @GRQRG-X01 @unit-level
  Scenario: Queries leave the graph as it was
    Given a graph
    When every query is run over it
    Then the graph is identical to its state before

  @GRQRG-B09 @unit-level
  Scenario: A node's unit codes come from its identity edges
    Given a spec PAYMT that specifies pay.go and is covered by pay.feature, which tests pay_test.go, and a spec SHARE that also specifies pay.go
    When the unit codes are asked
    Then pay.go has PAYMT and SHARE, pay_test.go has PAYMT, and a node no spec reaches has none
