# language: en
# @anchors
#   ref: C4CNC
#   updated_at: 2026-09-26
#   layer: feature

@C4CNC
Feature: C4Containers — the declared containers, each with the units that run in it, for the architecture page

  @C4CNC-B01 @unit-level
  Scenario: A container carries the specs of the layers it declares, in layer then code order
    Given a container "app" declaring the layers " SHARED " and "lambdas"
    And specs ROTAX in lambdas, UTILY and UTILX in shared
    When the containers are listed
    Then "app" carries ROTAX, UTILX, UTILY in that order

  @C4CNC-B02 @unit-level
  Scenario: An external container gets no level 3
    Given the containers app, api and an external database
    When the internal containers are listed
    Then all three are containers, and only app and api are internal

  @C4CNC-B03 @unit-level
  Scenario: A layer no container declares is named as orphan
    Given specs in screen, lambdas, shared and orfa, and containers declaring the first three
    When the orphan layers are asked
    Then the answer is exactly "orfa"

  @C4CNC-B04 @unit-level
  Scenario: No configuration means no containers and no orphans
    Given a compiler with no project configuration
    When the containers, internal containers and orphan layers are asked
    Then all three are empty

  @C4CNC-I01 @unit-level
  Scenario: Every spec layer is held by a container or named orphan
    Given containers declaring some of the project's layers
    When the containers' units and the orphan layers are collected
    Then every layer that has specs appears in one of the two

  @C4CNC-X01 @unit-level
  Scenario: The containers are only the declared ones
    Given a configuration declaring three containers over four layers of specs
    When the containers are listed
    Then exactly the three declared containers come back, by their declared names
