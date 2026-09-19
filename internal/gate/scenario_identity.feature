# language: en
# @anchors
#   ref: SCIDS
#   updated_at: 2026-09-19
#   layer: feature

@SCIDS
Feature: ScenarioIdentity — two scenarios of the same feature cannot share one code

  @SCIDS-B01 @unit-level
  Scenario: Two scenarios sharing one code are reported
    Given a feature whose happy path and alternative carry the same code
    When the gate confronts it
    Then it reports the repetition, because nothing links one of them to one test and the
      relational gates compare both titles against a single test

  @SCIDS-B02 @unit-level
  Scenario: The report names the repeated code
    Given a feature carrying one repeated scenario code among many distinct ones
    When the gate confronts it
    Then the message carries that code, so the reader does not have to scan the feature to
      find it

  @SCIDS-B03 @unit-level
  Scenario: The report says how many scenarios share the code
    Given a feature where three scenarios carry the same code
    When the gate confronts it
    Then the message says how many they are, which separates an accidental duplicate from a
      code borrowed across a whole rule

  @SCIDS-B04 @unit-level
  Scenario: The report teaches the way out with the project's own code
    Given a feature carrying a repeated scenario code
    When the gate confronts it
    Then the message shows the suffix notation applied to that very code, not to a generic
      placeholder

  @SCIDS-B05 @unit-level
  Scenario: The suffix gives each scenario its own identity
    Given two scenarios of one rule numbered with distinct suffixes
    When the gate confronts the feature
    Then it returns Pass, because numbered scenarios pair one-to-one with their tests

  @SCIDS-B06 @unit-level
  Scenario: Distinct codes pass
    Given a feature whose scenarios all carry codes of their own
    When the gate confronts it
    Then it returns Pass, because a rule with several scenarios is the common case

  @SCIDS-B07 @unit-level
  Scenario: The verdict is Pending and never a failure
    Given a feature carrying repeated scenario codes
    When the gate confronts it
    Then the verdict is Pending, because numbering is a migration and the gate is born over
      a base that did not know the notation

  @SCIDS-B08 @unit-level
  Scenario: An artifact that is not a feature leaves without a verdict
    Given a node whose kind is spec, code or test
    When the gate confronts it
    Then it returns Skip, because only a feature carries scenarios

  @SCIDS-B09 @unit-level
  Scenario: A feature with no coded scenario leaves without a verdict
    Given a feature whose scenarios carry no code at all
    When the gate confronts it
    Then it returns Skip, because there is nothing to confront and that absence is another
      gate's charge

  @SCIDS-B10 @unit-level
  Scenario: Several repeated codes are reported together in a stable order
    Given a feature carrying three different repeated scenario codes
    When the gate confronts it
    Then all three appear in one message, in an order that does not change between runs

  @SCIDS-B11 @unit-level
  Scenario: Long titles are shortened in the report
    Given a feature whose repeated scenarios carry very long titles
    When the gate confronts it
    Then the titles appear truncated, because the address is the code and the title only
      helps recognise which scenario is which

  @SCIDS-I01 @unit-level
  Scenario: Grouping is by the complete code, suffix included
    Given two scenarios sharing a code prefix and differing only in their suffixes
    When the gate confronts the feature
    Then it returns Pass, because grouping by the prefix alone would accuse exactly the
      projects that already did the migration the gate asks for

  @SCIDS-I02 @unit-level
  Scenario: The message is deterministic
    Given a feature carrying several repeated scenario codes
    When the gate confronts it twice
    Then both messages are identical, because the findings are ordered before being joined
      and the report must not churn between runs

  @SCIDS-X01 @unit-level
  Scenario: The gate does not judge whether the two scenarios describe different behaviours
    Given two scenarios sharing one code and describing plainly different behaviours
    When the gate confronts the feature
    Then the report says only that the code is repeated, because the ruler here is identity
      — one code, one scenario — and that is decidable by reading

  @SCIDS-X02 @unit-level
  Scenario: The gate does not look across features
    Given a feature whose codes are all distinct within it
    And another feature elsewhere reusing one of those codes
    When the gate confronts the first feature
    Then it returns Pass, because a code repeated across features is a different defect and
      seeing it would need a graph this gate does not take

  @SCIDS-X03 @unit-level
  Scenario: The gate does not charge the absence of a code on a scenario
    Given a feature mixing coded scenarios with scenarios carrying no code
    And every code present appearing exactly once
    When the gate confronts it
    Then it returns Pass, because charging the absence is the ruler of the gate that pairs
      scenarios to tests

  @SCIDS-X04 @unit-level
  Scenario: The gate does not renumber the scenarios
    Given a feature carrying a repeated scenario code
    When the gate confronts it
    Then the feature text is returned untouched, because the suffix carries meaning and a
      rewrite would invalidate every test already bound to the old code
