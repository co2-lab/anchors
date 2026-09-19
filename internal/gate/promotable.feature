# language: en
# @anchors
#   ref: PRGTP
#   updated_at: 2026-09-19
#   layer: feature

@PRGTP
Feature: PromotableGates — identifies clean informative gates ready for promotion to blocking

  @PRGTP-B01 @unit-level
  Scenario: Returns an empty list of Promotable gates when profile contains no gates
    Given an evaluation profile with an empty gate summary map
    When PromotableGates is invoked
    Then it returns an empty slice of promotable gates

  @PRGTP-B02 @unit-level
  Scenario: An informative gate with passes and zero failures is included
    Given an evaluation profile with an informative gate having positive passes and zero failures
    When PromotableGates is invoked
    Then the gate is returned in the promotable list

  @PRGTP-B03 @unit-level
  Scenario: An informative gate with failures is excluded from promotion
    Given an evaluation profile with an informative gate having one or more failures
    When PromotableGates is invoked
    Then the gate is excluded from the promotable list

  @PRGTP-B04 @unit-level
  Scenario: An informative gate with zero passes is excluded as having no data
    Given an evaluation profile with an informative gate having zero passes and zero failures
    When PromotableGates is invoked
    Then the gate is excluded from the promotable list

  @PRGTP-B05 @unit-level
  Scenario: A blocking gate is excluded from promotion suggestions
    Given an evaluation profile with a blocking gate having positive passes and zero failures
    When PromotableGates is invoked
    Then the blocking gate is excluded from the promotable list

  @PRGTP-B06 @unit-level
  Scenario: Returned Promotable populates Gate with the declared gate name
    Given an evaluation profile with a clean informative gate
    When PromotableGates is invoked
    Then each returned Promotable has Gate matching the declared gate name

  @PRGTP-B07 @unit-level
  Scenario: Returned Promotable records Passou equal to passed node count
    Given an evaluation profile with a clean informative gate having twelve passed nodes
    When PromotableGates is invoked
    Then the returned Promotable has Passou equal to twelve

  @PRGTP-B08 @unit-level
  Scenario: PromotableGates returns candidates sorted deterministically by name
    Given an evaluation profile with multiple clean informative gates in arbitrary order
    When PromotableGates is invoked
    Then the returned gates are ordered alphabetically by gate name

  @PRGTP-B09 @unit-level
  Scenario: Multiple clean informative gates are all collected
    Given an evaluation profile with three clean informative gates
    When PromotableGates is invoked
    Then all three gates are returned in the promotable list

  @PRGTP-I01 @unit-level
  Scenario: An informative gate is candidate if and only if non-blocking zero failures and positive passes
    Given various gate configurations covering combinations of blocking, pass count, and failure count
    When PromotableGates filters the gates
    Then only gates that are non-blocking with zero failures and positive passes are included

  @PRGTP-I02 @unit-level
  Scenario: A gate with zero passes is never classified as clean
    Given an informative gate that only recorded skips and zero passes
    When PromotableGates evaluates the profile
    Then the gate is omitted from promotion suggestions

  @PRGTP-X01 @unit-level
  Scenario: Does not automatically promote gates or modify anchors yaml
    Given an evaluation profile containing promotable gates
    When PromotableGates is invoked
    Then it produces candidate suggestions without modifying project configuration files

  @PRGTP-X02 @unit-level
  Scenario: Does not evaluate gate execution results directly from disk
    Given an evaluation profile constructed in memory
    When PromotableGates is invoked
    Then it derives candidates solely from the profile without reading disk artifacts
