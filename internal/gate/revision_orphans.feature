# language: en
# @anchors
#   ref: RVORP
#   updated_at: 2026-09-22
#   layer: feature

@RVORP
Feature: RevisionOrphans — the rules a revision changed the meaning of, without saying so

  @RVORP-B01 @unit-level
  Scenario: Non-spec artifacts skip confrontation
    Given an artifact whose kind is not a spec
    When the gate confronts it
    Then it returns Skip, because a revision is charged of the document that carries it

  @RVORP-B02 @unit-level
  Scenario: A spec with no revision has nothing to confront
    Given a spec declaring no revision
    When the gate confronts it
    Then it returns Skip

  @RVORP-B03 @unit-level
  Scenario: A revision naming no rule cannot be confronted
    Given a spec whose revision declares no Revises field
    When the gate confronts it
    Then it fails, because a revision that names nothing cannot be checked against any sibling

  @RVORP-B04 @unit-level
  Scenario: A revision naming a rule the spec does not define fails
    Given a spec whose revision names a rule code absent from the document
    When the gate confronts it
    Then it fails, reporting the unknown code

  @RVORP-B05 @unit-level
  Scenario: A sibling sharing vocabulary and left unmentioned is reported
    Given a spec whose revision rewrote a rule about the badge
    And a sibling invariant that still speaks of the badge
    And that invariant appears in neither Revises nor Checked
    When the gate confronts it
    Then it fails, naming the invariant and the shared terms

  @RVORP-B06 @unit-level
  Scenario: Every vocabulary-sharing sibling accounted for passes
    Given a spec whose revision names every sibling that shares its vocabulary
    When the gate confronts it
    Then it passes

  @RVORP-B07 @unit-level
  Scenario: Checked clears the accusation without asserting correctness
    Given a spec whose revision declares Checked for a vocabulary-sharing sibling
    When the gate confronts it
    Then that sibling leaves the accusation

  @RVORP-I01 @unit-level
  Scenario: A revised rule is never its own orphan
    Given a spec whose revision names a rule in Revises
    When the gate confronts it
    Then that rule is absent from the orphan list

  @RVORP-X01 @unit-level
  Scenario: Terms shared by the whole unit do not discriminate
    Given a spec where every rule title carries the same domain word
    When the gate confronts it
    Then that word does not by itself make a rule an orphan
