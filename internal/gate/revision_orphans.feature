# language: en
# @anchors
#   ref: RVORP
#   updated_at: 2026-09-30
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

  @RVORP-B08 @unit-level
  Scenario: A rule's title is its heading, not its usage row
    Given a spec whose rules read the same field in a usage table, and a revision of one of them
    When the revision is confronted
    Then the sibling reading the same field is not reported, since the titles are the headings

  @RVORP-B09 @unit-level
  Scenario: A revision that revised no rule says so
    Given a revision declaring it revised no rule with its reason, one declaring it with no reason, and one whose reason cites a code
    When each spec is confronted
    Then the first and the third pass, and the second is pending
