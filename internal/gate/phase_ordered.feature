# language: en
# @anchors
#   ref: PHORP
#   updated_at: 2026-09-19
#   layer: feature

@PHORP
Feature: PhaseOrdered — plan phases and phase dependencies must be ordered and consistent

  @PHORP-B01 @unit-level
  Scenario: The PlanPhases function extracts catalogued phase codes
    Given a plan containing catalogued phase headings
    When PlanPhases extracts the phase identifiers
    Then it returns the phase codes in the order they appear in the text

  @PHORP-B02 @unit-level
  Scenario: Non-plan artifacts skip phase ordering confrontation
    Given an artifact node whose kind is not plan
    When phase ordering confronts it
    Then it returns Skip, restricting internal phase checking to plans

  @PHORP-B03 @unit-level
  Scenario: Plans without phase headings skip confrontation
    Given a plan containing no phase headings or phase-like sections
    When phase ordering confronts it
    Then it returns Skip, treating phase structuring as optional for simple plans

  @PHORP-B04 @unit-level
  Scenario: Plans with phase-like sections lacking codes return Pending
    Given a plan dividing work into level-three sections without catalogued codes
    When phase ordering confronts it
    Then it returns Pending, reminding the author to give identity codes to phases

  @PHORP-B05 @unit-level
  Scenario: Plans declaring valid backward phase dependencies pass
    Given a plan whose phases depend only on preceding catalogued phases
    When phase ordering confronts it
    Then it returns Pass, confirming sequential prerequisite order

  @PHORP-B06 @unit-level
  Scenario: Plans with duplicate phase codes fail
    Given a plan containing two phases defined with the same identity code
    When phase ordering confronts it
    Then it returns Fail, rejecting ambiguous duplicate phase identifiers

  @PHORP-B07 @unit-level
  Scenario: Phases depending on uncatalogued phase codes fail
    Given a plan phase declaring a dependency on an unknown phase code
    When phase ordering confronts it
    Then it returns Fail, reporting the uncatalogued phase reference

  @PHORP-B08 @unit-level
  Scenario: Phases depending on themselves fail
    Given a plan phase declaring a dependency on its own code
    When phase ordering confronts it
    Then it returns Fail, rejecting self-dependency loops

  @PHORP-B09 @unit-level
  Scenario: Phases depending on future phases fail
    Given a plan phase declaring a dependency on a phase appearing later in the plan
    When phase ordering confronts it
    Then it returns Fail, rejecting forward prerequisite dependencies

  @PHORP-B10 @unit-level
  Scenario: Specifications declaring existing phase dependencies pass
    Given a specification declaring needs pointing to catalogued phases in plans
    When phase existence confronts it
    Then it returns Pass, confirming the phase dependencies exist

  @PHORP-B11 @unit-level
  Scenario: Specifications declaring missing phase dependencies fail
    Given a specification declaring needs pointing to a non-existent phase code
    When phase existence confronts it
    Then it returns Fail, naming the missing phase to unblock resolution

  @PHORP-B12 @unit-level
  Scenario: Artifacts declaring valid parents pass
    Given an artifact declaring an existing artifact code or phase as its parent
    When parent validation confronts it
    Then it returns Pass, confirming valid structural containment

  @PHORP-B13 @unit-level
  Scenario: Artifacts declaring invalid parents, self-parenting, or parent cycles fail
    Given artifacts declaring non-existent parents, self-parenting, or circular parent chains
    When parent validation confronts them
    Then it returns Fail, protecting tree traversal from broken containment

  @PHORP-I01 @unit-level
  Scenario: Phase dependencies are strictly acyclic and backward-directed
    Given plan phases with cyclical or forward-pointing prerequisite declarations
    When phase ordering confronts them
    Then it returns Fail, ensuring execution order flows forward

  @PHORP-I02 @unit-level
  Scenario: Phase and parent targets must exist in the map
    Given artifacts referencing phase or parent codes absent from the graph
    When the gates confront them
    Then they return Fail, preventing broken containment references

  @PHORP-I03 @unit-level
  Scenario: Parent chains are cycle-free and bounded
    Given an artifact whose parent hierarchy forms a circular loop
    When parent validation evaluates the chain
    Then it detects the cycle and returns Fail before infinite recursion occurs

  @PHORP-I04 @unit-level
  Scenario: Phase detection identifies level-three sections regardless of language
    Given plans formatted with level-three sections in any language
    When phase ordering evaluates section headers
    Then it detects structural phase candidates uniformly

  @PHORP-X01 @unit-level
  Scenario: Small plans are not required to catalog phases
    Given a concise single-phase implementation plan
    When phase ordering runs
    Then it skips confrontation without enforcing phase breakdown ceremony

  @PHORP-X02 @unit-level
  Scenario: Phase duration and calendar timing are not verified
    Given a plan defining phases with delivery date notes
    When phase ordering runs
    Then it validates dependency order without evaluating schedule duration

  @PHORP-X03 @unit-level
  Scenario: Both artifact codes and phase codes are accepted as parents
    Given artifacts declaring either artifact codes or phase codes as parent
    When parent validation runs
    Then it accepts both kinds of valid hierarchy roots symmetrically
