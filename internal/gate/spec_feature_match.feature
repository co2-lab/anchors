# language: en
# @anchors
#   ref: SFMSP
#   updated_at: 2026-09-26
#   layer: feature

@SFMSP
Feature: SpecFeatureMatch — every requirement the spec defines has at least one scenario

  @SFMSP-B01 @unit-level
  Scenario: A requirement no scenario tags is failed and named
    Given a spec defining two requirements
    And a linked feature carrying a scenario for only one of them
    When the gate confronts it
    Then it returns Fail naming the uncovered requirement, because otherwise it would
      cross the whole pipeline with nothing verifying it

  @SFMSP-B02 @unit-level
  Scenario: A requirement that has a scenario is not accused
    Given a spec defining one covered requirement and one uncovered requirement
    When the gate confronts it
    Then the verdict carries the uncovered one and never the covered one

  @SFMSP-B03 @unit-level
  Scenario: With every requirement tagged the gate passes
    Given a spec defining two requirements
    And a linked feature carrying one scenario tagged for each
    When the gate confronts it
    Then it returns Pass

  @SFMSP-B04 @unit-level
  Scenario: A code merely cited contracts no obligation
    Given a spec defining one requirement of its own
    And prose and a Dependency Table citing codes of other units
    And a feature covering only its own requirement
    When the gate confronts it
    Then it returns Pass, because a spec cites other units' codes all the time and
      without that distinction the gate would be a noise generator

  @SFMSP-B05 @unit-level
  Scenario: A per-requirement waiver with a written reason waives
    Given a spec whose second requirement carries the waiver marker followed by a reason
    And a feature covering only the first requirement
    When the gate confronts it
    Then it returns Pass, and the reason stays in the spec as the record that it was a
      decision rather than forgetfulness

  @SFMSP-B06 @unit-level
  Scenario: A bare per-requirement waiver does not waive
    Given a spec whose second requirement carries the waiver marker with nothing after it
    And a feature covering only the first requirement
    When the gate confronts it
    Then it returns Fail, because the waiver requires a written reason

  @SFMSP-B07 @unit-level
  Scenario: A whole-spec waiver drags the waiver to every requirement
    Given a spec declaring with a reason that it has no feature
    And two requirements neither of which carries a waiver of its own
    When the gate confronts it
    Then it returns Skip, because with no feature no requirement of it can have a
      scenario — one decision, one place

  @SFMSP-B08 @unit-level
  Scenario: Without the whole-spec waiver the same requirements keep failing
    Given the same two requirements and the same empty feature
    And no whole-spec waiver anywhere in the spec
    When the gate confronts it
    Then it returns Fail, because the drag cannot become a silent way of muting the gate

  @SFMSP-B09 @unit-level
  Scenario: A bare whole-spec waiver drags nothing
    Given a spec whose whole-spec waiver marker has nothing written after it
    And a feature carrying no scenario
    When the gate confronts it
    Then it returns Fail, because otherwise the marker would be a switch that turns the
      gate off without accounting for it

  @SFMSP-B10 @unit-level
  Scenario: A spec with no feature returns Skip
    Given a spec defining a requirement and no feature linked to it
    When the gate confronts it
    Then it returns Skip, because that absence is the ruler of the triad gate and
      accusing it here would print the same defect twice

  @SFMSP-B11 @unit-level
  Scenario: Requirements are looked for across every linked feature
    Given a spec covered by two features
    And each feature carrying the scenario of a different requirement
    When the gate confronts it
    Then it returns Pass, because the requirement only needs to be in some of them

  @SFMSP-B12 @unit-level
  Scenario: An artifact that is not a spec returns Skip
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only a spec defines requirements

  @SFMSP-B13 @unit-level
  Scenario: A rule alias needs no scenario of its own
    Given a spec whose failure rule is written as REF to a behaviour rule with a reason
    And a feature with a scenario for the behaviour rule only
    When the gate confronts the spec
    Then it returns Pass, and the same row without the alias is charged as uncovered

  @SFMSP-B14 @unit-level
  Scenario: An alias that stands for no rule fails
    Given aliases pointing at an undefined rule, at another alias, or with no reason
    When the gate confronts the spec
    Then it returns Fail naming the alias, even when the spec has no feature

  @SFMSP-I01 @unit-level
  Scenario: Every waiver requires a written reason
    Given in turn a bare per-requirement marker and a bare whole-spec marker
    When the gate confronts each of them
    Then both still fail, because a bare marker is a switch with no accounting and
      silence without a why is what the gate exists to end

  @SFMSP-I02 @unit-level
  Scenario: Each gate accuses one thing
    Given a spec with a defined requirement and no feature at all
    When the gate confronts it
    Then it skips rather than failing, because the missing feature belongs to the triad
      gate and reporting it here would print the same defect twice

  @SFMSP-X01 @unit-level
  Scenario: The gate does not judge whether the scenario proves the requirement
    Given a spec defining one requirement
    And a feature whose scenario carries the tag and asserts nothing at all
    When the gate confronts it
    Then it returns Pass, because the tag is deterministic and the shape of the
      assertion is the ruler of another gate

  @SFMSP-X02 @unit-level
  Scenario: A cited code produces no accusation
    Given a spec whose Dependency Table cites three codes of other units
    And a feature covering only the code this spec defines
    When the gate confronts it
    Then none of the cited codes appears in the verdict, because a gate that cries wolf
      gets switched off — which costs more than the defect it was catching

  @SFMSP-X03 @unit-level
  Scenario: The gate does not confront feature against test
    Given a spec whose every requirement has a scenario
    And no test binding any of those scenarios
    When the gate confronts it
    Then it returns Pass, because that edge already has its own watcher and this gate
      exists for the edge before it, which had none

  @SFMSP-E01 @unit-level
  Scenario: Without a built map the confrontation is pending
    Given a spec confronted with no map built
    When the gate confronts it
    Then it returns Pending with the no-map message

  @SFMSP-E02 @unit-level
  Scenario: A feature missing from disk does not hide the scenarios of the other features
    Given a spec linked to a feature that is gone from disk and to a feature that tags its requirement
    When the gate confronts it
    Then it returns Pass
