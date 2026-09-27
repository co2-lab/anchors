# language: en
# @anchors
#   ref: AGRLG
#   updated_at: 2026-09-27
#   layer: feature

@AGRLG
Feature: AgentRoles — who is who in a project, and the capabilities each role carries

  @AGRLG-B01 @unit-level
  Scenario: Every known role presents itself
    Given the known roles
    When they are listed twice
    Then both lists are the same seven roles, each with a title, a sentence of what it does and a capability

  @AGRLG-B02 @unit-level
  Scenario: Only the product owner and the architect decide the product
    Given every known role
    When each is asked whether it decides the product
    Then only product-owner and architect answer yes

  @AGRLG-B03 @unit-level
  Scenario: The structure is the architect's
    Given the architect and the product owner
    When each is asked whether it decides the structure
    Then the architect answers yes and the product owner no

  @AGRLG-B04 @unit-level
  Scenario: QA, reviewers and dev execute and do not decide
    Given the roles qa, reviewer, reviewer-security, reviewer-performance and dev
    When each is asked about deciding and reviewing
    Then none decides the product or the structure and all can review

  @AGRLG-B05 @unit-level
  Scenario: Each reviewing role has its own lens
    Given the roles qa, reviewer, reviewer-security and reviewer-performance, and the dev
    When their lenses are read
    Then the four lenses are non-empty and distinct, and the dev has none

  @AGRLG-B06 @unit-level
  Scenario: Typed roles are recognised by name and abbreviation
    Given the typed names "po", "arquiteto", "segurança", "perf", " DEV " and "tech-lead"
    When they are parsed
    Then they are product-owner, architect, reviewer-security, reviewer-performance, dev and nothing

  @AGRLG-B07 @unit-level
  Scenario: A role lists its capabilities in alphabetical order
    Given the architect, whose capabilities are declared out of alphabetical order
    When its capabilities are listed
    Then they come in alphabetical order, and so do every known role's

  @AGRLG-X01 @unit-level
  Scenario: An unknown role can do nothing
    Given the roles "", "tech-lead", "gerente" and "DEV"
    When each is asked about every capability
    Then none has any capability
