# language: en
# @anchors
#   ref: PSDPL
#   updated_at: 2026-09-19
#   layer: feature

@PSDPL
Feature: PlanSourceDeclared — a plan that names a source has to declare who builds it

  @PSDPL-B01 @unit-level
  Scenario: A source that lives only in the prose is failed
    Given a plan whose prose names a source in bold
    And another plan seeds that source's adapter
    And this plan's needs list does not name that other plan
    When the gate confronts it
    Then it returns Fail, because the dependency existed only in the prose and survived
      the day the other plan dropped the adapter

  @PSDPL-B02 @unit-level
  Scenario: The verdict names which source and where its adapter lives
    Given a plan naming a source whose adapter another plan seeds
    And that other plan absent from the needs list
    When the gate confronts it
    Then the verdict carries the source name and the plan that owns the adapter, so the
      fix is a declaration rather than an investigation

  @PSDPL-B03 @unit-level
  Scenario: With the owning plan declared in needs the gate passes
    Given a plan naming a source whose adapter another plan seeds
    And that other plan listed in this plan's needs
    When the gate confronts it
    Then it returns Pass, because the prose and the declaration now say the same thing

  @PSDPL-B04 @unit-level
  Scenario: Every source of the line is confronted on its own
    Given one source line naming two sources in bold
    And both adapters seeded by a plan this one does not declare
    When the gate confronts it
    Then it returns Fail naming both, because one declared source does not cover the rest

  @PSDPL-B05 @unit-level
  Scenario: The source name matches the adapter regardless of case
    Given a plan naming the source in lower case
    And the seeded adapter file spelling it in mixed case
    When the gate confronts it
    Then the two are matched and the undeclared dependency is charged

  @PSDPL-B06 @unit-level
  Scenario: A source whose adapter nobody seeds is not charged
    Given a plan naming a source
    And no plan in the map seeds an adapter for it
    When the gate confronts it
    Then it does not fail, because the source may belong to a plan that does not exist
      yet and the gate cannot invent a dependency

  @PSDPL-B07 @unit-level
  Scenario: The plan that seeds the adapter is not charged for itself
    Given a plan that names a source and itself seeds that source's adapter
    When the gate confronts it
    Then it does not fail, because a plan does not depend on itself

  @PSDPL-B08 @unit-level
  Scenario: A plan with no source line returns Skip
    Given a plan whose prose names no source at all
    When the gate confronts it
    Then it returns Skip, because there is nothing to confront and that is not approval

  @PSDPL-B09 @unit-level
  Scenario: An artifact that is not a plan returns Skip
    Given a spec whose text carries a source line
    When the gate confronts it
    Then it returns Skip, because the gate has jurisdiction over plans only

  @PSDPL-I01 @unit-level
  Scenario: What was not measured is never approved
    Given in turn a plan confronted with no graph built, and one whose map seeds no
      adapter at all
    When the gate confronts each of them
    Then neither returns Pass, because approving there would stamp a confrontation that
      never happened

  @PSDPL-I02 @unit-level
  Scenario: A seeded file off the naming pattern owns nothing
    Given a plan that seeds a file whose name does not end in the adapter suffix
    And another plan naming that same source in bold
    When the gate confronts the consumer
    Then nothing is charged, because a wrong accusation costs more than a missed one —
      it teaches the reader to ignore the gate

  @PSDPL-X01 @unit-level
  Scenario: The gate does not confront the order of the phases
    Given a plan that declares in needs the plan building its adapter
    And that owning plan scheduled in a later phase than this one
    When the gate confronts it
    Then it returns Pass, because ordering is the ruler of another gate and holding it
      in two places would let the two diverge

  @PSDPL-X02 @unit-level
  Scenario: The gate does not demand a needs pointing at nothing
    Given a plan naming a source no plan in the map builds
    When the gate confronts it
    Then it does not fail, because charging it would demand a declaration pointing at
      nothing — the gate would be asking for a lie instead of catching one

  @PSDPL-X03 @unit-level
  Scenario: The gate does not interpret what the source is for
    Given a plan whose prose names a source in bold only to say it was ruled out
    And another plan seeds that source's adapter
    When the gate confronts it
    Then it still fails, because the ruler is the bold name on the source line —
      deciding whether the plan really consumes it is interpretation, and interpretation
      is not what a blocking gate can hold
