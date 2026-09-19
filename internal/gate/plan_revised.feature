# language: en
# @anchors
#   ref: PLRVP
#   updated_at: 2026-09-19
#   layer: feature

@PLRVP
Feature: PlanRevised — mutual revision visibility between superseded and revising plans

  @PLRVP-B01 @unit-level
  Scenario: Non-plan artifacts skip confrontation
    Given an artifact whose kind is not a plan
    When the gate confronts it
    Then it returns Skip, avoiding confrontation on non-plan artifacts

  @PLRVP-B02 @unit-level
  Scenario: Confronting without a map graph returns pending
    Given a plan artifact and a nil map graph
    When the gate confronts it
    Then it returns Pending, because revision relationships cannot be determined without the map

  @PLRVP-B03 @unit-level
  Scenario: A plan with neither revisions nor revisers skips confrontation
    Given a plan that does not revise any plan and is not revised by any other plan
    When the gate confronts it
    Then it returns Skip, because there are no revision relationships to enforce

  @PLRVP-B04 @unit-level
  Scenario: Declaring a revision target that does not exist in the map fails
    Given a revising plan that declares a revision target not present in the map graph
    When the gate confronts it
    Then it returns Fail, citing the absent target plan

  @PLRVP-B05 @unit-level
  Scenario: A revising plan receives a pending reminder when the target lacks a top notice
    Given a revising plan that declares an existing target plan
    And the target plan file does not contain a top revision notice
    When the gate confronts the revising plan
    Then it returns Pending, reminding the author to write the revision notice on the target plan

  @PLRVP-B06 @unit-level
  Scenario: The pending reminder on a revising plan clears once the target carries the notice
    Given a revising plan that declares an existing target plan
    And the target plan file contains the required top revision notice
    When the gate confronts the revising plan
    Then it returns Skip, clearing the pending reminder

  @PLRVP-B07 @unit-level
  Scenario: A revised plan lacking a top revision notice fails
    Given a plan that is revised by another plan in the map graph
    And the plan content contains no top revision notice
    When the gate confronts the revised plan
    Then it returns Fail, reporting the revising plan identifier

  @PLRVP-B08 @unit-level
  Scenario: A revised plan placing the revision notice after line 40 fails
    Given a plan that is revised by another plan in the map graph
    And the top revision notice appears after the first 40 lines of content
    When the gate confronts the revised plan
    Then it returns Fail, because the notice was placed too late to inform top-down readers

  @PLRVP-B09 @unit-level
  Scenario: A revised plan with top notice but no section amendment markers returns pending
    Given a plan that is revised by another plan in the map graph
    And the plan contains a top revision notice within the first 40 lines
    And the plan contains no section amendment markers
    When the gate confronts the revised plan
    Then it returns Pending, indicating that specific affected sections must be marked

  @PLRVP-B10 @unit-level
  Scenario: A revised plan with top notice and marked section amendments passes
    Given a plan that is revised by another plan in the map graph
    And the plan contains a top revision notice in the first 40 lines
    And the plan contains section amendment markers in affected sections
    When the gate confronts the revised plan
    Then it returns Pass, confirming mutual revision visibility

  @PLRVP-B11 @unit-level
  Scenario: Markdown alerts and metadata directives are both accepted as valid markers
    Given revised plan content using markdown callouts or metadata revision directives
    When the gate confronts it
    Then it accepts both forms without enforcing a single syntax convention

  @PLRVP-I01 @unit-level
  Scenario: Revision notices must be placed within the first 40 lines
    Given a revised plan where the revision notice is delayed beyond line 40
    When the gate confronts it
    Then it returns Fail, upholding the invariant that top-down readers encounter warnings first

  @PLRVP-I02 @unit-level
  Scenario: Missing section amendment markers yield pending rather than failure
    Given a revised plan having a top notice but lacking section-level markings
    When the gate confronts it
    Then it returns Pending instead of Fail, allowing whole-plan revisions to proceed

  @PLRVP-I03 @unit-level
  Scenario: The pending reminder clears once the revised plan is notified
    Given a revising plan whose target plan file receives the top revision notice
    When the gate confronts the revising plan
    Then it clears the pending reminder, preventing permanent noise

  @PLRVP-X01 @unit-level
  Scenario: Language neutrality allows markdown alerts and metadata directives
    Given revision notices written using universal callout syntax or machine directives
    When the gate confronts the document
    Then it evaluates marker structure without demanding specific spoken languages

  @PLRVP-X02 @unit-level
  Scenario: Prose description quality accompanying revision markers is not evaluated
    Given revision markers accompanied by brief or varied prose text
    When the gate confronts the document
    Then it verifies marker presence without grading prose completeness

  @PLRVP-X03 @unit-level
  Scenario: Section amendment markers are not demanded on unrevised plans
    Given a plan that has no revising plans targeting it
    When the gate confronts it
    Then it skips section amendment checks
