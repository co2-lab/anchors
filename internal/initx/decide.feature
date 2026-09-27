# language: en
# @anchors
#   ref: INDCN
#   updated_at: 2026-09-26
#   layer: feature

@INDCN
Feature: InitDecisions — the pure decisions of init over the proposed configuration: code layers, tags and governs rules

  @INDCN-B01 @unit-level
  Scenario: The code layer names are listed sorted, and only code layers
    Given a configuration with spec, guide, mobile-code, backend-code and generated layers
    When the code layer names are listed
    Then the list is backend-code, generated, mobile-code

  @INDCN-B02 @unit-level
  Scenario: Pruning removes the code layers not kept and never an artifact layer
    Given the same configuration
    And the user keeps mobile-code and backend-code
    When the code layers are pruned
    Then generated is gone
    And mobile-code, backend-code, spec and guide remain

  @INDCN-B03 @unit-level
  Scenario: The candidate tags are the union of every layer's tags, deduplicated and sorted
    Given the same configuration
    When the tags are collected
    Then the tags are backend, frontend, generated, guide, mobile, spec

  @INDCN-B04 @unit-level
  Scenario: One governs rule per guide answered with a tag, ordered by guide, skipping the unanswered and none
    Given answers for four guides: frontend, backend, none and no answer
    When the governs rules are built
    Then there are two rules, the backend guide first and the frontend guide second

  @INDCN-B05 @unit-level
  Scenario: No answer gives no governs rule
    Given no governs answer
    When the governs rules are built
    Then there is no rule

  @INDCN-I01 @unit-level
  Scenario: After pruning, the code layer names are exactly the kept code layers
    Given the same configuration
    And a keep set of mobile-code, spec and a layer that does not exist
    When the code layers are pruned and then listed
    Then the list is exactly mobile-code
