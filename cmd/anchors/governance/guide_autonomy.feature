# language: en
# @anchors
#   ref: ATGDT
#   updated_at: 2026-09-26
#   layer: feature

@ATGDT
Feature: AutonomyGuide — what an agent does with what it does not know, by the role declared locally

  @ATGDT-B01 @unit-level
  Scenario: A role that decides the product is told to record its decisions
    Given a project whose local settings declare the role product-owner
    When the autonomy section is rendered
    Then it reads "Your role (Product Owner) decides the direction of this product."
    And it says the decision is written, through `--for-user`

  @ATGDT-B02 @unit-level
  Scenario: A declared role that does not decide is told who decides
    Given a project whose local settings declare the role dev
    When the autonomy section is rendered
    Then it reads "Your role (Dev) does NOT decide the direction of this product"
    And it names the `product-owner` or the `architect` as who decides

  @ATGDT-B03 @unit-level
  Scenario: With no role declared the guide is the closed one
    Given a project with no local settings
    When the autonomy section is rendered
    Then it reads "You did not declare a role" and never "Your role ("
    And it bans asking whoever is running the agent

  @ATGDT-B04 @unit-level
  Scenario: Every profile reads that preparing the environment asks no authorization
    Given each of the roles dev, architect, product-owner, qa and reviewer, and no role at all
    When the autonomy section is rendered
    Then it says preparing the environment does not ask for authorization, naming `doctor --fix`, `anchors settings role <role> --date` and `anchors map build`
    And it names REVERSIBILITY as the ruler of what DOES ask for authorization

  @ATGDT-B05 @unit-level
  Scenario: A role with a lens reads it
    Given a project whose local settings declare the role qa
    When the autonomy section is rendered
    Then it contains "This role's lens (QA)" followed by the qa lens
    And for the role dev, which has no lens, no lens section appears

  @ATGDT-B06 @unit-level
  Scenario: Whoever does not decide is told not to ask, to move on and what not to escalate
    Given a project that declares the role dev, and one that declares no role
    When the autonomy section is rendered for each
    Then both read "Do not ask whoever is running you.", "move on to the next card" and "What is NOT to be escalated"

  @ATGDT-B07 @unit-level
  Scenario: The old user_issues flag with no role is named as such, never as an empty role
    Given a project that declares only the old flag user_issues true, and one that declares it false
    When the autonomy section is rendered for each
    Then neither reads "Your role" nor "()"
    And the first reads that its declaration by the old user_issues flag says it decides the direction of this product
    And the second reads that it does NOT decide, and "Do not ask whoever is running you."

  @ATGDT-I01 @unit-level
  Scenario: An unreadable declaration reads as no role
    Given a project whose local settings file holds malformed content naming product-owner
    When the autonomy section is rendered
    Then it reads "You did not declare a role" and bans asking

  @ATGDT-X01 @unit-level
  Scenario: A role that decides the product is not forbidden to ask
    Given a project whose local settings declare the role product-owner
    When the autonomy section is rendered
    Then it does not contain "Do not ask whoever is running you."

  @ATGDT-B08 @unit-level
  Scenario: An unreadable settings file is named, and the guide stays closed
    Given a settings file that exists but cannot be parsed
    When the autonomy guide is rendered
    Then it names the read error instead of saying no role was declared, and keeps the closed instructions
