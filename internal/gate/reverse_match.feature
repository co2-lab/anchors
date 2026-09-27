# language: en
# @anchors
#   ref: RVMTR
#   updated_at: 2026-09-27
#   layer: feature

@RVMTR
Feature: ReverseMatch — every scenario still has its rule, and every proven code still has its scenario

  @RVMTR-B01 @unit-level
  Scenario: feature-spec-match skips a node that is not a feature, and a feature with no coded scenario
    Given a spec node, and in turn a feature of unit UNITX with no coded scenario
    When feature-spec-match confronts each
    Then both return Skip

  @RVMTR-B02 @unit-level
  Scenario: Without a map both gates are Pending
    Given the feature and the test of unit UNITX and no map
    When feature-spec-match and test-feature-match confront them
    Then both return Pending with the no-map message

  @RVMTR-B03 @unit-level
  Scenario: A feature no spec covers is Pending
    Given a map where no spec covers the feature of unit UNITX
    When feature-spec-match confronts it
    Then it returns Pending

  @RVMTR-B04 @unit-level
  Scenario: A covering spec that defines no requirement leaves feature-spec-match Pending
    Given the feature of unit UNITX covered by a spec with prose only
    When feature-spec-match confronts it
    Then it returns Pending, saying the spec defines no requirement

  @RVMTR-B05 @unit-level
  Scenario: A scenario whose rule the spec no longer declares fails, named as written
    Given a feature with scenarios UNITX-B01 and UNITX-B10, covered by a spec that defines only UNITX-B01, and in turn by one that defines both
    When feature-spec-match confronts it
    Then the first returns Fail naming UNITX-B10 and not UNITX-B01, and the second returns Pass

  @RVMTR-B06 @unit-level
  Scenario: A numbered variant is the same rule
    Given a spec defining UNITX-B01, and a feature with scenarios UNITX-B01#01 and UNITX-B01#02, and in turn one with UNITX-B10#01
    When feature-spec-match confronts each
    Then the variants of B01 pass, and the verdict for the other names UNITX-B10#01

  @RVMTR-B07 @unit-level
  Scenario: A data state is defined by its name
    Given a spec of unit TREXX defining the states DS-data-present and DS-data-empty in a table, and a feature citing those two and TREXX-DS-filter
    When feature-spec-match confronts it
    Then it returns Fail naming only TREXX-DS-filter

  @RVMTR-B08 @unit-level
  Scenario: The visual baseline is never charged
    Given a spec defining UNITX-B01, and a feature with UNITX-B01 and the visual baseline UNITX-VR
    When feature-spec-match confronts it
    Then it returns Pass

  @RVMTR-B09 @unit-level
  Scenario: The rules of every covering spec count together
    Given a feature with UNITX-B01 and UNITX-B10, covered by one spec defining B01 and another defining B10
    When feature-spec-match confronts it
    Then it returns Pass

  @RVMTR-B10 @unit-level
  Scenario: test-feature-match skips what is not a test, and a test no feature exercises is Pending
    Given a feature node, and in turn a test of unit UNITX that no feature exercises
    When test-feature-match confronts each
    Then the first returns Skip and the second Pending

  @RVMTR-B11 @unit-level
  Scenario: A linked feature with no coded scenario leaves test-feature-match Pending
    Given a test naming UNITX-B10, exercised by a feature whose only scenario has no code
    When test-feature-match confronts it
    Then it returns Pending, saying the feature declares no coded scenario

  @RVMTR-B12 @unit-level
  Scenario: A code the test names that no scenario declares fails
    Given a feature declaring UNITX-B01, and a test naming UNITX-B01 and UNITX-B10
    When test-feature-match confronts the test
    Then it returns Fail naming UNITX-B10

  @RVMTR-B13 @unit-level
  Scenario: A rule declared as a variant is a declared scenario for the test
    Given a feature declaring UNITX-B01#01, and a test naming UNITX-B01
    When test-feature-match confronts the test
    Then it returns Pass

  @RVMTR-B14 @unit-level
  Scenario: A revision code is not charged as a rule
    Given a feature declaring UNITX-B01, and a test naming UNITX-R0002 and UNITX-B01
    When test-feature-match confronts the test
    Then it returns Pass

  @RVMTR-E01 @unit-level
  Scenario: A linked spec or feature that cannot be read contributes nothing
    Given a feature covered by a spec file that does not exist, and a test exercised by a feature file that does not exist
    When each gate confronts its node
    Then both return Pending

  @RVMTR-I01 @unit-level
  Scenario: Neither gate answers Pass when it had nothing to match against
    Given no map, no linked origin, or an origin that declares nothing
    When each gate confronts its node
    Then every verdict is Pending

  @RVMTR-X01 @unit-level
  Scenario: Codes of another unit are never charged
    Given a feature of unit UNITX whose scenario also cites OTHER-B07, and a test of UNITX that builds a fixture with OTHER-B07
    When each gate confronts its node
    Then both return Pass

  @RVMTR-X02 @unit-level
  Scenario: A code named only in a comment is not a claim of proof
    Given a feature declaring UNITX-B01, and a test naming UNITX-B10 only in a comment
    When test-feature-match confronts the test
    Then it returns Pass

  @RVMTR-B15 @unit-level
  Scenario: A data state defined with the unit prefix is read at the code length the project declares
    Given a project that declares code length 7
    When a spec's table row defines "TREXXXX-DS-data-present"
    Then the spec defines exactly the data state "DS-data-present"

  @RVMTR-B16 @unit-level
  Scenario: A support file is not confronted as a test
    Given a support file of a test layer with no feature linked
    When test-feature-match judges it
    Then it skips saying the file is support

  @RVMTR-B17 @unit-level
  Scenario: A test of a declarative unit is not charged a feature
    Given tests with no feature linked, found beside their code by the derivation: one of a declarative util, one of a governed screen, and one with no unit
    When test-feature-match judges them
    Then the first is skipped saying its unit is declarative, and the others stay pending
