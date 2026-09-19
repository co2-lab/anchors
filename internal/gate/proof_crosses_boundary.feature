# language: en
# @anchors
#   ref: PCBPR
#   updated_at: 2026-09-19
#   layer: feature

@PCBPR
Feature: ProofCrossesBoundary — when a rule claims a relation, the proof must reach the other side

  @PCBPR-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only a spec catalogues the rules that claim relations

  @PCBPR-B02 @unit-level
  Scenario: A marked rule whose governed code does not import the cited unit fails
    Given a rule carrying the single-source mark and citing another unit's file
    And the governed code never imports that file
    When the gate confronts the spec
    Then it returns Fail, because the claim is prose while the proof stays local

  @PCBPR-B03 @unit-level
  Scenario: A citation living only in a comment does not satisfy the charge
    Given a marked rule citing another unit's file
    And the governed code names that file twice, in comments, and imports nothing
    When the gate confronts the spec
    Then it returns Fail, because counting the comment would approve exactly the case
      that motivated the gate

  @PCBPR-B04 @unit-level
  Scenario: Governed code that imports the cited unit passes
    Given a marked rule citing another unit's file
    And the governed code imports that file
    When the gate confronts the spec
    Then it returns Pass

  @PCBPR-B05 @unit-level
  Scenario: An import through a different alias still matches
    Given a marked rule citing the unit by its full repository path
    And the governed code imports the same module through a project alias, with no extension
    When the gate confronts the spec
    Then it returns Pass, because alias and extension vary per project and the module
      base name does not

  @PCBPR-B06 @unit-level
  Scenario: A declared waiver with a written reason is not charged
    Given a rule whose relation is declared non-importable with a written reason
    When the gate confronts the spec
    Then it returns Skip, and the verdict records that the relation was waived

  @PCBPR-B07 @unit-level
  Scenario: The owner stamp does not demand an import
    Given a rule carrying both the single-source mark and the owner stamp
    When the gate confronts the spec
    Then it returns Pass, because the owner does not mirror anybody — it is the source

  @PCBPR-B08 @unit-level
  Scenario: A relation claimed in prose, with no mark, is reported as a suspicion
    Given a rule saying it mirrors another unit's file, with no declared mark
    When the gate confronts the spec
    Then the verdict names that rule and its target, so the convention is taught instead
      of the delivery being barred on first encounter

  @PCBPR-B09 @unit-level
  Scenario: The mark with no target on the line is reported
    Given a rule carrying the single-source mark and citing no unit and no rule code
    When the gate confronts the spec
    Then the verdict reports the rule, because the mark says "I mirror something" and
      with no something there is nothing to confront

  @PCBPR-B10 @unit-level
  Scenario: A target declared by rule code is resolved through the map
    Given a marked rule naming another unit by its rule code instead of a path
    And the map resolves that code to a spec governing a code file
    And that code file is not imported by the governed code
    When the gate confronts the spec
    Then it returns Fail naming the resolved file, because the code is stable identity
      and the path is not

  @PCBPR-B11 @unit-level
  Scenario: A rule code that resolves to nothing is reported as unresolved
    Given a marked rule naming a rule code no unit of the map carries
    When the gate confronts the spec
    Then the verdict reports the unresolved target, because a rule pointing at nothing
      warns nobody

  @PCBPR-B12 @unit-level
  Scenario: A rule code of the unit itself is not a target
    Given a marked rule citing a rule code of its own unit
    When the gate confronts the spec
    Then no import is charged for it, because self-reference is not the other side of
      a boundary

  @PCBPR-B13 @unit-level
  Scenario: A rule with no relation claim leaves without a verdict
    Given a catalogued rule describing a local behaviour and naming no other unit
    When the gate confronts the spec
    Then it returns Skip, because there was nothing to charge

  @PCBPR-B14 @unit-level
  Scenario: Imports in other language shapes satisfy the charge
    Given a marked rule citing another unit's file
    And the governed code reaches it through a require, a from, a use or an include line
    When the gate confronts the spec
    Then it returns Pass, because the ruler is the import line and not one ecosystem's
      spelling of it

  @PCBPR-I01 @unit-level
  Scenario: The claim and the target must be on the same rule line
    Given a spec whose relation claim lives in a paragraph outside the rule table
    And the catalogued rules themselves name no other unit
    When the gate confronts the spec
    Then nothing is charged, because binding the demand to surrounding prose would bind
      it to no rule at all

  @PCBPR-I02 @unit-level
  Scenario: The import is charged on the governed code, never on the test
    Given a marked rule citing another unit's file
    And the test file imports it while the governed code does not
    When the gate confronts the spec
    Then it returns Fail, because a test may legitimately reach the function without
      importing the unit — what matters is the unit having no copy

  @PCBPR-I03 @unit-level
  Scenario: A satisfied charge does not swallow a pending suspicion
    Given one marked rule whose import is honoured
    And another rule claiming a relation in prose, with no mark
    When the gate confronts the spec
    Then the verdict still reports the prose rule, because a file may carry one rule
      that passed and another that nobody confronts

  @PCBPR-X01 @unit-level
  Scenario: The gate does not hunt duplicated concepts across the project
    Given two units defining the same table of values and no spec declaring the relation
    When the gate confronts either spec
    Then it returns Skip, because sweeping the cartesian product of the layers would
      demand judgement — the gate stops duplication from coming back once named

  @PCBPR-X02 @unit-level
  Scenario: The gate does not compare the values on the two sides
    Given a marked rule whose governed code imports the cited unit
    And the two sides hold values that disagree today
    When the gate confronts the spec
    Then it returns Pass, because the ruler is whether the proof reaches across —
      the import is what makes a divergence impossible to keep silent

  @PCBPR-X03 @unit-level
  Scenario: The gate does not decide whether an unmarked prose claim blocks
    Given a rule claiming a relation in prose only
    When the gate confronts the spec
    Then it emits the finding and leaves the decision to block to the gate's blocking
      setting in the project's Structure
