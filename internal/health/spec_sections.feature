# language: en
# @anchors
#   code: SSFSP
#   ref: SPSCS
#   updated_at: 2026-10-03
#   layer: feature

@SPSCS
Feature: SpecSections — the doctor tells when most specs of a layer lack a section a gate needs to see

  @SPSCS-B01 @unit-level
  Scenario: A screen spec without any recommended section is told about all six
    Given one screen spec with none of the recommended sections and no gate declared
    When the spec sections are checked
    Then the findings name navigation, data-contract, data-states, testids, domain and states

  @SPSCS-B02 @unit-level
  Scenario: Exactly half or a minority missing is silent
    Given four business logic specs of which one, and then two, lack the domain section
    When the spec sections are checked
    Then domain is not reported, and it is reported once three of the four lack it

  @SPSCS-B03 @unit-level
  Scenario: A declared gate that is blind gives a warning
    Given a screen spec without a data contract and the dependency-honored gate declared
    When the spec sections are checked
    Then the data-contract finding is a warning secao-ausente

  @SPSCS-B04 @unit-level
  Scenario: An undeclared gate gives an informational recommendation that asks to declare it
    Given a screen spec without navigation and no gate declared
    When the spec sections are checked
    Then the navigation finding is an informational secao-recomendada whose text names route-declared

  @SPSCS-B05 @unit-level
  Scenario: A title in another language counts
    Given a screen spec with a Navigation section titled in English
    When the spec sections are checked
    Then navigation is not reported

  @SPSCS-B06 @unit-level
  Scenario: The header layer wins over the map's
    Given a spec node whose map layer is spec and whose header declares screen
    When the spec sections are checked
    Then the screen-only data-contract section is demanded from it

  @SPSCS-B07 @unit-level
  Scenario: Screen sections are not demanded from other layers
    Given a business logic spec with none of the recommended sections
    When the spec sections are checked
    Then none of navigation, testids, data-contract, data-states or states is reported

  @SPSCS-B08 @unit-level
  Scenario: Three specs lacking a section give one finding with the numbers
    Given three screen specs without navigation
    When the spec sections are checked
    Then there is exactly one navigation finding and its text says 3 of 3

  @SPSCS-B09 @unit-level
  Scenario: A nil map or configuration gives nothing
    Given a nil map, or a nil configuration
    When the spec sections are checked
    Then there is no finding

  @SPSCS-I01 @unit-level
  Scenario: Adding the reported section removes the finding
    Given a screen spec reported for navigation
    When a Navigation section is added and the sections are checked again
    Then navigation is no longer reported

  @SPSCS-X01 @unit-level
  Scenario: No finding names a spec file
    Given three screen specs lacking sections
    When the spec sections are checked
    Then every finding's subject is a section key, never a spec path

  @SPSCS-E01 @unit-level
  Scenario: Specs missing on disk are left out of the counts
    Given a map whose only spec nodes are screen specs missing on disk
    When the spec sections are checked
    Then there is no finding
