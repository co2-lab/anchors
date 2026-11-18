# language: en
# @anchors
#   code: VRPRF
#   ref: VRPRS
#   updated_at: 2026-11-10
#   layer: feature

@VRPRS
Feature: VariantsProven — a rule that names its variants is proven once per variant

  @VRPRS-B01 @unit-level
  Scenario: A variant no scenario of its rule names fails, naming it
    Given a rule naming the variants audit and notification, and one scenario of it tagged for the audit
    When the gate confronts the spec, and again with an outline whose Examples name the notification
    Then the first fails naming the notification, and the second passes

  @VRPRS-B02 @unit-level
  Scenario: A rule that follows another's variants diverges while one of them is unproven
    Given an error rule following the variants of a rule proven for both, and proven itself only for the audit
    When the gate confronts the spec
    Then it diverges, saying the error rule is proven only without the notification

  @VRPRS-B03 @unit-level
  Scenario: A spec naming no variants, with no feature, or no spec at all, is skipped
    Given a spec whose rules name no variants, one with variants and no feature, and a code file
    When the gate confronts each
    Then each is skipped

  @VRPRS-E01 @unit-level
  Scenario: A rule following a rule that names no variants fails
    Given a rule following the variants of a rule that declares none
    When the gate confronts the spec
    Then it fails, naming both

  @VRPRS-E02 @unit-level
  Scenario: Without a map the gate is pending
    Given a spec with variants and no built map
    When the gate confronts it
    Then it is pending
