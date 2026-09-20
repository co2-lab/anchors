# language: en
# @anchors
#   ref: RTDCL
#   updated_at: 2026-09-20
#   layer: feature

@RTDCL
Feature: RouteDeclared — a screen declares how one arrives, and names its neighbours

  @RTDCL-B01 @unit-level
  Scenario: An artifact that is not a screen leaves without a verdict, and says why
    Given a spec whose layer is hook, business-logic, store or DAO
    When the gate confronts it
    Then it returns Skip naming the reason, because a bare indeterminate count would
      leave the reader wondering whether the silence is their problem

  @RTDCL-B02 @unit-level
  Scenario: A screen with no named route fails
    Given a screen spec whose header declares no route
    When the gate confronts it
    Then it returns Fail, because without the route the screen is a node nobody can reach

  @RTDCL-B03 @unit-level
  Scenario: A navigation row carrying a generic term fails
    Given a screen spec with a named route
    And a navigation table row naming a generic destination instead of a concrete screen
    When the gate confronts it
    Then it returns Fail, because an edge that names no concrete screen points nowhere

  @RTDCL-B04 @unit-level
  Scenario: A screen with a named route and concrete neighbours passes
    Given a screen spec declaring a named route
    And navigation rows naming concrete screens
    When the gate confronts it
    Then it returns Pass

  @RTDCL-B05 @unit-level
  Scenario: Route and navigation are recognised in either declared language
    Given two equivalent screen specs written in different declared languages
    When the gate confronts each of them
    Then both are measured alike, because a project writing in its own language must be
      measured and not silently approved

  @RTDCL-I01 @unit-level
  Scenario: The header's declared layer wins over the node's tags
    Given an artifact whose header declares one layer and whose tags say another
    When the gate resolves its identity
    Then the header wins, because it is what the author wrote on purpose while a tag can
      come from inference

  @RTDCL-I02 @unit-level
  Scenario: A generic term in prose is not accused
    Given a screen spec with a named route and concrete navigation rows
    And prose inside the navigation section mentioning a generic destination
    When the gate confronts it
    Then it returns Pass, because only table rows declare edges and a sentence is the
      author writing

  @RTDCL-X01 @unit-level
  Scenario: No layer other than screen is charged for a route
    Given specs for a hook, a store and a piece of business logic, none declaring a route
    When the gate confronts each of them
    Then none is accused, because charging them was the legacy validator's vice and a
      gate that cries over what cannot be fixed teaches the team to ignore it

  @RTDCL-X02 @unit-level
  Scenario: The declared route is not confronted against the real router
    Given a screen spec naming a route that no routing table defines
    When the gate confronts it
    Then it returns Pass, because the ruler here is the spec's internal coherence and
      confronting the name against the router is a different confrontation
