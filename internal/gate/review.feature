# language: en
# @anchors
#   ref: RVDUR
#   updated_at: 2026-10-02
#   layer: feature

@RVDUR
Feature: ReviewsDue — the targets of a reviewed gate that no review covers at their current revision

  @RVDUR-B01 @unit-level
  Scenario: The targets of a reviewed gate with no review at their revision are to review
    Given a gate on specs with review, and two specs, one reviewed at its revision
    When what is to review is listed
    Then only the other spec is listed, with the gate's review question

  @RVDUR-B02 @unit-level
  Scenario: A gate with no review has nothing to review, and a named gate lists its own
    Given a reviewed gate and a gate with no review over the same specs
    When what is to review is listed, and then for the reviewed gate alone, and for the other
    Then only the reviewed gate's targets are listed, and the other gate lists nothing

  @RVDUR-B03 @unit-level
  Scenario: A change to a reviewed target makes it to review again
    Given a spec reviewed at its revision
    When its revision changes
    Then it is to review again

  @RVDUR-I01 @unit-level
  Scenario: Listing what is to review changes no verdict
    Given a judged and reviewed gate over a spec
    When the gate runs before and after the spec is reviewed
    Then its judgment is asked both times

  @RVDUR-X01 @unit-level
  Scenario: The list never blocks
    Given a reviewed gate with targets to review
    When the list is built
    Then it carries no verdict and no blocking
