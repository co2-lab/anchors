# language: en
# @anchors
#   ref: UNTCP
#   updated_at: 2026-10-01
#   layer: feature

@UNTCP
Feature: UnitComplete — the pieces that realise a spec EXIST

  @UNTCP-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has a unit to demand

  @UNTCP-B02 @unit-level
  Scenario: A recognised layer leaves without a verdict
    Given a spec whose layer is declared with the declarative regime
    When the gate confronts it
    Then it returns Skip, because a recognised layer has neither spec nor unit by definition

  @UNTCP-B03 @unit-level
  Scenario: Without a map the verdict is undetermined
    Given a spec and no graph built
    When the gate confronts it
    Then it returns Pending, because approving without being able to look would assert
      what was never measured

  @UNTCP-B04 @unit-level
  Scenario: A spec with the three pieces linked passes
    Given a spec linked to its code, to its feature and to the test that proves it
    When the gate confronts it
    Then it returns Pass

  @UNTCP-B05 @unit-level
  Scenario: A spec missing a piece is failed, and the verdict names which
    Given a spec linked to its code and to nothing else
    When the gate confronts it
    Then it returns Fail
    And the verdict names the feature and the test, and says where each one is born

  @UNTCP-B06 @unit-level
  Scenario: The layer may waive a piece for every spec in it
    Given a layer declaring that the test edge is optional
    And a spec of that layer linked to its code and its feature only
    When the gate confronts it
    Then it returns Pass, because the waiver is declared in the Structure, in plain sight

  @UNTCP-B07 @unit-level
  Scenario: The unit may waive a piece in its own spec, with a written reason
    Given a spec carrying a waiver marker for the test, followed by the reason
    And the spec is linked to its code and its feature
    When the gate confronts it
    Then it returns Pass, because the decision belongs to the unit and is written where
      whoever reads the spec will see it

  @UNTCP-I01 @unit-level
  Scenario: The test is reached in two hops, through the feature
    Given a spec linked to a feature, and that feature linked to the test
    And no edge going straight from the spec to the test
    When the gate confronts it
    Then it returns Pass, because who points at the test is the FEATURE — checking it
      straight on the spec would report a missing test across the whole project

  @UNTCP-I02 @unit-level
  Scenario: Waiving the test while the feature carries a scenario is a contradiction
    Given a spec whose test is waived
    And a linked feature carrying one scenario
    When the gate confronts it
    Then it returns Fail, because either the scenario is real and someone must prove it,
      or it should not exist

  @UNTCP-I03 @unit-level
  Scenario: Waiving the test demands saying where the proof is, and the place must exist
    Given a spec whose test waiver points at a target that no file realises
    When the gate confronts it
    Then it returns Fail, because an orphan reference proves nothing

  @UNTCP-I04 @unit-level
  Scenario: A waiver covers only the piece it declares
    Given a spec that waives the feature and is linked to its code only
    When the gate confronts it
    Then it returns Fail naming the test, because waiving one piece never waives the others

  @UNTCP-X01 @unit-level
  Scenario: The gate does not confront whether the pieces MATCH one another
    Given a spec whose linked feature describes a behaviour the test does not prove
    And the three pieces exist and are linked
    When the gate confronts it
    Then it returns Pass, because matching is the work of the relational gates — this one
      exists precisely because they fail open when the piece is absent

  @UNTCP-B09 @unit-level
  Scenario: A per-rule waiver in a table row does not waive the unit
    Given a spec whose only waiver sits inside a rule's table row
    When the gate reads the unit's waivers
    Then no piece of the unit is waived

  @UNTCP-B08 @unit-level
  Scenario: A piece declared TO BE DEVELOPED leaves the verdict undetermined
    Given a spec declaring `@TBD` for a piece it has not written yet, with the reason
    When the gate confronts it
    Then it returns a divergence and never Pass, because `@TBD` is DEBT while `@no-<piece>`
      is a permanent waiver — treating them alike erased the divergence work from the radar
      for the honest declaration of whoever assumed it

  @UNTCP-X02 @unit-level
  Scenario: The gate does not judge the QUALITY of any piece
    Given a spec linked to a feature with no scenarios and to an empty test
    When the gate confronts it
    Then it returns Pass, because the ruler here is EXISTENCE — confronting the content
      belongs to another gate, and mixing the two would fail by a criterion this one
      cannot measure

  @UNTCP-E01 @unit-level
  Scenario: A test missing from disk does not orphan a reference another test resolves
    Given the map lists a test that is gone from disk, and a test on disk citing the code the @no-test points at
    When the gate confronts the spec
    Then it returns Pass

  @UNTCP-E02 @unit-level
  Scenario: A feature missing from disk does not hide the scenario of the covered feature
    Given the spec is covered by a feature that is gone from disk and by a feature with a scenario
    And the spec waives the test with @no-test
    When the gate confronts the spec
    Then it fails, naming the feature whose scenario contradicts the waiver
