# language: en
# @anchors
#   ref: TRCMT
#   updated_at: 2026-09-19
#   layer: feature

@TRCMT
Feature: TriadComplete — the pieces that realise a spec EXIST

  @TRCMT-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has a triad to demand

  @TRCMT-B02 @unit-level
  Scenario: A recognised layer leaves without a verdict
    Given a spec whose layer is declared with the declarative regime
    When the gate confronts it
    Then it returns Skip, because a recognised layer has neither spec nor triad by definition

  @TRCMT-B03 @unit-level
  Scenario: Without a map the verdict is undetermined
    Given a spec and no graph built
    When the gate confronts it
    Then it returns Pending, because approving without being able to look would assert
      what was never measured

  @TRCMT-B04 @unit-level
  Scenario: A spec with the three pieces linked passes
    Given a spec linked to its code, to its feature and to the test that proves it
    When the gate confronts it
    Then it returns Pass

  @TRCMT-B05 @unit-level
  Scenario: A spec missing a piece is failed, and the verdict names which
    Given a spec linked to its code and to nothing else
    When the gate confronts it
    Then it returns Fail
    And the verdict names the feature and the test, and says where each one is born

  @TRCMT-B06 @unit-level
  Scenario: The layer may waive a piece for every spec in it
    Given a layer declaring that the test edge is optional
    And a spec of that layer linked to its code and its feature only
    When the gate confronts it
    Then it returns Pass, because the waiver is declared in the Structure, in plain sight

  @TRCMT-B07 @unit-level
  Scenario: The unit may waive a piece in its own spec, with a written reason
    Given a spec carrying a waiver marker for the test, followed by the reason
    And the spec is linked to its code and its feature
    When the gate confronts it
    Then it returns Pass, because the decision belongs to the unit and is written where
      whoever reads the spec will see it

  @TRCMT-I01 @unit-level
  Scenario: The test is reached in two hops, through the feature
    Given a spec linked to a feature, and that feature linked to the test
    And no edge going straight from the spec to the test
    When the gate confronts it
    Then it returns Pass, because who points at the test is the FEATURE — checking it
      straight on the spec would report a missing test across the whole project

  @TRCMT-I02 @unit-level
  Scenario: Waiving the test while the feature carries a scenario is a contradiction
    Given a spec whose test is waived
    And a linked feature carrying one scenario
    When the gate confronts it
    Then it returns Fail, because either the scenario is real and someone must prove it,
      or it should not exist

  @TRCMT-I03 @unit-level
  Scenario: Waiving the test demands saying where the proof is, and the place must exist
    Given a spec whose test waiver points at a target that no file realises
    When the gate confronts it
    Then it returns Fail, because an orphan reference proves nothing

  @TRCMT-I04 @unit-level
  Scenario: A waiver covers only the piece it declares
    Given a spec that waives the feature and is linked to its code only
    When the gate confronts it
    Then it returns Fail naming the test, because waiving one piece never waives the others

  @TRCMT-X01 @unit-level
  Scenario: The gate does not confront whether the pieces MATCH one another
    Given a spec whose linked feature describes a behaviour the test does not prove
    And the three pieces exist and are linked
    When the gate confronts it
    Then it returns Pass, because matching is the work of the relational gates — this one
      exists precisely because they fail open when the piece is absent

  @TRCMT-B09 @unit-level
  Scenario: A per-rule waiver in a table row does not waive the unit
    Given a spec whose only waiver sits inside a rule's table row
    When the gate reads the unit's waivers
    Then no piece of the triad is waived

  @TRCMT-B08 @unit-level
  Scenario: A piece declared TO BE DEVELOPED leaves the verdict undetermined
    Given a spec declaring `@TBD` for a piece it has not written yet, with the reason
    When the gate confronts it
    Then it returns Pending and never Pass, because `@TBD` is DEBT while `@no-<piece>`
      is a permanent waiver — treating them alike erased the pending work from the radar
      for the honest declaration of whoever assumed it

  @TRCMT-X02 @unit-level
  Scenario: The gate does not judge the QUALITY of any piece
    Given a spec linked to a feature with no scenarios and to an empty test
    When the gate confronts it
    Then it returns Pass, because the ruler here is EXISTENCE — confronting the content
      belongs to another gate, and mixing the two would fail by a criterion this one
      cannot measure
