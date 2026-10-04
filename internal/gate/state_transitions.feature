# language: en
# @anchors
#   ref: VTRST
#   updated_at: 2026-10-03
#   layer: feature

@VTRST
Feature: StateTransitions — every change of a visual unit is proven through what the screen shows

  @VTRST-B01 @unit-level
  Scenario: What has nothing to confront leaves both gates without a verdict
    Given a spec node and a spec with no code
    When both gates confront each
    Then each leaves without a verdict

  @VTRST-B02 @unit-level
  Scenario: Every validation is the trigger of a transition
    Given validations LOGIN-V01 and LOGIN-V02 and presentation validation LOGIN-P01, and a State Flow naming LOGIN-V01 alone
    When validation-transitions confronts the unit's code
    Then it fails naming LOGIN-V02 and LOGIN-P01

  @VTRST-B03 @unit-level
  Scenario: A transition from or to an unknown state does not count
    Given a State Flow row from LOGIN-S01 to LOGIN-S09 triggered by LOGIN-V01, and no LOGIN-S09 state
    When validation-transitions confronts the unit's code
    Then it names LOGIN-V01 with both ends

  @VTRST-B04 @unit-level
  Scenario: A validation exempted with a reason is not asked
    Given LOGIN-V02 marked "@no-state: only trims spaces" and LOGIN-P01 marked "@no-state" with no reason
    When validation-transitions confronts the unit's code
    Then LOGIN-V02 is not asked, and LOGIN-P01 is named as an exemption with no reason

  @VTRST-B05 @unit-level
  Scenario: Every error names the message it shows
    Given errors citing LOGIN-M01, citing LOGIN-M09 that no message is, citing nothing, and one marked "@no-message: logged only"
    When error-message-declared confronts the unit's code
    Then it names the one citing nothing and the one citing LOGIN-M09

  @VTRST-E01 @unit-level
  Scenario: A code file with no spec beside it leaves without a verdict
    Given a code file with no spec beside it
    When both gates confront it
    Then each leaves without a verdict, saying so

  @VTRST-B06 @unit-level
  Scenario: Sections nested under another are found
    Given a spec whose "### Validações" and "### Erros / Falhas" sit under "## Rules (Regras de Negócio)"
    When both gates confront the unit's code
    Then they read the validations and the errors

